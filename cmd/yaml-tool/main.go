package main

import (
	"github.com/docopt/docopts/parser"
	"os"
)

func main() {
	err, ast := docopt_language.Load_ast_from_yaml(os.Args[1])
	if err != nil {
		panic(err)
	}
	docopt_language.Serialize_ast(ast, "")
}
