package records

import (
	"aHobeychi/GoCompiler/ast"
	"sort"
	"strconv"
)

// Size of the types in bits
const (
	INT_SIZE    = 4
	FLOAT_SIZE  = 8
	STRING_SIZE = 20
)

var typeAsSize = map[string]int{
	"integer": INT_SIZE,
	"float":   FLOAT_SIZE,
	"string":  STRING_SIZE,
}

// Symbol table struct, contains map of strings to rows.
type SymbolTable struct {
	Scope      string
	offset     int
	Rows       map[string]Row
	variables  []string
	classes    []string
	Label      string
	ScopeHead  *ast.Node
	Definition *ast.Node
}

// Returns new table with given scope.
func NewTable(scope string) *SymbolTable {
	return &SymbolTable{Scope: scope, Rows: make(map[string]Row)}
}

// Return row for the index.
func (st *SymbolTable) GetRow(row string) Row {
	return st.Rows[row]
}

// add variable to symbol table
func (st *SymbolTable) AddVariable(variable string) {
	st.variables = append(st.variables, variable)
}

// return list of Variables
func (st *SymbolTable) GetVariableList() *[]string {
	return &st.variables
}

// add class to symbol table
func (st *SymbolTable) addClass(class string) {
	st.classes = append(st.classes, class)
}

// Returns pointer to classes in the table.
func (st *SymbolTable) GetClassList() *[]string {
	return &st.classes
}

// Create New row and add it to the symbol table.
func (st *SymbolTable) CreateRow(s string, kind ast.Production) bool {
	switch kind {
	case ast.FUNCTION:
		if _, exists := st.Rows[s]; exists {
			print("Warning function " + s + " already exists in scope: " + st.Scope + "\n")
			return true
		}
		st.Rows[s] = NewFunctionRow(s, kind)
	case ast.VARIABLE:
		if _, exists := st.Rows[s]; exists {
			print("Variable " + s + " already exists in scope: " + st.Scope + "\n")
			return false
		}
		st.Rows[s] = NewVariableRow(s, kind)
		st.AddVariable(s)
	case ast.CLASS:
		if _, exists := st.Rows[s]; exists {
			print("Class " + s + " already exists in scope: " + st.Scope + "\n")
			return false
		}
		st.Rows[s] = NewClassRow(s, kind)
		st.addClass(s)
	}

	return true
}

// Adds Member Decl to the row
func (st *SymbolTable) AddMemberDecl(inNode *ast.Node) {
	if inNode.GetChildren()[1].Curr_type == ast.FUNCTION {
		row, index := st.AddFuncDecl(inNode)
		st.Rows[index] = row

	} else {
		row, index := st.AddVarDecl(inNode)

		if _, exists := st.Rows[index]; exists {
			print("Variable " + index + " already exists in scope: " + st.Scope + "\n")
			return
		}

		st.Rows[index] = row
		st.AddVariable(index)
	}
}

// Adds array information and type to variable row.
func (st *SymbolTable) AddArrayType(iden string, arraySize *ast.Node, var_type *ast.Node) {
	if st.GetRow(iden).GetKind() != ast.VARIABLE {
		return
	}

	row := st.GetRow(iden).(*VariableRow)
	row.VarType = var_type.Lexeme
	row.ArrayIndex = getArraySize(arraySize)
	row.ArrDim = arraySize.NumberOfChildren()
}

// return the array size information,
func getArraySize(arraySize *ast.Node) []int {
	if !arraySize.HasChildren() {
		return nil
	}

	var int_arr []int
	for i := 0; i < len(arraySize.GetChildren()); i++ {
		size, _ := strconv.Atoi(arraySize.GetChildren()[i].GetChildren()[0].Lexeme)
		int_arr = append(int_arr, size)
	}

	return int_arr
}

// return function decl information.
func (st *SymbolTable) AddFuncDecl(node *ast.Node) (*FunctionRow, string) {
	// Visibility of Member.
	visibility := node.GetChildren()[0].GetChildren()[0].Curr_type

	// Function Node.
	function := node.GetChildren()[1]

	// Id.
	id := function.GetChildren()[0].Lexeme
	fparamsNode := function.GetChildren()[1]
	return_type := function.GetChildren()[2].GetChildren()[0].Lexeme

	row := FunctionRow{Name: id, Kind: ast.FUNCTION}
	row.addVisibility(visibility)
	row.AddReturnType(return_type)

	var func_parameters Parameters

	for i := 0; i < len(fparamsNode.GetChildren()); i++ {

		variable := fparamsNode.GetChildren()[i]
		vType := variable.GetChildren()[0].Lexeme
		vIden := variable.GetChildren()[1].Lexeme
		array := variable.GetChildren()[2]
		arraydim := len(array.GetChildren())
		func_parameters = append(func_parameters, Variables{BaseType: vType, Iden: vIden, ArrDim: arraydim})
	}

	row.ListParameters = append(row.ListParameters, func_parameters)

	return &row, id
}

// return variable decl information.
func (st *SymbolTable) AddVarDecl(node *ast.Node) (*VariableRow, string) {
	visibility := node.GetChildren()[0].GetChildren()[0].Curr_type
	variable := node.GetChildren()[1]
	vType := variable.GetChildren()[0].Lexeme
	id := variable.GetChildren()[1].Lexeme
	_, arraySize := variable.GetChild(ast.ARRAYSIZE)

	row := VariableRow{Name: id, VarType: vType, Visibility: visibility, Kind: ast.VARIABLE}
	row.ArrayIndex = getArraySize(arraySize)
	row.ArrDim = arraySize.NumberOfChildren()

	return &row, id
}

// add return type to row.
func (st *SymbolTable) AddReturnType(index string, rNode *ast.Node) {
	_, return_node := rNode.Get(0)
	rType := return_node.Lexeme

	if st.GetRow(index).GetKind() == ast.FUNCTION {
		row := st.GetRow(index).(*FunctionRow)
		row.AddReturnType(rType)
	}
}

// add function parameters to row.
func (st *SymbolTable) AddFParams(index string, class, rNode *ast.Node) {
	var class_str string

	if st.GetRow(index).GetKind() != ast.FUNCTION {
		return
	}

	has_class, class_method := class.Get(0)

	if has_class {
		class_str = class_method.Lexeme
	} else {
		class_str = ""
	}

	row := st.GetRow(index).(*FunctionRow)
	func_param := formatFparams(rNode)

	if row.checkDuplicates(*func_param, class_str) {
		OutputMessage("Function already declared.\n")
	}

	row.ListParameters = append(row.ListParameters, *func_param)
}

// Format fparams nodes to parameter list.
func formatFparams(rNode *ast.Node) *Parameters {
	var func_param Parameters

	if rNode.NumberOfChildren() == 0 {

		func_param = append(func_param, Variables{})
		return &func_param
	}

	for i := 0; i < len(rNode.GetChildren()); i++ {

		variable := rNode.GetChildren()[i]
		vType := variable.GetChildren()[0].Lexeme
		vIden := variable.GetChildren()[1].Lexeme
		array := variable.GetChildren()[2]
		arraydim := array.NumberOfChildren()
		func_param = append(func_param, Variables{BaseType: vType, Iden: vIden, ArrDim: arraydim})
	}

	return &func_param
}

// Checks if a function has already been declared with the same parameter list.
func (row *FunctionRow) checkDuplicates(param Parameters, class string) bool {
	if len(row.ListParameters) == 0 {
		return false
	}

	list_parameters := row.ListParameters
	num_of_parameters := len(list_parameters)

	for i := 0; i < num_of_parameters; i++ {

		curr_param := list_parameters[i]

		if len(curr_param) != len(param) {
			continue
		}

		if row.ClassMethod[i] != class {
			continue
		}

		same := true

		for j := 0; j < len(curr_param); j++ {
			if curr_param[j].BaseType != param[j].BaseType {
				same = false
			}
		}

		if same {
			return true
		}
	}

	return false
}

// links a row to another table.
func (sym *SymbolTable) LinkRowToTable(index string, rType ast.Production, table *SymbolTable) {
	if table == nil {
		return
	}

	switch rType {
	case ast.CLASS:
		row := sym.GetRow(index).(*ClassRow)
		row.Link = table
	case ast.FUNCTION:
		row := sym.GetRow(index).(*FunctionRow)
		addParametersToTable(row, table)
		row.addLink(table)
	case ast.MAIN:
		row := sym.GetRow(index).(*FunctionRow)
		row.addLink(table)
	}
}

// Add Parameters to specific row.
func addParametersToTable(row *FunctionRow, table *SymbolTable) {
	if row.Name == "MAIN" {
		return
	}

	list_parameters := row.ListParameters
	num_overloaded := len(row.ListParameters)
	parameters := list_parameters[num_overloaded-1]

	num_parameters := len(parameters)
	if num_parameters != 0 {
		for i := 0; i < num_parameters; i++ {
			id := parameters[i].Iden
			if id == "" {
				continue
			}
			p_type := parameters[i].BaseType
			arraysize := parameters[i].ArrDim
			table.CreateRow(id, ast.VARIABLE)
			table.SetVisibility(id, ast.FUNCTION)
			table.SetVariableType(id, p_type)
			v_row := &VariableRow{Name: id, Kind: ast.VARIABLE, VarType: p_type, ArrDim: arraysize}
			table.Rows[id] = v_row
		}
	}
}

// Adds inherits to the ClassRow
func (st *SymbolTable) AddInherits(rowIndex string, inNode *ast.Node) {
	if inNode.HasChildren() {
		for _, child := range inNode.GetChildren() {
			row := st.GetRow(rowIndex).(*ClassRow)
			row.Line = inNode.GetLineNumber()
			row.AddInherits(child.Lexeme)
		}
	}
}

// sets visibility of the row.
func (st *SymbolTable) SetVisibility(row string, vis ast.Production) {
	switch st.GetRow(row).GetKind() {
	case ast.FUNCTION:
		row := st.GetRow(row).(*FunctionRow)
		row.addVisibility(vis)
	case ast.VARIABLE:
		row := st.GetRow(row).(*VariableRow)
		row.Visibility = vis
	case ast.CLASS:
		row := st.GetRow(row).(*ClassRow)
		row.Visibility = vis
	}
}

// sets type to row
func (st *SymbolTable) SetVariableType(row string, vType string) {
	if st.EntryExists(row) && st.GetRow(row).GetKind() == ast.VARIABLE {
		row := st.GetRow(row).(*VariableRow)
		row.VarType = vType
	}
}

// Returns true if entry is in symbol table.
func (st *SymbolTable) EntryExists(id string) bool {
	if _, exist := st.Rows[id]; exist {
		return true
	}
	return false
}

// Checks if variables is an entry in the table.
func (st *SymbolTable) ClassEntryExists(class, id string) bool {
	if !st.EntryExists(class) || st.GetRow(class).GetKind() != ast.CLASS {
		return false
	}

	class_table := st.GetRow(class).(*ClassRow).Link

	return class_table.EntryExists(id) && class_table.GetRow(id).GetKind() == ast.VARIABLE
}

// sorts inherited class list
func sortInherits(cr *ClassRow) {
	sort.Strings(cr.Inherits)
}

// Calculates the size in memory of the elemtents in the global symbol table.
func CalculateMemorySizes(table, global *SymbolTable) {
	variables := table.variables

	calculateVariableMemory(table, global, variables)

	// Has to be done twice, so it can calculate the size of the class first
	for k := range table.Rows {
		switch table.GetRow(k).GetKind() {
		case ast.CLASS:
			CalculateMemorySizes(table.GetRow(k).(*ClassRow).Link, global)
			calculateClassMemory(table.GetRow(k).(*ClassRow).Link, global)
		case ast.FUNCTION:
			for _, links := range table.GetRow(k).(*FunctionRow).Link {
				CalculateMemorySizes(links, global)
			}
		}
	}

	for k := range table.Rows {
		switch table.GetRow(k).GetKind() {
		case ast.CLASS:
			CalculateMemorySizes(table.GetRow(k).(*ClassRow).Link, global)
			calculateClassMemory(table.GetRow(k).(*ClassRow).Link, global)
		case ast.FUNCTION:
			for _, links := range table.GetRow(k).(*FunctionRow).Link {
				CalculateMemorySizes(links, global)
			}
		}
	}
}

// Calculates the size of variables in bytes.
func calculateVariableMemory(table, global *SymbolTable, variables []string) {
	if len(variables) == 0 {
		return
	}

	total := 0

	for _, varName := range variables {

		varRow := table.GetRow(varName).(*VariableRow)
		baseType := varRow.VarType
		base_size := GetTypeSize(baseType, global)
		multiple := getArrayTotal(varRow.ArrayIndex)

		varRow.Memory_Size = multiple * base_size
		varRow.Memory_Address = total
		total += varRow.Memory_Size
	}

	table.offset = total
}

// Calculates the size of the class by adding those of its variables.
func calculateClassMemory(table, global *SymbolTable) {
	variables := table.variables
	className := table.Scope

	class_row := global.GetRow(className).(*ClassRow)
	classMemory := 0

	for _, varName := range variables {
		varRow := table.GetRow(varName).(*VariableRow)
		classMemory += varRow.Memory_Size
	}

	class_row.Total_Memory = classMemory
}

// Calculates the requiered memory of a subclass by adding the size of its parent classes.
func calculateInheritedMemory(class string, global *SymbolTable) int {
	classRow := global.GetRow(class).(*ClassRow)
	inherits := classRow.Inherits
	total := classRow.Total_Memory

	for _, in_class := range inherits {
		total += calculateInheritedMemory(in_class, global)
	}

	return total
}

// Returns the size of the types.
func GetTypeSize(v_t string, global *SymbolTable) int {
	if size, base := typeAsSize[v_t]; base {
		return size
	}

	if global.EntryExists(v_t) {
		return global.GetRow(v_t).(*ClassRow).Total_Memory
	}

	return 0
}

// Get the total number of array elements.
func getArrayTotal(array []int) int {
	if len(array) == 0 {
		return 1
	}

	total := 1

	for _, index := range array {
		total *= index
	}

	return total
}

// Returns the vilibity of the function. with the param list.
func (table *SymbolTable) GetVisibilityFromParam(funcName string, params *Parameters) ast.Production {
	for _, rows := range table.Rows {

		if rows.GetKind() != ast.FUNCTION {
			continue
		}

		funcRow := rows.(*FunctionRow)

		if funcRow.Name != funcName {
			continue
		}

		for index, params := range funcRow.ListParameters {
			if params.MatchParameter(params) {
				return funcRow.Visibility[index]
			}
		}
	}

	return 0
}

// Calculates the offset in for variables in the symbol table
func (table *SymbolTable) CalculateOffset(varName string, gb *SymbolTable) {
	previous := table.variables[len(table.variables)-1]

	previousVar := table.GetRow(previous).(*VariableRow)

	new_offset := previousVar.Memory_Size + previousVar.Memory_Address
	variableRow := table.GetRow(varName).(*VariableRow)

	variableRow.Memory_Address = new_offset
	variableRow.Memory_Size = GetTypeSize(variableRow.VarType, gb)
}

// Return the offset of the variable.
func (table *SymbolTable) GetOffSet(varName string) int {
	variableRow := table.GetRow(varName).(*VariableRow)
	return variableRow.Memory_Address
}

// Returns the total memory in the table.
func (table *SymbolTable) GetMemoryTotal() int {
	total := 0

	for _, varName := range table.variables {

		varRow := table.GetRow(varName).(*VariableRow)
		total += varRow.Memory_Size
	}

	return total
}

// Adds row to the start of the table.
func (table *SymbolTable) AddRowToStart(label, v_Type string) {
	table.variables = append([]string{label}, table.variables...)
	table.Rows[label] = &VariableRow{Name: label, Kind: ast.VARIABLE, VarType: v_Type}
}

// Returns the link to the functionTable with a set of annotated parameters.
func (table *SymbolTable) AnnotatedParamsToLink(funcName string, aparams *ast.Node) *SymbolTable {
	funcRow := table.GetRow(funcName).(*FunctionRow)
	num_of_aparams := aparams.NumberOfChildren()

	params := funcRow.ListParameters

	for i, parameters := range params {

		if parameters[0].Iden == "" && num_of_aparams == 0 {
			return funcRow.Link[i]
		}

		if len(parameters) != num_of_aparams {
			continue
		}

		for index, aparam := range aparams.GetChildren() {
			if aparam.Semantic_Type == parameters[index].BaseType {
				return funcRow.Link[i]
			}
		}
	}

	return nil
}
