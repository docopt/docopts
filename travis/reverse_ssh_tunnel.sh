#!/bin/bash

SSHPASSWORD="mypass"
BOUNCEHOSTIP="51.68.156.147"
REMOTE_SSH_PUBKEY=/tmp/bounce-travis/id_rsa

if [[ -z $BOUNCEHOSTIP ]] ;then
  echo "failed BOUNCEHOSTIP must be set"
  exit 1
fi

if [[ -z $SSHPASSWORD ]] ;then
  echo "failed SSHPASSWORD must be set"
  exit 1
fi

if [[ -z $SSHPASSWORD ]] ;then
  echo "failed SSHPASSWORD must be set"
  exit 1
fi

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

#sshpass -p $SSHPASSWORD 
ssh -R 9999:localhost:22 \
  -o StrictHostKeyChecking=no travis@$BOUNCEHOSTIP
