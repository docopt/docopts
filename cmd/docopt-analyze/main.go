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
	if current_node.Token != nil {
		fmt.Printf("%s%s:%q\n", indent, current_node.Type, current_node.Token.Value)
	} else {
		fmt.Printf("%s%s\n", indent, current_node.Type)
	}
	for _, n := range current_node.Children {
		print_ast(n, indent_level+2)
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
		os.Exit(1)
	}
	ast := p.Parse()

	fmt.Printf("detected Prog_name:%s\n", p.Prog_name)
	print_ast(ast, 0)

	fmt.Printf("number of error: %d\n", p.Error_count)
	os.Exit(p.Error_count)
}
