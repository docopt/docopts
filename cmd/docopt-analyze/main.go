//
// docopt-analyze parse and generate parsing information about docopt language
//
package main

import (
	"fmt"
	docopt_language "github.com/docopt/docopts/parser"
	"os"
	// still use legacy embedded docopt lib
	"github.com/alecthomas/repr"
	"github.com/docopt/docopts/docopt-go"
	"strings"
)

var Usage string = `docopt language grammar analyzer

Usage:
  docopt-analyze [-y] [-r] [-s] FILENAME

Options:
  -y      Serialize yaml AST for unit testing
	-r      Print AST as repr
	-s      Simple_print_tree Print using minimalist tree
`

var Version string = "0.2"

// print_ast() visual print AST for user
func print_ast(current_node *docopt_language.DocoptAst, indent_level int) {
	var indent string
	for i := 0; i < indent_level; i++ {
		indent += " "
	}

	nb_children := len(current_node.Children)
	repeatable := ""
	if current_node.Repeat {
		repeatable = "REPEATABLE one or more"
	}

	if current_node.Token != nil {
		// final node with token
		if current_node.Type == docopt_language.Options_node {
			// Options_node are unmatched node
			fmt.Printf("%s- %s: %s\n", indent, current_node.Type, current_node.Token.GoString())
		} else {
			fmt.Printf("%s- %s: %q %s\n", indent, current_node.Type, current_node.Token.Value, repeatable)
		}
	} else {
		// syntax node without token
		switch nb_children {
		case 0:
			fmt.Printf("%s%s: []\n", indent, current_node.Type)
		case 1:
			fmt.Printf("%s%s: %s\n", indent, current_node.Type, repeatable)
		default:
			fmt.Printf("%s%s: (%d) %s\n", indent, current_node.Type, nb_children, repeatable)
		}
	}

	if current_node.Type == docopt_language.Prologue {
		// Prologue is printed merged (no nested level, only one level of Children)
		children := current_node.Children
		indent += "  "
		new_line := true
		out := ""
		for i := 0; i < nb_children; i++ {
			if new_line && children[i].Token.Type != docopt_language.NEWLINE {
				out = fmt.Sprintf("%s- %s: \"%s", indent, children[i].Type, children[i].Token.Value)
				new_line = false
			} else if children[i].Token.Type == docopt_language.NEWLINE {
				if new_line {
					continue
				}

				fmt.Printf("%s\"\n", out)
				new_line = true
			} else {
				out += children[i].Token.Value
			}
		}
	} else {
		for i := 0; i < nb_children; i++ {
			print_ast(current_node.Children[i], indent_level+2)
		}
	}
}

// nil_parent() put all parent at nil, helper for not poluting output with pointers with repr
func nil_parent(n *docopt_language.DocoptAst) {
	if n.Parent != nil {
		n.Parent = nil
	}
	for _, c := range n.Children {
		nil_parent(c)
	}
}

func main() {
	docopt_p := &docopt.Parser{
		OptionsFirst: true,
	}

	args, err := docopt_p.ParseArgs(Usage, nil, Version)
	if err != nil {
		msg := fmt.Sprintf("you're not suppose to get here: %v\n", err)
		panic(msg)
	}

	filename := args["FILENAME"].(string)
	print_repr := args["-r"].(bool)
	simple_print_tree := args["-s"].(bool)
	serialize := args["-y"].(bool)

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("error: fail to open '%s': $s\n", filename)
		os.Exit(1)
	} else {
		if !print_repr && !serialize {
			fmt.Printf("parsing: %s\n", filename)
		}
	}

	// get our grammar parser as p
	p, err := docopt_language.ParserInit(data)
	if err != nil {
		fmt.Println(err)
		os.Exit(999)
	}
	ast := p.Parse()

	if print_repr {
		nil_parent(ast)
		repr.Println(ast)
	} else if serialize {
		var out []string
		docopt_language.Serialize_DocoptAst(ast, "", nil, &out)
		fmt.Print(strings.Join(out, ""))
	} else if simple_print_tree {
		docopt_language.Simple_print_tree(ast, "")
	} else {
		fmt.Printf("Detected Prog_name:%s\n", p.Prog_name)

		fmt.Printf("============== AST ===============\n")
		print_ast(ast, 0)

		fmt.Printf("number of error: %d\n", p.Error_count)
		for _, e := range p.Errors {
			fmt.Println(e)
		}
		os.Exit(p.Error_count)
	}
}
