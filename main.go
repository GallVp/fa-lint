package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

var version = "dev"

func main() {

	// Command line flags
	threads := flag.Int("threads", 4, "Number of threads to use")
	versionFlag := flag.Bool("version", false, "Show version")
	fastaFile := flag.String("fasta", "", "Fasta file to process")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	flag.Parse()

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
		go validateRecord(recordChannel, &wg)
	}

	parseFasta(*fastaFile, recordChannel)

	log.Infof("Waiting for workers to finish...")
	wg.Wait()
	log.Infof("All workers finished processing")
	fmt.Printf("Fasta is valid: %s\n", *fastaFile)
}

type FastaRecord struct {
	ID           string
	Seq          string
	StartsAtLine int
	endsAtLine   int
}

func parseFasta(file string, recordChannel chan<- FastaRecord) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	defer close(recordChannel)

	var fastaIDs = map[string]struct{}{}
	var lineNumber int = 0
	var fastaID string
	var fastaIDLineNumber int

	var seqBuilder strings.Builder

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		newFastaID := validateHeader(lineNumber, line, fastaIDs)

		if newFastaID != "" && len(fastaIDs) > 1 { // Second fasta sequence starts
			recordChannel <- FastaRecord{ID: fastaID, Seq: seqBuilder.String(), StartsAtLine: fastaIDLineNumber, endsAtLine: lineNumber - 1}

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
		seqBuilder.WriteString(line)
	}

	recordChannel <- FastaRecord{ID: fastaID, Seq: seqBuilder.String(), StartsAtLine: fastaIDLineNumber, endsAtLine: lineNumber}

	log.Infof("Parsed %d lines from fasta file", lineNumber)
}

func validateHeader(lineNumber int, line string, fastaIDs map[string]struct{}) string {

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

	_, found := fastaIDs[fastaID]

	if !found {
		fastaIDs[fastaID] = struct{}{}
	} else {
		log.Fatalf("Duplicate fasta ID found near line #%d: %s", lineNumber, line)
	}

	log.Infof("Found fasta ID near line #%d: %s", lineNumber, fastaID)

	return fastaID
}

func validateRecord(recordChannel <-chan FastaRecord, wg *sync.WaitGroup) {
	defer wg.Done()
	for record := range recordChannel {
		log.Infof("Processing record %s with length %d", record.ID, len(record.Seq))
		record.validateRecordImpl()
		log.Infof("Finished processing record %s", record.ID)
	}
}

func (record FastaRecord) validateRecordImpl() {
	if len(record.Seq) == 0 {
		log.Errorf("Empty sequence for record near line #%d: %s", record.StartsAtLine, record.ID)
	}
}
