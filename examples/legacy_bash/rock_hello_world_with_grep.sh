#!/usr/bin/env bash
# Example from README.md
PATH=$PATH:../..

#? rock 0.1.0
#? Copyright (C) 200X Thomas Light
#? License RIT (Robot Institute of Technology)
#? This is free software: you are free to change and redistribute it.
#? There is NO WARRANTY, to the extent permitted by law.

##? Usage: rock [options] <argv>...
##?
##? Options:
##?       --help     Show help options.
##?       --version  Print program version.

help=$(grep "^##?" "$0" | cut -c 5-)
version=$(grep "^#?"  "$0" | cut -c 4-)
eval "$(docopts -h "$help" -V "$version" : "$@")"

for arg in "${argv[@]}"; do
    echo "$arg"
done
