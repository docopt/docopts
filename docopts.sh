#!/bin/bash 
# vim: set et sw=4 ts=4 sts=4:
#
# docopt for bash
#
# Usage: See API_proposal.md


# fetch Usage: from the given filname
# usually $0 in the main level script
docopt_get_help_string() {
    local myfname=$1
    # filter the block (all blocks) starting at a "# Usage:" and ending 
    # at an empty line, one level of comment markup is removed
    sed -n -e '/^# Usage:/,/^$/ s/^# \?//p' < $myfname
}

#
docopt() {
    #   docopts [options] -h <msg> : [<argv>...]

    # find python docopts
    local libexec=$(dirname $0)
    echo "echo \"dummy $libexec\""
}


docopt_find_docopts() {
    # docopts is the python wrapper using docopt.py
    # it is now embedded in docopts
    echo ../docopts
}
