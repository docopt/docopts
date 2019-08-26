package main

import (
  "os"
  "fmt"

  "github.com/alecthomas/repr"

  "github.com/alecthomas/participle"
  "github.com/alecthomas/participle/lexer"
)

// # docopt top level parser 3 sections.
// Docopt         = [ Prologue ] Usage  [ Options ] .
//
// Prologue       = Line_of_text EOL | EOL | { Prologue } .
// EOL            = "\n" | "\r\n" | "\r" .
// Line_of_text   = non_EOL { non_EOL } EOL .
// non_EOL        = "\u0000"…"\u0009" | "\u000B"…"\uffff" .
//
// Usage          = "Usage:"  Usage_content .
// Usage_content  = [ EOL ] | Usage_line .
// Usage_line     = Program_Usage | Indent Program_Usage .
// Program_Usage  = Line_of_text .
// Indent         = "  " { " " } | "\t" { "\t" } .
//
// Options        = "Options:" EOL Options_lines .
// Options_lines  = Options_line { Options_line } .
// Options_line   = Indent Line_of_text .

type Docopt struct {
  Prologue *Free_text  `@@?`
  Usage *Usage `@@`
	Options *Options `@@?`
  Free_text *Free_text  `@@?`
}

type Free_text struct {
	Pos lexer.Position

  Description []string `( @Indent? @Line_of_text "\n" | @"\n" )*`
}

type Usage struct {
	Pos lexer.Position

  Usage_first   *string         `  "Usage:" ( @Line_of_text "\n" )?`
  Usage_lines   []*Usage_line   `           @@+`
}

type Usage_line struct {
	Pos lexer.Position

  Usage_content  *string   `  Indent @Line_of_text "\n"`
  Comment        *string   `| ( @Line_of_text "\n" | @"\n"+ )`
}

type Options struct {
  Pos lexer.Position

  Options_lines []string `"Options:" "\n" ( Indent @Line_of_text "\n" | "\n" )+`
}

var (
  // A custom lexer for docopt input
  doctop_Lexer = lexer.Must(lexer.Regexp(
    `(?P<NewLine>\n)` +
  // catch multiple blank
    `|(?P<Indent>\s{2,})` +
  // skip single blank
		`|(\s)` +
    `|(?P<Section>^(Usage|Options):)` +
    `|(?P<Line_of_text>[^\n]+)`,
  ))

  parser = participle.MustBuild(&Docopt{},
    participle.UseLookahead(2),
    participle.Lexer(doctop_Lexer),
    //participle.Elide("Comment", "Whitespace"),
    )

)

func main() {
  filename := os.Args[1]
	f, err := os.Open(filename)
  if err != nil {
    fmt.Printf("error: fail to open %s\n", filename)
    return
  } else {
    fmt.Printf("parsing: %s\n", filename)
  }

  ast := &Docopt{}
  if err = parser.Parse(f, ast) ; err == nil {
    fmt.Println("no error")
    repr.Println(ast)
  } else {
    fmt.Println("Parse error")
    fmt.Println(err)
    fmt.Println("======================= partial AST ==========================")
    repr.Println(ast)
  }
}
