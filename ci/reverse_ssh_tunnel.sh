#!/usr/bin/env bash
#
# Usage: uncomment the reverse_ssh_tunnel lines in .github/workflows/ci.yml at the bottom of the file
# required: a bounce host MUST be set before!
# NOTE: ansible playbook to build the bounce_host not provided yet.

set -euo pipefail

###################################################################### Config

# The IP here is a temporay public cloud VM
# You MUST set the good IP on BOUNCEHOST or a valid IP.
# This will be the ssh machine that this script will connect too from travis instance
BOUNCEHOST="travis.opensource-expert.com"

###################################################################### Code

# Where to store localy the remote bounce host ssh key
TEMP_SSH_KEYS=/tmp/$USER-bounce-travis/id_rsa

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

fetch_ssh_keys()
{
  local tmp_ssh_dir=$(dirname $TEMP_SSH_KEYS)
  mkdir -p $tmp_ssh_dir
  wget -O $TEMP_SSH_KEYS http://$BOUNCEHOST/id_rsa
  wget -O ${TEMP_SSH_KEYS}.pub http://$BOUNCEHOST/id_rsa.pub
  chmod 600 $tmp_ssh_dir/*
}

cleanup()
{
  local tmp_ssh_dir=$(dirname $TEMP_SSH_KEYS)
  rm -rf $tmp_ssh_dir
  kill $SSH_AGENT_PID
}

bounce_host_remote_exec()
{
  NOOP_DELAY=30
  # our 30 minutes timeout
  MAX=$((30*60))
  echo 'bounce host connected'

  # generate the stop command
  echo 'pkill -f bounce_host_remote_exec; pkill -f sleep' > disconnect
  chmod a+x disconnect

  # message when disconnect kill is sent
  trap "echo 'bounce_host existing'; kill -9 $$" QUIT TERM EXIT

  local start=$SECONDS
  local elapsed
  while true
  do
      sleep $NOOP_DELAY
      elapsed=$(($SECONDS-$start))
      echo "noop ${elapsed}s"
      if [[ $elapsed -gt $MAX ]] ; then
        echo "timeout reached ${MAX}s"
        break
      fi
  done
  echo "bounce_host connection closed"
}

################################################ main

fail_if_empty BOUNCEHOST TEMP_SSH_KEYS
fetch_ssh_keys

# start ssh agent and add the private key
eval $(ssh-agent)
trap "cleanup" QUIT TERM EXIT
ssh-add $TEMP_SSH_KEYS

# allow ssh back to us with the same key
mkdir -p $HOME/.ssh
chmod 700 $HOME/.ssh
cp $TEMP_SSH_KEYS.pub $HOME/.ssh/authorized_keys
chmod 600 $HOME/.ssh/authorized_keys

# for 10 min without output Success
# https://travis-ci.org/Sylvain303/docopts/builds/540455090#L1295
ssh -R 9999:localhost:22 \
  -o StrictHostKeyChecking=no travis@$BOUNCEHOST \
  "$(typeset -f bounce_host_remote_exec); bounce_host_remote_exec"
