package main

import (
	"fmt"
	docopt_language "github.com/docopt/docopts/parser"
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

	p, _ := docopt_language.ParserInit(data)
	p.Parse()

	//print_ast(ast)

	fmt.Printf("number of error: %d\n", p.Error_count)
	os.Exit(p.Error_count)
}
