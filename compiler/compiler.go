// Package compiler contains the final struct.
// holds all different parts
package compiler

import (
	"aHobeychi/GoCompiler/ast"
	"aHobeychi/GoCompiler/lexer"
	"aHobeychi/GoCompiler/parser"
	"aHobeychi/GoCompiler/records"
	"aHobeychi/GoCompiler/reporting"
	"aHobeychi/GoCompiler/semantics"
	"aHobeychi/GoCompiler/translation"
	"os"
)

// Type compiler contains all necessary data types.
type Compiler struct {
	input      string
	output     string
	verbose    bool
	save_table bool
	save_ast   bool
	parser     *parser.Parser
	root       *ast.Node
	table      *records.SymbolTable
}

// Returns new compiler variable.
func New(input, output string) (*Compiler, bool) {

	lex := lexer.New(input, output)
	p := parser.New(lex, false, "", "")

	return &Compiler{parser: p, output: output, input: input}, true
}

// Compiles the input file and exports the instructions file.
func (comp *Compiler) Compile() {

	if comp.parser.IsErrorFree() && comp.runSemanticChecks() {

		file, _ := os.Create(comp.output + ".m")
		translation.SetOutputPath(file)
		translation.Translate(comp.root, comp.table)
		comp.printMemTable()

		defer file.Close()
	}

}

// Runs Semantic checks on symbol table and AST.
func (comp *Compiler) runSemanticChecks() bool {

	correct := false

	if comp.parser.IsErrorFree() {

		correct = semantics.CrawlTree(comp.table, comp.root)
		records.CalculateMemorySizes(comp.table, comp.table)
	}

	comp.printMemTable()

	return correct
}

//Parses The File
func (comp *Compiler) ParseFile() {

	var table_f *os.File
	var lex_f *os.File
	var prod_f *os.File

	if comp.save_table {

		table_f, _ = os.Create(comp.output + "_table")
		records.SetOutputPath(table_f)
	}

	if comp.verbose {

		lex_f, _ = os.Create(comp.output + "_outlextokens")
		reporting.SetTokenOutPath(lex_f)

		prod_f, _ = os.Create(comp.output + "_productions")
		reporting.SetProductPath(prod_f)
	}

	error_f, _ := os.Create(comp.output + "_errors")

	reporting.SetOutputPath(error_f)
	comp.root, comp.table = comp.parser.Parse()

	if comp.save_ast {

		comp.PrintAST()
	}

	records.PrintSymbolTable(comp.table)
	defer table_f.Close()
	defer lex_f.Close()
	defer prod_f.Close()
}

func (comp *Compiler) SetFlags(verbose, ast, symbol bool) {

	// Verbose
	comp.verbose = verbose
	comp.parser.SetVerbose(verbose)

	// Symbol Table
	comp.save_table = symbol

	// Abstract Symbol Table
	comp.save_ast = ast
}

// Prints AST to a file
func (comp *Compiler) PrintAST() {

	file, _ := os.Create(comp.output + "_ast")
	comp.root.TraverseToFile(0, file)

	defer file.Close()
}

// Prints Memory Table
func (comp *Compiler) printMemTable() {

	var table_f *os.File

	records.CalculateMemorySizes(comp.table, comp.table)

	if comp.save_table {

		table_f, _ = os.Create(comp.output + "_mem")
		records.SetOutputPath(table_f)
		records.PrintSymbolTable(comp.table)
		table_f.Close()
	}
}

// Closses all outputs in the parser.
func (comp *Compiler) KillOutputs() {
	comp.parser.CloseOutputs()
}
