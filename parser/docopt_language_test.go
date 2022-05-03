// vim: set ts=4 sw=4 sts=4 noet:
//
// unit test for docopt_language.go
//
package docopt_language

import (
	"github.com/docopt/docopts/grammar/lexer"
	// https://pkg.go.dev/github.com/stretchr/testify@v1.7.1/assert#pkg-functions
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
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

func helper_load_usage(t *testing.T, usage_filename string) (string, *DocoptParser, error) {
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
	filename, p, _ := helper_load_usage(t, "docopts.docopt")

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

	DocoptNodes_init_reverse_map()
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
		t.Errorf("%s: expected nb_children %d got %d", parsed.Type, nb_children, len(parsed.Children))
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
	_, p, _ := helper_load_usage(t, "docopts.docopt")

	options, err := p.transform_Options_section_to_map()
	assert := assert.New(t)
	assert.Nil(err)
	assert.Greater(len(options), 0, "options map must have elements")

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

	assert.Nil(options["-A"].Long)
}

func consume_me(p *DocoptParser) (Reason, error) {
	p.current_node.AddNode(NONE_node, nil)
	var dummy Reason_Value = 33
	return Reason{dummy, true}, nil
}

func check_first_child_type(t *testing.T, n *DocoptAst, expected DocoptNodeType) bool {
	if n.Type != expected {
		t.Errorf("Children[0] wrong type: got %s expected %s", n.Type, expected)
		return false
	}
	return true
}

func Test_Consume_loop(t *testing.T) {
	data := []byte("Usage: pipo molo")
	p, err := ParserInit(data)
	if err != nil {
		t.Errorf("ParserInit failed: %s", err)
	}

	p.CreateNode(Root, nil)
	if p.ast.Type != Root {
		t.Errorf("create Root: got %s expected %s", p.ast.Type, Root)
	}

	var our_def DocoptNodeType = Last_node_type + 1
	p.Parse_def[our_def] = &Consumer_Definition{
		create_self_node: true,
		create_node:      true,
		toplevel_node:    Usage_Expr,
		consume_func:     consume_me,
	}

	reason, err := p.Consume_loop(our_def)
	if err != nil {
		t.Errorf("Consume_loop returned err: %s", err)
	}

	if reason.Value != 33 {
		t.Errorf("Consume_loop returned reason value: got %d expected %d", reason.Value, 33)
	}

	c := p.ast.Children[0]
	if check_first_child_type(t, c, our_def) {
		c2 := c.Children[0]
		check_first_child_type(t, c2, Usage_Expr)
	}
}

// helper for DRY code
type Match_func func(*DocoptAst) (bool, error)

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func helper_ensure_matched(t *testing.T, f Match_func, node *DocoptAst) {
	matched, err := f(node)
	funcname := GetFunctionName(f)
	if err != nil {
		t.Errorf("%s: error %s", funcname, err)
	}
	if !matched {
		t.Errorf("%s: not matched node: %v", funcname, node)
	}
}

func Test_Match_Usage_node(t *testing.T) {
	// ============================================================== Usage_command
	node := &DocoptAst{
		Type: Usage_command,
		Token: &lexer.Token{
			Type:  IDENT,
			Value: "run",
			// not used in this context yet
			//Pos:        lexer.Position{Filename: "non-filename"},
			//Regex_name: "a regex",
			//State_name: "a state",
		},
	}

	if node.Type != Usage_command {
		t.Errorf("node Type error: got %s expected %s", node.Type, Usage_command)
	}

	command := "run"
	m := &MatchEngine{
		opts: DocoptOpts{},
		i:    0,
		argv: []string{command},
	}

	helper_ensure_matched(t, m.Match_Usage_node, node)
	if len(m.opts) != 1 {
		t.Errorf("Match_Usage_node: m.opts map wrong size, got %d expect %d", len(m.opts), 1)
	}
	if val, present := m.opts[command]; !present {
		t.Errorf("Match_Usage_node: map m.opts[%s] doesn't exists ", command)
	} else {
		if val != true {
			t.Errorf("Match_Usage_node: m.opts[%s] got %s expected true", command, val)
		}
	}
	if m.i != 1 {
		t.Errorf("Match_Usage_node: i should have increased got %d expected %d", m.i, 1)
	}

	// --------------------------------------- retest as Repeat-able argument
	node.Repeat = true
	m.i = 0
	// reset map
	m.opts = DocoptOpts{}
	helper_ensure_matched(t, m.Match_Usage_node, node)
	if len(m.opts) != 1 {
		t.Errorf("Match_Usage_node: m.opts map wrong size, got %d expect %d", len(m.opts), 1)
	}
	if val, present := m.opts[command]; !present {
		t.Errorf("Match_Usage_node: map m.opts[%s] doesn't exists ", command)
	} else {
		if val != 1 {
			t.Errorf("Match_Usage_node: m.opts[%s] got %s expected 1", command, val)
		}
	}
	if m.i != 1 {
		t.Errorf("Match_Usage_node: m.i should have increased got %d expected %d", m.i, 1)
	}

	// Repeat-able counted 2 times
	// another time (we rewind the m.argv index)
	m.i = 0
	helper_ensure_matched(t, m.Match_Usage_node, node)
	if val, present := m.opts[command]; !present {
		t.Errorf("Match_Usage_node: map m.opts[%s] doesn't exists ", command)
	} else {
		if val != 2 {
			t.Errorf("Match_Usage_node: m.opts[%s] got %v expected %d", command, val, 2)
		}
	}

	// ============================================================== Usage_argument
	name := "FILE"
	node = &DocoptAst{
		Type: Usage_argument,
		Token: &lexer.Token{
			Type:  IDENT,
			Value: name,
		},
	}

	m.i = 0
	// reset map
	m.opts = DocoptOpts{}
	helper_ensure_matched(t, m.Match_Usage_node, node)
	if val, present := m.opts[name]; !present {
		t.Errorf("Match_Usage_node: map m.opts[%s] doesn't exists ", name)
	} else {
		if val != m.argv[0] {
			t.Errorf("Match_Usage_node: m.opts[%s] got %v expected '%s'", name, val, m.argv[0])
		}
	}

	// -------------------------------- Repeat-able Usage_argument
	node.Repeat = true
	m.i = 0
	// reset map
	m.opts = DocoptOpts{}
	helper_ensure_matched(t, m.Match_Usage_node, node)
	if val, present := m.opts[name].([]string); !present {
		t.Errorf("Match_Usage_node: map m.opts[%s] doesn't exists ", name)
	} else {
		if len(val) != 1 {
			t.Errorf("Match_Usage_node: m.opts[%s] size got %d expected %d", name, len(val), 1)
		}
		if val[0] != m.argv[0] {
			t.Errorf("Match_Usage_node: m.opts[%s] => val[0] got %v expected %s", name, val[0], m.argv[0])
		}
	}
}

func Test_Match_Usage_node_Usage_long_option(t *testing.T) {
	assert := assert.New(t)

	option_name := "--myopt"
	m := &MatchEngine{
		opts: DocoptOpts{},
		i:    0,
		argv: []string{option_name},
	}

	// ========================================== node  Usage_long_option without child
	node := &DocoptAst{
		Type: Usage_long_option,
		Token: &lexer.Token{
			Type:  LONG,
			Value: option_name,
		},
	}
	helper_ensure_matched(t, m.Match_Usage_node, node)
	assert.Equal(true, m.opts[option_name].(bool), "Usage_long_option present must be true")

	// --------------------------------------- retest as Repeat-able Usage_long_option
	helper_test_Usage_option_Repeat_common(t, m, node, option_name)
}

// Call: helper_find_node(p, p.usage_node.Children[2], "LONG --speed")
// wrapper on Find_recursive_by_Token()
// `p` is used to revert Token Type from a string to a valid Token rune
func helper_find_node(p *DocoptParser, start *DocoptAst, node_desc string) (*DocoptAst, bool) {
	r := strings.Split(node_desc, " ")
	t := &lexer.Token{
		// Type is a rune from lexer's Symbols
		// map fail match will panic, OK
		Type:  p.all_symbols[r[0]],
		Value: r[1],
	}

	_, n, ok := start.Find_recursive_by_Token(t, -1)
	return n, ok
}

func helper_node_is_type(n *DocoptAst, t rune, v string) bool {
	if n.Token == nil {
		return false
	}
	return n.Token.Type == t && n.Token.Value == v
}

func Test_Match_Usage_node_Usage_short_option(t *testing.T) {
	assert := assert.New(t)

	option_name := "-m"
	m := &MatchEngine{
		opts: DocoptOpts{},
		i:    0,
		argv: []string{option_name},
	}

	// ========================================== node  Usage_short_option without child
	node := &DocoptAst{
		Type: Usage_short_option,
		Token: &lexer.Token{
			Type:  SHORT,
			Value: option_name,
		},
	}
	helper_ensure_matched(t, m.Match_Usage_node, node)
	assert.Equal(true, m.opts[option_name].(bool), "Usage_short_option present without alt must be true")

	// --------------------------------------- check for alternative long
	// we use a more complex tree from the parser
	// -m, --merge
	_, p, err := helper_load_usage(t, "test_input_short_option.docopt")
	assert.Nil(err)

	options, err := p.transform_Options_section_to_map()
	assert.Nil(err)
	assert.Greater(len(options), 0, "options map must have elements")

	// prepare MatchEngine
	m.options = &options
	m.i = 0
	m.opts = DocoptOpts{}
	long_option_name := "--merge"

	// look for -m node in usage_line 1 (Prog_name is at Children[0])
	if n, found := helper_find_node(p, p.usage_node.Children[1], "SHORT -m"); found {
		node = n
	} else {
		t.Errorf("node not found")
	}

	helper_ensure_matched(t, m.Match_Usage_node, node)
	assert.Equal(true, m.opts[long_option_name].(bool), "Usage_short_option present with alt must be true")
	// short option doesn't exist because we have a long alias in m.options
	_, exists := m.opts[option_name]
	assert.False(exists, "Usage_short_option -m with alt must not be set in m.opts")

	// --------------------------------------- retest as Repeat-able Usage_short_option (remove alternative long)
	m.options = nil
	helper_test_Usage_option_Repeat_common(t, m, node, option_name)

	// DISABLED
	//
	// // --------------------------------------- Usage_short_option that has no alternative
	// m.options = &options
	// option_name = "-A"
	// o, exists := m.Get_OptionRule(option_name)
	// assert.True(exists)
	// assert.Nil(o.Long)

	// m.i = 0
	// arg_value := "some_argument"
	// m.argv = Split_argv([]string{"-A", arg_value})
	// m.opts = DocoptOpts{}

	// if n, found := helper_find_node(p, p.usage_node.Children[2], "SHORT -A"); found {
	// 	node = n
	// } else {
	// 	t.Errorf("node not found for -A")
	// }
	// helper_ensure_matched(t, m.Match_Usage_node, node)
	// assert.Len(m.opts, 1)
	// opt, exists := m.opts[option_name].(string)
	// assert.True(exists, "Usage_short_option -A without alt must be set in m.opts")
	// assert.Equal(opt, arg_value)
}

func helper_test_Usage_option_Repeat_common(t *testing.T, m *MatchEngine, node *DocoptAst, option_name string) {
	assert := assert.New(t)
	// same testing code for Usage_short_option and Usage_long_option
	node.Repeat = true
	m.i = 0
	// reset map
	m.opts = DocoptOpts{}
	helper_ensure_matched(t, m.Match_Usage_node, node)

	assert.Len(m.opts, 1)
	_, exists := m.opts[option_name]
	assert.True(exists, "option_name key must exist in m.opts")
	assert.Equal(1, m.opts[option_name].(int), "Repeat-able: must be a counter")
	assert.Equal(1, m.i, "m.i invalid index")

	// Repeat-able counted 2 times
	// another time (we rewind the argument index)
	m.i = 0
	helper_ensure_matched(t, m.Match_Usage_node, node)
	assert.Equal(2, m.opts[option_name].(int), "2 times Repeat-able must be a counter")

	// ---------------------- option with argument
	option_value := "FILENAME"
	node.Children = []*DocoptAst{
		&DocoptAst{
			Type: Usage_argument,
			Token: &lexer.Token{
				Type:  ARGUMENT,
				Value: option_value,
			},
		},
	}
	node.Repeat = false

	option_argument := "some_file.txt"
	m.argv = append(m.argv, option_argument)
	m.i = 0
	helper_ensure_matched(t, m.Match_Usage_node, node)
	assert.Equal(2, m.i, "m.i invalid index")
	assert.Equal(option_argument, m.opts[option_name].(string), "option must have an argument")

	// ---------------------- Usage_short_option with argument Repeat-able
	node.Repeat = true
	m.argv = append(m.argv, option_argument, option_argument)
	assert.Len(m.argv, 4)
	m.i = 0
	helper_ensure_matched(t, m.Match_Usage_node, node)
	// move to one option
	assert.Equal(2, m.i, "m.i invalid index")
	assert.Equal([]string{option_argument}, m.opts[option_name].([]string), "Repeat-able option must have []string arguments")
}

func Test_Split_argv(t *testing.T) {
	input := []string{
		"pipo",
		"molo",
		"--",
		"=",
		"--opt",
		"--opt2=value",
		"--opt3=",
		"--opt4=\"value\"",
		"-S",
		"============",
		"pipo=molo",
		"--=",
	}

	expected := []string{
		"pipo",
		"molo",
		"--",
		"=",
		"--opt",
		"--opt2", "value",
		"--opt3=",
		"--opt4", "\"value\"",
		"-S",
		"============",
		"pipo=molo",
		"--=",
	}

	splited := Split_argv(input)
	assert.Equal(t, expected, splited, "argv spliting failed")
}

func Test_Match_empty_argv(t *testing.T) {
	assert := assert.New(t)
	_, p, err := helper_load_usage(t, "test_input_allow_empty_argv.docopt")
	assert.Nil(err)
	assert.NotNil(p)
	node := p.usage_node.Children[1]
	assert.Equal(Usage_line, node.Type)
	expr := node.Children[1]
	assert.Equal(Usage_Expr, expr.Type)

	m := &MatchEngine{}

	matched, err := m.Match_empty_argv(expr)
	assert.Nil(err)
	assert.True(matched)

	// test with Usage_line 2
	expr = p.usage_node.Children[2].Children[1]
	assert.Equal(Usage_Expr, expr.Type)
	assert.Len(expr.Children, 0)

	matched, err = m.Match_empty_argv(expr)
	assert.Nil(err)
	assert.True(matched)

	// test with Usage_line 3: fail
	expr = p.usage_node.Children[3].Children[1]
	assert.Equal(Usage_Expr, expr.Type)
	assert.Greater(len(expr.Children), 0, "Usage_Expr must have children")
	matched, err = m.Match_empty_argv(expr)
	assert.Nil(err)
	assert.False(matched, "Match_empty_argv must fail on mandatory argument")
}

func Test_Match_Usage_Expr(t *testing.T) {
	t.Skip()

	assert := assert.New(t)
	_, p, err := helper_load_usage(t, "docopts.docopt")
	assert.Nil(err)
	assert.NotNil(p)

	expr := p.usage_node.Children[1].Children[1]
	m := &MatchEngine{
		opts: DocoptOpts{},
		i:    0,
		argv: Split_argv([]string{"pipo", "molo"}),
	}
	matched, err := m.Match_Usage_Expr(expr)
	assert.Nil(err)
	assert.False(matched, "Match_Usage_Expr must fail on invalid argument")
	assert.Equal(0, m.i, "m.i should not change if matched is false && err == nil")

	m.argv = Split_argv([]string{"-h", "Usage: prog", ":"})
	m.i = 0
	matched, err = m.Match_Usage_Expr(expr)
	assert.Nil(err)
	assert.True(matched, "Match_Usage_Expr must succed")
}

func Test_Get_OptionRule(t *testing.T) {
	assert := assert.New(t)
	_, p, _ := helper_load_usage(t, "docopts.docopt")
	options, err := p.transform_Options_section_to_map()
	assert.Nil(err)

	m := &MatchEngine{}
	assert.Nil(m.options)

	// initialize
	m.options = &options
	assert.Greater(len(*m.options), 0, "OptionsMap must have elements")

	// an option without alternative long name
	o, ok := m.Get_OptionRule("-A")
	assert.True(ok)
	assert.Nil(o.Long)
}

// TODO:
// ensure one Usage section
// ensure Usage matched case insensitive
// check p.options_node pointing to Options_section:
