package main

import (
	"fmt"
	docopt_language "github.com/docopt/docopts/parser"
	"os"
)

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

func main() {
	filename := os.Args[1]
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("error: fail to open '%s': $s\n", filename)
		os.Exit(1)
	} else {
		fmt.Printf("parsing: %s\n", filename)
	}

	p, err := docopt_language.ParserInit(data)
	if err != nil {
		fmt.Println(err)
		os.Exit(999)
	}
	ast := p.Parse()

	fmt.Printf("Detected Prog_name:%s\n", p.Prog_name)

	fmt.Printf("============== AST ===============\n")
	print_ast(ast, 0)

	fmt.Printf("number of error: %d\n", p.Error_count)
	for _, e := range p.Errors {
		fmt.Println(e)
	}
	os.Exit(p.Error_count)
}
