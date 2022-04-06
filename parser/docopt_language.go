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

	lexer_state_changed bool
	run                 bool
}

type Consume_method func() error
type Consume_func struct {
	name    string
	consume Consume_method
}

var (
	NEWLINE    rune
	SECTION    rune
	PROG_NAME  rune
	USAGE      rune
	SHORT      rune
	LONG       rune
	ARGUMENT   rune
	PUNCT      rune
	IDENT      rune
	LONG_BLANK rune
	DEFAULT    rune
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
		current_token: nil,
		next_token:    nil,

		symbols_name: lexer.SymbolsByRune(states),
		all_symbols:  states.Symbols(),

		Error_count:         0,
		max_error:           10,
		lexer_state_changed: false,
		run:                 true,
	}

	// TODO: initialize token in token_docopt
	NEWLINE = p.all_symbols["NEWLINE"]
	SECTION = p.all_symbols["SECTION"]
	PROG_NAME = p.all_symbols["PROG_NAME"]
	USAGE = p.all_symbols["USAGE"]
	SHORT = p.all_symbols["SHORT"]
	LONG = p.all_symbols["LONG"]
	ARGUMENT = p.all_symbols["ARGUMENT"]
	PUNCT = p.all_symbols["PUNCT"]
	IDENT = p.all_symbols["IDENT"]
	LONG_BLANK = p.all_symbols["LONG_BLANK"]
	DEFAULT = p.all_symbols["DEFAULT"]

	return p, nil
}

func (p *DocoptParser) NextToken() {
	if p.next_token != nil && p.lexer_state_changed {
		p.s.Reject(p.next_token)
		p.lexer_state_changed = false
		p.next_token = nil
	}

	if p.next_token != nil {
		p.current_token = p.next_token
	} else {
		p.try_get_NextToken(&p.current_token)
	}

	p.try_get_NextToken(&p.next_token)

	if p.Error_count >= p.max_error {
		p.FatalError("too many error leaving")
	}
}

func (p *DocoptParser) try_get_NextToken(token_to_store **lexer.Token) error {
	if p.Error_count >= p.max_error {
		return fmt.Errorf("try_get_NextToken: already too many errors")
	}

	t, err := p.s.Next()
	if err == nil {
		*token_to_store = &t
	} else {
		// error collector
		p.AddError(err)

		if p.Error_count >= p.max_error {
			return fmt.Errorf("try_get_NextToken: too many errors")
		}

		p.s.Discard(1)
		return p.try_get_NextToken(token_to_store)
	}

	return nil
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
	p.run = false
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
		consumer("Consume_Free_Section", p.Consume_Free_section),
		consumer("Consume_Options", p.Consume_Options),
		consumer("Consume_Free_Section", p.Consume_Free_section),
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
func (p *DocoptParser) Print_all_token() error {
	for p.run {
		p.NextToken()
		fmt.Printf("%s:%q\n", p.symbols_name[p.current_token.Type], p.current_token.Value)
		if p.current_token.Type == lexer.EOF {
			break
		}
	}

	if p.run {
		return nil
	} else {
		return fmt.Errorf("Print_all_token: parser stoped")
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
	// we start parsing we are at Root node
	p.CreateNode(Prologue, nil)

	for p.run {
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
			return fmt.Errorf("EOF encountered will parsing Prologue, without 'Usage:' found")
		}
	}

	return fmt.Errorf("%s: parser stoped", p.current_node.Type)
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
	// current_node: Usage_line right after PROG_NAME has been matched

	for p.run {
		p.NextToken()
		if p.current_token.Type == lexer.EOF {
			// assert leaving condition are met
			return nil
		}

		// matching a PROG_NAME will start a new Usage_line
		if p.current_token.Type == PROG_NAME {
			if p.Prog_name != p.current_token.Value {
				return fmt.Errorf(
					"Consume_Usage_line:(%s) PROG_NAME encountered with a distinct value:%s, invalid Token: '%v' extracted with: %s",
					p.s.Current_state.State_name,
					p.Prog_name,
					p.current_token,
					p.current_token.State_name)
			}

			usage_line := p.current_node.Parent.AddNode(Usage_line, nil)
			usage_line.AddNode(Prog_name, p.current_token)
			p.current_node = usage_line
			continue
		}

		if p.current_token.Type == USAGE {
			return fmt.Errorf("Consume_Usage_line: USAGE invalid Token: %v", p.current_token)
		}

		if p.current_token.Type == NEWLINE {
			if p.next_token.Type == NEWLINE {
				// two consecutive NEWLINE
				// consume next NEWLINE
				p.NextToken()
				// leave Usage parsing
				return nil
			}

			// single NEWLINE skipping
			continue
		}

		// assign a token
		switch p.current_token.Type {
		case SHORT:
			n = Usage_short_option
		case LONG:
			n = Usage_long_option
		case ARGUMENT:
			n = Usage_argument
		case PUNCT:
			switch p.current_token.Value {
			case "[":
				n = Usage_optional_group
			case "(":
				n = Usage_required_group
			case "...":
				if err := p.Consume_ellipsis(); err != nil {
					return err
				}
				continue
			case "=":
				if err := p.Consume_assign(p.next_token); err != nil {
					return err
				}
				// consume ARGUMENT assigned
				p.NextToken()
				continue
			//case "|":
			//p.Ensure_node(p.current_node, Group_alternative)

			default:
				// unmatched PUNCT
				n = Usage_unmatched_punct
			}

			if n != Usage_unmatched_punct {
				// try to match a group required or optional
				if err := p.Consume_group(n); err != nil {
					return err
				}

				// assert
				if p.current_node.Type != Usage_line {
					p.FatalError(fmt.Sprintf("p.Consume_group(%s) did not restore current_node: %s",
						n,
						p.current_node.Type))
				}
				continue
			}

			// else: unmatched PUNCT will added to the AST
		case IDENT:
			n = Usage_command
		default:
			n = Unmatched_node
		}
		p.current_node.AddNode(n, p.current_token)
	}

	return fmt.Errorf("%s: parser stoped", p.current_node.Type)
}

func (p *DocoptParser) Consume_ellipsis() error {
	nb := len(p.current_node.Children)
	if nb > 0 {
		p.current_node.Children[nb-1].Repeat = true
	} else {
		return fmt.Errorf("%s: elipsis not expected on such node without Children, invalid Token: %v",
			p.current_node.Type, p.current_token)
	}
	return nil
}

func (p *DocoptParser) Consume_group(group_type DocoptNodeType) error {
	group := p.current_node.AddNode(group_type, nil)
	saved_current_node := p.current_node
	p.current_node = group
	var err error = nil
	var n DocoptNodeType
forLoop:
	for p.run {
		p.NextToken()
		switch p.current_token.Type {
		case lexer.EOF, PROG_NAME:
			err = fmt.Errorf("%s: %s unexpected, missing closing bracket ']'",
				p.current_node.Type,
				p.symbols_name[p.current_token.Type])
			break forLoop
		case USAGE:
			err = fmt.Errorf("%s: USAGE invalid Token: %v", p.current_node.Type, p.current_token)
			break forLoop
		case NEWLINE:
			if p.next_token.Type == NEWLINE {
				// two consecutive NEWLINE
				err = fmt.Errorf("%s: 2 consecutive NEWLINE invalid Token: %v", p.current_node.Type, p.current_token)
				break forLoop
			}
			continue
		case SHORT:
			n = Usage_short_option
		case LONG:
			n = Usage_long_option
		case ARGUMENT:
			n = Usage_argument
		case PUNCT:
			switch p.current_token.Value {
			case "[":
				if err = p.Consume_group(Usage_optional_group); err != nil {
					break forLoop
				}
				continue
			case "(":
				if err = p.Consume_group(Usage_required_group); err != nil {
					break forLoop
				}
				continue
			case "|":
				if p.current_node.Type != Group_alternative {
					// move actual Children to a new Group_alternative node
					alternative := &DocoptAst{
						Type:     Group_alternative,
						Token:    nil,
						Parent:   p.current_node,
						Children: p.current_node.Children,
					}
					for _, c := range alternative.Children {
						c.Parent = alternative
					}
					p.current_node.Children = []*DocoptAst{alternative}
					p.current_node = alternative
				}
				continue
			case "]":
				if p.current_node.Type == Group_alternative {
					p.current_node = p.current_node.Parent
				}

				if !(p.current_node.Type == Usage_optional_group || p.current_node.Parent.Type == Usage_optional_group) {
					err = fmt.Errorf("%s: closing bracket unexpected, invalid Token: %v", p.current_node.Type, p.current_token)
				}
				break forLoop
			case ")":
				if p.current_node.Type == Group_alternative {
					p.current_node = p.current_node.Parent
				}

				if !(p.current_node.Type == Usage_required_group || p.current_node.Parent.Type == Usage_required_group) {
					err = fmt.Errorf("%s: closing parenthese unexpected, invalid Token: %v", p.current_node.Type, p.current_token)
				}
				break forLoop
			case "=":
				if err = p.Consume_assign(p.next_token); err != nil {
					break forLoop
				}
				// consume ARGUMENT assigned
				p.NextToken()
				continue
			case "...":
				if err = p.Consume_ellipsis(); err != nil {
					break forLoop
				}
				continue
			}

			// unmatched PUNCT
			n = Usage_unmatched_punct

		case IDENT:
			n = Usage_command
		default:
			n = Unmatched_node
		}

		p.current_node.AddNode(n, p.current_token)
	}

	if p.run {
		p.current_node = saved_current_node
		return err
	} else {
		return fmt.Errorf("%s: parser stoped: %s", p.current_node.Type, err)
	}

}

func (p *DocoptParser) Consume_First_Program_Usage() error {
	// assert p.Prog_name == ""
	p.Change_lexer_state("state_First_Program_Usage")
	BLANK := p.all_symbols["BLANK"]
	// p.current_node has been set previously and must be Usage_section
	for p.run {
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

	return fmt.Errorf("%s: parser stoped", p.current_node.Type)
}

// This are section like part of the definition not yet used
// This basically allow more comment, but node are added to the ast
func (p *DocoptParser) Consume_Free_section() error {
	if p.s.Current_state.State_name != "state_Free" {
		// entering Free_section after: Usage_section or Options_section
		p.Change_lexer_state("state_Free")
		p.current_node = p.ast.AddNode(Free_section, nil)

	} else {
		// nested free section: we matched another SECTION token inside a Free_section
		p.current_node = p.ast.AddNode(Free_section, nil)
	}

	if p.current_token.Type == SECTION {
		p.current_node.AddNode(Section_name, p.current_token)
	}

	for p.run {
		p.NextToken()

		if p.current_token.Type == lexer.EOF {
			return nil
		}

		// leaving condition
		if p.current_token.Type == SECTION {
			if strings.EqualFold(p.current_token.Value, "options:") {
				return nil
			}

			if strings.EqualFold(p.current_token.Value, "usage:") {
				return fmt.Errorf("%s: Usage: token found outside Usage_section: %v",
					p.current_node.Type,
					p.current_token)
			}

			// test if the current section has already some content (was empty unamed section)
			nb := len(p.current_node.Children)
			if nb == 0 {
				p.current_node.AddNode(Section_name, p.current_token)
				continue
			}

			// nested Free_section
			// Free_section leaving condition are: EOF or SECTION == Options: or error
			return p.Consume_Free_section()
		}

		p.current_node.AddNode(Section_node, p.current_token)
	}

	return fmt.Errorf("%s: parser stoped", p.current_node.Type)
}

func (p *DocoptParser) Change_lexer_state(new_state string) error {
	p.lexer_state_changed = true
	return p.s.ChangeState(new_state)
}

func (p *DocoptParser) Consume_Options() error {
	section_node := p.ast.AddNode(Options_section, nil)
	// only start parsing Options if we start on a token SECTION == Options:
	if p.current_token.Type != SECTION || !strings.EqualFold(p.current_token.Value, "options:") {
		return nil
	}

	p.Change_lexer_state("state_Options")
	section_node.AddNode(Section_name, p.current_token)
	p.current_node = section_node

	var n DocoptNodeType
	var err error

	for p.run {
		p.NextToken()

		if p.current_token.Type == lexer.EOF {
			return nil
		}

		if p.current_token.Type == NEWLINE {
			if p.next_token.Type == NEWLINE {
				// two consecutive NEWLINE
				// consume next NEWLINE
				p.NextToken()
				// leave Usage parsing of Consume_Options
				return nil
			}

			// else: single NEWLINE
		}

		n = Options_node

		switch p.current_token.Type {
		case SECTION:
			return nil
		case LONG_BLANK:
			if p.next_token.Type == SHORT || p.next_token.Type == LONG {
				if err = p.Consume_option_line(); err != nil {
					return err
				}
			}
			continue
		case NEWLINE:
			continue
		}

		// unmatch Options_node
		p.current_node.AddNode(n, p.current_token)
	}

	return fmt.Errorf("%s: parser stoped: %s", p.current_node.Type, err)
}

// consume the next token which must be ARGUMENT as the argument of the last
// node added.
func (p *DocoptParser) Consume_assign(argument *lexer.Token) error {
	if argument.Type != ARGUMENT {
		return fmt.Errorf("%s: Consume_assign must be followed by ARGUMENT, invalid token: %v",
			p.current_node.Type, argument)
	}

	nb_children := len(p.current_node.Children)
	if nb_children == 0 {
		// Consume_assign must called after having assigned a option LONG in Usage_line
		// or any option in Options_line called with oe without equal sign
		return fmt.Errorf("Consume_assign: current_node must have an option child, invalid Token: %v", p.current_token)
	}

	prev_child := p.current_node.Children[nb_children-1]
	var node_type DocoptNodeType
	switch prev_child.Type {
	// only those kind of node can have assignment with ARGUMENT
	case Usage_long_option:
		node_type = Usage_argument
	case Option_long, Option_short:
		node_type = Option_argument
	default:
		return fmt.Errorf("Consume_assign: node %s cannot have assignment '=', invalid Token: %v",
			prev_child.Type, p.current_token)
	}

	prev_child.AddNode(node_type, argument)
	return nil
}

func (p *DocoptParser) Consume_option_alternative() error {
	// create the parent node on first call
	if p.current_node.Type != Option_alternative_group {
		nb := len(p.current_node.Children)
		if nb == 0 {
			return fmt.Errorf("%s: comma unexpected without alternative option name, invalid Token: %v", p.current_node.Type, p.current_token)
		}

		p.current_node = p.current_node.Replace_children_with_group(Option_alternative_group)
	}

	// eat next option alternative
	for p.run {
		p.NextToken()
		switch p.current_token.Type {
		case lexer.EOF, LONG_BLANK, NEWLINE:
			if len(p.current_node.Children) <= 1 {
				return fmt.Errorf("%s: %s unexpected without matchin alternative option name, invalid Token: %v",
					p.current_node.Type, p.symbols_name[p.current_token.Type], p.current_token)
			}
			// leaving condition OK
			return nil
		case SHORT:
			p.current_node.AddNode(Option_short, p.current_token)
		case LONG:
			p.current_node.AddNode(Option_long, p.current_token)
		case PUNCT:
			switch p.current_token.Value {
			case ",":
				continue
			case "=":
				if err := p.Consume_assign(p.next_token); err != nil {
					return err
				}
				// consume ARGUMENT assigned
				p.NextToken()
				continue
			default:
				return fmt.Errorf("%s: unexpected PUNC, invalid Token: %v", p.current_node.Type, p.current_token)
			}
		}
	}

	return fmt.Errorf("%s: parser stoped", p.current_node.Type)
}

func (p *DocoptParser) Consume_option_line() error {
	// we did look a head on token: p.current_token is an option LONG or SHORT
	// the option argument will be consumed during the first loop
	saved_node := p.current_node
	option_line := p.current_node.AddNode(Option_line, nil)
	p.current_node = option_line
	var err error = nil
forLoopOptionLine:
	for p.run {
		p.NextToken()

		switch p.current_token.Type {
		case lexer.EOF, NEWLINE:
			// leaving condition option without description
			if len(p.current_node.Children) == 0 {
				err = fmt.Errorf("%s: %s unexpected empty option, invalid Token: %s",
					p.current_node.Type, p.symbols_name[p.current_token.Type], p.current_token)
			}
			break forLoopOptionLine
		case LONG_BLANK:
			// LONG_BLANK in Consume_option_line occurs after options are comsumed
			// leaving condition of Consume_Usage_line
			err = p.Consume_option_description()
			break forLoopOptionLine
		case SHORT:
			p.current_node.AddNode(Option_short, p.current_token)
		case LONG:
			p.current_node.AddNode(Option_long, p.current_token)
		case ARGUMENT:
			if err = p.Consume_assign(p.current_token); err != nil {
				break forLoopOptionLine
			}
		case PUNCT:
			switch p.current_token.Value {
			case ",":
				continue
			case "=":
				if err := p.Consume_assign(p.next_token); err != nil {
					break forLoopOptionLine
				}
				// consume ARGUMENT assigned
				p.NextToken()
			default:
				err = fmt.Errorf("%s: unexpected PUNC, invalid Token: %v", p.current_node.Type, p.current_token)
				break forLoopOptionLine
			}
		default:
			err = fmt.Errorf("%s: Consume_option_line invalid Token: %v", p.current_node.Type, p.current_token)
			break forLoopOptionLine
		}
	} // end forLoopOptionLine

	if p.run {
		p.current_node = saved_node
		return err
	} else {
		return fmt.Errorf("%s: parser stoped", p.current_node.Type)
	}
}

// option description occurs after option has been parsed and can continue on multiple line
// indented by LONG_BLANK. The description is terminated when a new option SHORT or LONG
// is matched at the beginning of the line: NEWLINE LONG_BLANK (SHORT | LONG)
//
//                            Start consume description here
//                                |
// Options:                       v
//   -h <msg>, --help=<msg>        The help message in docopt format.
//                                 Without argument outputs this help.
//                                 If - is given, read the help message from
//                                 standard input.
//                                 If no argument is given, print docopts's own
//                                 help message and quit.
// => LONG_BLANK + option ==> leaving
func (p *DocoptParser) Consume_option_description() error {
	description := p.current_node.AddNode(Option_description, nil)
	current_line := 0

	for p.run {
		p.NextToken()

		switch p.current_token.Type {
		case NEWLINE:
			current_line++
			if p.next_token.Type == NEWLINE {
				// two consecutive NEWLINE
				description.AddNode(Description_node, p.current_token)

				// consume next NEWLINE
				p.NextToken()
				// leave Consume_option_description
				return nil
			}
			// else: single NEWLINE will be consumed as part of the description

		// all the following are leaving condition, other token will be collected as part of the description
		case lexer.EOF:
			return nil
		case LONG_BLANK:
			if current_line > 0 && (p.next_token.Type == SHORT || p.next_token.Type == LONG) {
				// LONG_BLANK need to be re extracted for starting the next Options_line
				p.s.Reject(p.current_token)
				p.next_token = nil
				p.current_token = nil
				return nil
			}
			// LONG_BLANK inside description
		case DEFAULT:
			return p.Consume_option_default()
		}

		description.AddNode(Description_node, p.current_token)
	}

	return fmt.Errorf("%s: parser stoped", p.current_node.Type)
}

func (p *DocoptParser) Consume_option_default() error {
	return nil
}
