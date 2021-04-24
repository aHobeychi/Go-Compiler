package translation

import (
	"aHobeychi/GoCompiler/ast"
	"aHobeychi/GoCompiler/records"
	"aHobeychi/GoCompiler/utilities"
	"strconv"
)

// Variables for the package. Easier than pushing them through functions.
var gb_table *records.SymbolTable
var tree_root *ast.Node
var memory *memory_stack
var registerPool *utilities.Stack
var topOffset string
var has_output bool

// begins the translation of the program.
func Translate(root *ast.Node, global_table *records.SymbolTable) {

	has_output = false
	gb_table = global_table
	tree_root = root

	memory = &memory_stack{label_offset: make(map[string]int),
		scope_offset: make(map[int]*records.SymbolTable),
		stack_head:   0}

	memory.addScopeToStack("MAIN", gb_table.GetRow("MAIN").(*records.FunctionRow).Link[0])

	registerPool = getNewRegisterPool()

	registerPool.Pop()                // r0 register will contain the value 0
	_, topOffset = registerPool.Pop() // r14 will contain the top address

	for _, child := range root.GetChildren() {

		if child.Curr_type == ast.MAIN {

			appendExecScope("entry", "", "")
			resetTopAddress() // Sets the top address
			translateFunction(child)
		}

		if child.Curr_type == ast.FUNCDEF {
			translateFunctionBlocks(child)
		}
	}

	appendEmptySpace()
	appendExecScope("hlt", "", "")

	if has_output {
		appendEmptySpace()
		appendExecScope("buf res 20", "", "")
		appendEmptySpace()
	}

	OutputMessage(exec_code)

}

// Translate the statements within a function.
func translateFunction(body *ast.Node) {

	has_body, func_body := body.GetChild(ast.FUNCTIONBODY)

	if !has_body {
		return
	}

	has_stats, stats := func_body.GetChild(ast.STATEMENT)

	if !has_stats {
		return
	}

	translateStatements(stats)

}

// Translate Statement Block
func translateStatements(root *ast.Node) {

	for _, child := range root.GetChildren() {
		translateStatement(child)
	}

}

// Translates statements - function body.
func translateStatement(root *ast.Node) {

	switch root.Curr_type {

	case ast.ASSIGN:
		translateAssign(root)
	case ast.WHILE:
		translateWhile(root)
	case ast.IF:
		translateIf(root)
	case ast.WRITE:
		translateWrite(root)
	case ast.READ:
		translateRead(root)
	case ast.IDENTIFIER:
		translateFreeId(root)
	case ast.RETURN:
		translateReturn(root)
	case ast.BREAK, ast.CONTINUE:
		translateTagJump(root)
	}

}

// returns the offset where the assignment will be stored.
func getStoreCode(root *ast.Node) string {

	off_register := getNestedOffset(root)
	exec_code := "-0(" + off_register + "),"

	root.LocalRegister = off_register
	return exec_code
}

// return string of the relative offset from the array index.
func getRelativeArrayOffset(previous, current *ast.Node) string {

	var varRow *records.VariableRow

	if gb_table.EntryExists(previous.Semantic_Type) && previous != current {

		parent_table := gb_table.GetRow(previous.Semantic_Type).(*records.ClassRow).Link
		varRow = parent_table.GetRow(current.Lexeme).(*records.VariableRow)

	} else {
		varRow = memory.getRowFromScope(current.Lexeme)
	}

	_, index := current.GetChild(ast.INDICEREP)
	arr_index := varRow.ArrayIndex
	vType := varRow.VarType

	_, register := registerPool.Pop()

	exec_code := register + ",r0,0"
	appendExec(addi, exec_code)

	for index, indices := range index.GetChildren() {

		multiple := calculateMultiple(index, arr_index, vType)

		_, index_register := registerPool.Pop()

		translateExpr(indices)

		exec_code = index_register + "," + indices.LocalRegister + "," + strconv.Itoa(multiple)

		appendExec(muli, exec_code)

		exec_code = register + "," + register + "," + index_register
		appendExec(add, exec_code)

		registerPool.Push(index_register)
		registerPool.Push(indices.LocalRegister)
	}

	return register
}

// Calculates the offset multiple for a given array index.
func calculateMultiple(index int, array []int, vType string) int {

	multiple := records.GetTypeSize(vType, gb_table)

	if index >= len(array) {
		return multiple
	}

	remainder_arr := 1
	for i := index + 1; i < len(array); i++ {
		remainder_arr *= array[i]
	}

	return remainder_arr * multiple
}

// resets the top address and empty the function call stack.
func resetTopAddress() {

	appendExec(addi, topOffset+","+r0+","+top)

	memory = &memory_stack{label_offset: make(map[string]int),
		scope_offset: make(map[int]*records.SymbolTable),
		stack_head:   0}

	memory.addScopeToStack("MAIN", gb_table.GetRow("MAIN").(*records.FunctionRow).Link[0])
}
