================================================================================
 docopts
================================================================================
--------------------------------------------------------------------------------
 shell interface for docopt, the CLI description language
--------------------------------------------------------------------------------
:Author:        `Lari Rasku <rasku@lavabit.com>`_
:Date:           2012-10-09
:Copyright:     `MIT <http://opensource.org/licenses/MIT>`_
:Version:        0.5.0
:Manual section: 1

SYNOPSIS
================================================================================
``docopts`` [*options*] *doc* *version* [-- *argv*...]

DESCRIPTION
================================================================================
``docopts`` parses the command line argument vector *argv* according to the
`docopt <http://docopt.org>`_ string *doc* and echoes the results to standard
output as a snippet of Bash source code.  Passing this snippet as an argument to
`eval(1) <http://man.cx/eval(1)>`_ is sufficient for handling the CLI needs of
most scripts.

If *argv* matches one of the usage patterns defined in *doc*, ``docopts``
generates code for storing the parsed arguments as Bash variables.  As most
command line argument names are not valid Bash identifiers, some name mangling
will take place:

* ``<angle-brackets>``: ``angle_brackets``
* ``UPPER_CASE``: ``upper_case``
* ``--long-option``: ``long_option``
* ``-S``: ``s``

If one of the argument names cannot be mangled into a valid Bash identifier,
or two argument names map to the same variable name, ``docopt`` will exit with
an error, and you should really rethink your CLI.  The ``--`` and ``-``
commands will not be stored.

Alternatively, ``docopts`` can be invoked with the ``-A <name>`` option, which
stores the parsed arguments as fields of a Bash 4 associative array called
``<name>`` instead.  However, as Bash does not natively support nested arrays,
they are faked for repeatable arguments with the following access syntax::

    ${args[ARG,#]} # the number of arguments to ARG
    ${args[ARG,0]} # the first argument to ARG
    ${args[ARG,1]} # the second argument to ARG, etc.

The arguments are stored as follows:

* Non-repeatable, valueless arguments: `true(1) <http://man.cx/true(1)>`_
  if found, `false(1) <http://man.cx/false(1)>`_ if not
* Repeatable valueless arguments: the count of their instances in *argv*
* Non-repeatable arguments with values: the value as a string if found,
  the empty string if not
* Repeatable arguments with values: a Bash array of the parsed values

Unless the ``--no-help`` option is given, ``docopts`` handles the ``--help``
and ``--version`` options and their possible aliases specially,
generating code for printing the relevant message to standard output and
terminating successfully if either option is encountered when parsing *argv*.
Note however that this also requires listing the relevant option in
*doc*, and in ``--version``'s case, invoking ``docopts`` with a non-empty
*version* string.

If *argv* does not match any usage pattern in *doc*, ``docopts`` will generate
code for exiting the program with status 64 (``EX_USAGE`` in
`sysexits(3) <http://man.cx/sysexits(3)>`_) and printing a diagnostic error
message.

ARGUMENTS
================================================================================
:doc:                           The help message in docopt format.  If ``-`` is
                                given, read the help message from standard
                                input.
:version:                       A version message.  If an empty argument is
                                given via ``''``, no version message is used.
                                If ``-`` is given, the version message is read
                                from standard input.  The version message is
                                read after the help message if both are given
                                via standard input.

OPTIONS
================================================================================
  -A <name>                     Export the arguments as a Bash 4.x associative
                                array called *name*.
  -s <sep>, --separator=<sep>   The string to use to separate *doc* from
                                *version* when both are given via standard
                                input [default: ``----``]
  -H, --no-help                 Don't handle ``--help`` and ``--version``
                                specially.
  -d, --debug                   Export the arguments as JSON.  The docopt
                                language agnostic test suite can be run for
                                docopts by passing ``docopts -d - '' --`` as
                                its argument.
  -h, --help                    Show help options.
  -V, --version                 Print program version.

EXAMPLES
================================================================================
Read *doc* and *version* from standard input::

    eval "$(docopts - - -- "$@" <<EOF
    Usage: rock [options] <argv>...
    
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

Parse *doc* and *version* from script comments and pass them as command line
arguments::

    ## rock 0.1.0
    ## Copyright (C) 200X Thomas Light
    ## License RIT (Robot Institute of Technology)
    ## This is free software: you are free to change and redistribute it.
    ## There is NO WARRANTY, to the extent permitted by law.
    
    ### Usage: rock [options] <argv>...
    ### 
    ###       --help     Show help options.
    ###       --version  Print program version.
    
    help=$(grep "^### " "$0" | cut -c 5-)
    version=$(grep "^## "  "$0" | cut -c 4-)
    eval "$(docopts "$help" "$version" -- "$@")"
    
    for arg in "${argv[@]}"; do
        # do something
    done

Using the associative array::

    eval "$(docopts -A args "$help" "" -- "$@")"
    
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
        echo "${args[<argument-with-multiple-values>,$i]}"
        i=$[$i+1]
    done

INSTALLATION
================================================================================
To install ``docopts`` for every user, extract the release archive and execute
the following command in it::

    python setup.py install

To install ``docopts`` just for you, use this instead::

    python setup.py install --user

Alternatively, you can simply copy the ``docopts`` file to anywhere on your
``PATH``; it is self-contained.

VERSIONING
================================================================================
The ``docopts`` version number always matches that of the
`docopt Python reference implementation <https://github.com/docopt/docopt>`_
version against which it was built.  As ``docopt`` follows semantic versioning,
``docopts`` should work with any ``docopt`` release it shares the major version
number with; however, as both ``docopts`` and ``docopt`` are in major version
number 0 at the moment of writing this (2012-10-09), ``docopts`` can only be
relied to work with an installation of ``docopt`` with the exact same version
number.
