package translation

import (
	"aHobeychi/GoCompiler/ast"
)

// Translate Nodes that can be contained within an expression.
func translateExpr(root *ast.Node) {

	switch root.Curr_type {

	case ast.ASSIGN:
		translateAssign(root)
	case ast.PLUS, ast.MINUS, ast.DIV, ast.MUL:
		translateArithOp(root)
	case ast.NEGATIVE:
		translateNegative(root)
	case ast.POSITIVE:
		translatePositve(root)
	case ast.EQ, ast.NOTEQ, ast.LT, ast.LTE, ast.GT, ast.GTE:
		translateRelOp(root)
	case ast.NOT:
		translateNot(root)
	case ast.AND:
		translateAnd(root)
	case ast.OR:
		translateOr(root)
	case ast.TERNARY:
		translateTernary(root)
	case ast.IDENTIFIER:
		translateId(root)
	case ast.INT_VALUE:
		translateInt(root)
	}

}

// Translates function in expression and places the return value into a register.
func translateFuncInExpr(root *ast.Node) {

	translateFreeFunc(root)

	_, register := registerPool.Pop()

	exec_code := register + ",-" + root.Offset + "(r14)"
	appendExec(load, exec_code)

	root.LocalRegister = register
}

// Translates nested id.
func translateNestedId(root *ast.Node) {

	off_register := getNestedOffset(root)

	_, final_register := registerPool.Pop()
	exec_code := final_register + ",-0(" + off_register + ")"
	appendExec(load, exec_code)

	registerPool.Push(off_register)
	root.LocalRegister = final_register
}

// Returns the memory address of a nested identifier.
func getNestedOffset(root *ast.Node) string {

	_, off_register := registerPool.Pop()
	offset := memory.getOffset(root.Lexeme)

	exec_code := off_register + ",r14," + offset
	appendExec(subi, exec_code)

	if is_arr, _ := root.GetChild(ast.INDICEREP); is_arr {

		arr_Reg := getRelativeArrayOffset(root, root)

		exec_code := off_register + "," + off_register + "," + arr_Reg
		appendExec(sub, exec_code)
		registerPool.Push(arr_Reg)
	}

	previous := root
	current := root

	for {

		if has_child, child := current.GetChild(ast.IDENTIFIER); has_child {

			previous = current
			current = child

		} else {
			break
		}

		offset = memory.getRelativeOffset(previous, current)

		exec_code := off_register + "," + off_register + "," + offset
		appendExec(subi, exec_code)

		if is_arr, _ := current.GetChild(ast.INDICEREP); is_arr {

			arr_Reg := getRelativeArrayOffset(previous, current)

			exec_code := off_register + "," + off_register + "," + arr_Reg
			appendExec(sub, exec_code)
			registerPool.Push(arr_Reg)
		}

	}

	return off_register
}

// Translates identifier and places the value into a register.
func translateId(root *ast.Node) {

	if is_func, _ := root.GetChild(ast.APARAMS); is_func {
		translateFuncInExpr(root)
		return
	}

	translateNestedId(root)
}
