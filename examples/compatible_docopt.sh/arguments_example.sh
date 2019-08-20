#!/usr/bin/env bash

DOC="Argument parser
Usage: arguments_example.sh [-vqrh] [FILE] ...
       arguments_example.sh (--left | --right) CORRECTION FILE

Process FILE and optionally apply correction to either left-hand side or
right-hand side.

Arguments:
  FILE        optional input file
  CORRECTION  correction angle, needs FILE, --left or --right to be present

Options:
  -h --help
  -v       verbose mode
  -q       quiet mode
  -r       make report
  --left   use left-hand side
  --right  use right-hand side"

main_arguments()
{
  # main function for this script

  # only display parsed arguments
  set | grep "^$DOCOPT_PREFIX"

  return 0
}

DOCOPT_PREFIX=ARGS_
case $DOCOPT_PARSER in
  docopts)
    if [[ -z $(type -p docopts) ]] ; then
      echo "docopts not found in PATH, use: source example_env.sh"
      exit 1
    fi
    # docopts append _ to prefix
    eval "$(docopts -G ${DOCOPT_PREFIX%_} --docopt_sh -h "$DOC" : "$@")"
    ;;
  docopt.sh)
    eval "$(docopt "$@")"
    ;;
  "")
    echo "DOCOPT_PARSER is undefined"
    exit 1
    ;;
  *)
    echo "DOCOPT_PARSER unsuported value: $DOCOPT_PARSER"
    exit 1
    ;;
esac

main_arguments "$@"
