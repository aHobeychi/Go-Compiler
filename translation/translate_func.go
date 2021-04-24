package translation

import (
	"aHobeychi/GoCompiler/ast"
	"aHobeychi/GoCompiler/records"
	"strconv"
)

func translateFreeId(call *ast.Node) {

	is_func, _ := call.GetChild(ast.APARAMS)

	if is_func {
		translateFreeFunc(call)
	}

}

// Translate Jump to free function.
func translateFreeFunc(call *ast.Node) {

	// Loads Parameters to the proper memory Address
	funcName := call.Lexeme
	_, aparams := call.GetChild(ast.APARAMS) // These aparams have been anotated already
	func_table := gb_table.AnnotatedParamsToLink(funcName, aparams)
	label := func_table.Label

	migrateParameters(call, func_table)

	// Change Add the called function to the stack.
	offset := memory.getStackTop()
	memory.addScopeToStack(funcName, func_table)

	// Code to jump to function scope.
	exec_code := "r14,r14,-" + strconv.Itoa(offset)
	appendExec(addi, exec_code)

	// Jump To Function Stack Frame
	exec_code = "r15," + label
	appendExec(j_link, exec_code)

	// Decrement The Stack Frame
	exec_code = "r14,r14,-" + strconv.Itoa(offset)
	appendExec(subi, exec_code)

	// load the returned value into a register
	_, register := registerPool.Pop()
	exec_code = register + ",-" + strconv.Itoa(offset) + "(r14)"
	appendExec(load, exec_code)

	memory.moveBackOneFrame()

	// store the value of the register into memory
	final_offset := memory.getOffset(call.Tag)
	exec_code = "-" + final_offset + "(r14)," + register
	appendExec(store, exec_code)

	call.Offset = final_offset

	registerPool.Push(register)
}

// Translate all funcitons.
func translateFunctionBlocks(root *ast.Node) {

	for _, function := range root.GetChildren() {
		translateFunctionBlock(function)
	}

}

// Translate function block and gives it a unique label.
func translateFunctionBlock(function *ast.Node) {

	_, head := function.GetChild(ast.FUNCTIONHEAD)
	_, id := head.GetChild(ast.IDENTIFIER)
	funcName := id.Lexeme

	funcLabel := generateLabel(funcName)
	// table := gb_table.GetRow(funcName).(*records.FunctionRow).Link[0] // temp
	table := getTableFromNode(function)
	table.Label = funcLabel

	memory.addReturnAndAddress(table) //
	records.CalculateMemorySizes(table, gb_table)

	memory.addScopeToStack(funcName, table)

	appendExecScope(funcLabel, "", "")

	exec_code := "-4(r14),r15"
	appendExec(store, exec_code)
	translateFunction(function)
	exec_code = "r15,-4(r14)"
	appendExec(load, exec_code)
	exec_code = "r15"
	appendExec(j_back, exec_code)
}

// Sends the parameters to the correct memory address.
func migrateParameters(funcCall *ast.Node, table *records.SymbolTable) {

	definition := table.ScopeHead
	_, aparams := funcCall.GetChild(ast.APARAMS)

	current_offset := memory.getStackTop()

	_, fparams := definition.GetChild(ast.FPARAMS)

	for index, parameter := range fparams.GetChildren() {

		_, iden := parameter.GetChild(ast.IDENTIFIER)
		varName := iden.Lexeme
		address := table.GetRow(varName).(*records.VariableRow).Memory_Address

		aparam := aparams.GetChildren()[index]

		translateExpr(aparam)

		to_save := address + current_offset
		exec_code := "-" + strconv.Itoa(to_save) + "(" + topOffset + ")," + aparam.LocalRegister
		appendExec(store, exec_code)
		registerPool.Push(aparam.LocalRegister)

	}

}

// Return the correct table for the function node.
func getTableFromNode(function *ast.Node) *records.SymbolTable {

	_, head := function.GetChild(ast.FUNCTIONHEAD)
	_, id := head.GetChild(ast.IDENTIFIER)
	funcName := id.Lexeme

	allTables := gb_table.GetRow(funcName).(*records.FunctionRow)

	for _, table := range allTables.Link {

		if table.Definition == function {
			return table
		}
	}

	return nil
}
