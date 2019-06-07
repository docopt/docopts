#!/usr/bin/env bash
#
# sample of script using auto-parser.
#
# Usage: docopts_auto_example.sh [cmd] FILE...
#
# do something
# do something
#
# Example:
#   ./docopts_auto_example.sh cmd file1 file2
#   ./docopts_auto_example.sh cmd file1 file2
#   ./docopts_auto_example.sh -h

# Auto parse needs an empty line after the top comment above ^^^

# if docopts is in PATH no need to change it.
PATH=..:$PATH

# auto parse this file.
# That's all!
source ../docopts.sh --auto "$@"

# main code based on $ARGS produced by --auto
for a in ${!ARGS[@]} ; do
    echo "$a = ${ARGS[$a]}"
done
