#!/usr/bin/env bats

load 'libs/bats-support/load'
load 'libs/bats-assert/load'

# Missing fasta
@test "No fasta argument should fail" {
  run ./bin/fa-lint
  echo "$output"
  assert_output --partial 'Please provide a fasta file using -fasta'
  [ "$status" -eq 1 ]
}

@test "No fasta file should fail" {
  run ./bin/fa-lint -fasta test.fasta
  echo "$output"
  assert_output --partial 'open test.fasta: no such file or directory'
  [ "$status" -eq 1 ]
}

# Testing a good fasta file
@test "Good fasta file should return 0" {
  run ./bin/fa-lint -fasta test/fasta/good.fasta
  echo "$output"
  assert_output --partial 'Fasta is valid'
  [ "$status" -eq 0 ]
}

@test "Good mixed fasta file should return 0" {
  run ./bin/fa-lint -fasta test/fasta/good_mixed.fasta
  echo "$output"
  assert_output --partial 'Fasta is valid'
  [ "$status" -eq 0 ]
}

@test "Good multiline fasta file should return 0" {
  run ./bin/fa-lint -fasta test/fasta/good_multiline.fasta
  echo "$output"
  assert_output --partial 'Fasta is valid'
  [ "$status" -eq 0 ]
}

# Invalid fasta file cases
@test "Empty fasta file should fail" {
  run ./bin/fa-lint -fasta test/fasta/empty.fasta
  echo "$output"
  assert_output --partial 'Fasta file is empty'
  [ "$status" -eq 1 ]
}

@test "Empty fasta file with a single empty line should fail" {
  run ./bin/fa-lint -fasta test/fasta/empty2.fasta
  echo "$output"
  assert_output --partial 'Fasta file is empty'
  [ "$status" -eq 1 ]
}

@test "Empty fasta file with a space should fail" {
  run ./bin/fa-lint -fasta test/fasta/empty3.fasta
  echo "$output"
  assert_output --partial 'Fasta file must start with a header line'
  [ "$status" -eq 1 ]
}

@test "Fasta file with no > on the first line should fail" {
  run ./bin/fa-lint -fasta test/fasta/no_first_line.fasta
  echo "$output"
  assert_output --partial 'Fasta file must start with a header line'
  [ "$status" -eq 1 ]
}

@test "Fasta file with duplicate ids should fail" {
  run ./bin/fa-lint -fasta test/fasta/duplicates_no_spaces.fasta
  echo "$output"
  assert_output --partial 'Duplicate fasta ID found near line'
  [ "$status" -eq 1 ]
}

@test "Fasta file with duplicate ids and spaces should fail" {
  run ./bin/fa-lint -fasta test/fasta/duplicates_with_spaces.fasta
  echo "$output"
  assert_output --partial 'Duplicate fasta ID found near line'
  [ "$status" -eq 1 ]
}

@test "Fasta file with spaces in the sequence should fail" {
  run ./bin/fa-lint -fasta test/fasta/space.fasta
  echo "$output"
  assert_output --partial 'Invalid sequence character near line'
  [ "$status" -eq 1 ]
}

@test "Fasta file with embedded PHP code should fail" {
  run ./bin/fa-lint -fasta test/fasta/code.fasta
  echo "$output"
  assert_output --partial 'Invalid sequence character near line #80'
  [ "$status" -eq 1 ]
}

@test "Fasta file with just an empty sequence should fail" {
  run ./bin/fa-lint -fasta test/fasta/empty_seq.fasta
  echo "$output"
  assert_output --partial 'Empty sequence for record near line #1: 1'
  [ "$status" -eq 1 ]
}

@test "Fasta file with an empty sequence should fail" {
  run ./bin/fa-lint -fasta test/fasta/empty_seq2.fasta
  echo "$output"
  assert_output --partial 'Empty sequence for record near line #3: 2'
  [ "$status" -eq 1 ]
}

@test "Fasta with non uniform line wrapping should fail" {
  run ./bin/fa-lint -fasta test/fasta/non_uniform_wrapping.fasta
  echo "$output"
  assert_output --partial 'Sequence near line #11 violates preceding line wrapping length'
  [ "$status" -eq 1 ]
}

@test "Fasta with a completely masked sequence should fail" {
  run ./bin/fa-lint -fasta test/fasta/completely_masked.fasta
  echo "$output"
  assert_output --partial "Sequence near line #4 is comprised entirely of N, n, '.' or '*': a2"
  [ "$status" -eq 1 ]
}

@test "Zipped Fasta with a completely masked sequence should fail" {
  run ./bin/fa-lint -fasta test/fasta/completely_masked.fasta.gz
  echo "$output"
  assert_output --partial "Sequence near line #4 is comprised entirely of N, n, '.' or '*': a2"
  [ "$status" -eq 1 ]
}

@test "Good zipped fasta file should return 0" {
  run ./bin/fa-lint -fasta test/fasta/good.fasta.gz
  echo "$output"
  assert_output --partial 'Fasta is valid'
  [ "$status" -eq 0 ]
}

@test "Fasta with a '.' stop codon should fail" {
  run ./bin/fa-lint -fasta test/fasta/stop_codon.fasta
  echo "$output"
  assert_output --partial 'Invalid sequence character near line #2'
  [ "$status" -eq 1 ]
}

@test "Fasta with a '.' stop codon and -s should return 0" {
  run ./bin/fa-lint -s -fasta test/fasta/stop_codon.fasta
  echo "$output"
  assert_output --partial 'Fasta is valid'
  [ "$status" -eq 0 ]
}

@test "Fasta with a '*' stop codon should fail" {
  run ./bin/fa-lint -fasta test/fasta/stop_s_codon.fasta
  echo "$output"
  assert_output --partial 'Invalid sequence character near line #2'
  [ "$status" -eq 1 ]
}

@test "Fasta with a '*' stop codon and -S should return 0" {
  run ./bin/fa-lint -S -fasta test/fasta/stop_s_codon.fasta
  echo "$output"
  assert_output --partial 'Fasta is valid'
  [ "$status" -eq 0 ]
}

@test "Fasta with a anywhere '.' stop codon should fail" {
  run ./bin/fa-lint -fasta test/fasta/anywhere_stop_codon.fasta
  echo "$output"
  assert_output --partial 'Invalid sequence character near line #2'
  [ "$status" -eq 1 ]
}

@test "Fasta with a anywhere '.' stop codon and -s, -a should return 0" {
  run ./bin/fa-lint -a -s -fasta test/fasta/anywhere_stop_codon.fasta
  echo "$output"
  assert_output --partial 'Fasta is valid'
  [ "$status" -eq 0 ]
}

@test "Fasta with a anywhere '*' stop codon should fail" {
  run ./bin/fa-lint -fasta test/fasta/anywhere_s_stop_codon.fasta
  echo "$output"
  assert_output --partial 'Invalid sequence character near line #2'
  [ "$status" -eq 1 ]
}

@test "Fasta with a anywhere '*' stop codon and -S, -a should return 0" {
  run ./bin/fa-lint -a -S -fasta test/fasta/anywhere_s_stop_codon.fasta
  echo "$output"
  assert_output --partial 'Fasta is valid'
  [ "$status" -eq 0 ]
}

@test "Fasta with all '.' stop codons should fail" {
  run ./bin/fa-lint -a -s -fasta test/fasta/all_stop_codon.fasta
  echo "$output"
  assert_output --partial "Sequence near line #1 is comprised entirely of N, n, '.' or '*': seq1_F"
  [ "$status" -eq 1 ]
}

@test "Fasta with '.' at the end of multiple lines should fail" {
  run ./bin/fa-lint -s -fasta test/fasta/multiline_stop_codon.fasta
  echo "$output"
  assert_output --partial 'Invalid sequence character near line #2'
  [ "$status" -eq 1 ]
}

@test "Fasta with '*' at the end of multiple lines should fail" {
  run ./bin/fa-lint -S -fasta test/fasta/multiline_s_stop_codon.fasta
  echo "$output"
  assert_output --partial 'Invalid sequence character near line #2'
  [ "$status" -eq 1 ]
}

@test "Fasta file with lax ID characters should pass non-strict mode" {
  run ./bin/fa-lint -fasta test/fasta/lax_id.fasta
  echo "$output"
  assert_output --partial 'Fasta is valid'
  [ "$status" -eq 0 ]
}

@test "Fasta file with lax ID characters should fail strict mode" {
  run ./bin/fa-lint -w -fasta test/fasta/lax_id.fasta
  echo "$output"
  assert_output --partial "Invalid FASTA ID"
  [ "$status" -eq 1 ]
}
