#!/usr/bin/env bash
#
# Update from tarball
#
# Steps:
# 1. change GOVER to the wanted version
# 2. Run ./update_go.sh

set -euo pipefail

localgo=/usr/local/go
gobin=$localgo/bin/go

set +e
actual_gover=$($gobin version | grep -o -E 'go[0-9.]+')
set -e

# change this go version for updating
# can be copy/paste from the email announcing the version
GOVER=go1.17.1

if [[ $actual_gover == $GOVER ]] ; then
  echo "$gobin version: $actual_gover is the same as $GOVER"
  exit 0
fi

gotar="${GOVER}.linux-amd64.tar.gz"
url="https://dl.google.com/go/$gotar"

if [[ ! -e  $gotar ]] ; then
  wget $url -O $gotar
else
  echo "using local file: $gotar"
fi

if [[ -d $localgo ]] ; then
  old=/usr/local/go.old
  if [[ -d  $old ]] ; then
    echo "old dest exists, remove it: $old"
    exit 1
  fi

  echo "moving go to $old"
  sudo mv $localgo $old
fi

if [[ ! -d $localgo ]] ; then
  dest=/usr/local
  echo "extracting $gotar to $dest"
  sudo tar -C $dest -xzf $gotar
fi

echo "Go updated:"
$localgo/bin/go version

echo "add $localgo/bin to your PATH"
