package docopt_language

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type AstNode struct {
	Node     string
	Children []*AstNode
	Token    *AstToken
}

type AstToken struct {
	Type  string
	Value string
}

// yaml serialize our AST
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

func Load_ast_from_yaml(filename string) (error, *AstNode) {
	ast := AstNode{}
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading %s: %v", filename, err), nil
	}
	err = yaml.Unmarshal([]byte(data), &ast)
	if err != nil {
		fmt.Errorf("error: %v", err)
	}
	return nil, &ast
}
