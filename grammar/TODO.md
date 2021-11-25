# docopt Grammar
// extract the Usage: section (case-insensitive) until the next empty line

## goal

* We start with our `lexer_state` before testing handcrafted optimized lexer/scanner
* parse actual grammar (to be determined by studying docopt-go)
* better error handling (See Also: [Error_reporting.md](Error_reporting.md))
* improve what needed for better support: `print ast` `analyze` `explain`

## Current work

(our branch grammar's PROGRESS.md)

`lexer_docopt.go`
* ~~dont change state inside the lexer, let the caller decide to change state~~ done in `lexer_docopt.go`
 * ~~remove participle dependancy (keep only our lexer)~~
 * ~~tokenize `usages/rally.docopt` with `Free_text` section after NEWLINE terminating Usage tokenizing~~
* ~~finish standalone lexer without error failure on valid usage input (See: compare)~~
* ~~document orgininal docopt-go lib behavior on some exmaple to compare with our grammar~~
* start building ast from our grammar: docopt_language.ebnf

## issues to read

* https://github.com/docopt/docopts/issues/35#issuecomment-516356005 explore docopt-ng

## Lexer

* ~~must handle syntax error (return ERROR token + state)~~
* ~~must progress one token at time~~
* ~~must stop scanning after a NEWLINE empty line with `Usage_line`~~ done within the parser calling our lexer
* ~~compare with `docopt-go/docopt.go` `ParseSection(name, source string)`~~ (See: compare/parseSection.go)

## example of parser / AST

* https://golangdocs.com/golang-parser-package
* https://ruslanspivak.com/lsbasi-part7/ - building an AST in python

## outputs

### `analyze`

```
Usage:
  docopts analyze [--full|--ast] [--only=<number>] <string_usage>
  docopts analyze [--full|--ast] [--only=<number>] -f <filename_usage>
```


Example:
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


