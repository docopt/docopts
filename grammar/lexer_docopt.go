package main

import (
	"fmt"
	"github.com/docopt/docopts/grammar/lexer"
	"github.com/docopt/docopts/grammar/lexer_state"
	"github.com/docopt/docopts/grammar/token_docopt"
	"github.com/gookit/color"
	"os"
	"strings"
)

func fail_if_error(err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

type Parser struct {
	s          *lexer_state.StateLexer
	last_token rune
	prog_name  string
}

func (p *Parser) change_to_state(t *lexer.Token, new_state string, color_fmt *color.Color) {
	if color_fmt == nil {
		fmt.Printf("%s change state: %s => %s\n", t.Regex_name, p.s.Current_state.State_name, new_state)
	} else {
		fmt.Printf("%s change state: %s => %s\n", t.Regex_name, p.s.Current_state.State_name, color_fmt.Render(new_state))
	}
	if err := p.s.ChangeState(new_state); err != nil {
		panic(err)
	}
	p.last_token = t.Type
}

func main() {
	// use like func
	red := color.FgRed
	green := color.FgGreen
	yellow := color.FgYellow
	cyan := color.FgCyan

	//stateDef, err := lexer_state.Parse_lexer_state("state_Prologue", state_Prologue)
	states, err := lexer_state.CreateStateLexer(token_docopt.All_states, "state_Prologue")

	if err == nil {
		fmt.Println(states)
	} else {
		fmt.Println(err)
	}

	var p Parser
	p.s = states
	p.last_token = 0

	filename := os.Args[1]
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("error: fail to open %s\n", filename)
		return
	} else {
		fmt.Printf("parsing: %s\n", green.Render(filename))
	}

	// initialize the Lexer with a string that will be read from file f
	states.State_auto_change = false
	lex, err := states.Lex(f)
	fail_if_error(err)
	// display all tokens with their type
	symbols_name := lexer.SymbolsByRune(states)
	all_symbols := states.Symbols()
	NEWLINE := all_symbols["NEWLINE"]
	SECTION := all_symbols["SECTION"]
	PROG_NAME := all_symbols["PROG_NAME"]

	// tokenize loop
	error_count := 0
	max_error := 10
	for {
		t, err := lex.Next()

		// token error handling
		if err != nil {
			fmt.Println(red.Render(err))
			error_count++
			states.Discard(err.(*lexer.Error).Pos, 1)
			if error_count < max_error {
				continue
			} else {
				fmt.Println(yellow.Render("too many error leaving"))
				break
			}
		}
		fmt.Printf("%s:Token{%s, %q}\n", states.Current_state.State_name, symbols_name[t.Type], t.Value)

		if t.Type == PROG_NAME && p.prog_name == "" {
			p.prog_name = t.Value
			fmt.Printf(green.Render("%s:Assign PROG_NAME: '%s' \n"), states.Current_state.State_name, p.prog_name)
			p.s.DynamicRuleUpdate("PROG_NAME", p.prog_name)
		}

		// two consecutive NEWLINE
		if p.last_token != 0 {
			if p.last_token == NEWLINE && t.Type == NEWLINE {
				fmt.Printf(cyan.Render("2 NEWLINE: %s\n"), states.Current_state.State_name)
				if states.Current_state.State_name == "state_Usage" {
					p.change_to_state(&t, "state_Free", &yellow)
					continue
				}
			}
		}

		if states.Current_state.State_name == "state_Free" && t.Type == SECTION && strings.EqualFold(t.Value, "Options:") {
			p.change_to_state(&t, "state_Options", &cyan)
			continue
		}

		// if we encounter a leave_token we change our lexer state
		if new_state, ok := states.Current_state.Leave_token[t.Regex_name]; ok {
			p.change_to_state(&t, new_state, nil)
			continue
		}

		if t.Type == lexer.EOF {
			break
		}
		p.last_token = t.Type
	}

	fmt.Printf("number of error: %d\n", error_count)
	os.Exit(error_count)
}
