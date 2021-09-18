#!/usr/bin/env bash
# vim: set et sw=4 ts=4 sts=4:
#
# docopts helper for bash, provides some functions for common usages
#
# Usage:
#   # raw usage
#   source path/to/docopts.sh
#   usage=$(docopt_get_help_string "$0")
#   version=$(docopt_get_version_string "$0")
#   eval $(docopts -A ARGS -h "$usage" -V $version : "$@")
#   myarray=( $(docopt_get_values ARGS FILE") )
#
#   # or for auto parsing the caller script comment (bash 4 associative array):
#   source path/to/docopts.sh --auto "$@"
#
#   # or for using globals variables (bash 3.2 compatible):
#   source path/to/docopts.sh --auto -G "$@"
#
# Conventions:
#   The prefix docopt_* is used to export globals and functions
#   docopt_auto_parse() modify $HELP and $ARGS or populate $ARGS_* globals.
#
# Code should work on bash 3.2 (mostly macOS) except where bash 4 mentioned
#
# For bash 3.2, you could use --auto -G for globals.
# Or simply source the helper if you need it, and use docopts directly:
#   source path/to/docopts.sh
#   docopts -G ARGS -h "$help" -V $version : "$@"

# Doc:
# fetch `Usage:` bloc from the given filename
# usually $0 in the main level script
docopt_get_help_string() {
    local myfname=$1
    # filter the first block starting at a "# Usage:" and ending at an empty line
    # one level of comment markup is removed.
    awk '
        BEGIN { u=0; l=0 }
        # we catch the first Usage: match
        /^# Usage:/ {
            if(u == 0)
            {
                u=1
            }
        }

        # match all lines. (Usage: is also matched)
        {
            if(u == 1) {
                # append to an array
                usage[l]=$0
                l++
            }
        }

        # empty line
        /^$/ {
            if(u == 1)
            {
                # stop parsing when empty line found
                u=2
            }
        }

        # display result and format output
        END {
            for(i=0; i<l; i++) {
                # remove comment (see issue #47) 
                sub("^# ", "", usage[i])
                sub("^#", "", usage[i])
                print usage[i]
            }
        }
        ' < "$myfname"
}

# Doc:
# Fetch version information from the given filename or string.
# Usually $0 in the main level script, or the help string extracted
# by docopt_get_help_string()
#
# Use standard delimiter ----
docopt_get_version_string() {
    if [[ -f "$1" ]] ; then
        # filter the block (all blocks) starting at a "# Usage:" and ending
        # at an empty line, one level of comment markup is removed
        sed -n -e '/^# ----/,/^$/ { 1d; s/^# \{0,1\}//; /----/ d; p; }' < "$1"
    else
        # use docopts --separator behavior
        echo "$1"
    fi
}

# Doc:
# convert a repeatable option parsed by docopts into a bash 4 ARRAY
#   ARGS['FILE,#']=3
#   ARGS['FILE,0']=somefile1
#   ARGS['FILE,1']=somefile2
#   ARGS['FILE,2']=somefile3
# Usage: myarray=( $(docopt_get_values ARGS FILE") )
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

# Doc:
# echo evaluable code to get all the values into a bash array
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

# Doc:
# Auto parser for the same docopts usage over scripts, for laziness.
# Used by --auto.
#
# bash 4:    convention use $ARGS as associative array
# bash 3.2:  convention use $ARGS_ prefix for globals variables
#
# It uses this convention:
#  - help string in: $HELP (modified in global scope)
#  - Usage is extracted by docopt_get_help_string at beginning of the script
#  - arguments are evaluated at global scope in the bash 4 assoc $ARGS (or globals with -G)
#  - no version information is handled
#
docopt_auto_parse() {
    local use_associative=true
    if [[ $1 == '-G' ]] ; then
        use_associative=false
        shift
    fi
    local script_fname=$1
    shift
    # $HELP in global scope
    HELP="$(docopt_get_help_string "$script_fname")"
    if $use_associative ;  then
        # $ARGS[] assoc array must be declared outside of this function
        # or its scope will be local, that's why we don't print it.
        docopts -A ARGS --no-declare -h "$HELP" : "$@"
        res=$?
    else
        # uses globals with ARGS_ prefix
        docopts -G ARGS -h "$HELP" : "$@"
        res=$?
    fi
    return $res
}

# Doc:
# Extract the raw value of a parsed docopts output.
# arguments:
#  - assoc: the docopts assoc name
#  - key:   the wanted key
#  - docopts_out: the full parsed output (before eval)
docopt_get_raw_value() {
    local assoc=$1
    local key="$2"
    local docopts_out="$3"
    local kstr=$(printf "%s['%s']" $assoc "$key")
    # split on '=', outputs the remaining for the matching $1
    awk -F= "\$1 == \"$kstr\" {sub(\"^[^=]+=\", \"\", \$0);print}" <<<"$docopts_out"
}

# Doc:
# Debug, prints env variable ARGS or $1 formatted
# Usage: docopt_print_ARGS [ASSOC_ARRAY_NAME]
#        docopt_print_ARGS -G [VARIABLE_PREFIX]
# with -G VARIABLE_PREFIX grep ${VARIABLE_PREFIX}_ variables from environment
docopt_print_ARGS() {
    local use_associative=true
    if [[ $# -ge 1 && $1 == '-G' ]] ; then
        use_associative=false
        shift
    fi
    # $1 can be the name of the global assoc array, or the prefix if -G is given
    local assoc="${1:-}"
    if [[ -z $assoc ]] ; then
        assoc=ARGS
    fi

    local a
    if $use_associative ; then
        # bash dark magic copying $assoc argument to a local myassoc array
        # inspired by:
        # https://stackoverflow.com/questions/6660010/bash-how-to-assign-an-associative-array-to-another-variable-name-e-g-rename-t#8881121
        declare -A myassoc
        eval $(typeset -A -p $assoc|sed "s/ $assoc=/ myassoc=/")

        # loop on keys
        echo "docopt_print_ARGS => $assoc"
        for a in ${!myassoc[@]} ; do
            printf "%20s = %s\n" $a "${myassoc[$a]}"
        done
    else
        echo "docopt_print_ARGS -G => ${assoc}_*"
        set | grep "^${assoc}_"
    fi
}

## main code if sourced with arguments
if [[ $# -ge 1 && "$1" == "--auto" ]] ; then
    if [[ $# -ge 2 && $2 == '-G' ]] ; then
        shift 2
        eval "$(docopt_auto_parse -G "${BASH_SOURCE[1]}" "$@")"
    else
        shift
        # declare must be used at global scope to be accessible at
        # global level anywhere in the caller script.
        declare -A ARGS
        eval "$(docopt_auto_parse "${BASH_SOURCE[1]}" "$@")"
    fi
fi
