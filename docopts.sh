#!/bin/bash
# vim: set et sw=4 ts=4 sts=4:
#
# docopts helper for bash
#
# Usage:
#   source path/to/docopts.sh
#   docopts -A ARGS -h "$help" -V $version : "$@"
#
# the prefix docopt_* is used to export globals and functions

# compute this file dirpath:
docopt_sh_me=$(readlink -f "${BASH_SOURCE[0]}")
docopt_sh_dir="$(dirname "$docopt_sh_me")"

# fetch Usage: from the given filename
# usually $0 in the main level script
docopt_get_help_string() {
    local myfname=$1
    # filter the block (/!\ all blocks) starting at a "# Usage:" and ending
    # at an empty line, one level of comment markup is removed
    #
    ## sed -n -e '/^# Usage:/,/\(^# \?----\|^$\)/ { /----/ d; s/^# \?//p }' rock_no-stdin_example.sh

    # -n : no print output
    # -e : pass sed code inline
    # /^# Usage:/,/^$/ : filter range blocks from '# Usage:' to empty line
    #  s/^# \?// : substitute comment marker and an optional space
    #  p : print
    sed -n -e '/^# Usage:/,/^$/ s/^# \?//p' < $myfname
}

# fetch version information from the given filename or string
# usually $0 in the main level script, or the help string extracted
# by docopt_get_help_string()
docopt_get_version_string() {
    if [[ -f "$1" ]] ; then
        # filter the block (all blocks) starting at a "# Usage:" and ending
        # at an empty line, one level of comment markup is removed
        sed -n -e '/^# ----/,/^$/ s/^# \?//p' < "$1"
    else
        # use docopts --separator behavior
        echo "$1"
    fi
}

## function wrapper
## Usage: same as docopts.py
#docopt() {
#    #   docopts [options] -h <msg> : [<argv>...]
#    # call python parser on embedded code
#    python <(sed -n -e '/^### EMBEDDED/,$ s/^#> // p' "$docopt_sh_me") "$@"
#}

# convert a repeatable option parsed by docopts into a bash ARRAY
#   ARGS['FILE,#']=3
#   ARGS['FILE,0']=somefile1
#   ARGS['FILE,1']=somefile2
#   ARGS['FILE,2']=somefile3
# Usage: myarray=( $(docopt_get_values ARGS --repeatable-option") )
docopt_get_values() {
    local opt=$2
    local ref="\${$1[$opt,#]}"
    local nb_val=$(eval echo "$ref")
    local i=0
    local vars=""
    while [[ $i -lt $nb_val ]] ; do
        ref="\${$1[$opt,$i]}"
        eval "vars+=\" $ref\""
        i=$(($i + 1))
    done
    echo $vars
}

# echo evaluable code to get alls the values into a bash array
# Usage: eval "$(docopt_get_eval_array ARGS FILE myarray)"
docopt_get_eval_array() {
    local ref="\${$1[$2,#]}"
    local nb_val=$(eval echo "$ref")
    local i=0
    local vars=""
    echo "declare -a $3"
    while [[ $i -lt $nb_val ]] ; do
        ref="\${$1[$2,$i]}"
        eval "echo \"$3+=( '$ref' )\""
        i=$(($i + 1))
    done
}

# Auto parser for the same docopts usage over scripts, for lazyness.
#
# It use this convention:
#  - help string in: $HELP
#  - Usage extracted by docopt_get_help_string at beginning of the script
#  - arguments are evaluated at global level in the assoc: $ARGS[]
#  - no version information
#
docopt_auto_parse() {
    local script_fname=$1
    shift
    # $HELP in global scope
    HELP="$(docopt_get_help_string "$script_fname")"
    # $ARGS[] assoc array must be declared outside of this function
    # or it's scope will be local, that's why we don't print it.
    docopts -A ARGS --no-declare -h "$HELP" : "$@"
    # returns the status of the docopts command, not grep status
    return ${PIPESTATUS[0]}
}

# Extract the raw value of a parsed docopts output.
# arguments:
#  - assoc: the docopts assoc name
#  - key:   the wanted key
#  - docopts_out: the full parsed output (before eval)
get_raw_value() {
    local assoc=$1
    local key="$2"
    local docopts_out="$3"
    local kstr=$(printf "%s['%s']" $assoc "$key")
    # split on '=', outputs the remaining for the matching $1
    awk -F= "\$1 == \"$kstr\" {sub(\"^[^=]+=\", \"\", \$0);print}" <<<"$docopts_out"
}

## main code
# --auto : don't forget to pass "$@"
# Usage: source docopts.sh --auto "$@"
if [[ "$1" == "--auto" ]] ; then
    shift
    # declare must be used at global scope to be accessible at
    # global level any were in the caller script.
    declare -A ARGS
    eval "$(docopt_auto_parse "${BASH_SOURCE[1]}" "$@")"
fi
