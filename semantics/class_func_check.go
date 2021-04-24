package semantics

import (
	"aHobeychi/GoCompiler/ast"
	"aHobeychi/GoCompiler/records"
)

// Checks if class method has been declared with the proper signature.
func class_func_check() bool {

	correct := true

	for k := range gb_table.Rows {

		row := gb_table.GetRow(k)

		if row.GetKind() == ast.CLASS || row.GetKind() == ast.VARIABLE {
			continue
		}

		funcrow := row.(*records.FunctionRow)
		func_name := funcrow.Name
		classmethod := funcrow.ClassMethod

		for _, class := range classmethod {

			if class == "" {
				continue
			}

			if !gb_table.EntryExists(class) {

				outputMessage("No class: " + class + ", for function " + func_name + ".\n")
				correct = false

				continue
			}

			if !gb_table.EntryExists(class) && gb_table.GetRow(class).GetKind() == ast.CLASS {

				outputMessage("No class: " + class + ", for function " + func_name + ".\n")
				correct = false

				continue
			}

			class_table := gb_table.GetRow(class).(*records.ClassRow).Link

			if !class_table.EntryExists(func_name) {

				outputMessage("Class: " + class + " has no function: " + func_name + ".\n")
				correct = false

				continue
			}

			if class_table.EntryExists(func_name) && class_table.GetRow(func_name).GetKind() != ast.FUNCTION {

				outputMessage("Class: " + class + " has no function: " + func_name + ".\n")
				correct = false

				continue
			}

			class_funcs := class_table.GetRow(func_name).(*records.FunctionRow)

			if !matchReturnAndParams(class_funcs, funcrow) {

				outputMessage("Class : " + class + " has no function: " + func_name)
				outputMessage(" with specific parameter list.\n")
				correct = false

				continue
			}
		}
	}

	return check_defined_functions() && correct
}

// Checks if all declared functions have definitions.
func check_defined_functions() bool {

	correct := true

	for k := range gb_table.Rows {

		row := gb_table.GetRow(k)

		if row.GetKind() == ast.FUNCTION || row.GetKind() == ast.VARIABLE {
			continue
		}

		classrow := row.(*records.ClassRow)
		classTable := classrow.Link

		if classTable == nil {
			continue
		}

		for classR := range classTable.Rows {

			class_row := classTable.GetRow(classR)

			if class_row.GetKind() == ast.CLASS || class_row.GetKind() == ast.VARIABLE {
				continue
			}

			final_func := class_row.(*records.FunctionRow)
			funcName := final_func.Name
			return_type := final_func.ReturnType

			if !gb_table.EntryExists(funcName) {

				outputMessage("Function " + funcName + " has not been defined in the global scope, ")
				outputMessage("with the correct classmethod for class " + k + ".\n")
				correct = false

				continue
			}

			if gb_table.EntryExists(funcName) && gb_table.GetRow(funcName).GetKind() != ast.FUNCTION {

				outputMessage("Function " + funcName + " has not been defined in the global scope, ")
				outputMessage("with classmethod for class " + k + ".\n")
				correct = false

				continue
			}

			gb_function := gb_table.GetRow(funcName).(*records.FunctionRow)

			matches_class_methods := false

			for index, classMethods := range gb_function.ClassMethod {

				if classMethods == k {
					if matchInArr(return_type, gb_function.ReturnType[index]) {

						matches_class_methods = true
					}
				}
			}

			if !matches_class_methods {

				outputMessage("Function " + funcName + " has not been defined in the global scope, ")
				outputMessage("with classmethod " + k + " and return type.\n")

				correct = false
				continue
			}

		}
	}

	return correct
}

func matchInArr(arr []string, str string) bool {

	for _, arr_str := range arr {
		if arr_str == str {
			return true
		}
	}

	return false
}

func matchReturnAndParams(first_func, second_func *records.FunctionRow) bool {

	first_return := first_func.ReturnType
	second_return := second_func.ReturnType

	for firstIndex, rType := range first_return {
		for secondIndex, secondType := range second_return {

			if rType == secondType {
				if matchParameters(first_func.ListParameters[firstIndex], second_func.ListParameters[secondIndex]) {

					return true
				}
			}
		}
	}

	return false
}

// func matchParameters(first_parameters parameters, sec_parameters parameters) bool {
func matchParameters(first_parameters, sec_parameters records.Parameters) bool {
	return first_parameters.MatchParameter(sec_parameters)
}
