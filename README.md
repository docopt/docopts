# `docopts` - shell interface for docopt, the CLI description language

## SYNOPSIS

    docopts <var> <doc> <version> -- [<argument>...]  
    docopts <var> <doc> -- [<argument>...]  
    docopts <var> -- [<argument>...]  
    docopts -- [<argument>...]  

## DESCRIPTION

`docopts` parses command line _arguments_ according to the _doc_ string and
echoes a [bash4(1)][] code snippet to standard output.  Passing this snippet to
[eval(1)][] will result in one of the four following effects:

- If one of the _arguments_ is `--help` or `-h` and _doc_ specifies such
  an option, the process echoes _doc_ to standard output and exits successfully.
- If one of the _arguments_ is `--version`, _doc_ specifies such an option,
  and `docopts` was invoked with a _version_ argument, the process echoes
  _version_ to standard output and exits successfully.
- If the _arguments_ do not match a valid usage pattern in _doc_, the process
  echoes an appropriate error message to standard error and exits with status
  64 (`EX_USAGE` in [sysexits(3)][].)
- If the _arguments_ match a valid usage pattern in _doc_, an associative
  array called _var_ is introduced to the process environment, mapping
  subcommand, argument and long option names defined in _doc_ to their
  parsed values.  The values are parsed as follows:
  
  - Subcommands and argumentless options will map to [true(1)][] if
    found and [false(1)][] if not.
  - Option-arguments and arguments accepting single values will map to
    their value if found and to the empty string if not.
  - Arguments accepting multiple values will be stored as fake nested arrays:
    ```
    ${args[ARG,#]} # the number of arguments to ARG
    ${args[ARG,0]} # the first argument to ARG
    ${args[ARG,1]} # the second argument to ARG, etc.
    ```

`docopts` expects _doc_ to be valid [docopt(7)][] text and _var_ to be a valid
[bash4(1)][] identifier.

If _doc_ is not given, it is read from standard input.  _version_ can also be
given via standard input by separating it from _doc_ with a sequence of four
dashes.  When both _doc_ and _version_ are given, whether as command line
arguments or via standard input, their order does not matter; `docopts`
considers the first string with a valid usage pattern to be _doc_.

If _var_ is not given, `docopts` is invoked in test mode, echoing
`"user-error"` on error and a [json(7)][] representation of the parsed
arguments on success, both to standard output.  The output is compatible
with [docopt(7)][]'s language agnostic test suite.

## OPTIONS

* `--help`:
  Show help options.
* `--version`:
  Print program version.

## EXAMPLES

Read _doc_ and _version_ from standard input:

    eval "$(docopts args -- $@ <<EOF
    rock 0.1.0
    Copyright (C) 200X  Thomas Light
    License RIT (Robot Institute of Technology)
    This is free software: you are free to change and redistribute it.
    There is NO WARRANTY, to the extent permitted by law.
    ----
    Usage: rock [options] _argument_...
    
          --help     Show help options.
          --version  Print program version.
    EOF
    )"

Parse _doc_ and _version_ from script comments and pass them as command line
arguments:

    ## rock 0.1.0
    ## Copyright (C) 200X  Thomas Light
    ## License RIT (Robot Institute of Technology)
    ## This is free software: you are free to change and redistribute it.
    ## There is NO WARRANTY, to the extent permitted by law.
    
    ### Usage: rock [options] _argument_...
    ### 
    ###       --help     Show help options.
    ###       --version  Print program version.
    
    help=$(grep "^### " "$0" | cut -c 5-)
    version=$(grep "^## "  "$0" | cut -c 4-)
    eval "$(docopts args "$help" "$version" -- $@)"

Using the associative array:

    eval "$(docopts args "$help" "$version" -- $@)"
    
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

## VERSIONING

The `docopts` version number always matches that of the [docopt(3)][] Python
reference implementation version which it was built against.  As [docopt(3)][]
follows semantic versioning, `docopts` should work with any [docopt(3)][]
release it shares the major version number with; however, as both `docopts` and
[docopt(3)][] are in major version number 0 at the moment of writing this
(2012-09-10), `docopts` can only be relied to work with the version of
[docopt(3)][] with the exact same version number.

## AUTHOR

[Lari Rasku](mailto:raskug@lavabit.com)

## REPORTING BUGS

Report bugs at <https://github.com/docopt/docopts/issues>.

## COPYRIGHT

Copyright (C) 2012 Vladimir Keleshev, Lari Rasku.
License [MIT](http://opensource.org/licenses/MIT).
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.

[bash4(1)]:    http://tldp.org/LDP/abs/html/bashver4.html
[docopt(3)]:   https://github.com/docopt/docopt
[docopt(7)]:   http://docopt.org
[json(7)]:     http://json.org
[sysexits(3)]: http://man.cx/sysexits
[eval(1)]:     http://man.cx/eval
[true(1)]:     http://man.cx/true
[false(1)]:    http://man.cx/false

