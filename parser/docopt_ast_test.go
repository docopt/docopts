package docopt_language

import (
	"testing"
	// https://pkg.go.dev/github.com/stretchr/testify@v1.7.1/assert#pkg-functions
	"github.com/docopt/docopts/grammar/lexer"
	"github.com/stretchr/testify/assert"
	"strings"
)

func Test_AddNode(t *testing.T) {
	root := &DocoptAst{
		Type: Root,
	}

	n := root.AddNode(Usage_section, nil)
	assert := assert.New(t)

	assert.Len(root.Children, 1)
	assert.Equal(n, root.Children[0])
	assert.Equal(root, n.Parent)
}

func Test_Detach_child(t *testing.T) {
	root := &DocoptAst{
		Type: Root,
	}
	assert := assert.New(t)

	// ------------------------------------- one child
	n := root.AddNode(Usage_section, nil)
	assert.Len(root.Children, 1)
	assert.Equal(n, root.Children[0])

	n2 := root.Detach_child(0)
	assert.Equal(n2, n)
	assert.Len(root.Children, 0)
	assert.Nil(n.Parent)
	assert.Nil(n2.Parent)

	n, n2 = nil, nil

	// ------------------------------------- with multiple Children
	n1 := root.AddNode(Usage_section, nil)
	n2 = root.AddNode(Usage_line, nil)
	n3 := root.AddNode(Prologue, nil)
	assert.Len(root.Children, 3)

	for _, c := range root.Children {
		assert.Equal(root, c.Parent)
	}

	// --------------- detach first element
	root.Detach_child(0)
	assert.Nil(n1.Parent)

	assert.Len(root.Children, 2)
	assert.Equal(root.Children[0], n2)
	assert.Equal(root.Children[1], n3)

	// --------------- detach last element
	n1, n2 = root.Children[0], root.Children[1]
	n3 = root.AddNode(Free_section, nil)
	root.Detach_child(2)

	assert.Len(root.Children, 2)
	assert.Equal(root.Children[0], n1)
	assert.Equal(root.Children[1], n2)
	assert.Nil(n3.Parent)

	// --------------- detach middle element
	n3 = root.AddNode(Free_section, nil)
	assert.Len(root.Children, 3)
	assert.NotNil(n3.Parent)
	root.Detach_child(1)

	assert.Len(root.Children, 2)
	assert.Nil(n2.Parent)
	assert.Equal(root.Children[0], n1)
	assert.Equal(root.Children[1], n3)
}

// create a tree from a string of nodes_type
func helper_nested_AddNode(nodes_type string, parent *DocoptAst) *DocoptAst {
	DocoptNodes_init_reverse_map()
	splited := strings.Split(nodes_type, " ")
	s := 0
	nb := len(splited)
	if parent == nil {
		parent = &DocoptAst{
			Type: DocoptNodes[splited[s]],
		}
		s++
	}

	current := parent
	for i := s; i < nb; i++ {
		current = current.AddNode(DocoptNodes[splited[i]], nil)
	}

	return parent
}

func Test_Detach_from_parent(t *testing.T) {
	root := &DocoptAst{
		Type: Root,
	}
	assert := assert.New(t)

	helper_nested_AddNode("Usage_section Usage_line Usage_Expr Usage_optional_group Usage_Expr", root)
	assert.Len(root.Children, 1)
	assert.Equal(root.Children[0].Type, Usage_section)
	g := root.Children[0].Children[0].Children[0].Children[0]
	assert.Equal(g.Type, Usage_optional_group)

	p := g.Parent
	assert.NotNil(p)
	assert.True(g.Detach_from_parent())
	assert.Nil(g.Parent)
	assert.Equal(root.Children[0].Children[0].Children[0], p)
	assert.Len(p.Children, 0)
}

func Test_Find_recursive_by_Token(t *testing.T) {
	assert := assert.New(t)
	_, p, err := helper_load_usage(t, "test_input_allow_empty_argv.docopt")
	assert.Nil(err)
	assert.NotNil(p)

	tok := &lexer.Token{
		Type:  IDENT,
		Value: "mandatory",
	}

	pos, n, ok := p.usage_node.Find_recursive_by_Token(tok, -1)
	assert.True(ok)
	assert.Equal(1, pos)
	assert.Equal("mandatory", n.Token.Value)
	assert.Equal(Usage_Expr, n.Parent.Type)

	tok = &lexer.Token{
		Type:  PROG_NAME,
		Value: "myprog",
	}
	pos, n, ok = p.usage_node.Find_recursive_by_Token(tok, -1)
	assert.True(ok)
	assert.Equal(0, pos)
	assert.Equal(PROG_NAME, n.Token.Type)
	assert.Equal(Usage_line, n.Parent.Type)
}

func Test_Deep_copy(t *testing.T) {
	assert := assert.New(t)

	root := &DocoptAst{
		Type: Root,
	}
	helper_nested_AddNode("Usage_section Usage_line Usage_Expr Usage_optional_group Usage_Expr", root)
	usage := root.Children[0]
	usage.AddNode(Usage_line, nil)
	usage.AddNode(Usage_line, nil)

	new_root := root.Deep_copy()
	assert.Equal(len(new_root.Children), len(root.Children))
	if new_root == root {
		t.Errorf("root pointer copy, must be different")
	}
	for i, c := range usage.Children {
		if c == new_root.Children[0].Children[i] {
			t.Errorf("children pointer copy, must be different: %d", i)
		}
		if c.Parent == new_root.Children[0].Children[i].Parent {
			t.Errorf("parent pointer copy, must be different: %d", i)
		}
		if new_root.Children[0].Children[i].Parent != new_root.Children[0] {
			t.Errorf("child copy, must point back to copied Parent: %d", i)
		}
	}
	assert.Equal(root, new_root)
}

func Test_Deep_copy_exclude(t *testing.T) {
	assert := assert.New(t)

	root := &DocoptAst{
		Type: Root,
	}
	helper_nested_AddNode("Usage_section Usage_line Usage_Expr Usage_optional_group Usage_Expr", root)
	usage := root.Children[0]
	usage.AddNode(Usage_line, nil)
	usage.AddNode(Usage_line, nil)

	new_root := root.Deep_copy_exclude(&[]DocoptNodeType{Root})
	assert.Nil(new_root)

	new_root = root.Deep_copy_exclude(&[]DocoptNodeType{Usage_optional_group})
	expr := new_root.Children[0].Children[0].Children[0]
	assert.Equal(Usage_Expr, expr.Type)
	assert.Len(expr.Children, 0)

	// ensure the original as a child
	expr1 := root.Children[0].Children[0].Children[0]
	assert.Equal(Usage_Expr, expr1.Type)
	assert.Len(expr1.Children, 1)
}

func Test_Has_Parent(t *testing.T) {
	assert := assert.New(t)
	root := &DocoptAst{
		Type: Root,
	}
	helper_nested_AddNode("Usage_section Usage_line Usage_Expr Usage_optional_group Usage_Expr", root)
	usage := root.Children[0]
	usage.AddNode(Usage_line, nil)
	usage.AddNode(Usage_line, nil)

	expr := root.Children[0].Children[0].Children[0]
	assert.Equal(Usage_Expr, expr.Type)

	assert.True(expr.Has_Parent(Root))
	assert.True(expr.Has_Parent(Usage_section))
	assert.True(expr.Has_Parent(Usage_line))

	assert.False(expr.Has_Parent(Options_section))
	assert.False(expr.Has_Parent(Usage_Expr))
}
