#!/usr/bin/env bash
#
# Only required for macos for unit testing advanced bash 4+ functionality and for bats.
#
# Our hack for macos bash version too old for unit testing

pathshow ()
{
    local var=${1:-PATH};
    eval "echo \$$var | tr : $'\n'"
}

if [[ "$RUNNER_OS" == "macOS" ]]; then
  if ((BASH_VERSINFO[0] <= 3)) ; then
    HOMEBREW_NO_AUTO_UPDATE=1 brew install bash # quick install, no brew update (c.f. https://apple.stackexchange.com/a/293252/167983)
  fi
fi

echo "======================= splited PATH"
pathshow PATH

echo "======================= new bash version"
hash -r bash
bash --version
type bash

# require sed macos comaptible regexp \+ doesn't exist ==> \{1,\}
MY_BASH_VERSINFO=$(bash --version | sed -n -e '1 s/^.*version \([0-9.]\{1,\}\).*/\1/ p')
if [[ ! $MY_BASH_VERSINFO =~ ^[4-9] ]] ; then
  echo "install bash5 failed"
  exit 1
fi
