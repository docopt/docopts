#!/usr/bin/env bash
# vim: set et ts=4 sw=4 sts=4 ft=sh:
#
# functional test for docopts
# run with bats
#

DOCOPTS_BIN=../docopts

@test "pass on single-dash '-' in -G mode" {
    run $DOCOPTS_BIN -G ARGS -h 'Usage: prog dump [-]' : dump -
    echo "$output"
    [[ $status -eq 0 ]]
    [[ "$output" =~ ARGS__= ]]
}

@test "fail Global mode on single-dash '-'" {
    run $DOCOPTS_BIN -h 'Usage: prog dump [-]' : dump -
    echo "$output"
    [[ $status -ne 0 ]]
    regexp='Print_bash_global:Mangling not supported'
    [[ "$output" =~ $regexp ]]
}

@test "pass on single-dash double-dash in --no-mangle mode" {
    usage="
Usage: prog dump [-]
       prog extra [-c] [--] <unparsed>...
"
    run $DOCOPTS_BIN --no-mangle -h "$usage" : extra -c -- arg1 '<value string>'
    echo "$output"
    [[ $status -eq 0 ]]
    expected_regexp='--=true'
    [[ "$output" =~ $expected_regexp ]]
    expected_regexp='-=false'
    [[ "$output" =~ $expected_regexp ]]
}

@test "skip double-dash in Global mode" {
    usage="
Usage: prog extra [-c] [--] <unparsed>...
"
    run $DOCOPTS_BIN -h "$usage" : extra -- arg1 '<value string>'
    echo "$output"
    #c=false
    #unparsed=('arg1' '<value string>')
    #extra=true
    [[ $status -eq 0 ]]
    expected_regexp='unparsed=\([^)]+'
    [[ "$output" =~ $expected_regexp ]]
    [[ $(echo "$output" | wc -l) -eq 3 ]]
}

@test "-A mode still parse all single-dash double-dash" {
    usage="
Usage: prog dump [-]
       prog extra [-c] [--] <unparsed>...
"
    run $DOCOPTS_BIN -A ARGS -h "$usage" : extra -c -- arg1 '<value string>'
    echo "$output"
    [[ $status -eq 0 ]]
    expected_regexp="ARGS\\['--'\\]=true"
    [[ "$output" =~ $expected_regexp ]]
    expected_regexp="ARGS\\['-'\\]=false"
    [[ "$output" =~ $expected_regexp ]]
    [[ $(echo "$output" | wc -l) -eq 9 ]]
}

# bug #53 non empty array
@test "multiple length argument returns a 0 sized array" {
    # Global mode
    run $DOCOPTS_BIN -h 'Usage: prog [NAME...]' :
    [[ $status -eq 0 ]]
    expected_regexp='NAME=\(\)'
    [[ "$output" =~ $expected_regexp ]]
    [[ ${#lines[@]} -eq 1 ]]

    # also with -G
    run $DOCOPTS_BIN -G ARGS -h 'Usage: prog [NAME...]' :
    [[ $status -eq 0 ]]
    expected_regexp='ARGS_NAME=\(\)'
    [[ "$output" =~ $expected_regexp ]]
    [[ ${#lines[@]} -eq 1 ]]
}
