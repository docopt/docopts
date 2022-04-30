package docopt_language

import (
	"fmt"
	"strings"
)

// same as Opts in legacy docopt-go
type DocoptOpts map[string]interface{}

// Split_argv() transforms an argv vector by spliting long option assigned:
// --long=ARGUMENT ==> "--long", "ARGUMENT"
func Split_argv(argv []string) []string {
	n_argv := len(argv)
	// at most we will double the size
	splited := make([]string, 0, 2*n_argv)
	for _, a := range argv {
		if e := strings.IndexByte(a, '='); e != -1 {
			n := len(a)
			if n > 3 && a[0] == '-' && a[1] == '-' && e < n-2 {
				splited = append(splited, a[0:e], a[e+1:n])
				continue
			}
		}
		// keep unchanged
		splited = append(splited, a)
	}
	return splited
}

// MatchArgs() associate argv (os.Args) to parsed Options / Argument
// algorithm derive from docopt.ParseArgs() docopt-go/docopt.go
func (p *DocoptParser) MatchArgs(argv []string) (args DocoptOpts, err error) {
	if p.ast == nil {
		err = fmt.Errorf("error: ast is nil")
		return
	}

	//if len(usageSections) == 0 {
	//	err = newLanguageError("\"usage:\" (case-insensitive) not found.")
	//	return
	//}
	//if len(usageSections) > 1 {
	//	err = newLanguageError("More than one \"usage:\" (case-insensitive).")
	//	return
	//}

	// options := parseDefaults(doc)
	// READY
	//options, err := p.transform_Options_section_to_map()
	//if err != nil {
	//	return nil, err
	//}

	// formal, err := FormalUsage(usage)
	// pat, err := ParsePattern(formal, &options)

	// loop over Usage_line to find one that match
	var found int = -1
	for _, l := range p.usage_node.Children {
		if l.Type == Usage_line {
			Match_Usage_line(l, &argv, 0)
		}
	}

	if err == nil && found > -1 {
		// success
		return
	}

	// no match found, argument parsing error
	if err != nil {
		// error previously caught
		return
	}

	err = fmt.Errorf("no match found for argv: %v", argv)
	return
}

func Match_Usage_line(u *DocoptAst, argv *[]string, i int) (matched bool, args DocoptOpts, err error) {
	// in Usage_line node structure:
	//   Usage_line [2]
	//     Prog_name "PROG_NAME"
	//     Usage_Expr [6]
	// So Usage_Expr is always the 2nd child

	m := &MatchEngine{
		opts: DocoptOpts{},
		i:    i,
		argv: Split_argv(*argv),
	}
	matched = false
	expr := u.Children[1]
	if expr.Type != Usage_Expr {
		err = fmt.Errorf("Match_Usage_line: AST error 2nd node is not Usage_Expr, got: %s", expr.Type)
		return
	}

	matched, err = m.Match_Usage_Expr(expr)
	return
}

func (m *MatchEngine) Match_Usage_Expr(expr *DocoptAst) (matched bool, err error) {
	// docopts [options] -h <msg> : [<argv>...]
	// docopts -h "Usage: myprog cat [-c COLOR] FILENAME" : cat pipo
	//   Usage_Expr: cat [-c COLOR] STRING
	//   argv: "pipo"
	//
	//  Usage_Expr: (3)
	//    - Usage_command: "cat"
	//    Usage_optional_group: (2)
	//      - Usage_short_option: "-c"
	//      - Usage_argument: "COLOR"
	//    - Usage_argument: "FILENAME"

	matched = false
	nb := len(m.argv)
	if nb == 0 {
		matched, err = m.Match_empty_argv(expr)
		return
	}

forLoopMatch_Usage_Expr:
	for _, c := range expr.Children {
		switch c.Type {
		case Usage_optional_group, Usage_required_group:
			matched, err = m.Match_Usage_Group(c)
			if err == nil {
				if c.Type == Usage_required_group {
					// required
					if matched {
						continue
					} else {
						err = fmt.Errorf("expected Usage_required_group, faild to match '%s'", m.argv[m.i])
						return
					}
				}
				// optional can continue
				continue
			}
			// some errors
			break forLoopMatch_Usage_Expr
		case Usage_command,
			Usage_short_option,
			Usage_long_option,
			Usage_argument:
			matched, err = m.Match_Usage_node(c)
			if err == nil {
				if matched {
					continue
				} else {
					err = fmt.Errorf("expected Usage_command %s, faild to match '%s'", c.Token.Value, m.argv[m.i])
				}
			}
			// some errors
			break forLoopMatch_Usage_Expr
		default:
			err = fmt.Errorf("unhandled ast node %s", c.Type)
			break forLoopMatch_Usage_Expr
		} // end switch c.Type

		// unmatched Type
		err = fmt.Errorf("you are not supposed to be here: forLoopMatch_Usage_Expr")
		break forLoopMatch_Usage_Expr
	}

	return
}

// Match_empty_argv() `expr` is pointing to an Usage_Expr, we validate that this Expr
// accept empty arg = no mandatory argument
func (m *MatchEngine) Match_empty_argv(expr *DocoptAst) (bool, error) {
	if expr.Type != Usage_Expr {
		return false, fmt.Errorf("Match_empty_argv only accept Usage_Expr node as arugment, got %s", expr.Type)
	}

	if len(expr.Children) == 0 {
		return true, nil
	}

	// should only have optional group
	for _, c := range expr.Children {
		if c.Type == Usage_optional_group {
			continue
		}
		return false, nil
	}

	return true, nil
}

// Match_options() check if the current argv[i] can be in Usage_options_shortcut
// the argument must be a option SHORT or LONG
// seek in OptionsMap => match
func (m *MatchEngine) Match_options() (bool, error) {
	a := m.argv[m.i]
	s := len(a)
	is_short := s == 2 && a[0] == '-' && a[1] != '-'
	start_with_2dash := a[0] == '-' && a[1] == '-'

	if start_with_2dash && s < 4 || !is_short {
		return false, nil
	}

	//if o, ok := m.Get_OptionRule(a); ok {
	//	m.opts[o.Long]
	//}

	return false, fmt.Errorf("Match_options not implemented")
}

func (m *MatchEngine) Match_Usage_Group(g *DocoptAst) (matched bool, err error) {
	if g.Type != Usage_optional_group && g.Type != Usage_required_group {
		err = fmt.Errorf("unsupported node, expecting Usage_optional_group or Usage_required_group, got %s", g.Type)
		return
	}

	nb := len(g.Children)
	if nb != 1 {
		// group surround a single node Usage_Expr
		err = fmt.Errorf("Expected only one Children got %d, for Group: %s", nb, g.Type)
		return
	}

	expr := g.Children[0]
	old_i := m.i

	// [options]
	if g.Type == Usage_options_shortcut {
		matched, err = m.Match_options()
	} else {
		matched, err = m.Match_Usage_Expr(expr)
	}
	if !matched {
		m.i = old_i
	}
	return
}

type MatchEngine struct {
	i       int
	argv    []string
	opts    DocoptOpts
	options *OptionsMap
}

func (m *MatchEngine) Get_OptionRule(k string) (*OptionRule, bool) {
	if m.options != nil {
		if o, ok := (*m.options)[k]; ok {
			return o, true
		}
	}
	return nil, false
}

type MachAssignType int

const (
	String_type MachAssignType = 1 + iota
	String_repeat
	Bool_type
	Bool_repeat
)

func (m *MatchEngine) Match_Usage_option(n *DocoptAst, a *string, k *string) (bool, error) {
	matched := false
	var t MachAssignType
	if len(n.Children) > 0 && n.Children[0].Type == Usage_argument {
		// option has a required argument
		if len(m.argv)-(m.i+1) > 0 {
			old_i := m.i
			// will also be moved +1 at the end eating 2 argv
			m.i++

			if n.Repeat {
				t = String_repeat
			} else {
				t = String_type
			}
			// we force the key assignment with the option's name k
			if err := m.Match_Assign(t, n.Children[0], k); err != nil {
				m.i = old_i
				return false, err
			}
			matched = true
		} else {
			// no more argument in argv[]
			return false, fmt.Errorf("option: %s require an argument", *k)
		}
	} else {
		// option has no argument (true or false)
		if n.Repeat {
			t = Bool_repeat
		} else {
			t = Bool_type
		}
		if err := m.Match_Assign(t, n, k); err != nil {
			return false, err
		}
		matched = true
	}
	return matched, nil
}

func (m *MatchEngine) Match_Assign(t MachAssignType, n *DocoptAst, force_key *string) error {
	var k *string
	if force_key != nil {
		k = force_key
	} else {
		k = &n.Token.Value
	}
	switch t {
	case String_type, String_repeat:
		a := m.argv[m.i]
		if n.Repeat || t == String_repeat {
			if val, present := m.opts[*k].([]string); present {
				m.opts[*k] = append(val, a)
			} else {
				m.opts[*k] = []string{a}
			}
		} else {
			// Single
			m.opts[*k] = a
		}
	case Bool_type, Bool_repeat:
		if n.Repeat || t == Bool_repeat {
			//  command take no value (as option without argument)
			// check key exists
			if val, present := m.opts[*k].(int); present {
				m.opts[*k] = val + 1
			} else {
				m.opts[*k] = 1
			}
		} else {
			// Single
			m.opts[*k] = true
		}
	default:
		return fmt.Errorf("Match_Assign: unsupported MachAssignType: %d", t)
	}

	return nil
}

func (m *MatchEngine) Match_Usage_node(n *DocoptAst) (matched bool, err error) {
	matched = false
	a := m.argv[m.i]
	k := n.Token.Value
	// TODO: handle option default value
	switch n.Type {
	case Usage_command:
		if a == k {
			err = m.Match_Assign(Bool_type, n, nil)
			if err != nil {
				return
			}
			matched = true
		} else {
			m.opts[k] = false
		}
	case Usage_argument:
		err = m.Match_Assign(String_type, n, nil)
		if err != nil {
			return
		}
		matched = true
	case Usage_long_option:
		start_with_2dash := a[0] == '-' && a[1] == '-'
		if start_with_2dash && a == k {
			matched, err = m.Match_Usage_option(n, &a, &k)
		} else {
			// not matched
			m.opts[k] = false
		}
	case Usage_short_option:
		is_short := len(a) == 2 && a[0] == '-' && a[1] != '-'
		if is_short && a == k {
			// replace short option by its long version if it exists
			if alternative, ok := m.Get_OptionRule(k); ok {
				// ok even if alternative.Long is nil: Match_Usage_option uses Match_Assign
				// which will use default option's short name (k) in such case
				matched, err = m.Match_Usage_option(n, &a, alternative.Long)
			} else {
				matched, err = m.Match_Usage_option(n, &a, &k)
			}
			matched = true
		} else {
			// not matched
			if alternative, ok := m.Get_OptionRule(k); ok && alternative.Long != nil {
				m.opts[*alternative.Long] = false
			} else {
				m.opts[k] = false
			}
		}
	default:
		err = fmt.Errorf("unhandled node Type: %s", n.Type)
	}

	if matched {
		// move argv index
		m.i++
	}

	return
}
