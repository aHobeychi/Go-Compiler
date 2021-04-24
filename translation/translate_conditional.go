package translation

import (
	"aHobeychi/GoCompiler/ast"
)

// Translates the relop expressions.
func translateRelOp(root *ast.Node) {

	_, leftNode := root.Get(0)
	_, rightNode := root.Get(1)

	translateExpr(leftNode)
	translateExpr(rightNode)

	command := relOpInstruction[root.Curr_type]

	leftRegister := leftNode.LocalRegister
	rightRegister := rightNode.LocalRegister

	_, register := registerPool.Pop()

	exec_code := register + ","
	exec_code += leftRegister + "," + rightRegister

	registerPool.Push(leftRegister)
	registerPool.Push(rightRegister)

	root.LocalRegister = register

	appendExec(command, exec_code)
}

// Translates ternary expression and sets the value into a register.
func translateTernary(root *ast.Node) {

	_, expr := root.Get(0)
	_, leftNode := root.Get(1)
	_, rightNode := root.Get(2)

	translateExpr(expr)
	translateExpr(leftNode)
	translateExpr(rightNode)

	offset := memory.getOffset(root.Tag)

	_, result := registerPool.Pop()

	end_label := generateLabel("e_t")
	rightTern := generateLabel("r_t")

	exec_code := expr.LocalRegister + "," + rightTern
	appendExec(branch, exec_code)

	exec_code = "-" + offset + "(r14)," + leftNode.LocalRegister
	appendExec(store, exec_code)

	exec_code = result + "," + r0 + "," + leftNode.LocalRegister
	appendExec(add, exec_code)
	appendExec(jump, end_label)

	exec_code = "-" + offset + "(r14)," + rightNode.LocalRegister
	appendExecScope(rightTern, store, exec_code)

	exec_code = result + "," + r0 + "," + rightNode.LocalRegister
	appendExec(add, exec_code)

	appendExecScope(end_label, "", "")

	registerPool.Push(expr.LocalRegister)
	registerPool.Push(leftNode.LocalRegister)
	registerPool.Push(rightNode.LocalRegister)

	root.LocalRegister = result
}

// Translate not expression and stores value into a register.
func translateNot(root *ast.Node) {

	_, expr := root.Get(0)

	not_label := generateLabel("not")
	end_not := generateLabel("endnot")

	translateExpr(expr)

	expr_register := expr.LocalRegister
	_, result := registerPool.Pop()

	exec_code := expr_register + "," + not_label
	appendExec(branch, exec_code)

	exec_code = result + "," + r0 + ",0"
	appendExec(addi, exec_code)

	offset := memory.getOffset(root.Tag)
	exec_code = "-" + offset + "(r14)," + result

	appendExec(store, exec_code)
	appendExec(jump, end_not)

	exec_code = result + "," + r0 + ", 1"
	appendExecScope(not_label, addi, exec_code)

	exec_code = "-" + offset + "(r14)," + result
	appendExec(store, exec_code)

	appendExecScope(end_not, "", "")

	registerPool.Push(expr_register)
	root.LocalRegister = result
}

// Translate and operator.
func translateAnd(root *ast.Node) {

	_, leftNode := root.Get(0)
	_, RightNode := root.Get(1)

	zero_label := generateLabel("zero")
	and_label := generateLabel("and")

	translateExpr(leftNode)
	translateExpr(RightNode)

	_, result := registerPool.Pop()

	exec_code := leftNode.LocalRegister + "," + zero_label
	appendExec(branch, exec_code)

	exec_code = RightNode.LocalRegister + "," + zero_label
	appendExec(branch, exec_code)

	exec_code = result + "," + r0 + ",1"
	appendExec(addi, exec_code)

	appendExec(jump, and_label)

	exec_code = result + "," + r0 + ",0"
	appendExecScope(zero_label, addi, exec_code)

	offset := memory.getOffset(root.Tag)

	exec_code = "-" + offset + "(r14)," + result
	appendExecScope(and_label, store, exec_code)

	registerPool.Push(leftNode.LocalRegister)
	registerPool.Push(RightNode.LocalRegister)

	root.LocalRegister = result
}

// Translates Or operator.
func translateOr(root *ast.Node) {

	_, leftNode := root.Get(0)
	_, RightNode := root.Get(1)

	not_z := generateLabel("not_z")
	or_label := generateLabel("or")

	translateExpr(leftNode)
	translateExpr(RightNode)

	_, result := registerPool.Pop()

	exec_code := leftNode.LocalRegister + "," + not_z
	appendExec(branch_nz, exec_code)

	exec_code = RightNode.LocalRegister + "," + not_z
	appendExec(branch_nz, exec_code)

	exec_code = result + "," + r0 + ",0"
	appendExec(addi, exec_code)

	appendExec(jump, or_label)

	exec_code = result + "," + r0 + ",1"
	appendExecScope(not_z, addi, exec_code)

	offset := memory.getOffset(root.Tag)

	exec_code = "-" + offset + "(r14)," + result
	appendExecScope(or_label, store, exec_code)

	registerPool.Push(leftNode.LocalRegister)
	registerPool.Push(RightNode.LocalRegister)

	root.LocalRegister = result
}
