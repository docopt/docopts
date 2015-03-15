#!/bin/bash
# vim: set ts=4 sw=4 sts=4 ft=sh
#
# unit test with bats
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

@test "docopt_find_docopts" {
    run docopt_find_docopts
    [[ ! -z "$output" ]]
    docopts=$output
    [[ -f $docopts ]]
    #$ python3 docopts
    # stderr
    #Usage:
    #  docopts [options] -h <msg> : [<argv>...]
    #$ echo $?
    #1
    run python3 $docopts
    #echo "$output" > log
    [[ ${#lines[@]} -eq 2 ]]
    [[ $status -eq 1 ]]
    regexp='^ *docopts \[options\]'
    [[ "${lines[0]}" == 'Usage:' ]]
    [[ "${lines[1]}" =~ $regexp ]]
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
