// File that contains all the pre-translation type checking.
package semantics

import (
	"aHobeychi/GoCompiler/ast"
	"aHobeychi/GoCompiler/records"
	"strconv"
)

var type_correct bool

// Runs the type checking.
func type_check(root *ast.Node, scope string) bool {

	type_correct = true

	for _, child := range root.GetChildren() {

		if child.Curr_type == ast.MAIN {

			check_function_types("MAIN", "", "", nil, child)

		} else if child.Curr_type == ast.FUNCDEF {

			if !child.HasChildren() {
				continue
			}

			for _, function := range child.GetChildren() {

				var classMethod string

				_, funcHead := function.GetChild(ast.FUNCTIONHEAD)
				_, funcName := funcHead.GetChild(ast.IDENTIFIER)
				_, funcClass := funcHead.GetChild(ast.CLASSMETHOD)

				if funcClass.HasChildren() {

					classMethod = funcClass.GetChildren()[0].Lexeme

				}

				_, fparams := funcHead.GetChild(ast.FPARAMS)
				_, funcBody := function.GetChild(ast.FUNCTIONBODY)

				check_function_types(funcName.Lexeme, classMethod, "", fparams, funcBody)
			}
		}
	}

	return type_correct
}

// checks the function body.
func check_function_types(funcName, class, p_t string, params, root *ast.Node) {

	func_table := gb_table.GetRow(funcName).(*records.FunctionRow)

	if funcName == "MAIN" {
		check_body_types(func_table.Link[0], func_table.Link[0], "", p_t, root)
		return
	}

	location := func_table.GetLinkForClass(class, params)

	check_body_types(location, location, class, p_t, root)
}

// Runs the checks on the inside of the body types.
func check_body_types(location, scope *records.SymbolTable, classMethod, p_t string, root *ast.Node) {

	cType := root.Curr_type

	if root.Curr_type == ast.VARIABLEDECL {
		return
	}

	if cType == ast.ASSIGN {
		check_assign(location, scope, classMethod, root)
		return
	}

	if cType == ast.IDENTIFIER {

		varType := get_expression_type(location, scope, classMethod, p_t, root)
		root.Semantic_Type = varType.BaseType

		if varType.BaseType == "invalid" {
			type_correct = false
		}

		return
	}

	if cType == ast.WRITE {

		varType := get_expression_type(location, scope, classMethod, p_t, root.GetChildren()[0])

		if varType.BaseType == "invalid" {
			type_correct = false
		}

		return
	}

	if cType == ast.IDENTIFIER && root.GetParent().Curr_type == ast.STATEMENT {

		function_type(location, scope, classMethod, p_t, root)

		return
	}

	if cType == ast.LT || cType == ast.LTE || cType == ast.GT ||
		cType == ast.GTE || cType == ast.EQ || cType == ast.NOTEQ {

		check_rel_op(location, scope, classMethod, root)
		return
	}

	if cType == ast.BREAK || cType == ast.CONTINUE {

		if !searchForWhile(root) {

			outputMessage("Cannot have a " + ast.TypesStrings[cType] + " statement outside a while loop. Line: " + root.GetLine() + ".\n")
			type_correct = false
		}
	}

	if root.HasChildren() {

		for _, child := range root.GetChildren() {
			check_body_types(location, scope, classMethod, "", child)
		}
	}

}

// Checks if the two sides of the rel-op match.
func check_rel_op(location, scope *records.SymbolTable, classMethod string, node *ast.Node) *records.Variables {

	left := node.GetChildren()[0]
	right := node.GetChildren()[1]

	lType := get_expression_type(location, scope, classMethod, "", left)
	rType := get_expression_type(location, scope, classMethod, "", right)

	if (lType.BaseType == "float" || lType.BaseType == "integer") &&
		(rType.BaseType == "float" || rType.BaseType == "integer") && !lType.CompareVar(rType) {

		outputMessage("Warning comparing float to integer at line: " + left.GetLine() + ".\n")
		type_correct = false
		return records.NewVariable("invalid")
	}

	if !lType.CompareVar(rType) || lType.BaseType == "invalid" || rType.BaseType == "invalid" {

		outputMessage("RelOp Type Mismatch! In func: " + scope.Scope + ". Line: " + node.GetLine() + ".\n")
		outputMessage("\tleft: " + left.Lexeme + " ➤ " + lType.ToString() + " \n")
		outputMessage("\tright: " + right.Lexeme + " ➤ " + rType.ToString() + " \n")

		type_correct = false

		return records.NewVariable("invalid")
	}

	left.Semantic_Type = lType.BaseType
	right.Semantic_Type = rType.BaseType

	// Helps Generate the intermediate offset that will be used for later processing.
	label := generateTag(ast.TypesStrings[node.Curr_type])

	scope.CreateRow(label, ast.VARIABLE)
	scope.SetVariableType(label, lType.BaseType)
	scope.CalculateOffset(label, gb_table)
	scope.SetVisibility(label, ast.FUNCTION)

	node.Semantic_Type = lType.BaseType
	node.Tag = label

	return lType
}

// Checks to see if the two types of the assignment match.
func check_assign(location, scope *records.SymbolTable, classMethod string, node *ast.Node) {

	left := node.GetChildren()[0]
	right := node.GetChildren()[1]

	lType := get_expression_type(location, scope, classMethod, "", left)
	rType := get_expression_type(location, scope, classMethod, "", right)

	if !lType.CompareVar(rType) {

		outputMessage("Assign Type Mismatch! In func: " + scope.Scope)

		if classMethod != "" {
			outputMessage("::" + classMethod)
		}

		outputMessage(". Line: " + left.GetLine() + ".\n")
		outputMessage("\tleft: " + left.Lexeme + " ➤ " + lType.ToString() + " \n")
		outputMessage("\tright: " + right.Lexeme + " ➤ " + rType.ToString() + " \n")
		type_correct = false
		return
	}

	node.Semantic_Type = lType.BaseType
}

// Handle the type of arith expression
func handle_arith_type(location, scope *records.SymbolTable, classMethod, p_t string, node *ast.Node) *records.Variables {

	left := node.GetChildren()[0]
	right := node.GetChildren()[1]

	lType := get_expression_type(location, scope, classMethod, p_t, left)
	rType := get_expression_type(location, scope, classMethod, p_t, right)

	if lType.CompareVar(rType) {

		// Helps Generate the intermediate offset that will be used for latter processig.
		label := generateTag(ast.TypesStrings[node.Curr_type])

		scope.CreateRow(label, ast.VARIABLE)
		scope.SetVariableType(label, lType.BaseType)
		scope.CalculateOffset(label, gb_table)
		scope.SetVisibility(label, ast.FUNCTION)

		node.Tag = label
		node.Semantic_Type = lType.BaseType

		return lType
	} else {
		return records.NewVariable("invalid")
	}
}

// Returns the expressiong type
func get_expression_type(location, scope *records.SymbolTable, classMethod, p_t string, node *ast.Node) *records.Variables {

	switch node.Curr_type {

	case ast.INT_VALUE:
		node.Semantic_Type = "integer"
		return records.NewVariable("integer")
	case ast.FLOAT_VALUE:
		node.Semantic_Type = "float"
		return records.NewVariable("float")
	case ast.STRINGLIT:
		node.Semantic_Type = "string"
		return records.NewVariable("string")
	case ast.IDENTIFIER:
		return handle_id_type(location, scope, classMethod, p_t, node)
	case ast.PLUS, ast.MINUS, ast.MUL, ast.DIV:
		return handle_arith_type(location, scope, classMethod, p_t, node)
	case ast.POSITIVE, ast.NEGATIVE:
		return handle_signed_value(location, scope, classMethod, p_t, node)
	case ast.TERNARY:
		return handle_ternary_expr(location, scope, classMethod, p_t, node)
	case ast.NOT:
		return handle_not(location, scope, classMethod, p_t, node)
	case ast.EQ, ast.LT, ast.LTE, ast.GT, ast.GTE, ast.AND, ast.OR:
		return check_rel_op(location, scope, classMethod, node)
	}

	return records.NewVariable("invalid")
}

// Handles ternary expression.
func handle_ternary_expr(location, scope *records.SymbolTable, classMethod, p_t string, node *ast.Node) *records.Variables {

	first := node.GetChildren()[0]
	middle := node.GetChildren()[1]
	last := node.GetChildren()[2]

	firstT := get_expression_type(location, scope, classMethod, p_t, first)
	lType := get_expression_type(location, scope, classMethod, p_t, middle)
	rType := get_expression_type(location, scope, classMethod, p_t, last)

	if !lType.CompareVar(rType) || firstT.CompareVar(records.NewVariable("invalid")) {

		outputMessage("Cannot have inconsistent types in ternary operation, ")
		outputMessage(lType.ToString() + " != ")
		outputMessage(rType.ToString() + "")
		outputMessage(". Line: " + node.GetLine() + "\n")

		return records.NewVariable("invalid")

	}

	label := generateTag("TERNARY")
	scope.CreateRow(label, ast.VARIABLE)
	scope.SetVariableType(label, lType.BaseType)
	scope.CalculateOffset(label, gb_table)
	scope.SetVisibility(label, ast.FUNCTION)

	node.Tag = label

	return rType
}

// Handles not statement.
func handle_not(location, scope *records.SymbolTable, classMethod, p_t string, node *ast.Node) *records.Variables {

	left := node.GetChildren()[0]
	lType := get_expression_type(location, scope, classMethod, p_t, left)

	label := generateTag("NOT")

	scope.CreateRow(label, ast.VARIABLE)
	scope.SetVariableType(label, lType.BaseType)
	scope.CalculateOffset(label, gb_table)
	scope.SetVisibility(label, ast.FUNCTION)

	node.Tag = label

	return lType
}

// Handle the type of arith expression
func handle_signed_value(location, scope *records.SymbolTable, classMethod, p_t string, node *ast.Node) *records.Variables {

	left := node.GetChildren()[0]

	lType := get_expression_type(location, scope, classMethod, p_t, left)

	left.Semantic_Type = lType.BaseType

	if lType.CompareVar(records.NewVariable("integer")) {
		return lType
	} else if lType.CompareVar(records.NewVariable("float")) {
		return lType
	} else {
		return records.NewVariable("invalid")
	}
}

// Handles the type of the identifier.
func handle_id_type(location, scope *records.SymbolTable, classMethod, p_t string, node *ast.Node) *records.Variables {

	if exists, _ := node.GetChild(ast.APARAMS); exists {
		return function_type(location, scope, classMethod, p_t, node)
	} else {
		return identifier_type(location, scope, classMethod, p_t, node)
	}
}

// Returns the type of the identifier.
func identifier_type(location, class *records.SymbolTable, classMethod, p_t string, node *ast.Node) *records.Variables {

	identifierName := node.Lexeme
	var varType string
	var varRow *records.VariableRow

	if location.EntryExists(identifierName) {

		row := location.GetRow(identifierName)
		varRow = row.(*records.VariableRow)
		varType = varRow.GetType()

	} else if class.EntryExists(identifierName) {

		row := class.GetRow(identifierName)
		varRow = row.(*records.VariableRow)
		varType = varRow.GetType()

	} else if gb_table.EntryExists(classMethod) && gb_table.InheritsSearch(classMethod, identifierName) {

		varRow = gb_table.GetInheritedVar(classMethod, identifierName)
		varType = varRow.GetType()

	} else if gb_table.EntryExists(p_t) && gb_table.InheritsSearch(p_t, identifierName) {

		varRow = gb_table.GetInheritedVar(p_t, identifierName)
		varType = varRow.GetType()

	} else {

		outputMessage("Variable " + identifierName + " not in current scope: " + class.Scope)

		if classMethod != "" {
			outputMessage(" with class method: " + classMethod)
		}

		outputMessage(", Line: " + node.GetLine() + ". \n")
		return records.NewVariable("invalid")
	}

	if !(varType == "integer" || varType == "float" || varType == "string") {
		if !gb_table.EntryExists(varType) {

			outputMessage("Variable type " + varType + " has not been defined in file")
			outputMessage(", line: " + node.GetLine() + ". \n")
			return records.NewVariable("invalid")
		}
	} else {

		has_child, _ := node.GetChild(ast.IDENTIFIER)

		if has_child {
			outputMessage("Variable type " + varType + ", Cannot have subType.")
			outputMessage(", line: " + node.GetLine() + ". \n")
			return records.NewVariable("invalid")
		}

	}

	has_child, child := node.GetChild(ast.IDENTIFIER)

	node.Semantic_Type = varType

	if has_child {

		new_class_table := gb_table.GetRow(varType).(*records.ClassRow).Link
		return handle_id_type(location, new_class_table, classMethod, varType, child)

	} else {

		has_indice, indice_rep := node.GetChild(ast.INDICEREP)

		if varRow.ArrDim == 0 && !has_indice {

			node.MigrateTagToParents(node.Lexeme)
			return &records.Variables{Iden: identifierName, BaseType: varType, ArrDim: 0}

		} else {

			var num_indices int

			if has_indice {

				num_indices = indice_rep.NumberOfChildren()

				for _, indices := range indice_rep.GetChildren() {

					indice_type := get_expression_type(location, class, classMethod, p_t, indices)

					if !indice_type.CompareVar(records.NewVariable("integer")) {

						outputMessage("Cannot use non-integer index: \"" + indices.Lexeme)
						outputMessage("\", Type: " + indice_type.ToString() + " Line: " + node.GetLine() + "\n")

						return records.NewVariable("invalid")
					}
				}

			} else {
				num_indices = 0
			}

			differences := varRow.ArrDim - num_indices

			if differences < 0 {

				outputMessage("Variable " + identifierName + " cannot be indexed. ")
				outputMessage("Var Dim: " + strconv.Itoa(varRow.ArrDim) + " Num of Indices: " + strconv.Itoa(num_indices))
				outputMessage(". Line: " + node.GetLine() + "\n")

				return records.NewVariable("invalid")
			}

			node.MigrateTagToParents(node.Lexeme)
			return &records.Variables{BaseType: varType, ArrDim: differences}
		}

	}

}

// Returns the return type of the function.
func function_type(location, scope *records.SymbolTable, classMethod, p_t string, node *ast.Node) *records.Variables {

	return_type := "invalid"
	funcName := node.Lexeme
	_, aparams := node.GetChild(ast.APARAMS)
	parameters := convertToFparams(location, scope, classMethod, aparams)

	if gb_table.EntryExists(funcName) && gb_table.GetRow(funcName).GetKind() == ast.FUNCTION {

		var tmp_class string

		if isClass(p_t) {
			tmp_class = p_t
		}

		funcRow := gb_table.GetRow(funcName).(*records.FunctionRow)
		return_type = funcRow.GetReturn(tmp_class, parameters)
	}

	if p_t != "" {

		if gb_table.EntryExists(p_t) {

			classRow := gb_table.GetRow(p_t).(*records.ClassRow)
			class_table := classRow.Link

			if class_table.GetVisibilityFromParam(funcName, parameters) == ast.PRIVATE && classMethod != p_t {
				outputMessage("Cannot reference private function, from non-classmethod context. Line: " + node.GetLine())
				outputMessage("\n")
				return records.NewVariable("invalid")
			}
		}
	}

	if return_type == "invalid" {

		outputMessage("No Function " + funcName + ":" + node.GetLine() + " with signature: ")

		for index, params := range *parameters {

			outputMessage(params.ToString())

			if index == len(*parameters)-1 {
				break
			}
			outputMessage(", ")
		}

		outputMessage("\n")
	}

	if return_type == "invalid" {
		return records.NewVariable("invalid")
	}

	if has_child, id := node.GetChild(ast.IDENTIFIER); has_child {

		if return_type == "integer" || return_type == "float" || return_type == "string" ||

			return_type == "void" {
			outputMessage("Error at line: " + node.GetLine() + ", ")
			outputMessage("Variable type " + return_type + ", Cannot have subType. " + "\n")

			return records.NewVariable("invalid")
		}

		new_class_table := gb_table.GetRow(return_type).(*records.ClassRow).Link
		return handle_id_type(location, new_class_table, classMethod, p_t, id)
	}

	label := generateTag(funcName)

	scope.CreateRow(label, ast.VARIABLE)
	scope.SetVariableType(label, return_type)
	scope.CalculateOffset(label, gb_table)
	scope.SetVisibility(label, ast.FUNCTION)

	node.Tag = label

	return records.NewVariable(return_type)
}

// converts Aparams to Fparams format to help with comparison
func convertToFparams(location, scope *records.SymbolTable, classMethod string, aparams *ast.Node) *records.Parameters {

	var func_param records.Parameters

	if aparams.NumberOfChildren() == 0 {

		func_param = append(func_param, records.Variables{})

		return &func_param
	}

	for _, param := range aparams.GetChildren() {

		variable := get_expression_type(location, scope, classMethod, "", param)
		param.Semantic_Type = variable.BaseType
		param.ArrDim = variable.ArrDim

		func_param = append(func_param, *variable)
	}

	return &func_param
}

// Returns whether string represents a class.
func isClass(parentType string) bool {
	return gb_table.EntryExists(parentType)
}
