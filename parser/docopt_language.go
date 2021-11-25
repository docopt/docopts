package docopt_language

import (
	"fmt"
	"github.com/docopt/docopts/grammar/lexer"
	"github.com/docopt/docopts/grammar/lexer_state"
	"github.com/docopt/docopts/grammar/token_docopt"
	"strings"
)

type DocoptAst struct {
	Type int
	Node *DocoptAst
}

type DocoptParser struct {
	s             *lexer_state.StateLexer
	prog_name     string
	current_token lexer.Token
	next_token    lexer.Token

	// map symbols <=> name
	symbols_name map[rune]string
	all_symbols  map[string]rune

	Error_count int
	max_error   int
	errors      []error
	ast         *DocoptAst
}

var (
	NEWLINE   rune
	SECTION   rune
	PROG_NAME rune
)

func ParserInit(source []byte) (*DocoptParser, error) {
	states, err := lexer_state.CreateStateLexer(token_docopt.All_states, "state_Prologue")
	if err != nil {
		return nil, fmt.Errorf("ParserInit: fail to init lexer_state: %v", err)
	}

	// initialize the Lexer with the source
	states.State_auto_change = false
	states.InitSource(source)

	p := &DocoptParser{
		s: states,

		prog_name:     "",
		current_token: lexer.EMPTY,
		next_token:    lexer.EMPTY,

		symbols_name: lexer.SymbolsByRune(states),
		all_symbols:  states.Symbols(),

		Error_count: 0,
		max_error:   10,
	}

	NEWLINE = p.all_symbols["NEWLINE"]
	SECTION = p.all_symbols["SECTION"]
	PROG_NAME = p.all_symbols["PROG_NAME"]

	p.NextToken()
	return p, nil
}

func (p *DocoptParser) NextToken() {
	t, err := p.s.Next()
	if err == nil {
		if p.next_token.Type != lexer.EMPTY.Type {
			p.current_token = p.next_token
		}
		p.next_token = t
	} else {
		p.AddError(err)

		if p.Error_count >= p.max_error {
			p.FatalError("too many error leaving")
			return
		}

		p.s.Discard(err.(*lexer.Error).Pos, 1)
		p.NextToken()
	}
}

func (p *DocoptParser) Eat(token_type rune) {
}

func (p *DocoptParser) FatalError(msg string) {
	for _, e := range p.errors {
		fmt.Println(e)
	}
	panic(msg)
}

func (p *DocoptParser) AddError(e error) {
	p.errors = append(p.errors, e)
	p.Error_count++
}

func (p *DocoptParser) Parse() *DocoptAst {
	for {
		p.NextToken()
		fmt.Printf("%s:%q\n", p.symbols_name[p.current_token.Type], p.current_token.Value)
		if p.current_token.Type == lexer.EOF {
			break
		}
	}
	//p.Consume_Prologue()
	//p.Consume_Usage()
	//p.Consume_Free_Section()
	//p.Consume_Options()
	//p.Consume_Free_Section()

	return p.ast
}

func (p *DocoptParser) CreateNode(node_name string) error {
	return nil
}

func (p *DocoptParser) Consume_Prologue() error {
	// State_Prologue = `
	// (?P<NEWLINE>\n)
	// |(?P<USAGE>[Uu][Ss][Aa][Gg][Ee]:) => state_First_Program_Usage
	// |(?P<WORD>\S+)
	// |(?P<BLANK>\s+)
	// `

	USAGE := p.all_symbols["USAGE"]
	p.CreateNode("Prologue")

	for {
		p.NextToken()

		if p.current_token.Type == USAGE {
			p.CreateNode("Usage")
			p.Eat(1)
			p.Change_lexer_state("state_Usage")
			return nil
		}

		p.Eat(1)

		if p.current_token.Type == lexer.EOF {
			break
		}
	}

	return fmt.Errorf("EOF encountered will parsing Prologue")
}

func (p *DocoptParser) Consume_Usage() error {
	if p.current_token.Type == PROG_NAME && p.prog_name == "" {
		p.prog_name = p.current_token.Value
		fmt.Printf("%s:Assign PROG_NAME: '%s' \n", p.s.Current_state.State_name, p.prog_name)
		p.s.DynamicRuleUpdate("PROG_NAME", p.prog_name)
	}
	// two consecutive NEWLINE
	if p.current_token.Type == NEWLINE && p.next_token.Type == NEWLINE {
		p.Eat(2)
	}
	return nil
}

func (p *DocoptParser) Consume_Free_Section() error {
	if p.s.Current_state.State_name == "state_Free" && p.current_token.Type == SECTION && strings.EqualFold(p.current_token.Value, "Options:") {
		p.Change_lexer_state("state_Options")
	}
	return nil
}

func (p *DocoptParser) Change_lexer_state(new_stage string) error {
	return nil
}
