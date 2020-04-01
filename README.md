# docopts

[docopt](http://docopt.org/) for shell - make beautiful CLI with ease.

Status: working.

`docopts` : the command line wrapper for bash.

Most concepts are documented in the `docopt` (without S) manual - see [docopt.org](http://docopt.org/).

Many examples use associative arrays in bash 4.x, but there is legacy support for bash 3.2 on macOS (OS X) or legacy
GNU/Linux OS.

[make README.md]: # (./docopts --version | get_version "This is a transitional release:")

```
This is a transitional release: v0.6.3-rc1
```

This release will be maintained for compatibility, only fixes will be provided. The 0.6.3 version is fully compatible with
the previous version of `docopts`.

## SYNOPSIS

[make README.md]: # (./docopts --help | get_usage)

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
generates code for storing the parsed arguments as Bash variables. As most
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

This is the verbatim output of the `--help`:

[make README.md]: # (./docopts --help | sed -n '/^Options/,$ p')

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
  -G <prefix>                   Don't use associative array but output
                                Bash 3.2 compatible GLOBAL variables assignment:
                                  <prefix>_{mangled_args}={parsed_value}
                                Can be used with numeric incompatible options
                                as well.  See also: --no-mangle
  --no-mangle                   Output parsed option not suitable for bash eval.
                                Full option names are kept. Rvalue is still
                                shellquoted. Extra parsing is required.
  --no-declare                  Don't output 'declare -A <name>', used only
                                with -A argument.
  --debug                       Output extra parsing information for debugging.
                                Output cannot be used in bash eval.
```

## COMPATIBILITY

Bash 4.x and higher is the main target.

In order to use `docopts` with bash 3.2 (for macOS and old GNU/Linux versions) by avoiding bash 4.x associative arrays,
you can:

* don't use the `-A` option
* use GLOBAL generated mangled variables
* use `-G` `<prefix>` option to generate GLOBAL with `prefix_`
* use `source docopts.sh --auto -G` (see [example](examples/legacy_bash/sshdiff_with_docopts.sh))

The [`docopts.sh`](docopts.sh) helper allows the use of `set -u`, which
[gives an error](https://www.gnu.org/software/bash/manual/html_node/The-Set-Builtin.html#The-Set-Builtin)
on undefined variables in your scripts.

[Unofficial strict mode](http://redsymbol.net/articles/unofficial-bash-strict-mode/) for bash
should also work with `docopts.sh`, please report any [issue](https://github.com/docopt/docopts/issues) with examples.

## `docopts.sh` helper

The helper has its own documentation here [docs/README.md](docs/README.md).

## EXAMPLES

Find more examples in [examples/ folder](examples/).

This example reads the help and version messages from standard input (`docopts` found in `$PATH`):

[make README.md]: # (include examples/legacy_bash/rock_hello_world.sh)

[source examples/legacy_bash/rock_hello_world.sh](examples/legacy_bash/rock_hello_world.sh)

```bash
#!/usr/bin/env bash
# Example from README.md

PATH=$PATH:../..
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

The following example, parses the help and version messages from script comments and pass them as
command line arguments:

[make README.md]: # (include examples/legacy_bash/rock_hello_world_with_grep.sh)

[source examples/legacy_bash/rock_hello_world_with_grep.sh](examples/legacy_bash/rock_hello_world_with_grep.sh)

```bash
#!/usr/bin/env bash
# Example from README.md
PATH=$PATH:../..

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

The next example shows how using the Bash 4.x associative array with `-A`:

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

`docopts` was first developed by Lari Rasku <rasku@lavabit.com> and was written in Python based on the
[docopt Python parser](https://github.com/docopt/docopt).

The current version is written in [go](https://golang.org/) and is 100% compatible with previous Python-based `docopts`.
Please report any non working code with [issue](https://github.com/docopt/docopts/issues) and examples.

## Roadmap: A new shell API is proposed

Starting at release: 0.7.0 a new lib API based on JSON will be introduced:

See and contribute on the [docopts Wiki](https://github.com/docopt/docopts/wiki).

## Install

You only have to drop the binary, and optionally also the `docopts.sh` lib helper, in a directory on your PATH.
The binary is standalone and statically linked, so it runs everywhere.

See build section.

With root privileges you could do:

```
cp docopts docopts.sh /usr/local/bin
```

### Pre-built binaries

Pre-built Go binaries for GNU/Linux (32 and 64 bits) are attached to [releases](https://github.com/docopt/docopts/releases).

We provide a download helper:

```bash
git clone https://github.com/docopt/docopts.git
cd docopts
./get_docopts.sh
```

Rename to `docopts` and put it in your `PATH`:

```bash
mv docopts_linux_amd64 docopts
cp docopts docopts.sh /usr/local/bin
```

The cloned repository is no more used at this stage. 

Learn more about [pre-built binaries](docs/pre_built_binaries.md).

## Compiling

We encourage you to build your own binary, which is easy once
you have Go installed. Or find a local golang developer that you
trust and ask her, in exchange for a beer or two, if she could build it for you. ;)

Requires a [Go workspace](https://golang.org/doc/code.html#Organization).

local build:
(also done with our Makefile default target: `make`)

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

or via Makefile:

```
cd src/github.com/docopt/docopts
make all
make test
```

Tested builds are built on: 

[make README.md]: # (go version)

```
go version go1.11.4 linux/amd64
```

## Features

Warning: may be not up-to-date feature list.

The [`docopts.sh`](docopts.sh) helper is an extra bash library that you can source in your shell script.
This library provides some bash helpers and is not required in order to use `docopts`. See [docopts.sh
documentation](docs/README.md).

`docopts` doesn't need a python interpreter anymore, so it works on any legacy system too.

As of 2019-05-18

* `docopts` is able to reproduce 100% of the python version.
* unit tests for Go are provided, so hack as you wish.
* 100% of `language_agnostic_tester.py` tests pass (GNU/Linux 64bits).
* `bats-core` unittests and fonctional testing are provided too. 

## Developers

Read the [doc for developer](docs/developer.md).
