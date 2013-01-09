#!/bin/bash
# 
# Testee script for docopts.  This script reads an arbitrary docstring from
# standard input and uses it to parse whatever arguments are passed to it
# into a Bash 4 associative array, which is then dumped in JSON format.
# 
# Pass this file as an argument to `language_agnostic_tester.py` to test
# a `docopts` file located in the same directory. As of 2013-01-09, docopts
# fails the Naval Fate test, as there is no way to determine from just the
# array keys if an option is repeatable or accepts an integer argument:
# both `--speed=2` and `--speed --speed` map to `"--speed": 2`.
# 
# There is currently no way to automatically test the operation mode of
# docopts that name-mangles elements into Bash variables, as this
# transformation cannot be deterministically reversed into a format
# language_agnostic_tester.py expects.
# 
# To test docopts with different Python versions, set the `PYTHON` variable:
# 
#     PYTHON=/usr/bin/python3.2 language_agnostic_tester.py testee.sh
# 
# Note that `language_agnostic_tester.py` itself is only compatible with
# Python 2.7.

script=$(${PYTHON:-python} ./docopts -A args -h - : "$@" < /dev/stdin)

if [[ $(tail -n 1 <<< "$script") =~ ^exit\ [0-9]+$ ]] ; then
    echo '"user-error"'; exit
fi

shopt -s extglob
eval "$script"

echo -n '{'
for key in "${!args[@]}" ; do
    # if the key is not part of a fake nested array,
    # print it as-is
    if [[ -z "${args[${key%,*},#]}" ]] ; then
        [[ -z $sep ]] && sep=, || echo $sep
        value=${args[$key]}
        case "$value" in
            '')         echo -n "\"$key\": null";;
            +([0-9]))   if [[ "${key^^}" == "$key" ]] || [[ "$key" == \<?*\> ]]
                        then
                            echo -n "\"$key\": \"$value\""
                        else
                            echo -n "\"$key\": $value"
                        fi;;
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
