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
	Usage_options_shortcut
	Usage_replaced
	Usage_old
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
	// This DocoptNodeType must be the last one
	Last_node_type
)

type DocoptAst struct {
	Type     DocoptNodeType
	Token    *lexer.Token
	Children []*DocoptAst
	Parent   *DocoptAst
	Repeat   bool
}

var DocoptNodes map[string]DocoptNodeType

// make a reverse map of nodes name to their DocoptNodeType
func DocoptNodes_init_reverse_map() {
	if len(DocoptNodes) > 0 {
		// already initialized
		return
	}
	DocoptNodes = make(map[string]DocoptNodeType)
	for t := Root; t < Last_node_type; t++ {
		DocoptNodes[t.String()] = t
	}
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
	// replace Children with a new slice without child_index
	n.Children = append(n.Children[:child_index], n.Children[child_index+1:]...)
	detached.Parent = nil
	return detached
}

func (n *DocoptAst) Detach_from_parent() bool {
	parent := n.Parent
	if parent == nil {
		return false
	}
	found := false
	for i, c := range parent.Children {
		if c == n {
			parent.Detach_child(i)
			found = true
			break
		}
	}
	return found
}

// Find_recursive_by_Token() Search recursively in the AST from `start_node` for a node
// than match `t` Type and Value.
// Returns: position int parent's Children index, node pointer found, and a bool for the search success
func (start_node *DocoptAst) Find_recursive_by_Token(t *lexer.Token, parent_index int) (int, *DocoptAst, bool) {
	if start_node.Token != nil && start_node.Token.Type == t.Type && start_node.Token.Value == t.Value {
		return parent_index, start_node, true
	}

	for i, n := range start_node.Children {
		if pos, n2, ok := n.Find_recursive_by_Token(t, i); ok {
			return pos, n2, true
		}
	}
	return -1, nil, false
}

func (parent *DocoptAst) AppendNode(n *DocoptAst, override_parrent *DocoptAst) {
	if override_parrent != nil {
		n.Parent = override_parrent
	}
	parent.Children = append(parent.Children, n)
}

func (n *DocoptAst) Deep_copy() *DocoptAst {
	new_node := *n
	new_node.Children = []*DocoptAst{}
	var child_copy *DocoptAst
	for _, c := range n.Children {
		if len(c.Children) > 0 {
			child_copy = c.Deep_copy()
		} else {
			// copy node without Children direcly
			node_copy := *c
			child_copy = &node_copy
		}
		(&new_node).AppendNode(child_copy, &new_node)
	}
	return &new_node
}

func (n *DocoptAst) Deep_copy_exclude(exclude *[]DocoptNodeType) *DocoptAst {
	if exclude == nil {
		return n.Deep_copy()
	}

	for _, x := range *exclude {
		if n.Type == x {
			return nil
		}
	}

	new_node := *n
	new_node.Children = []*DocoptAst{}
	var child_copy *DocoptAst
	for _, c := range n.Children {
		child_copy = c.Deep_copy_exclude(exclude)
		if child_copy != nil {
			(&new_node).AppendNode(child_copy, &new_node)
		}
	}
	return &new_node
}

func (n *DocoptAst) Has_Parent(t DocoptNodeType) bool {
	if n.Parent != nil {
		if n.Parent.Type == t {
			return true
		} else {
			return n.Parent.Has_Parent(t)
		}
	}
	return false
}
