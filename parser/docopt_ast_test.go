package docopt_language

import (
	"testing"
	// https://pkg.go.dev/github.com/stretchr/testify@v1.7.1/assert#pkg-functions
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
