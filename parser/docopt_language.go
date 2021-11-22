package docopt_language

import (
	"github.com/docopt/docopts/grammar/lexer"
	"github.com/docopt/docopts/grammar/lexer_state"
	"github.com/docopt/docopts/grammar/token_docopt"
)

type DocoptParser struct {
	s             *lexer_state.StateLexer
	prog_name     string
	current_token lexer.Token
	next_token    lexer.Token

	// map symbols <=> name
	symbols_name map[rune]string
	all_symbols  map[string]rune

	error_count int
	max_error   int
	errors      []error
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
		current_token: token_docopt.EMPTY,
		next_token:    token_docopt.EMPTY,

		symbols_name: lexer.SymbolsByRune(states),
		all_symbols:  states.Symbols(),

		error_count: 0,
		max_error:   10,
	}

	NEWLINE = p.all_symbols["NEWLINE"]
	SECTION = p.all_symbols["SECTION"]
	PROG_NAME = p.all_symbols["PROG_NAME"]

	return &p, nil
}

func (p *DocoptParser) NextToken() {
	t, err := p.s.Next()
	if err != nil {
		p.AddError(err)
		p.error_count++

		if p.error_count >= p.max_error {
			p.FatalError("too many error leaving")
			return
		}

		p.s.Discard(err.(*lexer.Error).Pos, 1)
	}

	p.current_token = t
}

func (p *DocoptParser) FatalError(msg string) {
	for _, e := range p.errors {
		fmt.Println(e)
	}
	panic(msg)
}

func (p *DocoptParser) AddError(e error) {
	p.errors = append(p.errors, e)
	p.error_count++
}
