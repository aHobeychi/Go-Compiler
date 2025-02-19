// File That will check if function's return type match.
package semantics

import (
	"aHobeychi/GoCompiler/ast"
	"aHobeychi/GoCompiler/records"
)

// Checks for the existance of return type and actual return of that type.
func return_check(root *ast.Node) bool {
	correct := true

	for _, child := range root.GetChildren() {
		if child.Curr_type == ast.MAIN {
			if check_main_has_return(child) {
				outputMessage("Main function cannot have return statement. Line: " + child.GetLine() + "\n")
				correct = false
				continue
			}
		} else if child.Curr_type == ast.FUNCDEF {

			if !child.HasChildren() {
				continue
			}

			for _, function := range child.GetChildren() {

				_, funcHead := function.GetChild(ast.FUNCTIONHEAD)
				_, fparams := funcHead.GetChild(ast.FPARAMS)
				_, funcBody := function.GetChild(ast.FUNCTIONBODY)
				_, name := funcHead.GetChild(ast.IDENTIFIER)
				_, return_node := funcHead.GetChild(ast.RETURNTYPE)
				_, classNode := funcHead.GetChild(ast.CLASSMETHOD)
				has_stats, stat := funcBody.GetChild(ast.STATEMENT)

				return_type := return_node.GetChildren()[0].Lexeme

				if !has_stats && return_type == "void" {
					continue
				} else if !has_stats {
					outputMessage("Function " + name.Lexeme + " requires a return statement. Line:" + funcHead.GetLine() + "\n")
					continue
				}

				has_return, returnNode := stat.GetChild(ast.RETURN)

				if !has_return && return_type == "void" {
					continue
				}

				if has_return && return_type == "void" {
					outputMessage("Cannot Have return statement when declared return type is void. Line: " + funcHead.GetLine() + "\n")
					correct = false
					continue
				} else if !has_return && return_type == "void" {
					continue
				} else if !has_return && return_type != "void" {
					outputMessage("Missing Return Statement! Line: " + funcHead.GetLine())
					outputMessage(", Expected Return Type: " + return_type + "\n")
					correct = false
					continue
				}

				if !check_for_multiple_return(stat) {
					outputMessage("Function: " + name.Lexeme + " has more than 1 return statement")
					outputMessage(". Line: ")
					outputMessage(funcHead.GetLine() + ".\n")
					correct = false
					continue
				}

				var classMethod string

				if classNode.HasChildren() {
					classMethod = classNode.GetChildren()[0].Lexeme
				} else {
					classMethod = ""
				}

				func_table := gb_table.GetRow(name.Lexeme).(*records.FunctionRow)
				link := func_table.GetLinkForClass(classMethod, fparams)

				_, return_expr := returnNode.Get(0)

				actualType := get_expression_type(link, link, classMethod, "", return_expr)

				if actualType.BaseType != return_type {
					outputMessage("Return Type Mismatch! Line " + funcHead.GetLine())
					outputMessage("\n\tActual Type: " + actualType.BaseType)
					outputMessage("\n\tDeclared Type: " + return_type + "\n")
					correct = false
				}

				if return_type == "integer" || return_type == "float" ||
					return_type == "string" || return_type == "void" {
					continue
				}

				if !gb_table.EntryExists(return_type) {
					outputMessage("Function: " + name.Lexeme + " has invalid return ")
					outputMessage("type: " + return_type + ". Line: ")
					outputMessage(funcHead.GetLine() + ".\n")
					correct = false
					continue
				}

				if gb_table.GetRow(return_type).GetKind() != ast.CLASS {
					outputMessage("Function: " + name.Lexeme + " has invalid return ")
					outputMessage("type: " + return_type + ". Line: ")
					outputMessage(funcHead.GetLine() + ".\n")
					correct = false
					continue
				}

			}
		}
	}

	return correct
}

// checks if there is multiple return statements in the same function
func check_for_multiple_return(stat *ast.Node) bool {
	num_of_return := 0

	for _, node := range stat.GetChildren() {
		if node.Curr_type == ast.RETURN {
			num_of_return += 1
		}
	}

	return num_of_return <= 1
}

// Returns true, if main function has return statement.
func check_main_has_return(main *ast.Node) bool {
	implemented, body := main.GetChild(ast.FUNCTIONBODY)

	if !implemented {
		return false
	}

	has_stats, stats := body.GetChild(ast.STATEMENT)

	if !has_stats {
		return false
	}

	for _, node := range stats.GetChildren() {
		if node.Curr_type == ast.RETURN {
			return true
		}
	}

	return false
}
