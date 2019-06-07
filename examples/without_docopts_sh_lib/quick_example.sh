#!/usr/bin/env bash

usage() {
  cat << EOU
Usage:
  quick_example.sh tcp <host> <port> [--timeout=<seconds>]
  quick_example.sh serial <port> [--baud=9600] [--timeout=<seconds>]
  quick_example.sh -h | --help | --version

Examples:
 ./quick_example.sh tcp remote-node 80 --timeout=120
 ./quick_example.sh serial 123 --timeout=120
EOU
}

# if docopts is in PATH, not needed.
PATH=../..:$PATH

version='0.1.1rc'
parsed=$(docopts -A myargs -h "$(usage)" -V $version : "$@")
eval "$parsed"

# main code
for a in ${!myargs[@]} ; do
    echo "$a = ${myargs[$a]}"
done
