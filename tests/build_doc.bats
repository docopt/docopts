#!/usr/bin/env bash
#
# unit test for helper for building the README.md
#


source ../build_doc.sh

output_split_lines() {
  local oldIFS=$IFS
  IFS=$'\n'
  local i=0
  local l
  for l in $output
  do
    lines[$i]=$l
    i=$((i+1))
  done
  IFS=$oldIFS
}

setup() {
  # https://github.com/docopt/docopts/issues/39
  if [[ "$OSTYPE" =~ darwin ]] ; then
    skip "build_doc.sh skipped on macOS"
  fi
}

@test "extract_markup" {
  run extract_markup input_for_build_doc.txt

  echo "$output"
  [[ ${lines[0]} == "2 include test_content.txt" ]]
  [[ ${lines[1]} == '14 echo "new outputed content"' ]]
}

@test "parse_input" {
  run parse_input input_for_build_doc.txt
  echo "$output"
  [[ ${lines[0]} == "2@6@include 'test_content.txt'" ]]
  [[ ${lines[1]} == '14@16@echo "new outputed content"' ]]
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

@test "find_end_content" {
  run find_end_content 6 input_for_build_doc.txt
  [[ $output -eq 10 ]]
}

@test "include" {
  # file doesn't exist
  run include "some file" 42
  echo "$output"
  [[ $status -ne 0 ]]
  run include "test_content.txt" 42
  echo "$output"
  [[ $status -eq 0 ]]
  # bats bug #224 (blank line missing in $lines)
  readarray -t lines <<<"$output"
  [[ ${lines[0]} == "[source test_content.txt](test_content.txt)" ]]
  [[ ${lines[1]} == "" ]]
  [[ ${lines[2]} == '```' ]]
  [[ ${lines[3]} == "some include test content"  ]]
}

populate_var_by_printf_v() {
  printf -v $1 "%s" $2
  test $var == "value"
  return $?
}

@test "bats bug fail to printf -v var" {
  run populate_var_by_printf_v var value
  [[ -z $var ]] # should be $var == "value"
}

test_eval_wrapper_helper() {
  eval_wrapper 2 6 "include test_content.txt" our_start our_filename
  echo $our_start
  echo $our_filename
}

@test "eval_wrapper" {
  # it seems that populating variables by printf -v inside function seems not visible by bats
  run test_eval_wrapper_helper
  #  so we read content from stdout
  begin_line=${lines[0]}
  filename=${lines[1]}

  [[ -n $begin_line ]]
  [[ -n $filename ]]
  [[ $begin_line -eq 4 ]]
}

@test "build_sed_cmd" {
  infile=input_for_build_doc.txt
  run build_sed_cmd $infile
  expect=" -e '3,10 d' -e '11 r /tmp/content.includetest_contenttxt' -e '15,20 d' -e '21 r /tmp/content.echonewoutputedcontent'"
  [[ $output == $expect ]]

  # apply
  eval "sed $output $infile" | diff expect_build_doc.txt -
}

@test "extract_markup free character" {
  in_text="[make README.md]: # (grep -F \"():\" \$input)"
  run extract_markup <<< "$in_text"
  echo $output
  [[ $output == '1 grep -F "():" $input' ]]
}
