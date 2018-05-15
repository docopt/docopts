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

@test "docopt_auto_parse testing internal" {
    # internal
    source ../docopts.sh
    [[ ! -z "$docopt_sh_me" ]]
    mktmp
    [[ ! -z "$tmp" ]]
    unset args
    unset help
    declare -A args
    # auto call help without argument
    run docopt_auto_parse $tmp
    [[ ! -z "$output" ]]
    regexp="^echo 'Usage:"
    [[ "${lines[0]}" =~ $regexp ]]
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

