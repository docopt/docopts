package main

import (
  "os"
  "fmt"

  "github.com/alecthomas/repr"

  "github.com/alecthomas/participle"
  "github.com/alecthomas/participle/lexer"
)

/*  grammar participle syntax ~ ebnf
Docopt =
  Prologue?
  Usage
	Options?
  Free_text?

Prologue       =  Free_text+
Free_text      =  INDENT? LINE_OF_TEXT "\n" | "\n"
INDENT         =  \s{2,}
LINE_OF_TEXT   =  [^\n]+
Usage          =  "Usage:" ( LINE_OF_TEXT "\n" )? Usage_line+
Usage_line     =  Usage_content | Comment
Usage_content  =  INDENT LINE_OF_TEXT "\n"
Comment        =  LINE_OF_TEXT "\n" | "\n"+
Options        =  "Options:" "\n" Options_line+
Options_line   =  INDENT LINE_OF_TEXT "\n" | "\n"
*/

// ================================ grammar ===============================
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
    repr.Println(ast)
    fmt.Println("Parse Success")
  } else {
    fmt.Println("Parse error")
    fmt.Println(err)
    fmt.Println("======================= partial AST ==========================")
    repr.Println(ast)
  }
}
