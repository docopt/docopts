#!/usr/bin/env bash

usage() {
  cat << EOU
Usage: counted_example.sh --help
       counted_example.sh -v...
       counted_example.sh go [go]
       counted_example.sh (--path=<path>)...
       counted_example.sh <file> <file>

Try: counted_example.sh -vvvvvvvvvv
     counted_example.sh go go
     counted_example.sh --path ./here --path ./there
     counted_example.sh this.txt that.txt
EOU
}

# if docopts is in PATH, not needed.
PATH=../..:$PATH
eval "$(docopts -A ARGS -h "$(usage)" : "$@")"

# docopt_auto_parse use ARGS bash 4 global assoc array
# main code
# on assoc array '!' before nane gike hash keys
for a in ${!ARGS[@]} ; do
    echo "$a = ${ARGS[$a]}"
done
