package translation

import (
	"aHobeychi/GoCompiler/ast"
	"aHobeychi/GoCompiler/records"
	"strconv"
)

const (
	STACK_LIMIT = 10000
)

type memory_stack struct {
	label_offset map[string]int
	scope_offset map[int]*records.SymbolTable
	stack_head   int
	labels       []string
	mem_total    int
}

// Add the current scope to the stack
func (mem *memory_stack) addScopeToStack(label string, table *records.SymbolTable) {

	label = generateLabel(label)
	mem.addLabel(label)
	mem.label_offset[label] = mem.mem_total
	mem.scope_offset[mem.mem_total] = table
	mem.stack_head = mem.mem_total
	mem.mem_total += table.GetMemoryTotal()
}

func (mem *memory_stack) addInterMediate(label, v_type string) {

	table := mem.scope_offset[mem.stack_head]

	table.CreateRow(label, ast.VARIABLE)
	table.SetVariableType(label, v_type)
	table.CalculateOffset(label, gb_table)
	table.SetVisibility(label, ast.FUNCTION)

	added_memory := records.GetTypeSize(v_type, gb_table)
	mem.mem_total += added_memory

}

func (mem *memory_stack) getTableAtHead() *records.SymbolTable {
	return mem.scope_offset[mem.stack_head]
}

func (mem *memory_stack) getRowFromScope(varName string) *records.VariableRow {

	table := mem.scope_offset[mem.stack_head]
	varRow := table.GetRow(varName).(*records.VariableRow)

	return varRow
}

func (mem *memory_stack) getOffset(varName string) string {

	table := mem.scope_offset[mem.stack_head]
	offset := table.GetOffSet(varName)

	return strconv.Itoa(offset)
}

// Adds Label to the stack
func (mem *memory_stack) addLabel(label string) {
	mem.labels = append(mem.labels, label)
}

func (mem *memory_stack) flushMemory(offset int) {
	mem.stack_head = offset
}

func (mem *memory_stack) addReturnAndAddress(table *records.SymbolTable) {
	table.AddRowToStart("address", "integer")
	table.AddRowToStart("return", "integer")
}

func (mem *memory_stack) getStackTop() int {

	table := mem.scope_offset[mem.stack_head]
	top := table.GetMemoryTotal()

	return top
}

func (mem *memory_stack) moveBackOneFrame() {
	mem.labels = mem.labels[:len(mem.labels)-1]
	mem.stack_head = mem.label_offset[mem.labels[len(mem.labels)-1]]
}

func (mem *memory_stack) getRelativeOffset(previous, current *ast.Node) string {

	previous_table := gb_table.GetRow(previous.Semantic_Type).(*records.ClassRow).Link
	offset := previous_table.GetOffSet(current.Lexeme)

	return strconv.Itoa(offset)

}
