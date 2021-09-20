# docopt Grammar
// extract the Usage: section (case-insensitive) until the next empty line

## goal

* We start with our `lexer_state` before testing handcrafted optimized lexer/scanner
* parse actual grammar
* better error handling (See Also: Error_reporting.md)

## Current work

(our branch PROGRESS.md)

* dont change state inside the lexer, let the caller decide to change state
 *  remove participle dependancy
* finish standalone lexer without error failure

## issues to read

* https://github.com/docopt/docopts/issues/35#issuecomment-516356005 explore docopt-ng

## Lexer

* must handle syntax error (return ERROR token + state)
* must progress one token at time
* must stop scanning after a NEWLINE empty line avec `Usage_line`
* ~~compare with `docopt-go/docopt.go` +205 `parseSection(name, source string)`~~

## Participle - Fix Symbols() for parser

`func (sl *stateLexer) Symbols() (symbols map[string]rune) `

Is not compliant with the parser TOKEN.
Normalize rune and TOKEN between lexer and `lexer_state`

