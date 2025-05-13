# FA-LINT

`fa-lint` is a Fasta linter/validator inspired by [py_fasta_validator](https://github.com/linsalrob/py_fasta_validator) and [SeqKit](https://bioinf.shenwei.me/seqkit). It adheres to the following rules,

1. The Fasta must not be empty.
2. The Fasta must start with `>`.
3. Each header line starts with a `>`. The header precedes the sequence.
4. The sequence identifier is the string of characters in the header line following the `>` and up to the first whitespace. Everything after the first whitespace is descriptive, and can be as long as you like.
5. Each sequence identifier must be unique within the fasta file.
6. Every other line is considered a sequence line.
7. Sequence lines may not contain whitespace, numbers, or non-sequence characters. In other words, they must only contain the characters [A-Z] and [a-z].
8. Sequence lines can end with a new line or return depending on whether you have edited this file on a mac, pc, or linux machine.
9. Sequence lines can not be empty.

## Usage

```bash
fa-lint:
  -fasta string
        Fasta file to process
  -threads int
        Number of threads to use (default 6)
  -verbose
        Enable verbose logging
  -version
        Show version
```

## Threads Benchmark

Performed on a Apple M1 Pro 10 CPU 32 GB Machine. Recommended threads 6.

| Fasta size | Threads | Time (sec) |
| :--------: | :-----: | :--------: |
|    55M     |    1    |    1.7     |
|    55M     |    2    |    1.2     |
|    55M     |    4    |    0.9     |
|    55M     |    6    |    0.8     |
|    850M    |    1    |     27     |
|    850M    |    2    |     14     |
|    850M    |    4    |     8      |
|    850M    |    6    |     6      |
