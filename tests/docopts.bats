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
