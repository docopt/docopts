#!/bin/bash
#
# You can fetch some binary directly from release on github
#
# We encourage to build your own version from source.
#

GIT_USER=Sylvain303
GIT_PROJECT=docopts
BASE_URL=https://github.com/$GIT_USER/$GIT_PROJECT/releases/download
RELEASE=v0.6.3-alpha1
DEST_BIN=docopts

if [[ -e $DEST_BIN ]]
then
  echo "file in the way: '$DEST_BIN' remove it."
  exit 1
fi

if [[ $(getconf LONG_BIT) == "64" ]]
then
    echo "I'm 64-bits"
    URL="$BASE_URL/$RELEASE/docopts"
else
    echo "I'm 32-bits"
    URL="$BASE_URL/$RELEASE/docopts-32bits"
fi

set -e
echo "Fetching from: $URL"
wget -q -O $DEST_BIN "$URL"
file $DEST_BIN
chmod a+x $DEST_BIN
