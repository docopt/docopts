# WARNING: a new shell API is comming.

See Proposal lib API: [on `docopts` Wiki](https://github.com/docopt/docopts/wiki)

The [develop branch](https://github.com/docopt/docopts/tree/develop) is abandoned

## docopts

Status: Draft - work in progress

This is the golang version  of `docopts` : the command line wrapper for bash.

See the Reference manual for command line `docopts` (old
[README.rst](old_README.rst))

## Install

TODO: upload binary to release to they could be fetched directly

## Build

cross compile for 32btis
```
env GOOS=linux GOARCH=386 go build docopts.go
```

local build
```
go build docopts.go
```

built on: `go version go1.10.2 linux/amd64`

## Features

Warning: may be not up-to-date feature list.

The command line tools `docopts` was written in python. It is unmaintained.
A new shell lib has been added `docopts.sh`.

The `docopts.sh` is an extra bash library that you can source in your CLI script.
It may be removed and integrated in `docopts` binary.

You don't need a python interpreter anymore.

As of 2018-05-22

* `docopts` is able to reproduce 100% of the python version.
* Language agnostic test is 100% (on adm64 Linux)

## Developpers

All python related stuff will be removed.

If you want to clone this repository and hack docopts:

Use `git clone --recursive`, to get submodules only required for testing

fetch the extra golang version of docopt (required for building `docopts`)

```
go get github.com/docopt/docopt-go
```

If you forgot `--recursive`, you can also run afterward:

~~~bash
git submodule init
git submodule update
~~~

Current folder structure:

~~~
.
├── API_proposal.md
├── build.sh - no more used yet
├── docopts.go - main source code
├── docopts.sh - library wrapper adn helpers
├── examples - all examples may not works, as they have been forked from python version
│   ├── calculator_example.sh
│   ├── cat-n_wrapper_example.sh
│   ├── docopts_auto_example.sh
│   ├── quick_example.sh
│   └── rock_stdin_example.sh
├── language_agnostic_tester.py - old python JSON tester still used with testee.sh
├── LICENSE-MIT
├── old_README.rst - copied README
├── PROGRESS.mda - what I'm working on
├── README.md
├── testcases.docopt - agnostic testcases copied from python's docopt
├── testee.sh - bash wrapper to convers docopts assoc array output to JSON
├── tests - unit and functional testing written in bats
│   ├── bats [...]
│   ├── bats.alias
│   ├── docopts-auto.bats
│   ├── docopts.bats
│   ├── exit_handler.sh
│   └── TODO.md
└── TODO.md - Some todo list on this golang version of docopts
~~~

## Tests

Some tests are coded along this code base.

- bats bash unit tests and functionnal testing
- language_agnostic_tester (old python wrapper, full docopt compatibily tests)
- See Also: docopt.go own tests in golang

### Runing tests

#### bats
```
cd ./tests
. bats.alias
bats .
```

#### language_agnostic_tester

```
cd PATH/TO/docopts
python language_agnostic_tester.py ./testee.sh
```

#### golang docopt.go

```
cd PATH/to/go/src/github.com/docopt/docopt-go/
go test -v .
```
