================================================================================
 docopts
================================================================================
--------------------------------------------------------------------------------
 shell interface for docopt, the CLI description language
--------------------------------------------------------------------------------
:Author:        `Lari Rasku <rasku@lavabit.com>`_
:Date:           2013-02-07
:Copyright:     `MIT <http://opensource.org/licenses/MIT>`_
:Version:        0.6.1+fix
:Manual section: 1

SYNOPSIS
================================================================================
``docopts`` [*options*] ``-h`` *msg* : [*argv*...]

DESCRIPTION
================================================================================
``docopts`` parses the command line argument vector *argv* according to the
`docopt <http://docopt.org>`_ string *msg* and echoes the results to standard
output as a snippet of Bash source code.  Passing this snippet as an argument to
`eval(1) <http://man.cx/eval(1)>`_ is sufficient for handling the CLI needs of
most scripts.

If *argv* matches one of the usage patterns defined in *msg*, ``docopts``
generates code for storing the parsed arguments as Bash variables.  As most
command line argument names are not valid Bash identifiers, some name mangling
will take place:

* ``<Angle_Brackets>``: ``Angle_Brackets``
* ``UPPER-CASE``: ``UPPER_CASE``
* ``--Long-Option``: ``Long_Option``
* ``-S``: ``S``

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

* Non-repeatable, valueless arguments: ``true`` if found, ``false`` if not
* Repeatable valueless arguments: the count of their instances in *argv*
* Non-repeatable arguments with values: the value as a string if found,
  the empty string if not
* Repeatable arguments with values: a Bash array of the parsed values

Unless the ``--no-help`` option is given, ``docopts`` handles the ``--help``
and ``--version`` options and their possible aliases specially,
generating code for printing the relevant message to standard output and
terminating successfully if either option is encountered when parsing *argv*.
Note however that this also requires listing the relevant option in
*msg* and, in ``--version``'s case, invoking ``docopts`` with the ``--version``
option.

If *argv* does not match any usage pattern in *msg*, ``docopts`` will generate
code for exiting the program with status 64 (``EX_USAGE`` in
`sysexits(3) <http://man.cx/sysexits(3)>`_) and printing a diagnostic error
message.

OPTIONS
================================================================================
  -h <msg>, --help=<msg>        The help message in docopt format.
                                If - is given, read the help message from
                                standard input.
                                If no argument is given, print docopts's own
                                help message and quit.
  -V <msg>, --version=<msg>     A version message.
                                If - is given, read the version message from
                                standard input.  If the help message is also
                                read from standard input, read it first.
                                If no argument is given, print docopts's own
                                version message and quit.
  -O, --options-first           Disallow interspersing options and positional
                                arguments: all arguments starting from the
                                first one that does not begin with a dash will
                                be treated as positional arguments.
  -H, --no-help                 Don't handle --help and --version specially.
  -A <name>                     Export the arguments as a Bash 4.x associative
                                array called <name>.
  -s <str>, --separator=<str>   The string to use to separate the help message
                                from the version message when both are given
                                via standard input. [default: ----]

EXAMPLES
================================================================================
Read the help and version messages from standard input::

    eval "$(docopts -V - -h - : "$@" <<EOF
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

Parse the help and version messages from script comments and pass them as
command line arguments::

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
    eval "$(docopts -h "$help" -V "$version" : "$@")"
    
    for arg in "${argv[@]}"; do
        echo "$arg"
    done

Using the associative array::

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
        echo "${args[<argument-with-multiple-values>,$i]}"
        i=$[$i+1]
    done

VERSIONING
================================================================================
The ``docopts`` version number always matches that of the
`docopt Python reference implementation <https://github.com/docopt/docopt>`_
version against which it was built.  As ``docopt`` follows
`semantic versioning <http://semver.org>`_, ``docopts`` should work with any
``docopt`` release it shares the major version number with; however, as both
``docopts`` and ``docopt`` are in major version number 0 at the moment of
writing this, ``docopts`` can only be relied to work with an installation of
``docopt`` with the exact same version number.
