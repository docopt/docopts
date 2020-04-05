#!/usr/bin/env bash
#
# Usage: arguments_example.sh [-vqrh] [FILE] ...
#        arguments_example.sh (--left | --right) CORRECTION FILE
#
# Process FILE and optionally apply correction to either left-hand side or
# right-hand side.
#
# Arguments:
#   FILE        optional input file
#   CORRECTION  correction angle, needs FILE, --left or --right to be present
#
# Options:
#   -h --help
#   -v       verbose mode
#   -q       quiet mode
#   -r       make report
#   --left   use left-hand side
#   --right  use right-hand side
#

# if docopts is in PATH, not needed.
# Note: docopts.sh is also found in PATH
PATH=..:$PATH
# auto parse the header above, See: docopt_get_help_string
source docopts.sh --auto "$@"

# docopt_auto_parse use ARGS bash 4 global assoc array
# main code
# on bash assoc array a '!' before name gives hash keys list
for a in ${!ARGS[@]} ; do
    echo "$a = ${ARGS[$a]}"
done
