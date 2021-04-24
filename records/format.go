package records

import (
	"aHobeychi/GoCompiler/ast"
	"os"
)

var outputPath *os.File

func SetOutputPath(path *os.File) {
	outputPath = path
}

func OutputMessage(message string) {
	outputPath.WriteString(message)
}

// prints the symbol table and its rows.
func PrintSymbolTable(table *SymbolTable) {

	if table == nil {
		return
	}

	OutputMessage("―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――\n")
	OutputMessage("| \t\t\t\t―――――" + table.Scope + "―――――\n")
	OutputMessage("-------------------------------------------------------------------------------\n")

	for _, variable := range table.variables {

		PrintRow(table.GetRow(variable))
		OutputMessage("-------------------------------------------------------------------------------\n")
	}

	for k := range table.Rows {

		if table.GetRow(k).GetKind() == ast.VARIABLE {
			continue
		}
		PrintRow(table.GetRow(k))
		OutputMessage("-------------------------------------------------------------------------------\n")
	}

	OutputMessage("―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――\n")
	OutputMessage("\n")

	for k := range table.Rows {

		switch table.GetRow(k).GetKind() {
		case ast.FUNCTION:
			for i := 0; i < len(table.GetRow(k).(*FunctionRow).Link); i++ {
				PrintSymbolTable(table.GetRow(k).(*FunctionRow).getLink(i))
			}
		case ast.CLASS:
			PrintSymbolTable(table.GetRow(k).(*ClassRow).Link)
		}
	}
}
