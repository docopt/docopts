# docopt Grammar

## goal

* parse actual grammar
* better error handling (See Also: Error_reporting.md)

## issues to read

* https://github.com/docopt/docopts/issues/35#issuecomment-516356005 explore docopt-ng

## Current work

PROGRESS.md ?

* finish standalone lexer without error failure
 * use our current `lexer_state` or switch technology?


## Lexer

* must handle syntax error (return ERROR token)

## Participle - Fix Symbols() for parser

`func (sl *stateLexer) Symbols() (symbols map[string]rune) `

Is not compliant with the parser TOKEN.
Normalize rune and TOKEN between lexer and `lexer_state`

