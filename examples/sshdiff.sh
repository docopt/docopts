#!/bin/bash
#
# Show file differences between 2 hosts.
# Usage: sshdiff.sh [-h] [-s] <host1> <host2> <file> [<lines_context>]
#
# If not specified, <lines_context> defaults to 3.
#
# Use colordiff if available.
#
# Options:
#     -h       display this help and exit
#     -s       use sort instead of cat to show remote <file>
#
# Examples:
#     sshdiff.sh server1 hostname2 /etc/hostname

PATH=../..:$PATH
source docopts.sh --auto "$@"

docopt_print_ARGS


