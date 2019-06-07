#!/usr/bin/env bash
#
# unit test for helper for building the README.md
#

source ../build_doc.sh

@test "get_build_doc" {
  run get_build_doc input_for_build_doc.txt
  echo "$output"
  [[ ${lines[0]} == "4@8@@include 'link/to/file' 2" ]]
  [[ ${lines[1]} == "14@18@@some command" ]]
}

@test "to_filename" {
  run to_filename "awk '{print \$1}'"
  echo "$output"
  [[ $output == 'awkprint1' ]]
}

@test "markdown_extract_link" {
  run markdown_extract_link "[text](some/path/to/file)"
  echo "$output"
  [[ $output == 'some/path/to/file' ]]
  run markdown_extract_link "[text(some/path/to/file)"
  echo "$output"
  [[ $output == 'no_match' ]]
}

@test "strpos" {
  run strpos "[text](some/path/to/file)" ']'
  echo "$output"
  [[ $output -eq 5 ]]
  run strpos "[text(some/path/to/file)" pipo
  echo "$output"
  [[ $output == '-1' ]]
}

tester_get_usage() {
  cat << EOT | get_usage
some text
Usage:
  usage get_usage line 1
  usage get_usage line 2

empty line above
some more text
EOT
}

@test "get_usage" {
  run tester_get_usage
  echo "$output"
  [[ $(echo "$output" | wc -l) -eq 2 ]]
  [[ ${lines[0]} == "  usage get_usage line 1" ]]
}

@test "include" {
  run include "some file"
  echo "$output"
  [[ $status -ne 0 ]]
  run include "no_match" 42
  echo "$output"
  [[ $status -eq 0 ]]
  [[ $output == "include: file not matched at line '42'" ]]
  run include input_for_build_doc.txt 0
  echo "$output"
  [[ $status -eq 0 ]]
  [[ ${lines[0]} == "line 1" ]]
}
