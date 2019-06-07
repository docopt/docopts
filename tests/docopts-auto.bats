#!/usr/bin/env bash
# vim: set et ts=4 sw=4 sts=4 ft=sh:
#
# fonctional tests with bats
# test sourcing and auto parse header
#

# create a sample script that uses docopts.sh
# modify $TMP
# Usage: mktmp [global|assoc]
mktmp() {
    case $1 in
    globals|global)
        # invoke auto parse with -G and generate Globals variables
        TMP=./tmp-docopt_auto_parse_globals.sh
        cat <<'EOF' >$TMP
#!/usr/bin/env bash
# Usage: sometest [--opt] [--output=FILE] INFILE

source ../docopts.sh --auto -G "$@"

echo "$ARGS_INFILE"
EOF
    ;;

    assoc|"")
        TMP=./tmp-docopt_auto_parse_assoc.sh
        cat <<'EOF' >$TMP
#!/usr/bin/env bash
# Usage: sometest [--opt] [--output=FILE] INFILE

source ../docopts.sh --auto "$@"

echo "${ARGS[INFILE]}"
EOF
    ;;
    esac
    chmod a+x $TMP
    echo $TMP
}

# to be sure to find docopts binary in the $PATH
PATH=..:$PATH

@test "docopt_auto_parse testing internal behavior" {
    # internal
    source ../docopts.sh
    [[ -n "$(typeset -f docopt_get_help_string)" ]]

    mktmp assoc
    [[ -n "$TMP" ]]
    unset ARGS
    unset HELP
    declare -A ARGS

    # auto call without argument (which is an error) and display help
    run docopt_auto_parse $TMP
    echo "$output"
    echo "status=$status"
    [[ -n "$output" ]]
    [[ $status == 1 ]]
    regexp="^echo 'error:"
    [[ "${lines[0]}" =~ $regexp ]]
    [[ "${lines[-1]}" == "exit 64" ]]

    # but runing with -h ==> exit 0
    run docopt_auto_parse $TMP -h
    echo "$output"
    echo "status=$status"
    [[ -n "$output" ]]
    [[ $status -eq 0 ]]
    [[ "${lines[-1]}" == "exit 0" ]]

    # with some options
    run docopt_auto_parse $TMP --opt afilename
    regexp='^ARGS\[[^]]+\]'
    [[ "${lines[0]}" =~ $regexp ]]

    run docopt_auto_parse $TMP afilename
    [[ "${lines[0]}" =~ $regexp ]]
    rm $TMP
}

@test "docopt_auto_parse -G for globals" {
    source ../docopts.sh
    mktmp globals
    [[ -x "$TMP" ]]
    # temporary script has the expected name
    regexp='_globals.sh'
    [[ "$TMP" =~ $regexp ]]
    unset HELP

    # auto call without argument (which is an error) and display help
    run docopt_auto_parse -G $TMP
    echo "$output"
    echo "status=$status"
    [[ -n "$output" ]]
    [[ $status == 1 ]]
    regexp="^echo 'error:"
    [[ "${lines[0]}" =~ $regexp ]]
    [[ "${lines[-1]}" == "exit 64" ]]

    # but runing with -h ==> exit 0
    run docopt_auto_parse -G $TMP -h
    echo "$output"
    echo "status=$status"
    [[ -n "$output" ]]
    [[ $status -eq 0 ]]
    [[ "${lines[-1]}" == "exit 0" ]]

    # with some options
    run docopt_auto_parse -G $TMP --opt afilename
    echo "$output"
    regexp='^ARGS_(INFILE|opt|output)'
    [[ "${lines[0]}" =~ $regexp ]]

    run docopt_auto_parse -G $TMP afilename
    [[ $status -eq 0 ]]
    echo "$output"
    unset ARGS_INFILE
    eval "$output"
    [[ "$ARGS_INFILE" == 'afilename' ]]
    rm $TMP
}

@test "docopt_auto_parse functional testing associative array" {
    mktmp assoc
    [[ -x $TMP ]]
    run $TMP prout
    [[ "$output" == prout ]]
    rm $TMP
}

@test "docopt_auto_parse -G functional testing globals variables" {
    mktmp global
    [[ -x $TMP ]]
    run $TMP prout_from_global
    [[ "$output" == prout_from_global ]]
    rm $TMP
}

@test "no source" {
    # test isolation is ok
    [[ -z "${ARGS[*]}" ]]
    [[ $(env | grep ARGS_ | wc -l) -eq 0 ]]
}

@test "global eval" {
   eval $(docopts -G ARGS -h "usage: p  FILE..." : one two three)
   [[ $? -eq 0 ]]
   # check content
   for f in ${ARGS_FILE[@]} ; do
       echo "======== $f"
   done
   [[ ${#ARGS_FILE[@]} -eq 3 ]]
   [[ ${ARGS_FILE[0]} == one ]]
   [[ ${ARGS_FILE[1]} == two ]]
   [[ ${ARGS_FILE[2]} == three ]]
}

@test "docopts error" {
    run docopts -h "usage: p [-9] FILE..." : -9 f pipo
    echo "status=$status"
    [[ $status -eq 1 ]]
    run docopts -G ARGS -h "usage: p [-9] FILE..." : -9 f pipo
    echo "status=$status"
    [[ $status -eq 0 ]]
}

@test "--no-declare" {
    run docopts --no-declare -h "usage: cat  FILE..." : file1 file2
    echo "status=$status"
    [[ $status -eq 1 ]]

    run docopts -A myargs -h "usage: cat  FILE..." : file1 file2
    echo "status=$status"
    [[ $status -eq 0 ]]
    [[ ${lines[0]} == "declare -A myargs" ]]

    run docopts -A myargs --no-declare -h "usage: cat  FILE..." : file1 file2
    echo "status=$status"
    echo "lines0=${lines[0]}"
    [[ $status -eq 0 ]]
    [[ ${lines[0]} != "declare -A myargs" ]]
}

