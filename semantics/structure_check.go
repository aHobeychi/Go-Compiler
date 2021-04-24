package semantics

import "aHobeychi/GoCompiler/ast"

// Checks if break and continue statements are within a while loop.
func searchForWhile(root *ast.Node) bool {

	parent := root.GetParent()

	if parent == nil {
		return false
	}

	if parent.Curr_type == ast.WHILE {
		return true
	}

	return searchForWhile(parent)
}
