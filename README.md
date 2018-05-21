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
├── build.sh
├── docopts.go
├── docopts.sh
├── examples
│   ├── calculator_example.sh
│   ├── cat-n_wrapper_example.sh
│   ├── docopts_auto_example.sh
│   ├── quick_example.sh
│   └── rock_stdin_example.sh
├── language_agnostic_tester.py
├── LICENSE-MIT
├── old_README.rst
├── PROGRESS.md
├── README.md
├── testcases.docopt
├── testee.sh
├── tests
│   ├── bats [...]
│   ├── bats.alias
│   ├── docopts-auto.bats
│   ├── docopts.bats
│   ├── exit_handler.sh
│   └── TODO.md
└── TODO.md
~~~

## Tests

Some tests are coded along this code base.

- bats bash unit tests and functionnal testing
- language_agnostic_tester (old python wrapper, docopts compatibily tests)
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

### golang docopt.go

```
cd PATH/to/go/src/github.com/docopt/docopt-go/
go test -v .
```

