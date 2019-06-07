#!/usr/bin/env bash
usage() {
  cat << EOU
Not a serious example.

Usage:
  calculator_example.sh <value> ( ( + | - | * | / ) <value> )...
  calculator_example.sh <function> <value> [( , <value> )]...
  calculator_example.sh (-h | --help)

Examples:
  calculator_example.sh 1 + 2 + 3 + 4 + 5
  calculator_example.sh 1 + 2 '*' 3 / 4 - 5    # note quotes around '*'
  calculator_example.sh sum 10 , 20 , 30 , 40

Options:
  -h, --help

Example:
  ./calculator_example.sh 30 + 23 - 22
EOU
}

# not needed if docopts is in PATH
PATH=../..:$PATH
version='0.1'
parsed=$(docopts -A args -h "$(usage)" -V $version : "$@")
echo "$parsed"
echo "================================================================================"
eval "$parsed"
