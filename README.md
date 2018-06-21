# docopts

docopt for shell - make beautifull CLI with ease.

Status: Alpha - work is done.

This is the golang version of `docopts` : the command line wrapper for bash.

See the Reference manual for how to use `docopts` [old README.rst](old_README.rst).

This is a transitional release: 0.6.2

This release will be maintained for compatibity, only fix will be provided.

The [develop branch](https://github.com/docopt/docopts/tree/develop) is abandoned!

The current go version is 100% compatible with python's `docopts`.
Please report any non working code with [issue](https://github.com/docopt/docopts/issues) and examples.

## A new shell API is proposed

Starting at release: 0.7.0 a new lib API based on JSON will be introduced:

See and contribute on the [docopts Wiki](https://github.com/docopt/docopts/wiki).

## Install

You only have to drop the binary, and eventually the `docopts.sh` lib helper, in your PATH.
The binary is standalone and staticaly linked. So it runs everywhere.

See build section.

With root privileges you could:

```
cp docopts docopts.sh /usr/local/bin
```

### pre-built binary

Pre-built binary are attached to [releases](https://github.com/Sylvain303/docopts/releases).

Download and rename it as `docopts` and put it in your `PATH`.

```
mv docopts-32bit docopts
cp docopts docopts.sh /usr/local/bin
```

You are strongly encouraged to build your own binary. Find a local golang developper in whom you trust and ask her, in exchange of a beer or two, if she could build it for you. ;)

## Usage

See [manual](old_README.rst) and [Examples](examples/)

## Compiling

With a go workspace.

local build

```
go get github.com/docopt/docopt-go
go get github.com/docopt/docopts
cd src/github.com/docopt/docopts
go build docopts.go
```

cross compile for 32btis

```
env GOOS=linux GOARCH=386 go build docopts.go
```

or via Makefile (generate 64bits, 32bits, arm and OSX-64bits version of docopts)

```
cd src/github.com/docopt/docopts
make all
make test
```

Tested built on: `go version go1.10.2 linux/amd64`

## Features

Warning: may be not up-to-date feature list.

The [`docopts.sh`](docopts.sh) is an extra bash library that you can source in your CLI script.
This library provides some bash helpers and is not required in order to use `docopts`.

You don't need a python interpreter anymore, so it works on any legacy system.

As of 2018-05-22

* `docopts` is able to reproduce 100% of the python version.
* unit test for go are provided, so hack as you wish.
* 100% `language_agnostic_tester.py` passed (GNU/Linux 64bits)

## Developpers

All python related stuff has been removed, excepted `language_agnostic_tester.py`.

If you want to clone this repository and hack `docopts`:

Use `git clone --recursive`, to get submodules only required for testing with `bats`.
[`bats`](https://github.com/sstephenson/bats) or [`bats-core`](https://github.com/bats-core/bats-core) installed in your
PATH should work too.

Fetch the extra golang version of `docopt-go` (required for building `docopts`)

```
go get github.com/docopt/docopt-go
```

If you forgot `--recursive`, you can also run it afterward:

~~~bash
git submodule init
git submodule update
~~~

Current folder structure:

~~~
.
├── docopts.go - main source code
├── docopts_test.go - go unit tests
├── docopts.sh - library wrapper and helpers
├── examples - many ported examples in bash, all must be working
├── language_agnostic_tester.py - old python JSON tester still used with testee.sh
├── LICENSE-MIT
├── old_README.rst - copied README
├── PROGRESS.mda - what I'm working on
├── README.md
├── testcases.docopt - agnostic testcases copied from python's docopt
├── testee.sh - bash wrapper to convert docopts output to JSON (now use docopts.sh)
├── tests - unit and functional testing written in bats (require submodule)
└── TODO.md - Some todo list on this golang version of docopts
~~~

## Tests

Some tests are coded along this code base.

- bats bash unit tests and functionnal testing.
- `language_agnostic_tester.py` (old python wrapper, full docopt compatibily tests)
- See Also: [docopt.go](https://github.com/docopt/docopts-go/) own tests in golang.
- `docopts_test.go` go unit test for `docopts.go`

### Runing tests

```
make test
```

#### bats

```
cd ./tests
. bats.alias
bats .
```

#### language_agnostic_tester

This script was provided with the original `docopts`. I fixed number/string output parsing failure with an extra function
for bash in [docopts.sh](https://github.com/docopt/docopts/blob/docopts-go/docopts.sh#L108)
`docopt_get_raw_value()`. This is a hack to get 100% pass, and it is not very efficient.

```
python language_agnostic_tester.py ./testee.sh
```

#### golang docopt.go (golang parser lib)

This lib is outside this project, but it is the base of the `docopt` parsing for this wrapper.

```
cd PATH/to/go/src/github.com/docopt/docopt-go/
go test -v .
```

#### golang docopts (our bash wrapper)

```
go test -v
```
