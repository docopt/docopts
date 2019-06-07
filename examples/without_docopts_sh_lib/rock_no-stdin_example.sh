#!/usr/bin/env bash

usage() {
  cat << EOU
Usage: rock [options] <argv>...

Options:
      --verbose  Generate verbose messages.
      --help     Show help options.
      --version  Print program version.
EOU
}

version="rock 0.1.0
Copyright (C) 200X Thomas Light
License RIT (Robot Institute of Technology)
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law."

# no need to change PATH if docopts are in your PATH
PATH=../..:$PATH

parsed=$(docopts -V "$version" -h "$(usage)" : "$@")
#echo "$parsed"
eval "$parsed"

echo "verbose=$verbose"

if $verbose ; then
    echo "Hello, world!"
else
    :
    # I'm silent
fi
