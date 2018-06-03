#!/bin/bash
# vim: set et ts=4 sw=4 sts=4 ft=sh:
#
# unit test for bash helpers in docopts.sh
# run with bats
#

source ../docopts.sh

@test "docopt_get_help_string" {
    tmp=./tmp_docopt_get_help_string
    cat <<EOF > $tmp
#!/bin/bash
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
