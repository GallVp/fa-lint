# FA-LINT

[![Build & Test](https://github.com/GallVp/fa-lint/actions/workflows/ci.yml/badge.svg)](https://github.com/GallVp/fa-lint/actions/workflows/ci.yml)
[![Static Badge](https://img.shields.io/badge/Biocontainer-quay.io-blue)](https://quay.io/repository/biocontainers/fa-lint?tab=tags)

[![Anaconda-Server Badge](https://anaconda.org/bioconda/fa-lint/badges/version.svg)](https://anaconda.org/bioconda/fa-lint)
[![Anaconda-Server Badge](https://anaconda.org/bioconda/fa-lint/badges/platforms.svg)](https://anaconda.org/bioconda/fa-lint)
[![Anaconda-Server Badge](https://anaconda.org/bioconda/fa-lint/badges/license.svg)](https://anaconda.org/bioconda/fa-lint)
[![Anaconda-Server Badge](https://anaconda.org/bioconda/fa-lint/badges/downloads.svg)](https://anaconda.org/bioconda/fa-lint)

`fa-lint` is a Fasta linter/validator inspired by [py_fasta_validator](https://github.com/linsalrob/py_fasta_validator) and [SeqKit](https://bioinf.shenwei.me/seqkit). It adheres to the following rules,

1. **Not Empty**: File must not be empty.
2. **Starts with `>`**: First line must begin with `>`.
3. **Header Lines**: Each sequence starts with a header line (`>`).
4. **Identifiers**: The ID is the word after `>` up to the first space. It must:
   - Be unique
   - When `-w` flag is set, it must:
     - Start with a letter (A–Z or a–z)
     - Contain only letters, digits, or underscores (`_`) when `-w` flag is set
5. **Descriptions**: Text after the first space in a header is optional and free-form.
6. **Sequence Lines**:
   - Follow header lines
   - Contain only \[A–Z, a–z] (no whitespace, digits, or other characters)
   - Must not be empty
   - Must use consistent line wrapping
7. **Line Endings**: Can be LF (`\n`) or CRLF (`\r\n`)
8. **Masking**: Sequences cannot be fully masked (i.e., all `N`/`n`)
9. **Stop Codons**:
   - Final `.`/`*` allowed if `-s`/`-S` is set
   - In-frame stops allowed if `-a` is set
   - Sequence must not be entirely stop codons

## Usage

```text
fa-lint:
  -S    Allow stop-codon denoted by '*' as the last character in a sequence
  -a    Allow stop-codons anywhere in the sequence. Use in combination with -s or -S
  -fasta string
        Fasta file to process
  -s    Allow stop-codon denoted by '.' as the last character in a sequence
  -threads int
        Number of threads to use (default 6)
  -verbose
        Enable verbose logging
  -version
        Show version
  -w    Enable strict alphanumeric FASTA ID validation (A-Za-z0-9_ only)
```

## Threads Benchmark

Performed on a Apple M1 Pro 10 CPU 32 GB Machine. Recommended threads 6.

| Fasta size | Threads | Time (sec) |
| :--------: | :-----: | :--------: |
|    850M    |    6    |     11     |
|    850M    |    4    |     14     |
|    850M    |    2    |     26     |
|    850M    |    1    |     47     |
