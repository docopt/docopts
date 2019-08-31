# docopt Grammar

## Fix Symbols() for parser

`func (sl *stateLexer) Symbols() (symbols map[string]rune) `

Is not compliant with the parser TOKEN.
Normalize rune and TOKEN between lexer and lexer_state

