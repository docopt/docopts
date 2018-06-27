
usage() {
  cat << EOU

cat -n all files

Usage: cat-n_wrapper_example.sh [--count=N] FILE...

Arguments:
  FILE     input file

Options:
  --count=N   limit the number of line to display

Examples:
   ./cat-n_wrapper.sh --count=3 cat-n_wrapper.sh  quick_example.sh 
EOU
}


# no PATH changes required if docopts binary is in the PATH already
PATH=../..:$PATH
help=$(usage)
version='0.1'

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
n=${args[FILE,#]}
for i in $(seq 0 $(($n - 1)))
do
    f="${args[FILE,$i]}"
    echo "----- $f -------"
    cat_limit  "$f"
done
