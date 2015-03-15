#!/bin/bash 
# vim: set et sw=4 ts=4 sts=4:
#
# docopt for bash
#
# Usage: See API_proposal.md
#   source path/to/docopt.sh
#   docopt -A args -h "$help" -v $version : "$@"
#
# the prefix docopt_* is used to export globals and functions

# compute this file dirpath:
docopt_sh_dir="$(dirname $(readlink -f "${BASH_SOURCE[0]}"))"

# fetch Usage: from the given filname
# usually $0 in the main level script
docopt_get_help_string() {
    local myfname=$1
    # filter the block (all blocks) starting at a "# Usage:" and ending 
    # at an empty line, one level of comment markup is removed
    sed -n -e '/^# Usage:/,/^$/ s/^# \?//p' < $myfname
}

# function wrapper
docopt() {
    #   docopts [options] -h <msg> : [<argv>...]
    # find python docopts
    local libexec="$docopt_sh_dir"
    echo "echo \"# libexec=$libexec\""
    # " fix vim hilight
    # call python parser, require docopt.py
    python3 "$libexec/docopts" "$@"
}

docopt_find_docopts() {
    # docopts is the python wrapper using docopt.py
    # it is now embedded in docopts
    echo ../docopts
    # will do: docopt_sh_dir="$.dirname $.readlink -f "${BASH_SOURCE[0]}"))"
}
