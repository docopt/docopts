// vim: set ts=4 sw=4 sts=4 noet:
//
// unit test for docopt_language.go
//
package docopt_language

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func parse_usage(filename string) (*DocoptAst, error) {
	// data is []byte
	data, err := os.ReadFile(filename)
	p, err := ParserInit(data)
	if err != nil {
		return nil, err
	}
	ast := p.Parse()

	return ast, err
}

var DocoptNodes map[string]DocoptNodeType

func init_DocoptNodes() {
	DocoptNodes = make(map[string]DocoptNodeType)
	for t := Root; t < Last_node_type; t++ {
		DocoptNodes[t.String()] = t
	}
}

func TestParseUsages(t *testing.T) {
	usage_dir := "../grammar/usages/valid"
	filename := usage_dir + "/naval_fate.docopt"
	ast, err := parse_usage(filename)
	if err != nil {
		t.Errorf("parse_usage failed")
	}

	if ast == nil {
		t.Errorf("ast is nil")
	}

	ast_dir := usage_dir + "/../ast"
	ast_file := ast_dir + "/" + strings.Replace(filepath.Base(filename), ".docopt", "_ast.yaml", 1)

	if _, err := os.Stat(ast_file); err != nil {
		t.Errorf("ast_file is missing: '%s'", ast_file)
	}

	ast_from_yaml, err := Load_ast_from_yaml(ast_file)
	if err != nil {
		t.Errorf("error reading ast yaml file: '%s'", ast_file)
	}

	init_DocoptNodes()
	Match_ast(t, ast_from_yaml, ast)
}

// Compare all node from AstNode and DocoptAst
func Match_ast(t *testing.T, n *AstNode, parsed *DocoptAst) bool {
	expect := DocoptNodes[n.Node]
	if parsed.Type != expect {
		t.Errorf("expected node '%s' got %v", expect, parsed.Type)
		return false
	}

	if n.Token != nil {
		if n.Token.Value != parsed.Token.Value {
			t.Errorf("expected token '%s' got %v", n.Token.Value, parsed.Token.Value)
			return false
		}
	}

	nb_children := len(n.Children)
	if nb_children != len(parsed.Children) {
		t.Errorf("expected nb_children %d got %d", nb_children, len(parsed.Children))
		return false
	}

	if nb_children > 0 {
		for i := 0; i < nb_children; i++ {
			if !Match_ast(t, n.Children[i], parsed.Children[i]) {
				return false
			}
		}
	}

	return true
}

// ensure one Usage section
// ensure Usage matched case insensitive
// check p.options_node pointing to Options_section:
