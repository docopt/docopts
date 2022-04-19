# docopt Grammar

// extract the Usage: section (case-insensitive) until the next empty line

## Goal

* We start with our `lexer_state` before testing handcrafted optimized lexer/scanner
* parse actual grammar (to be determined by studying docopt-go)
* better error handling (See Also: [`Error_reporting.md`](Error_reporting.md))
* improve what needed for better support: `print ast` `analyze` `explain` `convert`

## Choices

* hand crafted parser for better error handling

## Current work

(our branch grammar's PROGRESS.md)

=> start building ast from our grammar: `docopt_language.ebnf`
* matching engine
 * parse AST for match
 * explore: parseLong() from docopt-legacy
 * compare recieved os.Args to valid argument or option
* parse Options: default value
* add some automated testing and unit testing
* Refactor String() method to our token with a go generate? => can use GoString()
* full error handling support  and reporting
* Update ebnf language (add link to ebnf description) [`grammar/docopt_language.ebnf`](grammar/docopt_language.ebnf)
* add unit test

## TODO
* container/list (double linked list) handle parsed token list?
* code profiling
* fuzzy wrong option report suggestion
* bash completion (See cod parser)


## issues to read

* https://github.com/docopt/docopts/issues/35#issuecomment-516356005 explore docopt-ng
* https://github.com/docopt/docopt/issues/190
* https://github.com/docopt/docopt.go/issues/61
* https://github.com/docopt/docopts/issues/17

## Lexer

* ~~must handle syntax error (return ERROR token + state)~~
* ~~must progress one token at time~~
* ~~must stop scanning after a NEWLINE empty line with `Usage_line`~~ done within the parser calling our lexer
* ~~compare with `docopt-go/docopt.go` `ParseSection(name, source string)`~~ (See: compare/parseSection.go)

## example of parser / AST

* https://golangdocs.com/golang-parser-package
* https://ruslanspivak.com/lsbasi-part7/ - building an AST in python

## outputs proposal

some extra command for docopts:

### `analyze`

```
Usage:
  docopts analyze [--full|--ast] [--only=<usage_entry>] <string_usage>
  docopts analyze [--full|--ast] [--only=<usage_entry>] -f <filename_usage>
```

Example:

Here self help message parsing.
```
docopts analyze "$(docopts -h)"

5 alternatives 'Usages:' for 'docopts'

use --full to expand all details and shortcut.
use --only=<number> to display only one Usage.

1.  docopts [options] [--docopt_sh] -h <msg> : [<argv>...]

   [options]        Optional:   'shortcut' any options as defined in 'Options:' section
   [--docopt_sh]    Optional:   'boolean' default false
   -h <msg>         Mandatory:  'string' stored in <msg> argcount: 1
                                 Alt syntax: --help=<msg>
                                 documented: yes
   :                Mandatory:  'command'
   <argv>           Optional:   'argument' oneormore


2.  docopts [options] -G <prefix> [--docopt_sh] -h <msg> : [<argv>...]

   [options]        Optional:   'shortcut' any options as defined in 'Options:' section
   -G <prefix>      Mandatory:  'string' stored in <prefix> argcount: 1
                                 documented: yes
   [--docopt_sh]    Optional:   'boolean' default false
                                 documented: yes
   -h <msg>         Mandatory:  'string' stored in <msg> argcount: 1
                                 Alt syntax: --help=<msg>
                                 documented: yes
   :                Mandatory:  'command'
   <argv>           Optional:   'argument' oneormore

3.  docopts [options] --no-mangle  -h <msg> : [<argv>...]

   [options]        Optional:   'shortcut' any options as defined in 'Options:' section
   --no-mangle      Mandatory:  'boolean' default false
                                 documented: yes
   -h <msg>         Mandatory:  'string' stored in <msg> argcount: 1
                                 Alt syntax: --help=<msg>
                                 documented: yes
   :                Mandatory:  'command'
   <argv>           Optional:   'argument' oneormore

4.  docopts [options] [--no-declare] -A <name>   -h <msg> : [<argv>...]

   [options]        Optional:   'shortcut' any options as defined in 'Options:' section
   [--no-declare]   Optional:   'boolean' default false
                                 documented: yes
   -A <name>        Mandatory:  'string' stored in <name> argcount: 1
   -h <msg>         Mandatory:  'string' stored in <msg> argcount: 1
                                 documented: yes
   :                Mandatory:  'command'
   <argv>           Optional:   'argument' oneormore

5.  docopts [options] (--print-ast|--print-pat-fix|--print-parsed) -h <msg> : [<argv>...]

   [options]        Optional:   'shortcut' any options as defined in 'Options:' section
   --print-ast      Choice:     'boolean' default false
                                 documented: no
   --print-pat-fix  Choice:     'boolean' default false
                                 documented: no
   --print-parsed   Choice:     'boolean' default false
                                 documented: no
   -h <msg>         Mandatory:  'string' stored in <msg> argcount: 1
                                 documented: yes
   :                Mandatory:  'command'
   <argv>           Optional:   'argument' oneormore


Options:
  [...] TODO describe parsed options here

  -h <msg>, --help=<msg>   require argument: yes
                           documentation: "The help message in docopt format. Without[... 197 chars/6 lines]"
```

### ast

Some ast in textual format:
* https://docs.python.org/3/library/ast.html
* https://eleni.blog/2020/06/27/abstract-syntax-trees-asts-in-go/
* https://tech.ingrid.com/introduction-ast-golang/
* https://pypi.org/project/astboom/
* https://github.com/etetoolkit/ete/

our ast example for one Usage:
```
docopts [options] (--print-ast|--print-pat-fix|--print-parsed) -h <msg> : [<argv>...]
   optional: (optionsshortcut)
   required: (either)
                     (option) : "--print-ast"
                     (option) : "--print-pat-fix"
                     (option) : "--print-parsed"
   required: (option)  : "--help", argcount: 1
   required: (command) : ":"
   optional: oneormore (argument) : "<argv>"
```

http://lrv.bplaced.net/syntaxtree/?

```
[docopts
  [usage
   [
    [optional: optionsshortcut]
    [required: [either
                     [ option ["--print-ast" value: false]]
                     [ option ["--print-pat-fix" value: false]]
                     [ option ["--print-parsed" value: false]]
               ]]
   [required [option  "--help", argcount: 1]]
   [required [command ":", value: false]]
   [optional [oneormore [argument  "<argv>"]] ]
  ]
]
]
```

https://astexplorer.net/
```yaml
docopts:
  prologue: []
  usage:
    - optional: optionsshortcut
    - required:
        either:
          - option:
              name: "--print-ast"
          - option:
              name: "--print-pat-fix"
          - option:
              name: "--print-parsed"
    - required:
        option:
           name: "--help"
           argcount: 1
    - required:
        command:
           name: ":"
    - optional:
        oneormore:
          argument:
            name: "<argv>"
```

### `Simple_print_tree`

It doesn't display descriptive node content.

```vim
:r! cmd/docopt-analyze/docopt-analyze -s grammar/usages/valid/docopts.docopt
```

parsing: grammar/usages/valid/docopts.docopt
Root [5]
  Prologue [17]
  Usage_section [5]
    Usage "Usage:"
    Usage_line [2]
      Prog_name "docopts"
      Usage_Expr [5]
        Usage_optional_group [1]
          Usage_Expr [1]
            Usage_command "options"
        Usage_short_option "-h"
        Usage_argument "<msg>"
        Usage_argument ":"
        Usage_optional_group [1]
          Usage_Expr [1]
            Usage_argument "<argv>"...
    Usage_line [2]
      Prog_name "docopts"
      Usage_Expr [8]
        Usage_optional_group [1]
          Usage_Expr [1]
            Usage_command "options"
        Usage_optional_group [1]
          Usage_Expr [1]
            Usage_long_option "--no-declare"
        Usage_short_option "-A"
        Usage_argument "<name>"
        Usage_short_option "-h"
        Usage_argument "<msg>"
        Usage_argument ":"
        Usage_optional_group [1]
          Usage_Expr [1]
            Usage_argument "<argv>"...
    Usage_line [2]
      Prog_name "docopts"
      Usage_Expr [7]
        Usage_optional_group [1]
          Usage_Expr [1]
            Usage_command "options"
        Usage_short_option "-G"
        Usage_argument "<prefix>"
        Usage_short_option "-h"
        Usage_argument "<msg>"
        Usage_argument ":"
        Usage_optional_group [1]
          Usage_Expr [1]
            Usage_argument "<argv>"...
    Usage_line [2]
      Prog_name "docopts"
      Usage_Expr [6]
        Usage_optional_group [1]
          Usage_Expr [1]
            Usage_command "options"
        Usage_long_option "--no-mangle"
        Usage_short_option "-h"
        Usage_argument "<msg>"
        Usage_argument ":"
        Usage_optional_group [1]
          Usage_Expr [1]
            Usage_argument "<argv>"...
  Free_section
  Options_section [13]
    Section_name "Options:"
    Option_line [3]
      Option_short "-h" [1]
        Option_argument "<msg>"
      Option_long "--help" [1]
        Option_argument "<msg>"
      Option_description [17]
    Option_line [3]
      Option_short "-V" [1]
        Option_argument "<msg>"
      Option_long "--version" [1]
        Option_argument "<msg>"
      Option_description [17]
    Option_line [3]
      Option_short "-s" [1]
        Option_argument "<str>"
      Option_long "--separator" [1]
        Option_argument "<str>"
      Option_description [8]
    Options_node "----"
    Options_node "]"
    Option_line [3]
      Option_short "-O"
      Option_long "--options-first"
      Option_description [11]
    Option_line [3]
      Option_short "-H"
      Option_long "--no-help"
      Option_description [2]
    Option_line [2]
      Option_short "-A" [1]
        Option_argument "<name>"
      Option_description [5]
    Option_line [2]
      Option_short "-G" [1]
        Option_argument "<prefix>"
      Option_description [15]
    Option_line [2]
      Option_long "--no-mangle"
      Option_description [8]
    Option_line [2]
      Option_long "--no-declare"
      Option_description [5]
    Option_line [2]
      Option_long "--debug"
      Option_description [5]
  Free_section
