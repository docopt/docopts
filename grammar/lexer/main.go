package main

import (
  "github.com/alecthomas/participle/lexer"
  "fmt"
  "os"
)

var lexer_regxep = `
`

// // A custom lexer for docopt input
// var doctop_Lexer = lexer.Must(lexer.Regexp(
//     `(?P<NewLine>\n)` +
//   // catch multiple blank
//     `|(?P<LongBlank>\s{2,})` +
//   // skip single blank
//    `|(\s)` +
//     `|(?P<Section>^(Usage|Options):)` +
//   // Default also match default: Keyword
//     `|(?P<Default>default:\s*[^\]]+)` +
//   // single letter incluging -?
//   // also describe long option
//   // also --
//     `|(?P<Option>-[A-Za-z0-9?]|--[A-Za-z][A-Za-z0-9_-]+|^--$)` +
//    `|(?P<Argument><[a-z][a-z-]+>|[A-Z_]{2,})` +
//    `|(?P<Punct>[\][=,()|]|\.{3})` +
//    `|(?P<String>[^ \t\n\][=,()|]+)`,
// ))

// A custom lexer for docopt input
var doctop_Lexer = lexer.Must(lexer.Regexp(
    `(?P<NewLine>\n)` +
  // catch multiple blank
    `|(?P<Indent>\s{2,})` +
  // skip single blank
    `|(\s)` +
    `|(?P<Section>^(Usage|Options):)` +
    `|(?P<Line_of_text>[^\n]+)`,
))


func main() {
  filename := os.Args[1]
	f, err := os.Open(filename)
  if err != nil {
    fmt.Printf("error: fail to open %s\n", filename)
    return
  } else {
    fmt.Printf("parsing: %s\n", filename)
  }
  // extract symbols from the Lexer
  sym := lexer.SymbolsByRune(doctop_Lexer)
  // initialize the Lexer with a string
  lex, _ := doctop_Lexer.Lex(f)
  // extract all token
  tok, _ := lexer.ConsumeAll(lex)
  // display all tokens with their type
  for _, v := range tok {
    fmt.Printf("Token{%s, %q}\n", sym[v.Type], v.Value)
  }
}
