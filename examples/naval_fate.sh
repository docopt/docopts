#!/bin/bash
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

VERSION='Naval Fate 2.0'

# if docopts is in PATH, not needed.
# Note: docopts.sh is also found in PATH
PATH=..:$PATH
# auto parse the header above, See: docopt_get_help_string
source docopts.sh --auto "$@"

docopt_print_ARGS
