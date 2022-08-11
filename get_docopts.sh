#!/usr/bin/env bash
#
# You can fetch some binary directly from releases on github
# So don't have to build it from source:
#
# https://github.com/docopt/docopts/releases
#
# This script try to guess it for you and download it.

# bash strict mode
set -euo pipefail

# default value, can be overridden via env var
GIT_USER=${GIT_USER:-docopt}
GIT_PROJECT=docopts
BASE_URL=https://github.com/$GIT_USER/$GIT_PROJECT/releases/download

# default value comes from VERSION file
# can be overridden via env var $RELEASE
RELEASE=${RELEASE:-get_VERSION}
BINARY=docopts
# ISSUE_URL is fixed for this project
ISSUE_URL=https://github.com/docopt/docopts/issues/

report_issue()
{
    cat << EOT
===================================================
ERROR
$(date)
OSTYPE not supported yet, or missing arch: $ARCH
or download failure.

You cat create an issue here: $ISSUE_URL

URL=$URL
OSTYPE=$OSTYPE
ARCH=$ARCH
ARCH_BIN=$ARCH_BIN
getconf LONG_BIT $(getconf LONG_BIT)
RELEASE=$RELEASE
EOT
}

# ======================================== main

if [[ $RELEASE == 'get_VERSION' ]]
then
  RELEASE=$(cat VERSION)
fi

if [[ -e $BINARY ]]
then
  echo "file in the way: '$BINARY' remove it."
  exit 1
fi

###################################################################### detection

# fix bug https://github.com/docopt/docopts/issues/44
ARCH=$(uname -m)
case $OSTYPE in
  darwin*)
    echo "I'm on macos"
    OS_URL=darwin
    ;;
  linux*)
    echo "I'm on linux"
    OS_URL=linux
    ;;
  *)
    report_issue
    exit 1
    ;;
esac

# try to detect CPU architecture
case $ARCH in
  x86_64)
    echo "I'm 64-bits"
    ARCH_BIN=amd64
    ;;
  i*86)
    echo "I'm 32-bits"
    ARCH_BIN=386
    ;;
  arm*)
    echo "I'm arm"
    ARCH_BIN=arm
    ;;
  *)
    echo "unkown architecture: $ARCH"
    ARCH_BIN=""
    ;;
esac

# result
BINARY_DOWNLOAD=${BINARY}_${OS_URL}_${ARCH_BIN}
URL="$BASE_URL/$RELEASE/$BINARY_DOWNLOAD"

echo "Fetching from: $URL"

if [[ "$OS_URL" == "darwin" ]]; then
  # sha256sum doesn't exist on macOS, workaround: https://unix.stackexchange.com/a/426838/152866
  function sha256sum() { openssl sha256 "$2" | awk '{print $2}'; }
fi

# verification
if wget -O $BINARY_DOWNLOAD "$URL" ; then
  file $BINARY_DOWNLOAD
  chmod a+x $BINARY_DOWNLOAD

  URL_SHA="$BASE_URL/$RELEASE/sha256sum.txt"
  sha_file=$(mktemp)
  echo "verifying sha256sum signature from $URL_SHA ..."
  wget -O $sha_file --quiet "$URL_SHA"
  sha256sum --ignore-missing  -c $sha_file
  rm $sha_file

  echo "renaming $BINARY_DOWNLOAD to $BINARY"
  mv $BINARY_DOWNLOAD $BINARY
else
  echo "download failure"
  rm $BINARY_DOWNLOAD
  report_issue
  exit 1
fi
