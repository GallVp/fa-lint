# FA-LINT

`fa-lint` is a Fasta linter/validator inspired by [py_fasta_validator](https://github.com/linsalrob/py_fasta_validator) and [SeqKit](https://bioinf.shenwei.me/seqkit). It adheres to the following rules,

1. The Fasta must not be empty
2. The Fasta must start with `>`
3. Each sequence identifier must be unique within the fasta file

4. Each header line starts with a `>`. The header precedes the sequence.
5. Every other line is considered a sequence line
6. Sequence lines may not contain whitespace, numbers, or non-sequence characters. In other words, they must only contain the characters [A-Z] and [a-z]
7. Sequence lines can end with a new line or return depending on whether you have edited this file on a mac, pc, or linux machine.
8. Sequence lines can not be empty.
9. The sequence identifier is the string of characters in the header line following the `>` and up to the first whitespace. Everything after the first whitespace is descriptive, and can be as long as you like