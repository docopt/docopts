#!/usr/bin/env bash
#
# Example reading from stdin
#

# no need to change PATH if docopts binary is in PATH already
PATH=../..:$PATH

# read both verion and usage from stdin
# in global mode (option names arg mangled, See docopts --help
parsed="$(docopts -V - -h - : "$@" <<EOF
Usage: rock [options] <argv>...

Options:
      --verbose  Generate verbose messages.
      --help     Show help options.
      --version  Print program version.
----
rock 0.1.0
Copyright (C) 200X Thomas Light
License RIT (Robot Institute of Technology)
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
EOF
)"

eval "$parsed"

echo "$parsed"
echo "==========================="

if $verbose ; then
  # env var are in global scope
  typeset -p argv help version
fi

# docopts will fail without <argv>
if [[ -n $argv ]] ; then
  for body in ${argv[@]}
  do
    echo "let's rock $body"
  done
fi

