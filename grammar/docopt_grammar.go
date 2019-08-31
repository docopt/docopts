package main

import (
  "os"
  "fmt"

  "github.com/alecthomas/repr"
  "github.com/alecthomas/participle"
  "github.com/alecthomas/participle/lexer"

	"github.com/docopt/docopts/grammar/lexer_state"
)

/*  grammar participle syntax ~ ebnf
Docopt =
  Prologue?
  Usage_section
	Options_section?
  Free_Section*

Prologue            =  Free_text+
Free_text           =  INDENT? LINE_OF_TEXT "\n" | "\n"
INDENT              =  \s{2,}
LINE_OF_TEXT        =  [^\n]+
Usage_section       =  "Usage:" ( Usage_content "\n" )? Usage_line+
Usage_line          =  Usage_content | Comment
Usage_content       =  INDENT Usage_expr "\n"
Comment             =  LINE_OF_TEXT "\n" | "\n"+
Usage_expr          =  Seq  ( "|" Seq )*
Seq                 =  ( Atom "..."? )*
Atom                =    "(" Expr ")"
                       | "[" Expr "]"
                       | "options"
                       | Long_def
                       | Shorts_option
                       | ARGUMENT
                       | Command
Shorts_option       =  SHORT | SHORT ARGUMENT
Long_def            =  LONG | LONG "="? ARGUMENT
Options_section     =  "Options:" "\n" Options_line+
Options_line        =  INDENT Options_flag INDENT Option_description
Option_description  =  (INDENT LINE_OF_TEXT "\n")*
                       (INDENT LINE_OF_TEXT Defaulf_value "\n")?
Defaulf_value       =  "[" DEFAULT LINE_OF_TEXT "]"
Free_Section        = SECTION "\n" Free_text*
*/

// ================================ lexer ===============================
var (
  state_Prologue = `
  (?P<NEWLINE>\n)
  |(?P<SECTION>^Usage:) => state_Usage_Line
  |(?P<LINE_OF_TEXT>[^\n]+)
  `

  state_Usage = `
  (?P<NEWLINE>\n)
  |(?P<USAGE>^Usage:)
  |(?P<SECTION>^[A-Z][A-Za-z _-]+:) => state_Options
  |(?P<INDENT>\s{2,}) => state_Usage_Line
  # skip single blank
  |(\s)
  # Match some kind of comment when not preceded by LongBlank
  |(?P<LINE_OF_TEXT>[^\n]+)
  `

  state_Usage_Line = `
  (?P<NEWLINE>\n) => state_Usage
  |(\s+)
  |(?P<SHORT>-[A-Za-z0-9?])
  |(?P<LONG>--[A-Za-z][A-Za-z0-9_-]+|^--$)
  |(?P<ARGUMENT><[A-Za-z][A-Za-z0-9_-]+>|[A-Z_][A-Z0-9_-]+)
  # Punctuation doesn't accept comma but elipsis ...
  |(?P<PUNCT>[\][=()|]|\.{3})
  |(?P<IDENT>[A-Za-z][A-Za-z0-9_-]+)
  `

  state_Options = `
  (?P<NEWLINE>\n)
  # Options: is matched by state_Usage
  |(?P<SECTION>^[A-Z][A-Za-z _-]+:) => state_Free
  |(?P<DEFAULT>^default:\s)
  |(?P<INDENT>\s{2,})
  # skip single blank
  |(\s)
  |(?P<SHORT>-[A-Za-z0-9?])
  |(?P<LONG>--[A-Za-z][A-Za-z0-9_-]+|^--$)
  |(?P<ARGUMENT><[A-Za-z][A-Za-z0-9_-]+>|[A-Z_][A-Z0-9_-]+)
  # Punctuation differe from state_Usage accepts comma
  |(?P<PUNCT>[\][=,()|])
  # LINE_OF_TEXT not matching []
  |(?P<LINE_OF_TEXT>[^\n[\]]+)
  `

  state_Free = `
  (?P<NEWLINE>\n)
  |(?P<SECTION>^[A-Z][A-Za-z _-]+:)
  |(?P<LINE_OF_TEXT>[^\n]+)
  `

  all_states = map[string]string{
    "state_Prologue" : state_Prologue,
    "state_Usage" : state_Usage,
    "state_Usage_Line" : state_Usage_Line,
    "state_Options" : state_Options,
    "state_Free" : state_Free,
  }
)

// ================================ grammar ===============================
type Docopt struct {
  Prologue *Free_text  `@@?`
  Usage *Usage `@@`
	Options *Options `@@?`
  Free_Section *Free_Section  `@@*`
}

type Free_text struct {
	Pos lexer.Position

  Description []string `( @LINE_OF_TEXT "\n" | @"\n" )*`
}

type Free_Section struct {
  Pos lexer.Position

  Section_name string   `@SECTION`
  Free_text *Free_text  `@@*`
}

type Usage struct {
	Pos lexer.Position

  Usage_first   *string         `  "Usage:" ( @LINE_OF_TEXT "\n" )?`
  Usage_lines   []*Usage_line   `           @@+`
}

type Usage_line struct {
	Pos lexer.Position

  Usage_content  *string   `  INDENT @LINE_OF_TEXT "\n"`
  Comment        *string   `| ( @LINE_OF_TEXT "\n" | @"\n"+ )`
}

type Options struct {
  Pos lexer.Position

  Options_lines []string `"Options:" "\n" ( INDENT @LINE_OF_TEXT "\n" | "\n" )+`
}

func main() {
  filename := os.Args[1]
	f, err := os.Open(filename)
  if err != nil {
    fmt.Printf("error: fail to open %s\n", filename)
    return
  } else {
    fmt.Printf("parsing: %s\n", filename)
  }

  // A custom lexer for docopt input
  doctop_Lexer, err := lexer_state.StateLexer(all_states, "state_Prologue")
  if err != nil {
    fmt.Println(err)
    return
  }

  parser := participle.MustBuild(&Docopt{},
    participle.UseLookahead(2),
    participle.Lexer(doctop_Lexer),
    //participle.Elide("Comment", "Whitespace"),
    )

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
