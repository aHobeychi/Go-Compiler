// Package parser contains all the functions pertaining to parsing the lexical
// tokens to create an abstact syntax tree. Contains the Parser struct and
// Parses using the Parse() function.
package parser

import (
	"aHobeychi/GoCompiler/ast"
	"aHobeychi/GoCompiler/lexer"
	"aHobeychi/GoCompiler/records"
	"aHobeychi/GoCompiler/reporting"
	"aHobeychi/GoCompiler/semantics"
	"os"
)

// Parser Struct Contains Instance of Lexer.
type Parser struct {
	l         *lexer.Lexer
	thisToken lexer.Token
	nextToken lexer.Token
	prod_out  *os.File
	err_out   *os.File
	firstSet  ast.Set
	followSet ast.Set
	table     *records.SymbolTable
	verbose   bool
	err_free  bool
}

// Sets if parser will print out all of error information.
func (p *Parser) SetVerbose(verbose bool) {
	p.verbose = verbose
}

// Return a new Parser.
func New(l *lexer.Lexer, verbose bool, err string, prod string) *Parser {

	first := ast.GenerateFirstSet()
	follow := ast.GenerateFollowSet()

	if verbose {

		err_f, _ := os.Create(err)
		prod_out, _ := os.Create(prod)
		parser := &Parser{
			l: l, verbose: verbose,
			err_out: err_f, prod_out: prod_out, err_free: true,
			firstSet: first, followSet: follow}
		parser.readToken()
		parser.readToken()

		return parser

	} else {

		parser := &Parser{l: l, verbose: verbose, err_free: true,
			firstSet: first, followSet: follow}

		parser.readToken()
		parser.readToken()

		return parser

	}
}

// Reads The next Token
func (p *Parser) readToken() {
	p.thisToken = p.nextToken
	p.nextToken = p.l.NextToken()
}

// Public function to parse the file.
func (p *Parser) Parse() (*ast.Node, *records.SymbolTable) {
	return p.parseProgram()
}

// begins the parsing using:  <Prog> ::= <ClassDecl> <FuncDef> 'main' <FuncBody>
// Returns an Abstract Syntax Tree.
func (p *Parser) parseProgram() (*ast.Node, *records.SymbolTable) {

	if !p.skipErrors(ast.PROGRAM) {
		return nil, nil
	}

	root := ast.New(ast.PROGRAM, "")
	// Creates the symbol table
	p.table = records.NewTable("PROGRAM")

	p.printProductions("PROGRAM")

	// Parse Class Declerations
	root.AddChild(p.parseClassDecl())

	// Parse Function Declerations
	root.AddChild(p.parseFuncDef())

	// Parse Main Functiom
	main := p.parseMain()

	root.AddChild(main)

	if p.thisToken.Type != lexer.EOF {
		p.printMissingMessage(lexer.EOF)
		p.parseFromFunc()
		p.parseFromMain()
		return ast.New(ast.INVALID, "COULD NOT PARSE TILL THE END"), nil
	} else if !p.err_free {
		return ast.New(ast.INVALID, "CONTAINS ERROR IN PROGRAM"), nil
	} else {
		root.AddChild(ast.New(ast.EOF, "EOF"))
		p.printProductions("<EOF>")
	}

	if !semantics.Inheritance_check(p.table) {
		p.err_free = false
		return ast.New(ast.INVALID, "CONTAINS ERROR IN PROGRAM"), nil
	}

	records.CalculateMemorySizes(p.table, p.table)

	return root, p.table
}

// Parses the main function.
func (p *Parser) parseMain() *ast.Node {

	if p.thisToken.Type == lexer.MAIN {
		p.readToken()
	} else {
		p.readTillSymbol(lexer.LBRACE)
		p.printMissingMessage(lexer.MAIN)
	}

	main := ast.New(ast.MAIN, "MAIN")
	main.Line = p.getLineNumber()
	body, mainTable := p.parseFuncBody("MAIN")

	main.AddChild(body)

	p.table.CreateRow("MAIN", ast.FUNCTION)
	p.table.LinkRowToTable("MAIN", ast.FUNCTION, mainTable)

	return main
}

// If parsing the Class declarations fail will still parse the functions
// to atleast output misssing characters or lexical errors.
func (p *Parser) parseFromFunc() {

	for p.thisToken.Type != lexer.FUNCTION {

		if p.thisToken.Type == lexer.EOF {
			break
		}
		p.readToken()
	}

	p.parseFuncDef()

	if p.thisToken.Type == lexer.MAIN {
		p.readToken()
	}
	p.parseFuncDef()
}

// If parsing the Func declarations fail will still parse the main func.
// to atleast output misssing characters or lexical errors.
func (p *Parser) parseFromMain() {

	for p.thisToken.Type != lexer.MAIN {

		if p.thisToken.Type == lexer.EOF {
			break
		}
		p.readToken()
	}
	p.readToken()

	p.parseFuncDef()
}

// Parses Tokens related to array size.
func (p *Parser) parseArraySizeRept() *ast.Node {

	if !p.skipErrors(ast.ARRAYSIZE) {
		return nil
	}

	arraySpecifications := ast.New(ast.ARRAYSIZE, "")

	for {

		if p.thisToken.Type == lexer.LBRACKET {
			p.readToken()
		} else {
			break
		}

		p.printProductions("ArraySize:")
		arraySpecification := ast.New(ast.ARRAYSIZE, "")
		arraySpecification.AddChild(p.parseIntNum())
		arraySpecifications.AddChild(arraySpecification)

		if p.thisToken.Type == lexer.RBRACKET {
			p.readToken()
		} else {
			p.printMissingMessage(lexer.RBRACKET)
		}
	}

	return arraySpecifications
}

// Parses Tokens related Integer Values.
func (p *Parser) parseIntNum() *ast.Node {

	if p.thisToken.Type == lexer.INT_VALUE {

		p.printProductions("IntNum: " + p.thisToken.Lexeme)
		int_value := ast.New(ast.INT_VALUE, p.thisToken.Lexeme)
		p.readToken()
		return int_value

	} else if p.thisToken.Type == lexer.INVALIDINT {

		p.printMissingMessage(lexer.INT_VALUE)
		int_value := ast.New(ast.INVALID, p.thisToken.Lexeme)
		p.readToken()
		return int_value

	}

	return nil
}

// Parses Tokens related Float Values.
func (p *Parser) parseFloatNum() *ast.Node {

	if p.thisToken.Type == lexer.FLOAT_VALUE {

		p.printProductions("FloatNum: " + p.thisToken.Lexeme)
		float_value := ast.New(ast.FLOAT_VALUE, p.thisToken.Lexeme)
		p.readToken()
		return float_value

	} else if p.thisToken.Type == lexer.INVALIDFLOAT {

		float_value := ast.New(ast.INVALID, p.thisToken.Lexeme)
		p.readToken()
		return float_value

	}

	return nil
}

// Parses Tokens related to String Literals.
func (p *Parser) parseStringLit() *ast.Node {

	if p.thisToken.Type == lexer.STRINGLIT {

		p.printProductions("StringLit: " + p.thisToken.Lexeme)
		stringlit := ast.New(ast.STRINGLIT, p.thisToken.Lexeme)
		p.readToken()
		return stringlit

	}

	return nil
}

// Parses Tokens related to Visibility: Public, Private, Default
func (p *Parser) parseVisibility() *ast.Node {

	if !p.skipErrors(ast.VISIBILITY) {
		return nil
	}

	visibility := ast.New(ast.VISIBILITY, "")

	switch p.thisToken.Type {
	case lexer.PUBLIC:
		p.printProductions("Visibility: public")
		visibility.AddChild(ast.New(ast.PUBLIC, "public"))
		p.readToken()
		return visibility
	case lexer.PRIVATE:
		p.printProductions("Visibility: private")
		visibility.AddChild(ast.New(ast.PRIVATE, "private"))
		p.readToken()
		return visibility
	}

	p.printProductions("Visibility: default")
	visibility.AddChild(ast.New(ast.DEFAULT, "default"))
	return visibility
}

// Parses Tokens related to Variables <Variable> ::= 'id' <VariableIdnest>
func (p *Parser) parseVariable() *ast.Node {

	if !p.skipErrors(ast.VARIABLE) {
		return nil
	}

	var id *ast.Node

	if p.thisToken.Type == lexer.IDENT {

		p.printProductions("Variable: ")
		id = p.parseIdentifier()

	}

	id.AddChild(p.parseVariableIdnest())
	return id
}

// Parses Tokens <VariableIdnest> ::= <IndiceRep> <VariableIdnestTail>
func (p *Parser) parseVariableIdnest() *ast.Node {

	if !p.skipErrors(ast.VARIABLEIDNEST) {
		return nil
	}

	varIdNest := ast.New(ast.VARIABLEIDNEST, "")
	varIdNest.AddChild(p.parseIndiceRep())
	varIdNest.AddChild(p.parseVariableIdnestTail())

	return varIdNest

}

// Parses token related to the array indice
// <IndiceRep> ::= '[' <Expr> ']' <IndiceRep> | EPSILON.
func (p *Parser) parseIndiceRep() *ast.Node {

	if !p.skipErrors(ast.INDICEREP_PROD) {
		return nil
	}

	indiceRep := ast.New(ast.INDICEREP, "")
	indiceRep.Line = p.getLineNumber()
	p.printProductions("IndiceRep: ")

	if p.thisToken.Type != lexer.LBRACKET {
		return nil
	}

	for {

		if p.thisToken.Type != lexer.LBRACKET {
			break
		} else {
			p.readToken()
		}

		expr := p.parseExpression()
		expr.Line = p.getLineNumber()
		indiceRep.AddChild(expr)

		if p.thisToken.Type != lexer.RBRACKET {
			break
		} else {
			p.readToken()
		}

	}

	return indiceRep
}

// <VariableIdnestTail> ::= '.' 'id' <VariableIdnest> | EPSILON
func (p *Parser) parseVariableIdnestTail() *ast.Node {

	if !p.skipErrors(ast.VARIABLEIDNESTTAIL) {
		return nil
	}

	varIdNest := ast.New(ast.VARIABLEIDNEST, "")
	varIdNest.Line = p.getLineNumber()

	if p.thisToken.Type == lexer.DOT {

		p.readToken()

		if p.thisToken.Type == lexer.IDENT {
			p.printProductions("NestedId:")
			varIdNest.AddChild(p.parseIdentifier())
		}

	} else {
		return nil
	}

	varIdNest.AddChild(p.parseVariableIdnestTail())
	return varIdNest
}

// Parses Tokens related to the assign operation =.
func (p *Parser) parseAssignOp() *ast.Node {

	if !p.skipErrors(ast.ASSIGN) {
		return nil
	}

	p.printProductions("AssignOp: " + p.thisToken.Lexeme)

	if p.thisToken.Type == lexer.ASSIGN {

		assignOp := ast.New(ast.ASSIGN, p.thisToken.Lexeme)
		assignOp.Line = p.getLineNumber()
		p.readToken()

		return assignOp
	}

	return nil
}

// Parses Tokens related to the assign operation =.
func (p *Parser) parseSign() *ast.Node {

	if !p.skipErrors(ast.SIGN) {
		return nil
	}

	p.printProductions("Sign: " + p.thisToken.Lexeme)

	switch p.thisToken.Type {
	case lexer.MINUS:
		p.readToken()
		return ast.New(ast.NEGATIVE, "-")
	case lexer.PLUS:
		p.readToken()
		return ast.New(ast.POSITIVE, "+")
	}

	return nil
}

// Parses Tokens that represents identifiers.
func (p *Parser) parseIdentifier() *ast.Node {

	p.printProductions("IDENTIFIER: " + p.thisToken.Lexeme)
	identifier := ast.New(ast.IDENTIFIER, p.thisToken.Lexeme)
	identifier.Line = p.getLineNumber()
	p.readToken()

	return identifier
}

// Parses Token related to Add Operations +, -.
func (p *Parser) parseAddOp() *ast.Node {

	if !p.skipErrors(ast.PLUS) {
		return nil
	}

	p.printProductions("AddOP: " + p.thisToken.Lexeme)

	switch p.thisToken.Type {
	case lexer.PLUS:
		p.readToken()
		return ast.New(ast.PLUS, "+")
	case lexer.MINUS:
		p.readToken()
		return ast.New(ast.MINUS, "-")
	case lexer.OR:
		p.readToken()
		return ast.New(ast.OR, "or")
	}

	return nil
}

// Parses tokens related to Relational operations ==, !=. <, >, <=, >=.
func (p *Parser) parseRelOp() *ast.Node {

	if !p.skipErrors(ast.LT) {
		return nil
	}

	p.printProductions("RelOp: " + p.thisToken.Lexeme)

	switch p.thisToken.Type {
	case lexer.EQUAL:
		p.readToken()
		return ast.New(ast.EQ, "==")
	case lexer.NOTEQ:
		p.readToken()
		return ast.New(ast.NOTEQ, "!=")
	case lexer.LT:
		p.readToken()
		return ast.New(ast.LT, "<")
	case lexer.GT:
		p.readToken()
		return ast.New(ast.GT, ">")
	case lexer.LTE:
		p.readToken()
		return ast.New(ast.LTE, "<=")
	case lexer.GTE:
		p.readToken()
		return ast.New(ast.GTE, ">=")
	}

	return nil
}

// Checks if the current token is a Sign
func (p *Parser) tokenIsSign() bool {

	switch p.thisToken.Type {
	case lexer.PLUS:
		return true
	case lexer.MINUS:
		return true
	default:
		return false
	}
}

// Checks if the current token is an Add operation.
func (p *Parser) tokenIsAddOp() bool {

	switch p.thisToken.Type {
	case lexer.PLUS:
		return true
	case lexer.MINUS:
		return true
	case lexer.OR:
		return true
	}
	return false
}

// Checks if the current token is a relational operation.
func (p *Parser) tokenIsRelOp() bool {

	switch p.thisToken.Type {
	case lexer.EQUAL:
		return true
	case lexer.NOTEQ:
		return true
	case lexer.LT:
		return true
	case lexer.GT:
		return true
	case lexer.LTE:
		return true
	case lexer.GTE:
		return true
	}
	return false
}

// Checks if the current token is a type: string, integer, float, Ident.
func (p *Parser) tokenIsType() bool {

	switch p.thisToken.Type {
	case lexer.INTEGER:
		return true
	case lexer.FLOAT:
		return true
	case lexer.STRING:
		return true
	case lexer.IDENT:
		return true
	}
	return false
}

// Checks if the current token is in the First set of the
// Expression Production rule.
func (p *Parser) tokenIsExpression() bool {

	switch p.thisToken.Type {
	case lexer.INT_VALUE:
		return true
	case lexer.FLOAT_VALUE:
		return true
	case lexer.STRINGLIT:
		return true
	case lexer.LPAREN:
		return true
	case lexer.NOT:
		return true
	case lexer.QUESTION:
		return true
	case lexer.IDENT:
		return true
	case lexer.PLUS:
		return true
	case lexer.MINUS:
		return true
	}
	return false

}

// Checks if the current token is in the First set of the
// Statement Production rule.
func (p *Parser) tokenIsStatement() bool {

	switch p.thisToken.Type {
	case lexer.IF:
		return true
	case lexer.WHILE:
		return true
	case lexer.READ:
		return true
	case lexer.WRITE:
		return true
	case lexer.RETURN:
		return true
	case lexer.BREAK:
		return true
	case lexer.CONTINUE:
		return true
	case lexer.IDENT:
		return true
	}
	return false
}

// Parses the tokens related to the Multiply operation: *, /, &.
func (p *Parser) parseMultOp() *ast.Node {

	if !p.skipErrors(ast.MUL) {
		return nil
	}

	p.printProductions("MultOp: " + p.thisToken.Lexeme)

	switch p.thisToken.Type {
	case lexer.MUL:
		p.readToken()
		return ast.New(ast.MUL, "*")
	case lexer.DIV:
		p.readToken()
		return ast.New(ast.DIV, "/")
	case lexer.AND:
		p.readToken()
		return ast.New(ast.AND, "&")
	}

	return nil
}

// Checks if the current token is in the first set of the Multiply
// First set.
func (p *Parser) tokenIsMultOp() bool {

	switch p.thisToken.Type {
	case lexer.MUL:
		return true
	case lexer.DIV:
		return true
	case lexer.AND:
		return true
	}
	return false
}

// Parses tokens related to the types: integer, string, float , Ident.
func (p *Parser) parseType() *ast.Node {

	if !p.skipErrors(ast.TYPE) {
		return nil
	}

	p.printProductions("Type: " + p.thisToken.Lexeme)

	switch p.thisToken.Type {
	case lexer.INTEGER:
		varType := ast.New(ast.TYPE, p.thisToken.Lexeme)
		p.readToken()
		return varType
	case lexer.FLOAT:
		varType := ast.New(ast.TYPE, p.thisToken.Lexeme)
		p.readToken()
		return varType
	case lexer.STRING:
		varType := ast.New(ast.TYPE, p.thisToken.Lexeme)
		p.readToken()
		return varType
	case lexer.IDENT:
		varType := ast.New(ast.TYPE, p.thisToken.Lexeme)
		p.readToken()
		return varType
	}

	return nil
}

// Return Line of the current token.
func (p *Parser) getLine() string {
	return p.thisToken.GetLine()
}

func (p *Parser) getLineNumber() int {
	return p.thisToken.Line
}

// Returns Message Of Missing Token
func (p *Parser) printMissingMessage(tType lexer.TokenType) {

	p.err_free = false

	if p.verbose {
		reporting.OutputMessage("Problem at: " + p.getLine() + ", Maybe: " +
			lexer.TokenAsString[tType] + "\n")
	}
}

// Print Productions
func (p *Parser) printProductions(mess string) {

	if p.verbose {
		reporting.OutputProd("" + mess + " ")
	}
}

// Close Outputs
func (p *Parser) CloseOutputs() {
	p.err_out.Close()
	p.prod_out.Close()
}

// Sets Error Output to be the same as all other phases.
func (p *Parser) SetErrorOut(path *os.File) {
	p.err_out = path
}

// Reads Till Semi Colon, Used to recover from error
func (p *Parser) readTillSymbol(symbol lexer.TokenType) {

	for {

		if p.thisToken.Type == symbol {
			break
		}
		p.readToken()
	}
}

// Skips Errors - Implementation of Don't Panic Technique
func (p *Parser) skipErrors(prod ast.Production) bool {

	if p.firstSet[prod][p.thisToken.Type] {

		return true

	} else if p.firstSet[prod][lexer.EMPTY] &&

		p.followSet[prod][p.thisToken.Type] {
		return true
	}

	p.printMissingMessage(lexer.EMPTY)

	for !(p.firstSet[prod][p.thisToken.Type] ||
		p.followSet[prod][p.thisToken.Type]) {

		p.readToken()

		if p.firstSet[prod][lexer.EMPTY] && p.followSet[prod][p.thisToken.Type] {
			return false
		}
	}

	return true
}

// Returns if parser has found an error, pre-semantic check.
func (p *Parser) IsErrorFree() bool {
	return p.err_free
}
