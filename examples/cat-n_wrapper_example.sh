#!/bin/bash
#
# cat -n all files
#
# Usage: cat-n_wrapper_example.sh [--count=N] FILE...
#
# Arguments:
#   FILE     input file
# 
# Options:
#   --count=N   limit the number of line to display
#
# Examples:
#    ./cat-n_wrapper.sh --count=3 cat-n_wrapper.sh  quick_example.sh 
#

source ../docopts.sh
help=$(docopt_get_help_string $0)
version='0.1'

parsed=$(docopt -A args -h "$help" -V $version : "$@")
echo "$parsed"
eval "$parsed"

cat_limit() {
    if [[ -z "${args[--count]}" ]] ; then
        cat -n "$1"
    else
        cat -n "$1" | head -"${args[--count]}"
    fi
}

# current docopts multiple argument wrapper
n=${args[FILE,#]}
for i in $(seq 0 $(($n - 1)))
do
    f="${args[FILE,$i]}"
    echo "----- $f -------"
    cat_limit  "$f"
done
