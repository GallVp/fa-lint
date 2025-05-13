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

# @test "Fasta file with spaces in the sequence should fail" {
#   run ./bin/fa-lint -fasta test/fasta/space.fasta
#   echo "$output"
#   [ "$status" -ne 0 ]
# }

# @test "Fasta file with duplicate ids and spaces should fail" {
#   run ./bin/fa-lint -fasta test/fasta/duplicates_with_spaces.fasta
#   echo "$output"
#   [ "$status" -ne 0 ]
# }

# @test "Fasta file with embedded PHP code should fail" {
#   run ./bin/fa-lint -fasta test/fasta/code.fasta
#   echo "$output"
#   [ "$status" -ne 0 ]
# }

# @test "Fasta file with one empty sequence should fail" {
#   run ./bin/fa-lint -fasta test/fasta/empty_seq.fasta
#   echo "$output"
#   [ "$status" -ne 0 ]
# }

# @test "Fasta file with an empty sequence should fail" {
#   run ./bin/fa-lint -fasta test/fasta/empty_seq2.fasta
#   echo "$output"
#   [ "$status" -ne 0 ]
# }
