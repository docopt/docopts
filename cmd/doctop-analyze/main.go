package main

import (
	"fmt"
	"github.com/docopt/docopts/parse/docopt_language"
	"os"
)

func main() {
	filename := os.Args[1]
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("error: fail to open '%s': $s\n", filename)
		os.Exit(1)
	} else {
		fmt.Printf("parsing: %s\n", filename)
	}

	p := docopt_language.ParserInit(data)
	for {
		p.NextToken()

		if p.current_token.Type == PROG_NAME && p.prog_name == "" {
			p.prog_name = p.current_token.Value
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
