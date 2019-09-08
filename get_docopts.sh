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

# default value comes from VERSION file, can be overridden via env var
RELEASE=${RELEASE:-VERSION}
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

if [[ $RELEASE == 'VERSION' ]]
then
  RELEASE=$(cat VERSION)
fi

if [[ -e $BINARY ]]
then
  echo "file in the way: '$BINARY' remove it."
  exit 1
fi

# try to detect CPU architecture
ARCH=$(arch)
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


case $OSTYPE in
  darwin*)
    echo "I'm on macos"
    URL="$BASE_URL/$RELEASE/${BINARY}_darwin_${ARCH_BIN}"
    ;;
  linux*)
    echo "I'm on linux"
    URL="$BASE_URL/$RELEASE/${BINARY}_linux_${ARCH_BIN}"
    ;;
  *)
    report_issue
    exit 1
    ;;
esac

echo "Fetching from: $URL"
if wget -O $BINARY "$URL" ; then
  file $BINARY
  chmod a+x $BINARY
else
  echo "download failure"
  rm $BINARY
  report_issue
  exit 1
fi
