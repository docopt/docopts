#!/usr/bin/env bash
# vim: st=2 sw=2 sts=2 et:

usage() {
  cat << EOU
Test for docopts double-dash handling

Usage:
  $0 cmd [options] [--] [<unparsed_options>...]
  $0 diff [-u] BEFORE AFTER
  $0 -h | --help

Options:
  -p --platform=<platform>  Platform to configure
  --trace                   Full trace output from bash
  -u                        unified diff
EOU
}

# global mode
# not needed if docopts already in PATH
PATH=../..:$PATH
parsed="$(docopts -G ARGS -h "$(usage)" : "$@")"
echo "$parsed"
eval "$parsed"
