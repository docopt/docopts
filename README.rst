``docopts``: shell interface for docopt, the CLI description language
=====================================================================

Synopses::

    docopts <var> <doc> <version> -- [<argument>...]
    docopts <var> <doc> -- [<argument>...]
    docopts <var> -- [<argument>...]
    docopts -- [<argument>...]

``docopts`` parses command line *argument*\s according to the *doc* string and
echoes a `Bash 4.x`_ code snippet to standard output.  Passing this snippet to
`eval`_ will result in one of the four following effects:

- If one of the *argument*\s is ``--help`` or ``-h`` and *doc* specifies such
  an option, the process echoes *doc* to standard output and exits successfully.
- If one of the *argument*\s is ``--version``, *doc* specifies such an option,
  and ``docopts`` was invoked with a *version* argument, the process echoes
  *version* to standard output and exits successfully.
- If the *argument*\s do not match a valid usage pattern in *doc*, the process
  echoes an appropriate error message to standard error and exits with status
  64 (``EX_USAGE`` in `sysexits.h`_.)
- If the *argument*\s match a valid usage pattern in *doc*, an associative
  array called *var* is introduced to the process environment, mapping
  subcommand, argument and long option names defined in *doc* to their
  parsed values.  The values are parsed as follows:
  
  - Subcommands and argumentless options will map to the program `true`_ if
    found and `false`_ if not.
  - Option-arguments and arguments accepting single values will map to
    their value if found and to the empty string if not.
  - Arguments accepting multiple values will be stored as fake nested arrays::
    
        ${args[ARG,#]} # the number of arguments to ARG
        ${args[ARG,1]} # the first argument to ARG
        ${args[ARG,2]} # the argument to ARG, etc.

``docopts`` expects *doc* to be valid `docopt`_ text and *var* to be a valid
`Bash 4.x`_ identifier.

If *doc* is not given, it is read from standard input.  *version* can also be
given via standard input by separating it from *doc* with a sequence of four
dashes.  When both *doc* and *version* are given, whether as command line
arguments or via standard input, their order does not matter; ``docopts``
considers the first string with a valid usage pattern to be *doc*.

If *var* is not given, ``docopts`` is invoked in test mode, echoing
``"user-error"`` on error and a `JSON`_ representation of the parsed
arguments on success, both to standard output.  The output is compatible
with `docopt`_'s language agnostic test suite.

Options
-------
--help     Show help options.
--version  Print program version.

Examples
--------
Read *doc* and *version* from standard input::

    eval "$(docopts args -- $@ <<EOF
    rock 0.1.0
    Copyright (C) 200X  Thomas Light
    License RIT (Robot Institute of Technology)
    This is free software: you are free to change and redistribute it.
    There is NO WARRANTY, to the extent permitted by law.
    ----
    Usage: rock [options] <argument>...
    
          --help     Show help options.
          --version  Print program version.
    EOF
    )"

Parse *doc* and *version* from script comments and pass them as command line
arguments::

    ## rock 0.1.0
    ## Copyright (C) 200X  Thomas Light
    ## License RIT (Robot Institute of Technology)
    ## This is free software: you are free to change and redistribute it.
    ## There is NO WARRANTY, to the extent permitted by law.
    
    ### Usage: rock [options] <argument>...
    ### 
    ###       --help     Show help options.
    ###       --version  Print program version.
    
    help=$(grep "^### " "$0" | cut -c 5-)
    version=$(grep "^## "  "$0" | cut -c 4-)
    eval "$(docopts args "$help" "$version" -- $@)"

Using the associative array::

    eval "$(docopts args "$help" "$version" -- $@)"
    
    if ${args[subcommand]} ; then
        echo "subcommand was given"
    fi
    
    if [ -n "${args[--long-option-with-argument]}" ] ; then
        echo "${args[--long-option-with-argument]}"
    else
        echo "--long-option-with-argument was not given"
    fi
    
    for i in $(seq 1 ${args[<argument-with-multiple-values>,#]}) ; do
        echo "${args[<argument-with-multiple-values>,$i]}"
    done

Installing
----------
To install from source, execute the following command in the release directory::

    python setup.py install

Versioning
----------
The ``docopts`` version number always matches that of the `docopt`_ Python
reference implementation version which it was built against.  As `docopt`_
follows semantic versioning, ``docopts`` should work with any `docopt`_ release
it shares the major version number with; however, as both ``docopts`` and
`docopt`_ are in major version number 0 at the moment of writing this
(2012-08-08), ``docopts`` can only be relied to work with the version of
`docopt`_ with the exact same version number.



.. _Bash 4.x:   http://tldp.org/LDP/abs/html/bashver4.html
.. _docopt:     http://docopt.org
.. _JSON:       http://json.org
.. _sysexits.h: http://man.cx/sysexits
.. _eval:       http://man.cx/eval
.. _true:       http://man.cx/true
.. _false:      http://man.cx/false
