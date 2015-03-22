#!/bin/bash
# Usage: doit [cmd] FILE...
#
# do somethingdo something

source ../docopts.sh --auto "$@"

for a in ${!args[@]} ; do
    echo "$a = ${args[$a]}"
done
