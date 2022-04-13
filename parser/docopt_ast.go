//go:generate stringer -type=DocoptNodeType
package docopt_language

import (
	"github.com/docopt/docopts/grammar/lexer"
)

type DocoptNodeType int

// ast nodes types
const (
	// NONE is a value for not to be matched DocoptNodeType
	Unmatched_node DocoptNodeType = -1
	NONE_node      DocoptNodeType = 0 + iota
	Root
	Prologue
	Prologue_node
	Usage_section
	Usage
	Usage_line
	Prog_name
	Usage_short_option
	Usage_long_option
	Usage_argument
	Usage_unmatched_punct
	Usage_command
	Usage_optional_group
	Usage_required_group
	Usage_Expr
	Group_alternative
	Free_section
	Section_name
	Section_node
	Options_section
	Options_node
	Option_line
	Option_short
	Option_long
	Option_argument
	Option_alternative_group
	Option_description
	Description_node
)

type DocoptAst struct {
	Type     DocoptNodeType
	Token    *lexer.Token
	Children []*DocoptAst
	Parent   *DocoptAst
	Repeat   bool
}

func (n *DocoptAst) AddNode(node_type DocoptNodeType, t *lexer.Token) *DocoptAst {
	new_node := &DocoptAst{
		Type:   node_type,
		Token:  t,
		Parent: n,
		Repeat: false,
	}
	n.Children = append(n.Children, new_node)
	return new_node
}

// Replace_children_with_group move all actual children of the node, to a new node of type `node_type`.
// Current node: `parent` becomes the new parent, and all its children becomes
// children of the new node `group_node`.
// returns the new group node recreated
func (parent *DocoptAst) Replace_children_with_group(node_type DocoptNodeType) *DocoptAst {
	group_node := &DocoptAst{
		Type:     node_type,
		Token:    nil,
		Parent:   parent,
		Children: parent.Children,
	}

	// move actual Children to new node
	for _, c := range group_node.Children {
		c.Parent = group_node
	}

	// create anew parent Children[] array with group_node only as new sole children
	parent.Children = []*DocoptAst{group_node}
	return group_node
}

func (n *DocoptAst) Detach_child(child_index int) *DocoptAst {
	detached := n.Children[child_index]
	n.Children = append(n.Children[:child_index], n.Children[child_index+1:]...)
	detached.Parent = nil
	return detached
}
