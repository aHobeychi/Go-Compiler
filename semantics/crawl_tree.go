package semantics

import (
	"aHobeychi/GoCompiler/ast"
	"aHobeychi/GoCompiler/records"
	"aHobeychi/GoCompiler/reporting"
)

var gb_table *records.SymbolTable

func outputMessage(message string) {
	print(message)
	reporting.OutputMessage(message)
}

// Runs all the tests and returns if the input file passes.
func CrawlTree(table *records.SymbolTable, root *ast.Node) bool {
	gb_table = table

	decl := declaration_check(root)
	return_res := return_check(root)
	class_func := class_func_check()

	if !decl || !return_res || !class_func {
		return false
	}

	type_check := type_check(root, "PROGRAM")

	correct := return_res &&
		class_func && type_check && decl

	if correct {
		mergeTable(root)
	}

	return correct
}
