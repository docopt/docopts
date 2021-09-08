#!/bin/bash

docopts --debug -h - : "$@" <<EOF
Test for docopts double-dash handling

Usage:
  $0 --platform=<platform> [options] [--] [<unparsed_options>...]
  $0 -h | --help

Options:
  -p --platform=<platform> Platform to configure
  --trace                  Full trace output from bash
EOF
