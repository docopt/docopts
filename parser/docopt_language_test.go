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

func parse_usage(filename string) (*DocoptParser, error) {
	// data is []byte
	data, err := os.ReadFile(filename)
	p, err := ParserInit(data)
	if err != nil {
		return p, err
	}
	// AST is available from p.ast
	p.Parse()
	return p, err
}

var DocoptNodes map[string]DocoptNodeType

func init_DocoptNodes() {
	DocoptNodes = make(map[string]DocoptNodeType)
	for t := Root; t < Last_node_type; t++ {
		DocoptNodes[t.String()] = t
	}
}

func load_usage(t *testing.T, usage_filename string) (string, *DocoptParser, error) {
	usage_dir := "../grammar/usages/valid"
	filename := usage_dir + "/" + usage_filename
	if _, err := os.Stat(filename); err != nil {
		t.Errorf("doctop file is missing: '%s'", filename)
		return filename, nil, err
	}

	p, err := parse_usage(filename)
	if err != nil {
		t.Errorf("parse_usage failed for: %s", filename)
	} else if p.ast == nil {
		t.Errorf("ast is nil")
	}
	return filename, p, err
}

func TestParseUsages(t *testing.T) {
	filename, p, _ := load_usage(t, "docopts.docopt")

	usage_dir := filepath.Dir(filename)
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
	Match_ast(t, ast_from_yaml, p.ast)
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

func Test_transform_Options_section_to_map(t *testing.T) {
	_, p, _ := load_usage(t, "docopts.docopt")

	options, err := p.transform_Options_section_to_map()
	if err != nil {
		t.Errorf("transform_Options_section_to_map error: %v", err)
	}

	if len(options) == 0 {
		t.Errorf("transform_Options_section_to_map: options map has no element")
	}

	if options["-s"].Arg_count != 1 {
		t.Errorf("transform_Options_section_to_map: options[\"-s\"] (separator) as not 1 mandatory argument count")
	}

	expected := "<str>"
	if *options["--separator"].Argument_name != expected {
		t.Errorf("transform_Options_section_to_map: options[\"--separator\"].Argument_name got: %q expected %q",
			expected,
			*options["--separator"].Argument_name)
	}

	if options["-s"] != options["--separator"] {
		t.Errorf("transform_Options_section_to_map: options[\"-s\"] != options[\"--separator\"]")
	}
}

// ensure one Usage section
// ensure Usage matched case insensitive
// check p.options_node pointing to Options_section:
