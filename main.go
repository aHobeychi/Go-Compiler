// Main File, will create compiler an run it on flagged files.
package main

import (
	"aHobeychi/GoCompiler/compiler"
	"flag"
)

// Flag Variables.
var (
	i       *string
	o       *string
	verbose *bool
	symbol  *bool
	ast     *bool
	h       *bool
)

// Initializes the flags.
func init() {
	i = flag.String("i", "none", "input file")
	o = flag.String("o", "sample/", "output file")
	verbose = flag.Bool("verbose", false, "print verbose")
	symbol = flag.Bool("symbol", false, "print symbol table")
	ast = flag.Bool("ast", false, "print symbol table")
	h = flag.Bool("h", false, "Display help message")
}

func main() {

	usage_message :=

		`To Run the compiler, Please run the program with the program using the following flags:

            -i:         followed by the input file.
            -o:         followed by the ouput path.
            -h:         to print the help message.
            -verbose:   to print out all error message.
            -ast:       to print out the abstract syntax tree to a file.
            -symbol:    to print out the symbol table to a file.
            `

	flag.Parse()

	if *h == true {
		print("\n" + usage_message + "\n")
	}

	if *i == "none" {
		print("\nMissing an input file. Use -h for help message.\n\n")
		return
	}

	if *o == "." {
		print("\nMissing Output File, Using the current working directory.")
		print(" Use -h for help message..\n\n")
		return
	}

	comp, _ := compiler.New(*i, *o)
	comp.SetFlags(*verbose, *ast, *symbol)
	comp.ParseFile()
	comp.Compile()
}
