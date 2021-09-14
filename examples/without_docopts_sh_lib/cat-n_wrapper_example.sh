
usage() {
  cat << EOU

cat -n all files

Usage: cat-n_wrapper_example.sh [--count=N] FILE...

Arguments:
  FILE     input file, - is converted to stdin

Options:
  --count=N   limit the number of line to display

Examples:
   ./cat-n_wrapper_example.sh --count=3 cat-n_wrapper.sh  quick_example.sh
   # read from standard input with -
   ./cat-n_wrapper_example.sh --count=3 - < cat-n_wrapper_example.sh
EOU
}

# no PATH changes required if docopts binary is in the PATH already
PATH=../..:$PATH
help=$(usage)
version='0.2'

parsed=$(docopts -A args -h "$help" -V $version : "$@")
#echo "$parsed"
eval "$parsed"

cat_limit() {
    if [[ -z "${args[--count]}" ]] ; then
        cat -n "$1"
    else
        cat -n "$1" | head -"${args[--count]}"
    fi
}

# current docopts multiple argument wrapper

# use intermediate variable FILE_DASH to help vim syntax highlighting
# you can put FILE,# as assoc key too.
FILE_DASH='FILE,#'
n=${args[$FILE_DASH]}
for i in $(seq 0 $(($n - 1)))
do
    f="${args[FILE,$i]}"
    if [[ $f == '-' ]] ; then
      f=/dev/stdin
    fi
    echo "----- $f -------"
    cat_limit  "$f"
done
