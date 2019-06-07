#!/usr/bin/env bash
# Naval Fate.
#
# Usage:
#   naval_fate.sh ship new <name>...
#   naval_fate.sh ship <name> move <x> <y> [--speed=<kn>]
#   naval_fate.sh ship shoot <x> <y>
#   naval_fate.sh mine (set|remove) <x> <y> [--moored|--drifting]
#   naval_fate.sh -h | --help
#   naval_fate.sh --version
#
# Options:
#   -h --help     Show this screen.
#   --version     Show version.
#   --speed=<kn>  Speed in knots [default: 10].
#   --moored      Moored (anchored) mine.
#   --drifting    Drifting mine.
#


# if docopts is in PATH, not needed.
# Note: docopts.sh is also found in PATH
PATH=../..:$PATH

VERSION='Naval Fate 2.0'
source docopts.sh
# no version support in docopt_auto_parse() so we call docopts directly
usage=$(docopt_get_help_string "$0")
parsed="$(docopts -G ARGS -V "$VERSION" -h "$usage" : "$@")"
echo "============ parsed output"
echo "$parsed"
# now vars are populated at global scope
eval "$parsed"

echo "============== docopt_print_ARGS"
docopt_print_ARGS -G
