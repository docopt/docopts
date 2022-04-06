# grammar development and exploration

some main code:

## native golang parser

Work-in-Progress

following an EBNF (like) grammar and keeping strong error handling in mind, the code is here:

* proposed grammar: `docopt_language.ebnf` (based on actual docopt parser behavior study)
* code lib: `../parser/docopt_language.go`
* code runner: `../cmd/docopt-analyze/main.go`

```
cd ../cmd/docopt-analyze
for d in ../../grammar/usages/*.docopt; do   echo "=============== $d";   go run ./main.go $d; done | batcat
```

## `lexer_docopt.go`

is it DEPRECATED?

Our `lexer_docopt` as a stand alone lexer with state

```
go run lexer_docopt.go usages/rally.docopt
```

## `docopt_grammar.go`

DEPRECATED

Use participle and `lexer_state.StateLexer`
Some test grammar, abandonned too much missing features (for debugging and explaning grammar at runtime)
Don't compile anymore with this version of `lexer_state/lexer_state.go`, the dependancy with participle/lexer is
broken.

```
go run docopt_grammar.go usages/rally.docopt
```

## `docopt_top_grammar.go`

DEPRECATED

use participle and try to parse only top level grammar as the grammar changes for sub-element:

```
Docopt =
  Prologue?
  Usage
	Options?
  Free_text?
```
