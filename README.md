# WARNING: a new shell API is comming.

See Proposal lib API: [The Wiki](https://github.com/docopt/docopts/wiki)
The [develop branch](https://github.com/docopt/docopts/tree/develop) is abandoned

## docopts

Status: Draft

This is the golang version of command line wrapper for bash of docopts

See the Reference manual for command line `docopts` (old
[README.rst](https://github.com/docopt/docopts/blob/master/old_README.rst))

## Install

TODO: upload binary to release to they could be fetched directly

## Build

cross compile for 32btis
```
env GOOS=linux GOARCH=386 go build docopts.go
```

## Features

The command line tools `docopts` was written in python. It is unmaintained.
A new shell lib has been added `docopts.sh`.

The `docopts.sh` is an extra bash library you can source in your CLI script.
It may be removed and integrated in `docopts` binary.

You don't need a python interperter anymore.

See `examples/` for details.

## Developpers

All python related stuff will be removed.

If you want to clone this repository and hack docopts.

Use `git clone --recursive`, to get submodules.

fetch the extra golang version of docopt

```
go get github.com/docopt/docopt-go
```

If you forgot, you can also run:

~~~bash
git submodule init
git submodule update
~~~

Folder structure:

~~~
.
├── API_proposal.md - doc See wiki - to be removed
├── build.sh - build the embedded docopts.py into docopts.sh
├── docopt.py - original copy of docopt.py
├── docopts - current python wrapper - almost unmodified
├── docopts.py - copy of docopts, See build.sh
├── docopts.sh - bash lib - already embed both docopt.py + docopts.py
├── examples
│   ├── calculator_example.sh
│   ├── cat-n_wrapper_example.sh
│   ├── docopts_auto_examples.sh
│   └── quick_example.sh
├── language_agnostic_tester.py
├── setup.py
├── testcases.docopt
├── testee.sh
└── tests
    ├── bats/ - git submodules
    ├── bats.alias - source it to have bats working
    ├── docopts-auto.bats - unit test for --auto
    ├── docopts.bats - unit test docopts.sh
    └── exit_handler.sh - helper
~~~
