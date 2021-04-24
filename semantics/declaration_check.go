// Will check that all used identifiers are within the current scope.
package semantics

import (
	"aHobeychi/GoCompiler/ast"
	"aHobeychi/GoCompiler/records"
)

var decl_valid bool

// Checks if all referenced variables have already been
// declared or exists in the scope.
func declaration_check(root *ast.Node) bool {

	decl_valid = true

	check_types_exist()

	for _, child := range root.GetChildren() {

		if child.Curr_type == ast.MAIN {

			verifyFunction("MAIN", "", nil, child)

		} else if child.Curr_type == ast.FUNCDEF {

			if !child.HasChildren() {
				continue
			}

			for _, function := range child.GetChildren() {

				var classMethod string

				_, funcHead := function.GetChild(ast.FUNCTIONHEAD)
				_, funcName := funcHead.GetChild(ast.IDENTIFIER)
				_, funcClass := funcHead.GetChild(ast.CLASSMETHOD)
				_, fparams := funcHead.GetChild(ast.FPARAMS)
				_, funcBody := function.GetChild(ast.FUNCTIONBODY)

				if funcClass.HasChildren() {

					classMethod = funcClass.GetChildren()[0].Lexeme

				} else {
					classMethod = ""
				}

				verifyFunction(funcName.Lexeme, classMethod, fparams, funcBody)
			}
		}
	}

	return decl_valid
}

func check_types_exist() {

	for k := range gb_table.Rows {

		row := gb_table.GetRow(k)

		if row.GetKind() == ast.VARIABLE {
			continue
		}

		if row.GetKind() == ast.FUNCTION {

			funcrow := row.(*records.FunctionRow)
			check_function_contents(funcrow)

		} else {

			classrow := row.(*records.ClassRow)
			check_class_contents(classrow.Link)
		}

	}

}

func check_function_contents(funcrow *records.FunctionRow) {

	for _, tables := range funcrow.Link {

		for k := range tables.Rows {

			row := tables.GetRow(k)

			if row.GetKind() != ast.VARIABLE {
				continue
			}

			varRow := row.(*records.VariableRow)
			varType := varRow.VarType

			if varType == "integer" || varType == "float" || varType == "string" {
				continue
			}

			if !gb_table.EntryExists(varType) {

				outputMessage("Func: " + tables.Scope + ", contains variable: ")
				outputMessage(varRow.Name + " of type: " + varType + " that has not been defined.\n")
				decl_valid = false

				continue
			}
		}
	}
}

func check_class_contents(class_table *records.SymbolTable) {

	for k := range class_table.Rows {

		row := class_table.GetRow(k)

		if row.GetKind() != ast.VARIABLE {
			continue
		}

		varRow := row.(*records.VariableRow)
		varType := varRow.VarType

		if varType == "integer" || varType == "float" || varType == "string" {
			continue
		}

		if !gb_table.EntryExists(varType) {

			outputMessage("Class: " + class_table.Scope + ", contains variable: ")
			outputMessage(varRow.Name + " of type: " + varType + " that has not been defined.\n")
			decl_valid = false

			continue
		}

	}

}

// // will branch to diferrent functions, to check their func bodies.
func verifyFunction(func_name, class string, params, root *ast.Node) {

	func_table := gb_table.GetRow(func_name).(*records.FunctionRow)

	if func_name == "MAIN" {

		verifyFunctionBody(func_table.Link[0], func_table.Link[0], "", root)
		return
	}

	location := func_table.GetLinkForClass(class, params)

	verifyFunctionBody(location, location, class, root)
}

// Will check all the variables inside the function body to see if they been declared.
func verifyFunctionBody(table, scope *records.SymbolTable, classMethod string, root *ast.Node) {

	if root.Curr_type == ast.VARIABLEDECL {

		handleVariableDec(root, table)
		return
	}

	if root.Curr_type == ast.IDENTIFIER {

		handleIdentifer(table, scope, classMethod, "", root)
		return
	}

	if root.HasChildren() {

		for _, child := range root.GetChildren() {

			verifyFunctionBody(table, scope, classMethod, child)
		}
	}
}

func handleIdentifer(location, class *records.SymbolTable, classMethod, pt string, node *ast.Node) {

	if exists, _ := node.GetChild(ast.APARAMS); exists {

		handleFunction(location, class, classMethod, pt, node)

	} else {

		handleVariable(location, class, classMethod, pt, node)
	}
}

func handleVariable(location, class *records.SymbolTable, classMethod, pt string, node *ast.Node) {

	var varType string
	var varRow *records.VariableRow

	identifierName := node.Lexeme

	if location.EntryExists(identifierName) {

		row := location.GetRow(identifierName)
		varRow = row.(*records.VariableRow)
		varType = varRow.GetType()

	} else if class.EntryExists(identifierName) {

		row := class.GetRow(identifierName)
		varRow = row.(*records.VariableRow)
		varType = varRow.GetType()

	} else if pt != "" && gb_table.EntryExists(pt) && gb_table.InheritsSearch(pt, identifierName) {

		varRow = gb_table.GetInheritedVar(pt, identifierName)
		varType = varRow.GetType()

	} else if gb_table.EntryExists(classMethod) && gb_table.InheritsSearch(classMethod, identifierName) {

		varRow = gb_table.GetInheritedVar(classMethod, identifierName)
		varType = varRow.GetType()

		if gb_table.InheritsFromParen(classMethod, identifierName) {

			outputMessage("Error at line: " + node.GetLine() + ", ")
			outputMessage("Cannot reference private member " + identifierName)
			outputMessage(" in current scope: " + location.Scope)
			outputMessage(".\n")

			decl_valid = false

			return
		}

	} else {

		outputMessage("Error at line: " + node.GetLine() + ", ")
		outputMessage("Variable " + identifierName + " not in current scope: " + class.Scope)

		if classMethod != "" {
			outputMessage("::" + classMethod + "")
		}

		outputMessage(".\n")

		decl_valid = false

		return
	}

	// CHECKS IF THE VISIBILITY OF THE VARIABLE FITS THE CALL
	if !gb_table.ClassEntryExists(classMethod, identifierName) && varRow.Visibility == ast.PRIVATE {

		outputMessage("Error at line: " + node.GetLine() + ", ")
		outputMessage("Cannot reference private member " + identifierName + " of class " + pt)
		outputMessage(" in current scope: " + location.Scope)
		outputMessage(".\n")

		decl_valid = false

		return
	}

	has_child, child := node.GetChild(ast.IDENTIFIER)

	if has_child {

		if varType == "integer" || varType == "float" || varType == "string" {

			outputMessage("Error at line: " + node.GetLine() + ", ")
			outputMessage("Dot operator used on non-class type: " + varType + "\n")
			decl_valid = false

			return
		}

		new_class_table := gb_table.GetRow(varType).(*records.ClassRow).Link
		handleIdentifer(location, new_class_table, classMethod, varType, child)

	}

}

func handleFunction(location, class *records.SymbolTable, classMethod, pt string, node *ast.Node) {

	return_type := "invalid"

	funcName := node.Lexeme
	_, aparams := node.GetChild(ast.APARAMS)
	parameters := convertToFparams(location, class, classMethod, aparams)

	var funcRow *records.FunctionRow

	if gb_table.EntryExists(funcName) && gb_table.GetRow(funcName).GetKind() == ast.FUNCTION {

		funcRow = gb_table.GetRow(funcName).(*records.FunctionRow)
		return_type = funcRow.GetReturn(class.Scope, parameters)

	}

	if has_child, id := node.GetChild(ast.IDENTIFIER); has_child {

		if return_type == "integer" || return_type == "float" || return_type == "string" ||

			return_type == "void" || return_type == "invalid" {
			outputMessage("Error at line: " + node.GetLine() + ", ")
			outputMessage("Variable type " + return_type + ", Cannot have subType. " + "\n")

			return
		}

		new_class_table := gb_table.GetRow(return_type).(*records.ClassRow).Link
		handleIdentifer(location, new_class_table, classMethod, return_type, id)
	}

	return
}

func handleVariableDec(root *ast.Node, location *records.SymbolTable) {

	if !root.HasChildren() {
		return
	}

	for _, variable := range root.GetChildren() {

		_, variable_type := variable.GetChild(ast.TYPE)

		v_t := variable_type.Lexeme

		if !(v_t == "integer" || v_t == "float" || v_t == "string") {

			if !gb_table.EntryExists(v_t) {

				outputMessage("No Type: " + v_t + ", in var decl of func: ")
				outputMessage(location.Scope + ", Line: " + variable.GetLine() + "\n")
				decl_valid = false

				continue

			} else if gb_table.EntryExists(v_t) && gb_table.GetRow(v_t).GetKind() != ast.CLASS {

				outputMessage("No Type: " + v_t + ", in var decl of func: ")
				outputMessage(location.Scope + ", Line: " + variable.GetLine() + "\n")
				decl_valid = false

				continue
			}
		}
	}

}
