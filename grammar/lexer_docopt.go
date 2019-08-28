package main

import (
	"github.com/docopt/docopts/grammar/lexer_state"
  "fmt"
  // "os"
)

var state_Prologue = `
(?P<NewLine>\n)
|(?P<Section>^Usage:) => state_Usage
|(?P<Line_of_text>[^\n]+)
`
var state_Usage = `
(?P<NewLine>\n)
|(?P<Usage>^Usage:)
|(?P<Section>^Options:) => state_Options
|(?P<LongBlank>\s{2,})
#     skip single blank
|(\s)
|(?P<ShortOption>-[A-Za-z0-9?])
|(?P<LongOption>--[A-Za-z][A-Za-z0-9_-]+|^--$)
|(?P<Argument><[A-Za-z][A-Za-z0-9_-]+>|[A-Z_][A-Z0-9_-]+)
# Punctuation doesn't accept comma but elipsis ...
|(?P<Punct>[\][=()|]|\.{3})
|(?P<Ident>[A-Za-z][A-Za-z0-9_-]+)
|(?P<Line_of_text>[^\n]+)
`
var state_Options = `
(?P<NewLine>\n)
|(?P<Options>^Options:)
|(?P<Section>^[A-Z][A-Za-z _-]+:) => state_Free
|(?P<LongBlank>\s{2,})
#     skip single blank
|(\s)
|(?P<ShortOption>-[A-Za-z0-9?])
|(?P<LongOption>--[A-Za-z][A-Za-z0-9_-]+|^--$)
|(?P<Argument><[A-Za-z][A-Za-z0-9_-]+>|[A-Z_][A-Z0-9_-]+)
# Punctuation differe from state_Usage accepts coma
|(?P<Punct>[\][=,()|])
# Default also match default: Keyword
|(?P<Default>default:\s*[^\]]+)
|(?P<Line_of_text>[^\n]+)
`
var state_Free = `
(?P<NewLine>\n)
|(?P<Line_of_text>[^\n]+)
`

var all_states = map[string]string{
  "state_Prologue" : state_Prologue,
  "state_Usage" : state_Usage,
  "state_Options" : state_Options,
  "state_Free" : state_Free,
}

func main() {
  stateDef, err := lexer_state.Parse_lexer_state("state_Prologue", state_Prologue)

  if err == nil {
    fmt.Println(stateDef)
  } else {
    fmt.Println(err)
  }

  //sl := lexer_state.StateLexer(all_states, "state_Prologue")
  //fmt.Println(sl)

  //doctop_Lexer := lexer.Must(lexer.Regexp(lexer_Usage.pattern))
  //// extract symbols from the Lexer
  //sym := lexer.SymbolsByRune(doctop_Lexer)
  //// initialize the Lexer with a string
	//lex, _ := doctop_Lexer.Lex(os.Stdin)
  //// extract all token
  //tok, _ := lexer.ConsumeAll(lex)
  //// display all tokens with their type
  //for _, v := range tok {
  //  fmt.Printf("Token{%s, %q}\n", sym[v.Type], v.Value)
  //}
}
