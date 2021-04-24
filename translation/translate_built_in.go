package translation

import (
	"aHobeychi/GoCompiler/ast"
	"strconv"
)

// Translates If Stament and add the code to execution.
func translateIf(root *ast.Node) {

	elseLabel := generateLabel("else")
	endifLabel := generateLabel("endif")

	var exec_code string

	_, conditional := root.Get(0)

	translateExpr(conditional)
	register := conditional.LocalRegister

	exec_code = register + "," + elseLabel

	appendExec(branch, exec_code)

	registerPool.Push(register)

	_, firstBlock := root.Get(1)
	translateStatements(firstBlock)

	exec_code = endifLabel
	appendExec(jump, exec_code)

	has_else, elseBlock := root.Get(2)

	appendExecScope(elseLabel, "", "")

	if has_else {
		translateStatements(elseBlock)
	}

	appendExecScope(endifLabel, "", "")
}

// Translates while loop.
func translateWhile(root *ast.Node) {

	gowhile := generateLabel("while")
	endwhile := generateLabel("endwhile")

	searchForBreak(root, endwhile)
	searchForContinue(root, gowhile)

	appendExecScope(gowhile, "", "")

	_, conditional := root.Get(0)
	translateExpr(conditional)
	register := conditional.LocalRegister

	exec_code := register + "," + endwhile
	appendExec(branch, exec_code)

	_, firstBlock := root.Get(1)

	translateStatements(firstBlock)

	exec_code = gowhile
	appendExec(jump, exec_code)

	appendExecScope(endwhile, "", "")
}

// Translate Write Statement.
func translateWrite(root *ast.Node) {

	has_output = true

	_, node := root.Get(0)
	translateExpr(node)

	// Frame to top to use subroutine
	offset := memory.getStackTop()

	stored_address := offset + 8
	exec_code := "-" + strconv.Itoa(stored_address) + "(r14)," + node.LocalRegister
	appendExec(store, exec_code)

	registerPool.Push(node.LocalRegister)

	_, register := registerPool.Pop()
	exec_code = register + ",r0,buf"
	appendExec(addi, exec_code)

	buff_address := offset + 12
	exec_code = "-" + strconv.Itoa(buff_address) + "(r14)," + register
	appendExec(store, exec_code)

	exec_code = "r14,r14,-" + strconv.Itoa(offset)
	appendExec(addi, exec_code)

	// Jump To Function Stack Frame
	exec_code = "r15," + "intstr"
	appendExec(j_link, exec_code)

	exec_code = "-8(r14),r13"
	appendExec(store, exec_code)

	exec_code = "r15," + "putstr"
	appendExec(j_link, exec_code)

	// Decrement The Stack Frame
	exec_code = "r14,r14,-" + strconv.Itoa(offset)
	appendExec(subi, exec_code)

	registerPool.Push(register)

	// Prints empty line after command
	_, space_reg := registerPool.Pop()

	exec_code = space_reg + ",r0,10"
	appendExec(addi, exec_code)
	appendExec(out, space_reg)

	registerPool.Push(space_reg)
}

// Translates read statement.
func translateRead(root *ast.Node) {

	_, node := root.Get(0)
	_, register := registerPool.Pop()

	offset := getStoreCode(node)

	exec_code := register
	appendExec(in, exec_code)

	exec_code = offset + register
	appendExec(store, exec_code)

	registerPool.Push(register)
	registerPool.Push(node.LocalRegister)
}

// Translates return statement in function.
func translateReturn(root *ast.Node) {

	_, expr := root.Get(0)
	translateExpr(expr)

	exec_code := "-0(" + topOffset + ")," + expr.LocalRegister
	appendExec(store, exec_code)

	registerPool.Push(expr.LocalRegister)

	// Returns to previous function
	exec_code = "r15,-4(r14)"
	appendExec(load, exec_code)
	exec_code = "r15"
	appendExec(j_back, exec_code)

}

// Search for break statement and adds the label that it will jump to.
func searchForBreak(root *ast.Node, label string) {

	if root.Curr_type == ast.BREAK {
		root.Tag = label
	}

	for _, child := range root.GetChildren() {
		searchForBreak(child, label)
	}

}

// Search for continue statement and adds the label that it will jump to.
func searchForContinue(root *ast.Node, label string) {

	if root.Curr_type == ast.CONTINUE {
		root.Tag = label
	}

	for _, child := range root.GetChildren() {
		searchForContinue(child, label)
	}

}

// Used for break and continue statement, will jump to node tag.
func translateTagJump(root *ast.Node) {

	exec_code := root.Tag
	appendExec(jump, exec_code)
}
