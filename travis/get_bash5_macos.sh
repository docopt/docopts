#!/usr/bin/env bash
#
# Only required for macos for unit testing advanced bash 4+ functionality and for bats.
#
# Our hack for macos bash version too old for unit testing

# was in .travis.yml
# brew update takes very long time
# according to https://docs.travis-ci.com/user/reference/osx#homebrew
# homebrew is already updated, but it's still really slow
#- if [[ "$TRAVIS_OS_NAME" == "osx" ]]; then brew update ; fi
#- if [[ "$TRAVIS_OS_NAME" == "osx" ]]; then brew install bash; fi

pathshow ()
{
    local var=${1:-PATH};
    eval "echo \$$var | tr : $'\n'"
}

if [[ "$TRAVIS_OS_NAME" == "osx" ]]; then
  if ((BASH_VERSINFO[0] <= 3)) ; then

    # The following hack kept a bash5 binary in our repository
    # only for our speedup pupose.

    # we simply install it
    bash_bin=/usr/local/bin/bash
    gzip -dc ./bash-5.0.16_x86_64-apple-darwin17.7.0.gz | \
      sudo bash -c "cat > $bash_bin && chmod a+x $bash_bin"
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
  # ./travis/reverse_ssh_tunnel.sh
  exit 1
fi
