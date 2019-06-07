#!/usr/bin/env bash
# Not a serious example.
#
# Usage:
#   calculator_example.sh <value> ( ( + | - | * | / ) <value> )...
#   calculator_example.sh <function> <value> [( , <value> )]...
#   calculator_example.sh (-h | --help)
#
# Examples:
#   calculator_example.sh 1 + 2 + 3 + 4 + 5
#   calculator_example.sh 1 + 2 '*' 3 / 4 - 5    # note quotes around '*'
#   calculator_example.sh sum 10 , 20 , 30 , 40
#
# Options:
#   -h, --help
#
# Example:
#   ./calculator_example.sh 30 + 23 - 22
#

source ../docopts.sh
# not needed if docopts is in PATH
PATH=..:$PATH

help=$(docopt_get_help_string $0)
version='0.1'
parsed=$(docopts -A args -h "$help" -V $version : "$@")
echo "$parsed"
eval "$parsed"
