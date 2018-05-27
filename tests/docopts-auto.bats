#!/bin/bash
# vim: set et ts=4 sw=4 sts=4 ft=sh:
#
# unit test with bats
# test sourcing and auto parse header
#

mktmp() {
    tmp=./tmp-docopt_auto_parse.sh
    cat <<'EOF' >$tmp
#!/bin/bash
# Usage: sometest [--opt] [--output=FILE] INFILE

source ../docopts.sh --auto "$@"

echo "${ARGS[INFILE]}"
EOF
    echo $tmp
}

# to be sure to find docopts binray in the $PATH
PATH=..:$PATH

@test "docopt_auto_parse testing internal behavior" {
    # internal
    source ../docopts.sh
    [[ ! -z "$docopt_sh_me" ]]
    mktmp
    [[ ! -z "$tmp" ]]
    unset ARGS
    unset help
    declare -A ARGS

    # auto call help without argument (which is an error and display help)
    run docopt_auto_parse $tmp
    echo "$output"
    echo "status=$status"
    [[ ! -z "$output" ]]
    [[ $status == 1 ]]
    regexp="^echo 'error:"
    [[ "${lines[0]}" =~ $regexp ]]
    [[ "${lines[-1]}" == "exit 64" ]]

    # but runing with -h ==> exit 0
    run docopt_auto_parse $tmp -h
    echo "$output"
    echo "status=$status"
    [[ ! -z "$output" ]]
    [[ $status == 0 ]]
    [[ "${lines[-1]}" == "exit 0" ]]

    # with some options
    run docopt_auto_parse $tmp --opt afilename
    regexp='^ARGS\[[^]]+\]'
    [[ "${lines[0]}" =~ $regexp ]]

    run docopt_auto_parse $tmp afilename
    [[ "${lines[0]}" =~ $regexp ]]
    rm $tmp
}

@test "docopt_auto_parse functionnal testing" {
    mktmp
    [[ -f $tmp ]]
    chmod a+x $tmp
    run $tmp prout
    # echo "$output" >> log
    [[ "$output" == prout ]]
    rm $tmp
}

@test "no source" {
    # test isolation is ok
    [[ -z "${ARGS[*]}" ]]
    [[ -z "$docopt_sh_me" ]]
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

@test "get_raw_value" {
    source ../docopts.sh
    # --num arg is handled as a string
    run docopts -A myargs -h "usage: count --num=<counter>  FILE" : --num=2 file1
    [[ $status -eq 0 ]]
    run get_raw_value myargs --num "$output"
    echo "status=$status"
    echo "lines0=${lines[0]}"
    [[ $status -eq 0 ]]
    [[ ${lines[0]} == "'2'" ]]
    # -v is a counter it is handled as an integer
    run docopts -A myargs -h "Usage: prog -v ..." : -vvv
    run get_raw_value myargs -v "$output"
    [[ $status -eq 0 ]]
    [[ ${lines[0]} == "3" ]]
}
