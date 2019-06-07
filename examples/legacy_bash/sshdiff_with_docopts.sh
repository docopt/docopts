#!/usr/bin/env bash
#
# Show file differences between 2 hosts.
# Usage: sshdiff.sh [-h] [-s] HOST1 HOST2 FILENAME [LINES_CONTEXT]
# 
# Use colordiff if available.
# 
# Options:
#     -h   display this help and exit
#     -s   use sort instead of cat to show remote FILENAME
# 
# If not specified, LINES_CONTEXT defaults to 3.
# 
# Environment variable:
#   EXTRA_DIFF  pass extra argument to diff command
#   SSH_USER    use this user for remote ssh user (default: $USER)
# 
# Examples:
#     sshdiff.sh server1 server2 /etc/crontab
#     # ignore changes on blank lines
#     EXTRA_DIFF=-b sshdiff.sh webserver1 webserver2 /etc/apache2/apache2.conf

# The line above must be empty
PATH=../..:$PATH
source docopts.sh --auto -G "$@"

docopt_print_ARGS -G

# converted main code
FILTER="cat"
if [ -z "$AGRS_LINES_CONTEXT" ]; then
  ARGS_LINES_CONTEXT=3
fi

# selecting diff program
DIFF="$(which diff)"
COLORDIFF="$(which colordiff)"
if [ -x "$COLORDIFF" ]; then
  DIFF="$COLORDIFF"
fi

if [ -z "$SSH_USER" ]; then
  SSH_USER=$USER
fi

# complete final command
"$DIFF" $EXTRA_DIFF -U "$ARGS_LINES_CONTEXT" \
    <(ssh $SSH_USER@"$ARGS_HOST1" $FILTER "$ARGS_FILENAME") \
    <(ssh $SSH_USER@"$ARGS_HOST2" $FILTER "$ARGS_FILENAME")

