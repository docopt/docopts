package docopt_language

import (
	"fmt"
	"github.com/docopt/docopts/grammar/lexer"
	"github.com/docopt/docopts/grammar/lexer_state"
	"github.com/docopt/docopts/grammar/token_docopt"
	"strings"
)

type DocoptParser struct {
	s             *lexer_state.StateLexer
	Prog_name     string
	current_token *lexer.Token
	next_token    *lexer.Token

	// map symbols <=> name
	symbols_name map[rune]string
	all_symbols  map[string]rune

	Error_count  int
	max_error    int
	errors       []error
	ast          *DocoptAst
	current_node *DocoptAst
}

type Consume_method func() error
type Consume_func struct {
	name    string
	consume Consume_method
}

var (
	NEWLINE   rune
	SECTION   rune
	PROG_NAME rune
	USAGE     rune
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

		Prog_name:     "",
		current_token: &lexer.EMPTY,
		next_token:    &lexer.EMPTY,

		symbols_name: lexer.SymbolsByRune(states),
		all_symbols:  states.Symbols(),

		Error_count: 0,
		max_error:   10,
	}

	NEWLINE = p.all_symbols["NEWLINE"]
	SECTION = p.all_symbols["SECTION"]
	PROG_NAME = p.all_symbols["PROG_NAME"]
	USAGE = p.all_symbols["USAGE"]

	// initialize next_token with tokenizer in mode state_Prologue
	p.NextToken()
	return p, nil
}

func (p *DocoptParser) NextToken() {
	t, err := p.s.Next()
	if err == nil {
		if p.next_token.Type != lexer.EMPTY.Type {
			p.current_token = p.next_token
		}
		p.next_token = &t
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

func (p *DocoptParser) Eat(f Consume_method) error {
	if err := f(); err != nil {
		return err
	}
	return nil
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

func consumer(name string, method Consume_method) Consume_func {
	return Consume_func{
		name:    name,
		consume: method,
	}
}

func (p *DocoptParser) Parse() *DocoptAst {
	p.CreateNode(Root, nil)

	// list parsing_step
	parse := []Consume_func{
		consumer("Consume_Prologue", p.Consume_Prologue),
		consumer("Consume_Usage", p.Consume_Usage),
		//p.Consume_Free_Section()
		//p.Consume_Options()
		//p.Consume_Free_Section()
	}

	for _, step := range parse {
		if err := step.consume(); err != nil {
			fmt.Printf("error: %s: %v\n", step.name, err)
			p.Error_count++
		}
	}

	return p.ast
}

// simple call to our tokenizer for testing debuging purpose
func (p *DocoptParser) Print_all_token() {
	for {
		p.NextToken()
		fmt.Printf("%s:%q\n", p.symbols_name[p.current_token.Type], p.current_token.Value)
		if p.current_token.Type == lexer.EOF {
			break
		}
	}
}

func (p *DocoptParser) CreateNode(node_type DocoptNodeType, token *lexer.Token) {
	if p.current_node == nil {
		p.current_node = &DocoptAst{
			Type:  node_type,
			Token: token,
		}
	} else {
		p.current_node = p.current_node.AddNode(node_type, token)
	}

	if p.ast == nil {
		p.ast = p.current_node
	}
}

func (p *DocoptParser) Consume_Prologue() error {
	p.CreateNode(Prologue, nil)

	for {
		p.NextToken()

		if p.current_token.Type == USAGE {
			// leaving Prologue
			usage_node := p.ast.AddNode(Usage_section, nil)
			usage_node.AddNode(Usage, p.current_token)
			p.current_node = usage_node
			return nil
		}

		p.current_node.AddNode(Prologue_node, p.current_token)

		if p.current_token.Type == lexer.EOF {
			// Prologue must leave on an Usage token
			break
		}
	}

	return fmt.Errorf("EOF encountered will parsing Prologue")
}

func (p *DocoptParser) Consume_Usage() error {
	// Usage   = USAGE , First_Program_Usage | { Program_Usage } ;
	// First_Program_Usage = PROG_NAME , [ Expr ] ;
	// (*
	//  PROG_NAME is catched at first definition and stay the same literal for all the parsing
	//  Program_Usage can be break multi-line: Indent + PROG_NAME will start a new Program_Usage
	//
	//  Usage: ./my_program.sh [-h] [--lovely-option] FILENAME
	//         ./my_program.sh another LINE OF --usage
	//         my_program      will continue [the] [--above] <usage-definition>
	//
	//  PROG_NAME  on first usage parsing it becomes: "./my_program.sh"
	// *)
	// PROG_NAME = ? any non space characters ? ;
	// Program_Usage = Indent , PROG_NAME  [ Expr ] ;

	if err := p.Consume_First_Program_Usage(); err != nil {
		return err
	}

	if err := p.Consume_Usage_line(); err != nil {
		return err
	}

	return nil
}

func (p *DocoptParser) Consume_Usage_line() error {
	p.Change_lexer_state("state_Usage_Line")
	var n DocoptNodeType
	SHORT := p.all_symbols["SHORT"]
	LONG := p.all_symbols["LONG"]
	ARGUMENT := p.all_symbols["ARGUMENT"]
	PUNCT := p.all_symbols["PUNCT"]
	IDENT := p.all_symbols["IDENT"]
	// current_node: Usage_line
	for {
		p.NextToken()

		// matching a PROG_NAME will start a new Usage_line
		if p.current_token.Type == PROG_NAME {
			if p.Prog_name != p.current_token.Value {
				return fmt.Errorf(
					"Consume_Usage_line: PROG_NAME encountered with a distinct value:%s, invalid Token: (%s)'%v'",
					p.Prog_name,
					p.s.Current_state.State_name,
					p.current_token)
			}

			usage_line := p.current_node.parent.AddNode(Usage_line, nil)
			usage_line.AddNode(Prog_name, p.current_token)
			p.current_node = usage_line
			continue
		}

		if p.current_token.Type == USAGE {
			return fmt.Errorf("Consume_Usage_line: USAGE invalid Token: %v", p.current_token)
		}

		if p.current_token.Type == NEWLINE && p.next_token.Type == NEWLINE {
			// two consecutive NEWLINE
			// consume next NEWLINE
			p.NextToken()
			// leave Usage parsing
			return nil
		}

		if p.current_token.Type == lexer.EOF {
			return nil
		}

		switch p.current_token.Type {
		case SHORT:
			n = Usage_short_option
		case LONG:
			n = Usage_long_option
		case ARGUMENT:
			n = Usage_argument
		case PUNCT:
			n = Usage_punct
		case IDENT:
			n = Usage_indent
		default:
			n = Unmatched_node
		}
		p.current_node.AddNode(n, p.current_token)
	}

}

func (p *DocoptParser) Consume_First_Program_Usage() error {
	// assert p.Prog_name == ""
	p.Change_lexer_state("state_First_Program_Usage")
	BLANK := p.all_symbols["BLANK"]
	// p.current_node has been set previously and must be Usage_section
	for {
		p.NextToken()

		if p.current_token.Type == PROG_NAME {
			p.Prog_name = p.current_token.Value
			p.s.DynamicRuleUpdate("PROG_NAME", p.Prog_name)

			usage_line := p.current_node.AddNode(Usage_line, nil)
			usage_line.AddNode(Prog_name, p.current_token)
			p.current_node = usage_line
			return nil
		}

		if p.current_token.Type == BLANK {
			continue
		}

		if p.current_token.Type == NEWLINE {
			if p.next_token.Type == NEWLINE {
				// two consecutive NEWLINE
				if p.Prog_name == "" {
					return fmt.Errorf("Consume_First_Program_Usage: PROG_NAME not defined while leaving on 2 consecutive NEWLINE: %v", p.current_token)
				}
				// consume next NEWLINE
				p.NextToken()
				// leave
				return nil
			}

			continue
		}

		return fmt.Errorf("Consume_First_Program_Usage: expecting PROG_NAME, got: %s", p.symbols_name[p.current_token.Type])
	}
}

func (p *DocoptParser) Consume_Free_Section() error {
	if p.s.Current_state.State_name == "state_Free" && p.current_token.Type == SECTION && strings.EqualFold(p.current_token.Value, "Options:") {
		p.Change_lexer_state("state_Options")
	}
	return nil
}

func (p *DocoptParser) Change_lexer_state(new_state string) error {
	return p.s.ChangeState(new_state)
}
