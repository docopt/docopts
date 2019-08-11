#!/usr/bin/env bash
#
# Helper to setup environment for this examples
#

pathadd() {
  local after=0
  if [[ "$1" == "after" ]] ; then
    after=1
    shift
  fi

  local p

  for p in $*
  do
    if [ -d "$p" ] && ! echo $PATH | grep -E -q "(^|:)$p($|:)" ; then
      if [[ $after -eq 1 ]]
      then
        PATH="$PATH:${p%/}"
      else
        PATH="${p%/}:$PATH"
      fi
    fi
  done
}

# comput current dir
# reandlink -f use GNU readlink, available on macos via: brew install coreutils
EXAMPLES_DIR=$(dirname $(readlink -f $BASH_SOURCE))
DOCOPTS_BIN=$(type -p docopts)

if [[ -z $DOCOPTS_BIN ]] ; then
  echo "adding to PATH: $EXAMPLES_DIR/../.."
  pathadd $EXAMPLES_DIR/../..
  DOCOPTS_BIN=$(type -p docopts)
  if [[ -z $DOCOPTS_BIN ]] ; then
    echo "ERROR: docopts not found in PATH, get a binary or compile it."
  fi
else
  echo "using docopts in PATH: $DOCOPTS_BIN"
fi
