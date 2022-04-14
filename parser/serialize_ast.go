//
// serialization functions for Ast to YAML, loading and printing
//
package docopt_language

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// AstNode is a simpler structure of DocoptAst + Token
type AstNode struct {
	Node     string
	Children []*AstNode
	Token    *AstToken
}

type AstToken struct {
	Type  string
	Value string
}

// ==================================================
//
// Serialization methods use the same algorithm to walk through the nodes
// and they must me maintained together
//
// ==================================================

// serialize DocoptAst to YAML to stdout
func Serialize_DocoptAst(n *DocoptAst, indent string) {
	if n.Type == Root {
		fmt.Printf("---\n")
		fmt.Printf("%snode: %s\n", indent, n.Type)
	} else {
		fmt.Printf("%s- node: %s\n", indent, n.Type)
		indent += "  "
	}

	if n.Token != nil {
		fmt.Printf("%stoken: { type: %s, value: %q }\n", indent, n.Token.Regex_name, n.Token.Value)
	}

	nb_children := len(n.Children)
	if nb_children > 0 {
		fmt.Printf("%schildren:\n", indent)
		for i := 0; i < nb_children; i++ {
			Serialize_DocoptAst(n.Children[i], indent)
		}
	}
}

// yaml serialize our AstNode to stdout
func Serialize_ast(n *AstNode, indent string) {
	if n.Node == "Root" {
		fmt.Printf("---\n")
		fmt.Printf("%snode: %s\n", indent, n.Node)
	} else {
		fmt.Printf("%s- node: %s\n", indent, n.Node)
		indent += "  "
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
