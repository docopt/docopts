//
// serialization functions for Ast to YAML, loading and printing
//
package docopt_language

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

// AstNode is a simpler structure of DocoptAst + Token
type AstNode struct {
	Node             string
	Children         []*AstNode
	Token            *AstToken
	Usage_line_input *string
	Repeat           *string
}

type AstToken struct {
	Type  string
	Value string
}

// ==================================================
//
// Serialization methods use the same algorithm to walk through the nodes
// and they MUST be maintained together
//
// ==================================================

type pending_usage_string struct {
	insert_index int
	collected    []string
}

// Serialize_DocoptAst() serialize DocoptAst to YAML to stdout
// (for cmd/docopt-analyze/main.go -y)
// output a YAML suitable for Load_ast_from_yaml()
//
// call:
//   var out []string
//   docopt_language.Serialize_DocoptAst(ast, "", nil, &out)
//
// Always call from the Root with empty `indent` and nil `usage_string`
// Result of the serialization is stored in `out`
func Serialize_DocoptAst(n *DocoptAst, indent string, usage_string *pending_usage_string, out *[]string) {
	if n.Type == Root {
		*out = append(*out, fmt.Sprintf("---\n"))
		*out = append(*out, fmt.Sprintf("%snode: %s\n", indent, n.Type))
	} else {
		*out = append(*out, fmt.Sprintf("%s- node: %s\n", indent, n.Type))
		indent += "  "
	}

	if n.Repeat {
		*out = append(*out, fmt.Sprintf("%srepeat: true\n", indent))
	}

	if n.Type == Usage_line {
		usage_string = &pending_usage_string{
			insert_index: len(*out),
			collected:    make([]string, 0, 10),
		}
		// reserve the place holder
		*out = append(*out,
			fmt.Sprintf("%susage_line_input: %s\n", indent,
				fmt.Sprintf("USAGE WILL GO HERE: %d", usage_string.insert_index)))
	}

	nb_children := len(n.Children)

	if n.Token != nil {
		*out = append(*out, fmt.Sprintf("%stoken: { type: %s, value: %q }\n", indent, n.Token.Regex_name, n.Token.Value))
		if usage_string != nil {
			// repeat-able long option with argument the elipsis will be put at the end of the children loop
			if n.Repeat && nb_children == 0 {
				usage_string.collected = append(usage_string.collected, n.Token.Value+"...")
			} else {
				usage_string.collected = append(usage_string.collected, n.Token.Value)
			}
		}
	}

	if nb_children > 0 {
		close_token := ""
		if n.Type == Usage_optional_group {
			usage_string.collected = append(usage_string.collected, "[")
			close_token = "]"
		}
		if n.Type == Usage_required_group {
			usage_string.collected = append(usage_string.collected, "(")
			close_token = ")"
		}
		if n.Type == Usage_long_option {
			usage_string.collected = append(usage_string.collected, "=")
		}

		*out = append(*out, fmt.Sprintf("%schildren:\n", indent))
		for i := 0; i < nb_children; i++ {
			Serialize_DocoptAst(n.Children[i], indent, usage_string, out)
			if close_token != "" && i < nb_children-1 {
				usage_string.collected = append(usage_string.collected, "|")
			}
		}

		if close_token != "" {
			// group Repeat
			if n.Repeat {
				usage_string.collected = append(usage_string.collected, close_token+"...")
			} else {
				usage_string.collected = append(usage_string.collected, close_token)
			}
		}
		// repeat-able long option with argument put the elipsis at the end
		if n.Type == Usage_long_option && n.Repeat {
			usage_string.collected = append(usage_string.collected, close_token+"...")
		}
	}

	if n.Type == Usage_line {
		(*out)[usage_string.insert_index] = fmt.Sprintf(
			"%susage_line_input: \"%s\"\n",
			indent,
			strings.Join(usage_string.collected, " "))
	}
}

// Yaml serialize our AstNode to stdout (for cmd/yaml-tool/main.go)
// for validating purpose.
// Was the same walking tree as Serialize_DocoptAst() but on YAML AstNode input.
// See Load_ast_from_yaml()
// ./yaml-tool ../../grammar/usages/ast/docopts_ast.yaml | diff -u ../../grammar/usages/ast/docopts_ast.yaml -
// The above command must return no diff
func Serialize_ast(n *AstNode, indent string) {
	if n.Node == "Root" {
		fmt.Printf("---\n")
		fmt.Printf("%snode: %s\n", indent, n.Node)
	} else {
		fmt.Printf("%s- node: %s\n", indent, n.Node)
		indent += "  "
	}

	if n.Usage_line_input != nil {
		fmt.Printf("%susage_line_input: %q\n", indent, *n.Usage_line_input)
	}

	if n.Repeat != nil {
		fmt.Printf("%srepeat: true\n", indent)
	}
	if n.Token != nil {
		fmt.Printf("%stoken: { type: %s, value: %q }\n", indent, n.Token.Type, n.Token.Value)
	}

	nb_children := len(n.Children)
	if nb_children > 0 {
		fmt.Printf("%schildren:\n", indent)
		for i := 0; i < nb_children; i++ {
			Serialize_ast(n.Children[i], indent)
		}
	}
}

func Load_ast_from_yaml(filename string) (*AstNode, error) {
	ast := AstNode{}
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %v", filename, err)
	}
	err = yaml.Unmarshal([]byte(data), &ast)
	if err != nil {
		fmt.Errorf("error: %v", err)
	}
	return &ast, nil
}

// light AST print to stdout omiting descending nodes:
// Prologue, Free_section, Option_description
func Simple_print_tree(n *DocoptAst, indent string) {
	nb_children := len(n.Children)
	fmt.Printf("%s%s", indent, n.Type)

	if n.Token != nil {
		fmt.Printf(" %q", n.Token.Value)
		if n.Repeat {
			fmt.Printf("...")
		}
	}

	if nb_children > 0 {
		if n.Repeat {
			fmt.Printf(" [%d]...\n", nb_children)
		} else {
			fmt.Printf(" [%d]\n", nb_children)
		}
		if n.Type == Prologue || n.Type == Free_section || n.Type == Option_description {
			// skip those nodes
			return
		} else {
			for i := 0; i < nb_children; i++ {
				Simple_print_tree(n.Children[i], indent+"  ")
			}
		}
	} else {
		fmt.Printf("\n")
	}
}
