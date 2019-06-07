#!/usr/bin/env bash

usage() {
  cat << EOU
Naval Fate.

Usage:
  naval_fate.sh ship new <name>...
  naval_fate.sh ship <name> move <x> <y> [--speed=<kn>]
  naval_fate.sh ship shoot <x> <y>
  naval_fate.sh mine (set|remove) <x> <y> [--moored|--drifting]
  naval_fate.sh -h | --help
  naval_fate.sh --version

Options:
  -h --help     Show this screen.
  --version     Show version.
  --speed=<kn>  Speed in knots [default: 10].
  --moored      Moored (anchored) mine.
  --drifting    Drifting mine.EOU
EOU
}

# copied from docopts.sh
# Debug, prints env varible ARGS or $1 formated as a bash 4 assoc array
docopt_print_ARGS() {
    local assoc="$1"
    if [[ -z $assoc ]] ; then
        assoc=ARGS
    fi

    # bash dark magic copying $assoc argument to a local myassoc array
    # inspired by:
    # https://stackoverflow.com/questions/6660010/bash-how-to-assign-an-associative-array-to-another-variable-name-e-g-rename-t#8881121
    declare -A myassoc
    eval $(typeset -A -p $assoc|sed "s/ $assoc=/ myassoc=/")

    # loop on keys
    echo "docopt_print_ARGS => $assoc"
    local a
    for a in ${!myassoc[@]} ; do
        printf "%20s = %s\n" $a "${myassoc[$a]}"
    done
}

VERSION='Naval Fate 2.0'

# if docopts is in PATH, not needed.
PATH=../..:$PATH
# auto parse the header above, See: docopt_get_help_string
eval "$(docopts -A ARGS -h "$(usage)" -V "$VERSION" : "$@")"

docopt_print_ARGS
