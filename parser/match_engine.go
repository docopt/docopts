package docopt_language

import (
	"fmt"
)

// same as Opts in legacy docopt-go
type DocoptOpts map[string]interface{}

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

	matched = false
	expr := u.Children[1]
	if expr.Type != Usage_Expr {
		err = fmt.Errorf("Match_Usage_line: AST error 2nd node is not Usage_Expr, got: %s", expr.Type)
		return
	}

	matched, args, err = Match_Usage_Expr(expr, argv, i)
	return
}

func Match_Usage_Expr(expr *DocoptAst, argv *[]string, i int) (matched bool, args DocoptOpts, err error) {
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
	nb := len(*argv)
	if nb == 0 {
		matched, err = Match_empty_argv(expr)
		return
	}

forLoopMatch_Usage_Expr:
	for _, c := range expr.Children {
		switch c.Type {
		case Usage_optional_group, Usage_required_group:
			matched, args, err = Match_Usage_Group(c, argv, i)
			if err == nil {
				if c.Type == Usage_optional_group {
					// optional
					if matched {
						i++
					}
					continue
				} else {
					// required
					if matched {
						i++
						continue
					} else {
						err = fmt.Errorf("expected Usage_command %s, faild to match '%s'", c.Token.Value, (*argv)[i])
					}
				}
			}
			// some errors
			break forLoopMatch_Usage_Expr
		case Usage_command,
			Usage_short_option,
			Usage_long_option,
			Usage_argument:
			matched, err = Match_Usage_node(c, argv, &i, &args)
			if err == nil {
				if matched {
					i++
					continue
				} else {
					err = fmt.Errorf("expected Usage_command %s, faild to match '%s'", c.Token.Value, (*argv)[i])
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

func Match_empty_argv(expr *DocoptAst) (bool, error) {
	return false, fmt.Errorf("unfinished method Match_empty_argv")
}

func Match_Usage_Group(g *DocoptAst, argv *[]string, i int) (matched bool, args DocoptOpts, err error) {
	err = fmt.Errorf("unfinished method Match_Usage_Group")
	return
}

func Match_Usage_node(n *DocoptAst, argv *[]string, i *int, opts *DocoptOpts) (matched bool, err error) {
	matched = false
	a := (*argv)[*i]
	k := n.Token.Value
	switch n.Type {
	case Usage_command:
		if a == k {
			matched = true
			if n.Repeat {
				//  command take no value (as option without argument)
				// check key exists
				if val, present := (*opts)[a].(int); present {
					(*opts)[k] = val + 1
				} else {
					(*opts)[k] = 1
				}
			} else {
				// Single
				(*opts)[k] = true
			}
		} else {
			(*opts)[k] = false
		}
	case Usage_argument:
		matched = true
		if n.Repeat {
			if val, present := (*opts)[a].([]string); present {
				(*opts)[k] = append(val, a)
			} else {
				(*opts)[k] = []string{a}
			}
		} else {
			// Single
			(*opts)[k] = a
		}
	case Usage_short_option:
		err = fmt.Errorf("unhandled node Type: %s", n.Type)
		//if (*argv)[*i] == n.Token.Value {
		//	matched = true
		//	(*opts)[a] = true
		//	if len(n.Children) > 0 {
		//		c := n.Children[0]
		//		if c.Type == Usage_argument {
		//			// check key exists
		//			if _, present := (*opts)[c.Token.Value]; ! present {
		//				(*opts)[c.Token.Value] = true
		//			} else {
		//				err = fmt.Errorf("key '%s' already exists for")
		//		}
		//	}
		//} else {
		//	(*opts)[a] = false
		//}
	case Usage_long_option:
		err = fmt.Errorf("unhandled node Type: %s", n.Type)
	default:
		err = fmt.Errorf("unhandled node Type: %s", n.Type)
	}

	if matched {
		// move argv index
		*i++
	}

	return
}
