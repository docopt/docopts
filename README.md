# docopts

[docopt](http://docopt.org/) for shell - make beautiful CLI with ease.

Status: Alpha - work is done.

`docopts` : the command line wrapper for bash.

Most concepts are documented in the `docopt` (without S) manual - see [docopt.org](http://docopt.org/).

Many examples use associative arrays in bash 4.x, but there is legacy support for bash 3.2 on macOS (OSX) or legacy
GNU/Linux OS.

This is a transitional release: 0.6.2

This release will be maintained for compatibility, only fixes will be provided. The 0.6.2 version is fully compatible with
the previous version of `docopts`.

## SYNOPSIS

```
  docopts [options] -h <msg> : [<argv>...]
  docopts [options] [--no-declare] -A <name>   -h <msg> : [<argv>...]
  docopts [options] -G <prefix>  -h <msg> : [<argv>...]
  docopts [options] --no-mangle  -h <msg> : [<argv>...]
```

## DESCRIPTION

`docopts` parses the command line argument vector `<argv>` according to the
[docopt](http://docopt.org) string `<msg>` and echoes the results to standard
output as a snippet of Bash source code.  Passing this snippet as an argument to
[`eval(1)`](http://man.cx/eval(1)) is sufficient for handling the CLI needs of
most scripts.

If `<argv>` matches one of the usage patterns defined in `<msg>`, `docopts`
generates code for storing the parsed arguments as Bash variables.  As most
command line argument names are not valid Bash identifiers, some name mangling
will take place:

* `<Angle_Brackets>` ==> `Angle_Brackets`
* `UPPER-CASE` ==> `UPPER_CASE`
* `--Long-Option` ==> `Long_Option`
* `-S` ==> `S`
* `-4` ==> **INVALID** (without -G)

If one of the argument names cannot be mangled into a valid Bash identifier,
or two argument names map to the same variable name, `docopts` will exit with
an error, and you should really rethink your CLI, or use `-A` or `-G`.
The `--` and `-` commands will not be stored.

Note: You can use `--no-mangle` if you still want full input, this wont
produce output suitable for bash `eval(1)` but can be parsed by your own
code.

Alternatively, `docopts` can be invoked with the `-A <name>` option, which
stores the parsed arguments as fields of a Bash 4 associative array called
`<name>` instead.  However, as Bash does not natively support nested arrays,
they are faked for repeatable arguments with the following access syntax:

```
    ${args[ARG,#]} # the number of arguments to ARG
    ${args[ARG,0]} # the first argument to ARG
    ${args[ARG,1]} # the second argument to ARG, etc.
```

The arguments are stored as follows:

* Non-repeatable, valueless arguments: `true` if found, `false` if not
* Repeatable valueless arguments: the count of their instances in `<argv>`
* Non-repeatable arguments with values: the value as a string if found,
  the empty string if not
* Repeatable arguments with values: a Bash array of the parsed values

Unless the `--no-help` option is given, `docopts` handles the `--help`
and `--version` options and their possible aliases specially,
generating code for printing the relevant message to standard output and
terminating successfully if either option is encountered when parsing `<argv>`.

Note however that this also requires listing the relevant option in
`<msg>` and, in `--version`'s case, invoking `docopts` with the `--version`
option.

If `<argv>` does not match any usage pattern in `<msg>`, `docopts` will generate
code for exiting the program with status 64 [`EX_USAGE` in `sysexits(3)`](http://man.cx/sysexits(3))
and printing a diagnostic error message.

Note that due to the above, `docopts` can't be used to parse shell function
arguments: [`exit(1)`](http://man.cx/exit(1)) quits the entire interpreter,
not just the current function.

## OPTIONS

This is the verbatim output of `docopts --help`:

```
Options:
  -h <msg>, --help=<msg>        The help message in docopt format.
                                Without argument outputs this help.
                                If - is given, read the help message from
                                standard input.
                                If no argument is given, print docopts's own
                                help message and quit.
  -V <msg>, --version=<msg>     A version message.
                                If - is given, read the version message from
                                standard input.  If the help message is also
                                read from standard input, it is read first.
                                If no argument is given, print docopts's own
                                version message and quit.
  -s <str>, --separator=<str>   The string to use to separate the help message
                                from the version message when both are given
                                via standard input. [default: ----]
  -O, --options-first           Disallow interspersing options and positional
                                arguments: all arguments starting from the
                                first one that does not begin with a dash will
                                be treated as positional arguments.
  -H, --no-help                 Don't handle --help and --version specially.
  -A <name>                     Export the arguments as a Bash 4.x associative
                                array called <name>.
  -G <prefix>                   As with -A, but outputs Bash 3.x compatible
                                GLOBAL variables assignment, using the given
                                <prefix>_{option}={parsed_option}. Can be used
                                with numerical incompatible option as well.
                                See also: --no-mangle
  --no-mangle                   Output parsed option not suitable for bash eval.
                                As without -A but full option names are kept.
                                Rvalue is still shellquoted.
  --no-declare                  Don't output 'declare -A <name>', used only
                                with -A argument.
  --debug                       Output extra parsing information for debuging.
                                Output cannot be used in bash eval.
```

## EXAMPLES

More examples in [examples/ folder](examples/).

Read the help and version messages from standard input (`docopts` found in `$PATH`):

[examples/legacy_bash/rock_hello_world.sh](examples/legacy_bash/rock_hello_world.sh)

```bash
eval "$(docopts -V - -h - : "$@" <<EOF
Usage: rock [options] <argv>...

Options:
      --verbose  Generate verbose messages.
      --help     Show help options.
      --version  Print program version.
----
rock 0.1.0
Copyright (C) 200X Thomas Light
License RIT (Robot Institute of Technology)
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
EOF
)"

if $verbose ; then
    echo "Hello, world!"
fi
```

Parse the help and version messages from script comments and pass them as
command line arguments:

[examples/legacy_bash/rock_hello_world_with_grep.sh](examples/legacy_bash/rock_hello_world_with_grep.sh)

```bash
#? rock 0.1.0
#? Copyright (C) 200X Thomas Light
#? License RIT (Robot Institute of Technology)
#? This is free software: you are free to change and redistribute it.
#? There is NO WARRANTY, to the extent permitted by law.

##? Usage: rock [options] <argv>...
##?
##? Options:
##?       --help     Show help options.
##?       --version  Print program version.

help=$(grep "^##?" "$0" | cut -c 5-)
version=$(grep "^#?"  "$0" | cut -c 4-)
eval "$(docopts -h "$help" -V "$version" : "$@")"

for arg in "${argv[@]}"; do
    echo "$arg"
done
```

Using the Bash 4.x associative array with `-A`:

```bash
help="
Usage: example [--long-option-with-argument=value] <argument-with-multiple-values>...
"
eval "$(docopts -A args -h "$help" : "$@")"

if ${args[subcommand]} ; then
    echo "subcommand was given"
fi

if [ -n "${args[--long-option-with-argument]}" ] ; then
    echo "${args[--long-option-with-argument]}"
else
    echo "--long-option-with-argument was not given"
fi

i=0
while [[ $i -lt ${args[<argument-with-multiple-values>,#]} ]] ; do
    echo "$i: ${args[<argument-with-multiple-values>,$i]}"
    i=$[$i+1]
done
```


## History

`docopts` was first developed by Lari Rasku <rasku@lavabit.com> and was written in Python based on the Python parser.

The current version is written in [go](https://golang.org/) and is 100% compatible with previous Python-based `docopts`.
Please report any non working code with [issue](https://github.com/docopt/docopts/issues) and examples.

## Roadmap: A new shell API is proposed

Starting at release: 0.7.0 a new lib API based on JSON will be introduced:

See and contribute on the [docopts Wiki](https://github.com/docopt/docopts/wiki).

## Install

You only have to drop the binary, and eventually the `docopts.sh` lib helper, in your PATH.
The binary is standalone and staticaly linked. So it runs everywhere.

See build section.

With root privileges you could do:

```
cp docopts docopts.sh /usr/local/bin
```

### Pre-built binaries

Pre-built Go binaries for GNU/Linux (32 and 64 bits) are attached to [releases](https://github.com/Sylvain303/docopts/releases).

Download and rename to `docopts` and put it in your `PATH`:

```
mv docopts-32bit docopts
cp docopts docopts.sh /usr/local/bin
```

You are strongly encouraged to build your own binary. Find a local golang developper in whom you trust and ask her, in exchange for a beer or two, if she could build it for you. ;)


## Compiling

Requires a directory to use as [Go workspace](https://golang.org/doc/code.html#Organization).

local build:

```
go get github.com/docopt/docopt-go
go get github.com/docopt/docopts
cd src/github.com/docopt/docopts
go build docopts.go
```

cross compile for 32 bit:

```
env GOOS=linux GOARCH=386 go build docopts.go
```

or via Makefile (generates 64 bit, 32 bit, arm and macOS-64bit versions of docopts)

```
cd src/github.com/docopt/docopts
make all
make test
```

Tested builds are built on: `go version go1.11.4 linux/amd64`

## Features

Warning: may be not up-to-date feature list.

The [`docopts.sh`](docopts.sh) is an extra bash library that you can source in your CLI script.
This library provides some bash helpers and is not required in order to use `docopts`.

You don't need a python interpreter anymore, so it works on any legacy system.

As of 2019-05-18

* `docopts` is able to reproduce 100% of the python version.
* unit test for go are provided, so hack as you wish.
* 100% `language_agnostic_tester.py` passed (GNU/Linux 64bits)

## Developers

All python related stuff has been removed, excepted `language_agnostic_tester.py`.

If you want to clone this repository and hack `docopts`:

Use `git clone --recursive`, to get submodules - these are only required for testing with `bats`.
[`bats`](https://github.com/sstephenson/bats) or [`bats-core`](https://github.com/bats-core/bats-core) installed in your
PATH should work too.

Fetch the extra golang version of `docopt-go` (required for building `docopts`)

```
go get github.com/docopt/docopt-go
```

If you forgot `--recursive`, you can also run it afterwards:

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
├── PROGRESS.md - what I'm working on
├── README.md
├── testcases.docopt - agnostic testcases copied from python's docopt
├── testee.sh - bash wrapper to convert docopts output to JSON (now uses docopts.sh)
├── tests - unit and functional testing written in bats (requires submodule)
└── TODO.md - Some todo list on this golang version of docopts
~~~

## Tests

Some tests are coded along with this code base:

- bats - bash unit tests and functional testing.
- `language_agnostic_tester.py` - old python wrapper, full docopt compatibility tests.
- See also: [docopt.go](https://github.com/docopt/docopt.go) has its own tests in golang.
- `docopts_test.go` - go unit test for `docopts.go`

### Running tests

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

Run these tests from top of repo:
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
cd PATH/to/go/src/github.com/docopt/docopts
go test -v
```
