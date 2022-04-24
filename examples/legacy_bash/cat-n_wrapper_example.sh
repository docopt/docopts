#!/usr/bin/env bash
#
# cat -n all files
#
# Usage: cat-n_wrapper_example.sh [--count=N] FILE...
#
# Arguments:
#   FILE     input file, if FILE equal - stdin is used instead.
#
# Options:
#   --count=N   limit the number of line to display
#
# Examples:
#    ./cat-n_wrapper.sh --count=3 cat-n_wrapper.sh  quick_example.sh
#

# no PATH changes required if docopts binary is in the PATH already
PATH=../..:$PATH
source ../../docopts.sh
help=$(docopt_get_help_string $0)
version='0.1'

parsed=$(docopts -G args -h "$help" -V $version : "$@")
# Show parsed arguments
#echo "$parsed"
eval "$parsed"

cat_limit() {
  local filename=$1

  if [[ -z "$args_count" ]] ; then
      cat -n "$filename"
  else
      cat -n "$filename" | head -"$args_count"
  fi
}

# array len in bash
n=${#args_FILE[@]}
for i in $(seq 0 $(($n - 1)))
do
    f="${args_FILE[$i]}"
    echo "----- $f ------- $((i+1)) / $n"
    if [[ $f == '-' ]] ; then
      f=/dev/stdin
    fi
    if [[ -f $f ]]
    then
      cat_limit  "$f"
    else
      echo "file not found: $f"
    fi
done
