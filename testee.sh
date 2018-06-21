#!/usr/bin/env bash
#
# Testee script for docopts.  This script reads an arbitrary docopts usage
# from standard input and uses it to parse whatever arguments are passed to it
# into a Bash 4 associative array, which is then dumped in JSON format.
#
# Pass this file as an argument to `language_agnostic_tester.py` to test
# a `docopts` binary located in the same directory.
#
# As of 2018-05-22, docopts may fails the Naval Fate test, as there it is
# difficult to determine from just the array value if an option is repeatable
# counter or accepts an integer argument:
#   both `--speed=2` and `--speed --speed` map to `"--speed": 2`.
# A trick is to read the outputed value of docopts and not evaled result.
# See: docopt_get_raw_value()
#
# There is currently no way to automatically test the operation mode of
# docopts that name-mangles elements into Bash variables, as this
# transformation cannot be deterministically reversed into a format
# language_agnostic_tester.py expects.
#
# Usage:
#  python language_agnostic_tester.py ./testee.sh [ID_OF_THE_TEST]
#  # usage on stdin, args on command line
#  echo "usage: prog (go <direction> --speed=<km/h>)..." | ./testee.sh go left --speed=5  go right --speed=9
#
# To get ID_OF_THE_TEST:
#  grep -E '^\{|"user' testcases.docopt | cat -n | less
#
# Note that `language_agnostic_tester.py` is only compatible with
# Python 2.7.

source ./docopts.sh
script=$(./docopts -A args -h - : "$@" < /dev/stdin)

if [[ $(tail -n 1 <<< "$script") =~ ^exit\ [0-9]+$ ]] ; then
    echo '"user-error"'; exit
fi

shopt -s extglob
eval "$script"

# start JSON
echo -n '{'
regexp="^'[0-9]+'$"
for key in "${!args[@]}" ; do
    # if the key is not part of a fake nested array,
    # print it as-is
    if [[ -z "${args[${key%,*},#]}" ]] ; then
        [[ -z $sep ]] && sep=, || echo $sep
        value=${args[$key]}
        case "$value" in
            '')         echo -n "\"$key\": null";;
            +([0-9]))
              # For numeric value, the JSON is distinct if it is a counter
              # (no quote) or a string (quoted value). But bash can't distiguish
              # any. So we look at the outputed value as text
              if [[ $(docopt_get_raw_value args "$key" "$script") =~ $regexp ]]
              then
                  echo -n "\"$key\": \"$value\""
              else
                  echo -n "\"$key\": $value"
              fi
            ;;
            true|false) echo -n "\"$key\": $value";;
            *)          echo -n "\"$key\": \"$value\"";;
        esac
    # if the key is the length key of a fake nested array,
    # print the whole array
    elif [[ "${key: -2:2}" == ',#' ]] ; then
        [[ -z $sep ]] && sep=, || echo $sep
        key=${key%,*}
        n=${args[$key,#]}
        i=0
        echo -n "\"$key\": ["
        while [[ $i -lt $n ]] ; do
            [[ $i -gt 0 ]] && echo -n ', '
            echo -n "\"${args[$key,$i]}\""
            i=$[$i+1]
        done
        echo -n ']'
    fi
done
echo '}'
