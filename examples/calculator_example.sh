#!/bin/bash
# Not a serious example.
# 
# Usage:
#   calculator_example.py <value> ( ( + | - | * | / ) <value> )...
#   calculator_example.py <function> <value> [( , <value> )]...
#   calculator_example.py (-h | --help)
# 
# Examples:
#   calculator_example.py 1 + 2 + 3 + 4 + 5
#   calculator_example.py 1 + 2 '*' 3 / 4 - 5    # note quotes around '*'
#   calculator_example.py sum 10 , 20 , 30 , 40
# 
# Options:
#   -h, --help
# 
# Example:
#   ./calculator_example.sh 30 + 23 - 22
#

source ../docopts.sh

help=$(docopt_get_help_string $0)
version='0.1'
parsed=$(docopt -A args -h "$help" -V $version : "$@")
echo "$parsed"
eval "$parsed"
