// Package that defines the abstract syntax tree and functions
// related to its manipulation and traversal.
package ast

import (
	"os"
	"strconv"
)

// Node struct, all members of the AST will be of this type.
type Node struct {
	Line          int
	Lexeme        string
	LocalRegister string
	Tag           string
	Semantic_Type string
	Semantic_Call string
	ArrDim        int
	Offset        string
	Curr_type     Production
	parent        *Node
	children      []*Node
}

// Append a child to the Node
func (n *Node) AddChild(child *Node) {
	if child != nil {
		child.addParent(n)
		n.children = append(n.children, child)
	}
}

// Sets the Parent of the node
func (n *Node) addParent(parent *Node) {
	n.parent = parent
}

// Return a new Node.
func New(node Production, lexeme string) *Node {
	return &Node{Curr_type: node, Lexeme: lexeme, children: make([]*Node, 0)}
}

func (n *Node) Get(index int) (bool, *Node) {
	if n.NumberOfChildren() <= index {
		return false, nil
	}

	return true, n.children[index]
}

func (n *Node) GetChild(prod Production) (bool, *Node) {
	if !n.HasChildren() {
		return false, nil
	}
	for _, child := range n.children {
		if child.Curr_type == prod {
			return true, child
		}
	}
	return false, nil
}

// Returns the number of children of the node.
func (n *Node) NumberOfChildren() int {
	if n.children == nil {
		return 0
	}

	return len(n.children)
}

// Checks if the node has a specific child and returns it.
func (n *Node) CheckHasChild(prod Production) (bool, *Node) {
	for _, child := range n.GetChildren() {
		if child.Curr_type == prod {
			return true, child
		}
	}

	return false, nil
}

// Returns if the current node has any children.
func (n *Node) HasChildren() bool {
	return n.NumberOfChildren() != 0
}

// Returns a Pointers to the array containing all the nodes children.
func (n *Node) GetChildren() []*Node {
	return n.children
}

// Returns pointer to the parent node.
func (n *Node) GetParent() *Node {
	return n.parent
}

// Must Use Preorder Traversal
func (n *Node) Traverse(depth int) {
	for i := 0; i < depth; i++ {
		print("   ")
	}
	print(TypesStrings[n.Curr_type])
	print(": " + n.Lexeme + "\n")

	if n.HasChildren() {
		for _, child := range n.children {
			child.Traverse(depth + 1)
		}
	}
}

// Traverses the AST and Prints the production to a file.
func (n *Node) TraverseToFile(depth int, file *os.File) {
	for i := 0; i < depth; i++ {
		file.WriteString("   ")
	}

	file.WriteString(TypesStrings[n.Curr_type])
	file.WriteString(": " + n.Lexeme + "\n")

	if n.HasChildren() {
		for _, child := range n.children {
			child.TraverseToFile(depth+1, file)
		}
	}
}

// Returns line number as string, useful for error messaging.
func (n *Node) GetLine() string {
	return strconv.Itoa(n.Line)
}

// Returns line number as string, useful for error messaging.
func (n *Node) GetLineNumber() int {
	return n.Line
}

// Migrates the tag to the parents, used to help with code translation.
func (n *Node) MigrateTagToParents(tag string) {
	n.Semantic_Call = tag

	if n.parent.Curr_type != IDENTIFIER {
		return
	}

	n.parent.MigrateTagToParents(tag)
}
