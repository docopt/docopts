#!/bin/bash 
#
# Usage:
#  quick_example.sh tcp <host> <port> [--timeout=<seconds>]
#  quick_example.sh serial <port> [--baud=9600] [--timeout=<seconds>]
#  quick_example.sh -h | --help | --version

libpath=../
source $libpath/docopts.sh

help=$(get_help_string $0)
version='0.1.1rc'

parsed=$(docopt -A arguments "$help" $version)
echo "$parsed"
eval "$parsed"
