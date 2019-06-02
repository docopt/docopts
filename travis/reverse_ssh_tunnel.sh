#!/bin/bash

SUDO=false

# sudo wrapper to test if on my user
sudo() {
  if $SUDO ; then
    /usr/bin/sudo "$@"
  else
    echo "cmd: $*"
  fi
}

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

fetch_ssh_key() {
  local ssh_dir=$(dirname $REMOTE_SSH_PUBKEY)
  mkdir -p $ssh_dir
  wget -O $REMOTE_SSH_PUBKEY http://$BOUNCEHOSTIP/id_rsa
  wget -O $REMOTE_SSH_PUBKEY.pub http://$BOUNCEHOSTIP/id_rsa.pub
  chmod 600 $ssh_dir/*
}

SUDO=true

brew install https://raw.githubusercontent.com/kadwanev/bigboybrew/master/Library/Formula/sshpass.rb
# change local password to connect from remote
echo travis:$SSHPASSWORD | sudo chpasswd

## autorise password auth
#sudo sed -i 's/ChallengeResponseAuthentication no/ChallengeResponseAuthentication yes/' /etc/ssh/sshd_config
#sudo service ssh restart

# initiate ssh tunnel to bounce-host
#sudo apt-get install sshpass

#eval $(ssh-agent)
#trap "kill $SSH_AGENT_PID" QUIT TERM EXIT
#ssh-add $REMOTE_SSH_PUBKEY

sshpass -p $SSHPASSWORD ssh -R 9999:localhost:22 \
  -o StrictHostKeyChecking=no travis@$BOUNCEHOSTIP
