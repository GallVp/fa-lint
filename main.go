package main

import (
	"bufio"
	"compress/gzip"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

var version = "dev"

func main() {

	// Command line flags
	threads := flag.Int("threads", 6, "Number of threads to use")
	versionFlag := flag.Bool("version", false, "Show version")
	fastaFile := flag.String("fasta", "", "Fasta[.gz] file to process")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	allowStop := flag.Bool("s", false, "Allow stop-codon denoted by '.' as the last character in a sequence")
	allowSStop := flag.Bool("S", false, "Allow stop-codon denoted by '*' as the last character in a sequence")
	allowAnyStop := flag.Bool("a", false, "Allow stop-codons anywhere in the sequence. Use in combination with -s or -S")
	wStrict := flag.Bool("w", false, "Enable strict alphanumeric FASTA ID validation (A-Za-z0-9_ only)")
	flag.Parse()

	if *allowStop && *allowSStop {
		log.Fatal("-s and -S are mutually exclusive. Specify -s or -S")
	}

	if *allowAnyStop && !(*allowStop || *allowSStop) {
		log.Fatal("-a can only be used in combination with -s or -S")
	}

	if *versionFlag {
		fmt.Println(version)
		return
	}

	// Set up logging
	log.SetLevel(log.ErrorLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	// Validate input
	if *fastaFile == "" {
		log.Fatal("Please provide a fasta file using -fasta")
	}

	if *threads < 1 {
		log.Fatal("Number of threads must be at least 1")
	}

	if *threads > runtime.NumCPU() {
		log.Warnf("Requested %d threads, but only %d are available. Using %d threads.", *threads, runtime.NumCPU(), runtime.NumCPU())
		*threads = runtime.NumCPU()
	}

	log.Infof("Using %d threads", *threads)
	log.Infof("Processing fasta file: %s", *fastaFile)

	// Set the number of threads
	runtime.GOMAXPROCS(*threads)
	numWorkers := *threads

	recordChannel := make(chan FastaRecord, 100)

	var wg sync.WaitGroup
	for range numWorkers {
		wg.Add(1)
		go validateRecord(*allowStop, *allowSStop, *allowAnyStop, recordChannel, &wg)
	}

	parseFasta(*fastaFile, *wStrict, recordChannel)

	log.Infof("Waiting for workers to finish...")
	wg.Wait()
	log.Infof("All workers finished processing")
	fmt.Printf("Fasta is valid: %s\n", *fastaFile)
}

type FastaRecord struct {
	ID                 string
	Seq                string
	StartsAtLine       int
	lineWrappingLength int
}

func parseFasta(file string, wStrict bool, recordChannel chan<- FastaRecord) {
	var err error
	var fHandle *os.File

	fHandle, err = os.Open(file)
	errorOutIf(err)
	defer fHandle.Close()

	var ioReader io.Reader = fHandle

	if strings.HasSuffix(file, ".gz") {
		log.Infof("Detected gzipped fasta file: %s", file)

		var gzReader *gzip.Reader
		gzReader, err = gzip.NewReader(fHandle)
		errorOutIf(err)
		defer gzReader.Close()

		ioReader = gzReader
	}

	defer close(recordChannel)

	var fastaIDs = map[string]struct{}{}
	var lineNumber int = 0
	var fastaID string
	var fastaIDLineNumber int
	var lineWrappingLength int = 0

	var seqBuilder strings.Builder

	reader := bufio.NewReaderSize(ioReader, math.MaxInt32)
	for {

		line, isPrefix, err := reader.ReadLine()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			log.Fatalf("Error reading fasta file: %v", err)
		}

		if isPrefix {
			log.Fatalf("Fasta contains a line longer than %d. Please split it into smaller lines", math.MaxInt32)
		}

		lineStr := string(line)

		lineNumber++

		newFastaID := validateHeader(lineNumber, lineStr, fastaIDs, wStrict)

		if newFastaID != "" && len(fastaIDs) > 1 { // Second fasta sequence starts
			recordChannel <- FastaRecord{ID: fastaID, Seq: seqBuilder.String(), StartsAtLine: fastaIDLineNumber, lineWrappingLength: lineWrappingLength}

			fastaID = newFastaID
			fastaIDLineNumber = lineNumber

			seqBuilder.Reset()
			continue
		}

		if newFastaID != "" { // First fasta sequence starts
			fastaID = newFastaID
			fastaIDLineNumber = lineNumber
			continue
		}

		// Fasta sequence continues
		lineWrappingLength = max(lineWrappingLength, len(line))
		if seqBuilder.Len() < 1 {
			seqBuilder.WriteString(lineStr)
			continue
		}
		seqBuilder.WriteString("\n" + lineStr)
	}

	if lineNumber == 0 {
		log.Fatal("Fasta file is empty")
	}

	recordChannel <- FastaRecord{ID: fastaID, Seq: seqBuilder.String(), StartsAtLine: fastaIDLineNumber, lineWrappingLength: lineWrappingLength}

	log.Infof("Parsed %d lines from fasta file", lineNumber)
}

func validateHeader(lineNumber int, line string, fastaIDs map[string]struct{}, wStrict bool) string {

	if !strings.HasPrefix(line, ">") && len(fastaIDs) == 0 {
		log.Fatal("Fasta file must start with a header line")
	}

	if !strings.HasPrefix(line, ">") {
		return ""
	}

	fastaIDField := strings.Fields(line)[0]

	if len(fastaIDField) < 2 {
		log.Fatalf("Fasta header does not contain a valid ID near line #%d: %s", lineNumber, line)
	}

	fastaID := fastaIDField[1:]

	strictIDPattern := regexp.MustCompile(`^[A-Za-z]{1}\w+$`)
	if wStrict && !strictIDPattern.MatchString(fastaID) {
		log.Fatalf("Invalid FASTA ID '%s' near line #%d: IDs must match [A-Za-z][A-Za-z0-9_]", fastaID, lineNumber)
	}

	_, found := fastaIDs[fastaID]

	if !found {
		fastaIDs[fastaID] = struct{}{}
	} else {
		log.Fatalf("Duplicate fasta ID found near line #%d: %s", lineNumber, line)
	}

	log.Infof("Found fasta ID near line #%d: %s", lineNumber, fastaID)

	return fastaID
}

func validateRecord(allowStop bool, allowSStop bool, allowAnyStop bool, recordChannel <-chan FastaRecord, wg *sync.WaitGroup) {
	defer wg.Done()
	for record := range recordChannel {
		log.Infof("Processing record %s with length %d", record.ID, len(record.Seq))
		record.validateRecordImpl(allowStop, allowSStop, allowAnyStop)
		log.Infof("Finished processing record %s", record.ID)
	}
}

func (record FastaRecord) validateRecordImpl(allowStop bool, allowSStop bool, allowAnyStop bool) {
	if len(record.Seq) == 0 {
		log.Fatalf("Empty sequence for record near line #%d: %s", record.StartsAtLine, record.ID)
	}

	seqLines := strings.Split(record.Seq, "\n")

	log.Infof("Sequence %s has %d line(s)", record.ID, len(seqLines))

	var numberOfInvalidLines = 0
	for i, line := range seqLines {

		isLastLine := (i == len(seqLines)-1)

		if len(line) == 0 {
			log.Fatalf("Empty sequence line near line #%d", record.StartsAtLine+i+1)
		}

		if !isValidSequence(line, isLastLine, allowStop, allowSStop, allowAnyStop) {
			log.Fatalf("Invalid sequence character near line #%d", record.StartsAtLine+i+1)
		}

		if i > 0 && i < len(seqLines)-1 && record.lineWrappingLength != len(line) && len(seqLines[i+1]) != 0 {
			log.Fatalf("Sequence near line #%d violates preceding line wrapping length", record.StartsAtLine+i+1)
		}

		if isAllInvalidSequence(line) {
			numberOfInvalidLines += 1
		}
	}

	if len(seqLines) == numberOfInvalidLines {
		log.Fatalf("Sequence near line #%d is comprised entirely of N, n, '.' or '*': %s", record.StartsAtLine, record.ID)
	}
}

func isValidSequence(seq string, isLastLine bool, allowStop bool, allowSStop bool, allowAnyStop bool) bool {

	var pattern string

	if allowStop && allowAnyStop {
		pattern = `^[A-Za-z\.]+$`
	} else if allowSStop && allowAnyStop {
		pattern = `^[A-Za-z\*]+$`
	} else if allowStop && isLastLine {
		pattern = `^[A-Za-z]+\.?$`
	} else if allowSStop && isLastLine {
		pattern = `^[A-Za-z]+\*?$`
	} else {
		pattern = `^[A-Za-z]+$`
	}

	re := regexp.MustCompile(pattern)
	return re.MatchString(seq)
}

func isAllInvalidSequence(seq string) bool {
	pattern := `^[Nn\.\*]+$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(seq)
}

func errorOutIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
