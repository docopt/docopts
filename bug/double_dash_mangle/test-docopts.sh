#!/bin/bash

docopts_input="
Test for docopts double-dash handling

Usage:
  $0 --platform=<platform> [options] [--] [<unparsed_options>...]
  $0 -h | --help

Options:
  -p --platform=<platform> Platform to configure
  --trace                  Full trace output from bash
"

# read from the path
docopts_bin=$(type -p docopts)

# argument hack for choosing the binary to compare behavior
case $1 in
  branch)
    # localy compiled
    docopts_bin=../../docopts
    shift
    ;;
  dana)
    # from inside the container (or cloned localy)
    docopts_bin=./patched/docopts
    shift
    ;;
  current)
    # in the current folder (container or local)
    # from wget our latest release wget https://github.com/docopt/docopts/releases/download/v0.6.3-rc2/docopts_linux_amd64
    docopts_bin=./docopts_linux_amd64
    shift
    ;;
  input)
    # just returns the $docopts_bin and leave to be reused externaly
    echo "$docopts_input"
    exit
esac

if [[ ! -x $docopts_bin ]]
then
  echo "docopts_bin not found: $docopts_bin"
  exit 1
fi
$docopts_bin --version | head -1

$docopts_bin --debug -h "$docopts_input" : "$@"
