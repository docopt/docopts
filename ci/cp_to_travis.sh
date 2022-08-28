#!/usr/bin/env bash
#
# Usage: ./ci/cp_to_travis.sh LOCAL_FILENAME
#
# require: a bounce host MUST be set before!

set -euo pipefail

###################################################################### Config

# The IP here is a temporay public cloud VM
# You MUST set the good IP on BOUNCEHOST or a valid IP
BOUNCEHOST="travis.opensource-expert.com"

###################################################################### Code

fail_if_empty()
{
  local varname
  local v
  # allow multiple check on the same line
  for varname in $*
  do
    eval "v=\$$varname"
    if [[ -z "$v" ]] ; then
      echo "error: $varname empty or unset at ${BASH_SOURCE[1]}:${FUNCNAME[1]} line ${BASH_LINENO[0]}"
      exit 1
    fi
  done
}

# will be run on travis_host
remote_copy()
{
  local local_copy="$1"
  # dest MUST include parent project name too
  # /home/sylvain/code/go/src/github.com/docopt/docopts ==> parent docopts/
  local dest="$2"
  local dest_dir="$(dirname "$dest")"

  # find our base on travis
  local travis_base=$(ls -Fd /Users/travis/build/* | grep -E '/$')
  cd $travis_base || { echo "error: can't cd to '$travis_base'" ; exit 1; }
  cd "$dest_dir"

  # actually move the file to the final destination on travis_host
  mv $HOME/$local_copy "$(basename "$dest")"
}

# on bounce_host I scp the file first on travis_host
# using stdin bring some strange error.
# we are using a local tempfile
bounce_host_cp()
{
  local arg="$1"
  local file="$(basename "$arg")"
  local t="$(mktemp /dev/shm/$file.XXXXXX)"
  cat > "$t"
  scp -P 9999 "$t" localhost:
  # the temporay filename will be used on travis_host
  echo "$(basename $t)"
}

################################################ main

arg=$1
fail_if_empty BOUNCEHOST arg

arg=$(readlink -f "$1")
# remove our base
local_base=/home/sylvain/code/go/src/github.com/docopt/
dest=${arg#$local_base}
echo "$arg ==> $dest"

# Send file content on stdin.
# Our functions above will be executed on bounce_host + travis host
# The second typeset -f is executed on bounce_host to resolve
# quote escaping nested level.
cat "$arg" |
  ssh travis@$BOUNCEHOST "
    $(typeset -f bounce_host_cp remote_copy);
    t=\$(bounce_host_cp '$arg')
    ssh -p 9999 localhost \"\$(typeset -f remote_copy); remote_copy '\$t' '$dest'\"
    "
