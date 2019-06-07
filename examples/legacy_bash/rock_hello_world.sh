#!/usr/bin/env bash
# Example from README.md

PATH=$PATH:../..
eval "$(docopts -V - -h - : "$@" <<EOF
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

if $verbose ; then
    echo "Hello, world!"
fi
