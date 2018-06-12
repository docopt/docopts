# A new shell API is proposed

See Proposal lib API: [on `docopts` Wiki](https://github.com/docopt/docopts/wiki)

The [develop branch](https://github.com/docopt/docopts/tree/develop) is abandoned!

The current go version is 100% compatible with python's docopts.
Please report any non working code with [issue](https://github.com/docopt/docopts/issues) and examples.

## docopts

Status: Alpha - work is done

This is the golang version  of `docopts` : the command line wrapper for bash.

See the Reference manual for command line `docopts` (old
[README.rst](old_README.rst))

## Install

You have to drop the binary, and eventually the `docopts.sh` lib helper in your PATH.
The binary is standalone and staticaly linked. So it runs everywhere.

See build section.

next as root

```
cp docopts docopts.sh /usr/local/bin
```

### pre-built binary

pre-built binary are attached to [releases](https://github.com/Sylvain303/docopts/releases)

download and rename it as `docopts` and put in your `PATH`

```
mv docopts-32bit docopts
cp docopts docopts.sh /usr/local/bin
```

You are strongly encouraged to build your own binary. Find a local golang developper in whom you trust and ask her, for a beer or two, if she could build it for you. ;)

## Usage

See manual and [Examples](examples/)

## Compiling

With a go workspace.

cross compile for 32btis
```
env GOOS=linux GOARCH=386 go build docopts.go
```

local build
```
go get github.com/docopt/docopt-go
go get github.com/Sylvain303/docopts
cd src/github.com/Sylvain303/docopts
go build docopts.go
```

or via Makefile (generate 64bits, 32bits, arm and OSX-64bits version of docopts)

```
cd src/github.com/Sylvain303/docopts
make all
```

Tested built on: `go version go1.10.2 linux/amd64`

## Features

Warning: may be not up-to-date feature list.

The command line tools `docopts` was written in python. It is unmaintained.
A new shell lib has been added [`docopts.sh`](docopts.sh).

The `docopts.sh` is an extra bash library that you can source in your CLI script.
This library provides some bash helpers and is not required in order to use docopts.

You don't need a python interpreter anymore.

As of 2018-05-22

* `docopts` is able to reproduce 100% of the python version.
* unit test for go, hack as you wish.
* 100% `language_agnostic_tester.py` passed (GNU/Linux 64bits)

## Developpers

All python related stuff has been removed, excepted `language_agnostic_tester.py`.

If you want to clone this repository and hack docopts:

Use `git clone --recursive`, to get submodules only required for testing with `bats`.

Fetch the extra golang version of `docopt` (required for building `docopts`)

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
├── API_proposal.md - compatibility link to the proposed API.
├── build.sh - no more used yet
├── docopts.go - main source code
├── docopts.sh - library wrapper and helpers
├── examples - all examples are working
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
├── testee.sh - bash wrapper to convert docopts output to JSON
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
- `language_agnostic_tester.py` (old python wrapper, full docopt compatibily tests)
- See Also: docopt.go own tests in golang
- `docopts_test.go` go unit test for `docopts.go`

### Runing tests

#### bats
```
cd ./tests
. bats.alias
bats .
```

#### language_agnostic_tester

```
python language_agnostic_tester.py ./testee.sh
```

#### golang docopt.go (golang parser lib)

```
cd PATH/to/go/src/github.com/docopt/docopt-go/
go test -v .
```

#### golang docopts (our bash wrapper)

```
go test -v
```
