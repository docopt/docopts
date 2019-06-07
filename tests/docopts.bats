#!/usr/bin/env bash
# vim: set et ts=4 sw=4 sts=4 ft=sh:
#
# unit test for bash helpers in docopts.sh
# run with bats
#

source ../docopts.sh

@test "docopt_get_help_string" {
    tmp=./tmp_docopt_get_help_string
    cat <<EOF > $tmp
#!/usr/bin/env bash
# some test
# Usage: docopt_get_help_string $0
#   docopt_get_help_string prout
#   docopt_get_help_string -

# an empty line above
EOF
    [[ -f $tmp ]]
    run docopt_get_help_string $tmp
    regexp='Usage:'
    [[ "${lines[0]}" =~ $regexp ]]
    [[ ${#lines[@]} -eq 3 ]]

    rm -f $tmp
}

@test "docopt_get_help_string with Options (this test may fail if bats fix their issue #224)" {
    tmp=./tmp_docopt_get_help_string_with_options
    cat <<EOF > $tmp
#!/usr/bin/env bash
# rock
#
# Usage: rock [options] <argv>...
#
# Options:
#       --verbose  Generate verbose messages.
#       --help     Show help options.
#       --version  Print program version.

# an empty line above
EOF
    [[ -f $tmp ]]
    run docopt_get_help_string $tmp

    # truncate log
    # :>log
    #echo '--- output ---' >> log
    #echo "$output" >> log
    #echo '--- lines ---' >> log
    #for line in "${lines[@]}"; do
    #    echo "$line" >> log
    #done
    #echo '---' >> log
    #echo "Line count = ${#lines[@]}" >> log

    # bats bug missing empty line: https://github.com/bats-core/bats-core/issues/224
    # we can rely on bats to split HELP correctly
    line_count=$(wc -l <<<"$output")
    [[ $line_count -eq 6 ]]
    regexp='Options:'
    # index 1 must be 2 when bats issue 224 will be corrected
    [[ "${lines[1]}" =~ $regexp ]]

    rm -f $tmp
}

@test "docopt_get_values" {
    declare -A args
    args['FILE,#']=3
    args['FILE,0']=somefile1
    args['FILE,1']=somefile2
    args['FILE,2']=somefile3

    run docopt_get_values args FILE
    [[ ${#lines[@]} -eq 1 ]]
    [[ "$output" == "somefile1 somefile2 somefile3" ]]
    array=( $output )
    [[ ${#array[@]} -eq 3 ]]
    [[ ${array[2]} == 'somefile3' ]]
}

@test "docopt_get_eval_array" {
    declare -A args
    args['FILE,#']=4
    args['FILE,0']=somefile1
    args['FILE,1']=somefile2
    args['FILE,2']=somefile3
    args['FILE,3']="somefile4 with space inside"

    run docopt_get_eval_array args FILE myarray
    # echo "$output" >> log
    [[ ${#lines[@]} -eq 5 ]]
    eval "$output"
    [[ ${#myarray[@]} -eq 4 ]]
    [[ ${myarray[2]} == 'somefile3' ]]
    [[ ${myarray[3]} == "somefile4 with space inside" ]]
}

@test "docopt_get_raw_value" {
    PATH=..:$PATH
    # --num arg is handled as a string
    run docopts -A myargs -h "usage: count --num=<counter>  FILE" : --num=2 file1
    [[ $status -eq 0 ]]
    run docopt_get_raw_value myargs --num "$output"
    echo "status=$status"
    echo "lines0=${lines[0]}"
    [[ $status -eq 0 ]]
    [[ ${lines[0]} == "'2'" ]]
    # -v is a counter it is handled as an integer
    run docopts -A myargs -h "Usage: prog -v ..." : -vvv
    run docopt_get_raw_value myargs -v "$output"
    [[ $status -eq 0 ]]
    [[ ${lines[0]} == "3" ]]
}

@test "docopt_print_ARGS" {
    declare -A ARGS
    ARGS['FILE,#']=4
    ARGS['FILE,0']=somefile1
    ARGS['FILE,1']=somefile2
    ARGS['FILE,2']=somefile3
    ARGS['FILE,3']="somefile4 with space inside"
    run docopt_print_ARGS
    echo "output=$output"
    grep -q -E 'FILE,3' <<< "$output"

    # with a named assoc
    declare -A ourargs
    ourargs['FILE,#']=4
    ourargs['FILE,0']=somefile1
    ourargs['FILE,1']=somefile2
    ourargs['FILE,2']=somefile3
    ourargs['FILE,3']="somefile4 with space inside"
    run docopt_print_ARGS ourargs
    echo "output=$output"
    grep -q -E 'FILE,3' <<< "$output"
    grep -q -E 'ourargs' <<< "$output"
}

@test "docopt_print_ARGS -G" {
    ARGS_FILE=( somefile1 somefile2 somefile3 "somefile4 with space inside" )
    run docopt_print_ARGS -G
    echo "output=$output"
    grep -q -E 'ARGS_FILE' <<< "$output"

    unset ARGS_FILE

    # with a named prefix
    someprefix_FILE=( somefile1 somefile2 somefile3 "somefile4 with space inside" )
    run docopt_print_ARGS -G someprefix
    echo "output=$output"
    grep -q -E 'someprefix_FILE' <<< "$output"
}

@test "docopt_get_version_string from \$0" {
    tmp=./tmp_docopt_get_version_string
    cat <<EOF > $tmp
#!/usr/bin/env bash
#
# Usage: rock [options] <argv>...
# 
# Options:
#       --verbose  Generate verbose messages.
#       --help     Show help options.
#       --version  Print program version.
# ----
# rock 0.1.0
# Copyright (C) 200X Thomas Light
# License RIT (Robot Institute of Technology)
# This is free software: you are free to change and redistribute it.
# There is NO WARRANTY, to the extent permitted by law.

# an empty line above
EOF
    [[ -f $tmp ]]
    run docopt_get_version_string $tmp
    regexp='^rock'
    echo "${lines[0]}"
    [[ "${lines[0]}" =~ $regexp ]]
    [[ ${#lines[@]} -eq 5 ]]

    rm -f $tmp
}
