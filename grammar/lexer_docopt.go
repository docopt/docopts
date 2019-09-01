package main

import (
	"github.com/docopt/docopts/grammar/lexer_state"
	"github.com/docopt/docopts/grammar/token_docopt"
  "fmt"
  "os"
  "github.com/alecthomas/participle/lexer"
)

func fail_if_error(err error) {
  if err != nil {
    fmt.Printf("error: %v\n", err)
    os.Exit(1)
  }
}

func main() {
  //stateDef, err := lexer_state.Parse_lexer_state("state_Prologue", state_Prologue)
  states, err := lexer_state.StateLexer(token_docopt.All_states, "state_Prologue")

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
	lex, err := states.Lex(f)
  fail_if_error(err)
  // display all tokens with their type
  symbols := lexer.SymbolsByRune(states)
  // extract all token
  for {
    t, err := lex.Next()
    fail_if_error(err)
    fmt.Printf("%s:Token{%s, %q}\n", states.Current_state.State_name, symbols[t.Type], t.Value)
    if t.Type == lexer.EOF {
      break
    }
  }
}
