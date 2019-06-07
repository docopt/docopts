#!/bin/bash

# a bounce host must be set before
BOUNCEHOSTIP="51.68.156.147"
REMOTE_SSH_PUBKEY=/tmp/bounce-travis/id_rsa

fail_if_empty() {
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

fail_if_empty BOUNCEHOSTIP REMOTE_SSH_PUBKEY

fetch_ssh_keys() {
  local ssh_dir=$(dirname $REMOTE_SSH_PUBKEY)
  mkdir -p $ssh_dir
  wget -O $REMOTE_SSH_PUBKEY http://$BOUNCEHOSTIP/id_rsa
  wget -O ${REMOTE_SSH_PUBKEY}.pub http://$BOUNCEHOSTIP/id_rsa.pub
  chmod 600 $ssh_dir/*
}

fetch_ssh_keys

# start ssh agent and add the private key
eval $(ssh-agent)
trap "kill $SSH_AGENT_PID" QUIT TERM EXIT
ssh-add $REMOTE_SSH_PUBKEY

# allow ssh back to us with the same key
mkdir $HOME/.ssh
chmod 700 $HOME/.ssh
cp $REMOTE_SSH_PUBKEY.pub $HOME/.ssh/authorized_keys
chmod 600 $HOME/.ssh/authorized_keys

# for 10 min without output Success
# https://travis-ci.org/Sylvain303/docopts/builds/540455090#L1295
ssh -R 9999:localhost:22 \
  -o StrictHostKeyChecking=no travis@$BOUNCEHOSTIP
