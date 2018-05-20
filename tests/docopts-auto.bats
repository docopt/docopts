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

echo "${args[INFILE]}"
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
    unset args
    unset help
    declare -A args

    # auto call help without argument (which is an error and display help)
    run docopt_auto_parse $tmp
    echo "$output"
    echo "status=$status"
    [[ ! -z "$output" ]]
    [[ $status == 1 ]]
    regexp="^echo 'error:"
    [[ "${lines[0]}" =~ $regexp ]]
    [[ "${lines[2]}" == "exit 64" ]]

    # but runing with -h ==> exit 0
    run docopt_auto_parse $tmp -h
    echo "$output"
    echo "status=$status"
    [[ ! -z "$output" ]]
    [[ $status == 0 ]]
    [[ "${lines[1]}" == "exit 0" ]]

    # with some options
    run docopt_auto_parse $tmp --opt afilename
    regexp='^args\[[^]]+\]'
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
    [[ -z "${args[*]}" ]]
    [[ -z "$docopt_sh_me" ]]
}

