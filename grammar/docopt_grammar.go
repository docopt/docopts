package main

import (
  "os"
  "fmt"

  "github.com/alecthomas/repr"
  "github.com/alecthomas/participle"
  "github.com/alecthomas/participle/lexer"

	"github.com/docopt/docopts/grammar/lexer_state"
	"github.com/docopt/docopts/grammar/token_docopt"
)

/*  grammar participle syntax ~ ebnf
Docopt =
  Prologue?
  Usage_section
	Options_section?
  Free_Section*

Prologue            =  Free_text+
Free_text           =  LONG_BLANK? LINE_OF_TEXT "\n" | "\n"
Usage_section       =    "Usage:" Usage_expr "\n" Usage_line*
                       | "Usage:" "\n" Usage_line+
Usage_line          =  ( LONG_BLANK Usage_expr | Comment ) "\n"
Comment             =  LINE_OF_TEXT | "\n"+
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
Options_line        =  LONG_BLANK Options_flag LONG_BLANK Option_description
Option_description  =  (LONG_BLANK LINE_OF_TEXT "\n")*
                       (LONG_BLANK LINE_OF_TEXT Defaulf_value "\n")?
Defaulf_value       =  "[" DEFAULT LINE_OF_TEXT "]"
Free_Section        = SECTION "\n" Free_text*
*/

// ================================ grammar ===============================
type Docopt struct {
  Prologue *Free_text  `@@?`
  Usage *Usage `@@`
	Options *Options `@@?`
  Free_Section *Free_Section  `@@*`
}

type Free_text struct {
	Pos lexer.Position

  Description []string `@( LINE_OF_TEXT "\n" | "\n" )*`
}

type Free_Section struct {
  Pos lexer.Position

  Section_name string   `@SECTION "\n"`
  Free_text *Free_text  `@@`
}

type Usage struct {
	Pos lexer.Position

  Usage_first       *Usage_def      `( "Usage:" @@ "\n"`
  Usage_next_lines  []*Usage_line   `           @@*`
  Usage_lines       []*Usage_line   `| "Usage:" "\n"  @@+ )`
}

type Usage_line struct {
	Pos lexer.Position

  Usage_def  *Usage_def   `   LONG_BLANK @@ "\n"`
  Comment    *string      `| @( LINE_OF_TEXT "\n" | "\n"+ )`
}

type Usage_def struct {
	Pos lexer.Position

  Prog_name    string     `@IDENT`
  Pattern      *Pattern   `   @@?`
}

type Pattern struct {
	Pos lexer.Position

  Seq           *Seq  `@@ `
  Exlusiv_Seq  []Seq  ` ("|" @@)*`
}

type Seq struct {
	Pos lexer.Position

  Atom           []Atom   ` (  @@`
  One_or_more    *string  `    @"..."? )*`
}

type Atom struct {
	Pos lexer.Position

	Group            *Pattern   `  "(" @@ ")"`
	Optional         *Pattern		`| "[" @@ "]"`
  Options_shortcut *string    `| "options"`
  Option           *Option    `| @@`
  Argument         *string    `| @ARGUMENT`
  Command          *string    `| @IDENT`
}

type Option struct {
	Pos lexer.Position

  Short      *string    `( @SHORT`
  Short_arg  *string    `  @ARGUMENT?`
  Long       *string    `| @LONG`
  Argument   *string    `  "="? @ARGUMENT? )`
}


type Options struct {
  Pos lexer.Position

  Options_lines []Options_line `"Options:" "\n" @@+`
}

type Options_line struct {
  Pos lexer.Position

	// eat all, not yet parsed
  Text []string `@("\n"|LONG_BLANK|DEFAULT|PUNCT|SHORT|LONG|ARGUMENT|LINE_OF_TEXT)`
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
  doctop_Lexer, err := lexer_state.StateLexer(token_docopt.All_states, "state_Prologue")
  if err != nil {
    fmt.Println(err)
    return
  }

  parser := participle.MustBuild(&Docopt{},
    participle.UseLookahead(1),
    participle.Lexer(doctop_Lexer),
    //participle.Elide("Comment", "Whitespace"),
    )

  ast := &Docopt{}
  if err = parser.Parse(f, ast) ; err == nil {
    //repr.Println(ast)
		repr.Println(ast, repr.Hide(&lexer.Position{}))
    fmt.Println("Parse Success")
  } else {
    fmt.Println("Parse error")
    fmt.Println("======================= partial AST ==========================")
    repr.Println(ast)
    fmt.Println("================ end of partial AST ==========================")
    fmt.Println(err)
  }
}
