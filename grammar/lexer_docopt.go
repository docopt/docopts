package main

import (
	"github.com/docopt/docopts/grammar/lexer_state"
  "fmt"
  "os"
  "github.com/alecthomas/participle/lexer"
)

var state_Prologue = `
(?P<NewLine>\n)
|(?P<Section>^Usage:) => state_Usage_Line
|(?P<Line_of_text>[^\n]+)
`
var state_Usage = `
(?P<NewLine>\n)
|(?P<Usage>^Usage:)
|(?P<Section>^Options:) => state_Options
|(?P<LongBlank>\s{2,}) => state_Usage_Line
# skip single blank
|(\s)
|(?P<Line_of_text>[^\n]+)
`

var state_Usage_Line = `
(?P<NewLine>\n) => state_Usage
|(\s+)
|(?P<ShortOption>-[A-Za-z0-9?])
|(?P<LongOption>--[A-Za-z][A-Za-z0-9_-]+|^--$)
|(?P<Argument><[A-Za-z][A-Za-z0-9_-]+>|[A-Z_][A-Z0-9_-]+)
# Punctuation doesn't accept comma but elipsis ...
|(?P<Punct>[\][=()|]|\.{3})
|(?P<Ident>[A-Za-z][A-Za-z0-9_-]+)
`
var state_Options = `
(?P<NewLine>\n)
# Options: is matched by state_Usage
|(?P<Options>^Options:)
# Default also match default: Keyword
|(?P<Default>^default:\s)
|(?P<Section>^[A-Z][A-Za-z _-]+:) => state_Free
|(?P<LongBlank>\s{2,})
#     skip single blank
|(\s)
|(?P<ShortOption>-[A-Za-z0-9?])
|(?P<LongOption>--[A-Za-z][A-Za-z0-9_-]+|^--$)
|(?P<Argument><[A-Za-z][A-Za-z0-9_-]+>|[A-Z_][A-Z0-9_-]+)
# Punctuation differe from state_Usage accepts comma
|(?P<Punct>[\][=,()|])
# Line_of_text not matching []
|(?P<Line_of_text>[^\n[\]]+)
`
var state_Free = `
(?P<NewLine>\n)
|(?P<Line_of_text>[^\n]+)
`

var all_states = map[string]string{
  "state_Prologue" : state_Prologue,
  "state_Usage" : state_Usage,
  "state_Usage_Line" : state_Usage_Line,
  "state_Options" : state_Options,
  "state_Free" : state_Free,
}

func main() {
  //stateDef, err := lexer_state.Parse_lexer_state("state_Prologue", state_Prologue)
  states, err := lexer_state.StateLexer(all_states, "state_Prologue")

  if err == nil {
    fmt.Println(states)
  } else {
    fmt.Println(err)
  }

  filename := os.Args[1]
	f, err := os.Open(filename)
  if err != nil {
    fmt.Printf("error: fail to open %s\n", filename)
    return
  } else {
    fmt.Printf("parsing: %s\n", filename)
  }

  // initialize the Lexer with a string
	lex, _ := states.Lex(f)
  // extract all token
  tokens, _ := lexer.ConsumeAll(lex)
  // display all tokens with their type
  for _, v := range tokens {
    fmt.Printf("Token{%s, %q}\n", states.Symbol(v.Type), v.Value)
  }
}
