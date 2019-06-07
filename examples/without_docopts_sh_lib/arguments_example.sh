#!/usr/bin/env bash

usage() {
  cat << EOU
Usage: arguments_example.sh [-vqrh] [FILE] ...
          arguments_example.sh (--left | --right) CORRECTION FILE

Process FILE and optionally apply correction to either left-hand side or
right-hand side.

Arguments:
  FILE        optional input file
  CORRECTION  correction angle, needs FILE, --left or --right to be present

Options:
  -h --help
  -v       verbose mode
  -q       quiet mode
  -r       make report
  --left   use left-hand side
  --right  use right-hand side
EOU
}

# not needed if docopts already in PATH
PATH=../..:$PATH
eval "$(docopts -A ARGS -h "$(usage)" : "$@")"

# main code
# on assoc array '!' before nane gike hash keys
for a in ${!ARGS[@]} ; do
    echo "$a = ${ARGS[$a]}"
done
