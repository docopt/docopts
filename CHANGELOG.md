# Changelog

## bug fix v0.6.4-with-no-mangle-double-dash

Tag: `v0.6.4-with-no-mangle-double-dash`

This is a bug fix release.

features changes:
  - fix error refusing mangling double dash in global mode [#52]
  - still refuse to mangle single dash `[-]` in global mode (not compatible with v0.6.1)
  - now can mangle single-dash and double-dash in `-G` mode
  - fix output bash empty array #53

internal changes:
  - use Go 1.17.1 for compiling
  - `language_agnostic_tester.py` re-write for python3
  - sort argument key in `Print_bash_args()` `Print_bash_global()`
  - sort input keys in `docopts_test.go`
  - add PR #52 test to `testcases.docopt`
  - completed developer documentation
  - #52 ignore `[--]` double-dash option in `Print_bash_global()`
  - reformat of Go code with `gofmt`; indent change Space ==> Tab
  - add `tests/functional_tests_docopts.bats`

## docopts binary transitional v0.6.3-rc2

Tag: `v0.6.3-rc2`

This is a transitional release.

It supports all the previous command line API plus some extra options.
Fully compatible with previous 0.6.2 python code for Bash.

See: https://github.com/docopt/docopts/tree/v0.6.1%2Bfix
based on master branch

features changes:
  - `docopts.sh` function `docopt_get_help_string()` now uses awk to extract only first `Usage:`

internal changes:
  - use Go 1.14 for compiling
  - more pre-built binaries, removed darwin/386
  - fixed #44 `get_docopts.sh` for macOS + functional tests
  - removed `bats` git submodule
  - use [bats-core 1.2-dev](https://github.com/bats-core/bats-core) as testing framework from travis
  - `deploy.sh` removed, now uses its [own repository](https://github.com/opensource-expert/deploy.sh)
  - updated Makefile to use `get_ldflags.sh` from `deploy.sh`
  - [travis hack](travis/get_bash5_macos.sh) to get faster build on macos with our embedded bash5 binary

## docopts binary transitional v0.6.3-rc1

Tag: `v0.6.3-rc1`

This is a transitional release.

It supports all the previous API plus some extra command line options.
Fully compatible with previous 0.6.2 Python code for Bash.

See: https://github.com/docopt/docopts/tree/v0.6.1%2Bfix
based on master branch

changes:
  - more test for macOS
  - Bash 3.2 support is more documented and fixed
  - use `bats-core` as testing framework
  - updated README merged from old README.rst
  - now documentation introduce `docopts.sh` See [docs](https://github.com/docopt/docopts/tree/v0.6.3-rc1/docs/)
  - added `Makefile`
  - added `build_doc.sh` PoC markdown preprocessor

all examples written for docopts:
  - shebang conversion `#!/bin/bash` ==> `#!/usr/bin/env bash`
  - legacy example completed
  - example from README extracted a file, and merged in README via `build_doc.sh`
  - `sshdiff` full example coded
  - added examples with `--auto -G`

docopts.sh helper:
  - is more documented in the code
  - has documenation in [docs/README.md](https://github.com/docopt/docopts/tree/v0.6.3-rc1/docs/README.md)
  - now supports bash strict mode (`set -euo pipefail`)
  - now supports `--auto -G` to auto parse with global vars (doesn't use bash 4 associative array)
`docopts` behavior sould be unchanged:
  - add mangled name collision detection

## docopts for Bash first release in golang

Tag: `v0.6.3-alpha1`

This is a transitional release.

It is a complete rewrite of the Python code in Go.
It supports all of the previous API plus some extra command line options.

Fully compatible with previous 0.6.2 python code for Bash.
based on https://github.com/docopt/docopts/tree/packaging-debian

## first release in Go

Tag: `v0.6.2`

This is a transitional release.

Docopt for shell.
It supports all of the previous API.

Fully compatible with previous 0.6.1 Python code for Bash.
See: https://github.com/docopt/docopts/tree/v0.6.1%2Bfix
based on master branch

changes:
  - now written in Go, no more Python dependency
