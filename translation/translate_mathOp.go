package translation

import "aHobeychi/GoCompiler/ast"

// Translates assign expression.
func translateAssign(root *ast.Node) {

	_, leftNode := root.Get(0)
	_, rightNode := root.Get(1)

	translateExpr(rightNode)
	offset := getStoreCode(leftNode)

	exec_code := offset + rightNode.LocalRegister
	appendExec(store, exec_code)

	registerPool.Push(leftNode.LocalRegister)
	registerPool.Push(rightNode.LocalRegister)
}

// Translates arithmic expression, standard way of doing it.
func translateArithOp(root *ast.Node) {

	_, rightNode := root.Get(1)

	if rightNode.Curr_type == ast.INT_VALUE {
		translateSimpliedArith(root)
		return
	}

	_, leftNode := root.Get(0)

	translateExpr(leftNode)
	translateExpr(rightNode)

	command := arithOpInstruction[root.Curr_type]

	_, register := registerPool.Pop()

	exec_code := register + "," + leftNode.LocalRegister + "," + rightNode.LocalRegister
	appendExec(command, exec_code)

	root.LocalRegister = register

	registerPool.Push(leftNode.LocalRegister)
	registerPool.Push(rightNode.LocalRegister)
}

// Translates arithmic expression using integer instructions, saves one instruction.
func translateSimpliedArith(root *ast.Node) {

	_, leftNode := root.Get(0)
	_, rightNode := root.Get(1)

	translateExpr(leftNode)

	command := simplifiedArith[root.Curr_type]

	_, register := registerPool.Pop()
	int_value := rightNode.Lexeme

	exec_code := register + "," + leftNode.LocalRegister + "," + int_value
	appendExec(command, exec_code)

	root.LocalRegister = register
	registerPool.Push(leftNode.LocalRegister)
}

// Translates int and puts it into register.
func translateInt(root *ast.Node) {

	_, int_reg := registerPool.Pop()
	exec_code := int_reg + ",r0," + root.Lexeme

	appendExec(addi, exec_code)

	root.LocalRegister = int_reg
}

// Translates Negative expression.
func translateNegative(root *ast.Node) {

	_, expr := root.Get(0)
	translateExpr(expr)

	exec_code := expr.LocalRegister + "," + r0 + "," + expr.LocalRegister
	appendExec(sub, exec_code)
	root.LocalRegister = expr.LocalRegister
}

// Translates Positive expression statement.
func translatePositve(root *ast.Node) {
	_, expr := root.Get(0)
	translateExpr(expr)
	root.LocalRegister = expr.LocalRegister
}
