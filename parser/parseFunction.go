// File that contains all functions related to parsing functions.
package parser

import (
	"aHobeychi/GoCompiler/ast"
	"aHobeychi/GoCompiler/lexer"
	"aHobeychi/GoCompiler/records"
)

// Parses the
// <FuncDecl> ::= 'func' 'id' '(' <FParams> ')' ':' <FuncDeclTail> ';' Prod.
func (p *Parser) parseFuncDecl() *ast.Node {
	if !p.skipErrors(ast.FUNCDECL) {
		return nil
	}

	p.printProductions("<FUNCDEL>")

	functionNode := ast.New(ast.FUNCTION, "")
	if p.thisToken.Type == lexer.FUNCTION {
		p.readToken()
	}

	var id *ast.Node
	if p.thisToken.Type == lexer.IDENT {
		id = p.parseIdentifier()
	}
	functionNode.AddChild(id)

	if p.thisToken.Type == lexer.LPAREN {
		p.readToken()
	}

	fparams := p.parseFParams()
	functionNode.AddChild(fparams)

	if p.thisToken.Type == lexer.RPAREN {
		p.readToken()
	}

	if p.thisToken.Type == lexer.COLON {
		p.readToken()
	}

	tail := p.parseFuncDeclTail()
	functionNode.AddChild(tail)

	if p.thisToken.Type == lexer.SEMICOLON {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.SEMICOLON)
	}

	return functionNode
}

// Parses the <FuncDeclTail> ::= <Type> | 'void' Production.
func (p *Parser) parseFuncDeclTail() *ast.Node {
	p.printProductions("<FUNCDECLTAIL>")

	if p.thisToken.Type == lexer.VOID {

		p.readToken()

		fReturn := ast.New(ast.RETURNTYPE, "")
		fReturn.AddChild(ast.New(ast.VOID, "void"))

		return fReturn
	}

	if !p.skipErrors(ast.FUNCDECLTAIL) {
		return nil
	}

	fReturn := ast.New(ast.RETURNTYPE, "")
	fReturn.AddChild(p.parseType())

	return fReturn
}

// Parses the Function Parameter productions.
// <FParams> ::= <Type> 'id' <ArraySizeRept> <FParamsTail> | EPSILON
// <FParamsTail> ::= ',' <Type> 'id' <ArraySizeRept> <FParamsTail> | EPSILON
func (p *Parser) parseFParams() *ast.Node {
	if !p.skipErrors(ast.FPARAMS) {
		return nil
	}

	p.printProductions("<FPARAMS>")
	fparams := ast.New(ast.FPARAMS, "")

	for p.tokenIsType() {

		variable := ast.New(ast.VARIABLE, "Variable")

		variable.AddChild(p.parseType())

		if p.thisToken.Type == lexer.IDENT {
			variable.AddChild(p.parseIdentifier())
		} else {
			p.printMissingMessage(lexer.IDENT)
		}

		variable.AddChild(p.parseArraySizeRept())
		fparams.AddChild(variable)

		if !p.skipErrors(ast.FPARAMSTAIL) {
			return nil
		}

		if p.thisToken.Type != lexer.COMMA {
			break
		} else {
			p.printProductions("<FPARAMSTAIL>")
			p.readToken()
		}
	}

	return fparams
}

// Parses the <FuncDef> ::= <Function> <FuncDef> | EPSILON Production.
func (p *Parser) parseFuncDef() *ast.Node {
	funcDef := ast.New(ast.FUNCDEF, "")

	for p.thisToken.Type == lexer.FUNCTION {
		funcDef.AddChild(p.parseFunction())
	}

	if !funcDef.HasChildren() {
		return nil
	}
	return funcDef
}

func (p *Parser) parseFunction() *ast.Node {
	if !p.skipErrors(ast.FUNCTION) {
		return nil
	}

	functionNode := ast.New(ast.FUNCTION, "")
	p.printProductions("<FUNCTION>\n")

	head, headIden := p.parseFuncHead()
	functionNode.AddChild(head)
	p.printProductions("\n")

	funcBody, table := p.parseFuncBody(headIden)

	if funcBody != nil {
		functionNode.AddChild(funcBody)
	}

	p.printProductions("\n")

	// Link Table to Function Definition && FuncHead
	if table != nil {
		p.table.LinkRowToTable(headIden, ast.FUNCTION, table)
		table.Definition = functionNode
		table.ScopeHead = head
	}

	return functionNode
}

// Parses the
// <FuncHead> ::= 'func' 'id' <ClassMethod> '(' <FParams> ')”:'<FuncDeclTail>.
func (p *Parser) parseFuncHead() (*ast.Node, string) {
	if !p.skipErrors(ast.FUNCTIONHEAD) {
		return nil, ""
	}

	functionHead := ast.New(ast.FUNCTIONHEAD, "")
	functionHead.Line = p.getLineNumber()

	p.printProductions("<FUNHEAD>")

	if p.thisToken.Type == lexer.FUNCTION {
		p.readToken()
	}

	var iden *ast.Node
	class := p.parseClassMethod()

	if p.thisToken.Type == lexer.IDENT {
		iden = p.parseIdentifier()
	} else {
		p.printMissingMessage(lexer.IDENT)
	}

	functionHead.AddChild(iden)

	functionHead.AddChild(class)

	if p.thisToken.Type == lexer.LPAREN {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.LPAREN)
	}

	fparams := p.parseFParams()
	functionHead.AddChild(fparams)

	if p.thisToken.Type == lexer.RPAREN {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.RPAREN)
	}

	if p.thisToken.Type == lexer.COLON {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.COLON)
	}

	var return_type *ast.Node = p.parseFuncDeclTail()
	functionHead.AddChild(return_type)

	p.table.CreateRow(iden.Lexeme, ast.FUNCTION)
	p.table.AddReturnType(iden.Lexeme, return_type)
	p.table.AddFParams(iden.Lexeme, class, fparams)
	p.table.AddClassMethod(iden.Lexeme, class)

	return functionHead, iden.Lexeme
}

// Parses the <FuncBody> ::= '{' <MethodBodyVar> <StatementList> '}' Prod.
func (p *Parser) parseFuncBody(scope string) (*ast.Node, *records.SymbolTable) {
	if !p.skipErrors(ast.FUNCTIONBODY) {
		return nil, nil
	}

	p.printProductions("<FUNCBODY>")
	functionBody := ast.New(ast.FUNCTIONBODY, "")

	if p.thisToken.Type == lexer.LBRACE {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.LBRACE)
		p.readTillSymbol(lexer.RBRACE)
		return nil, nil
	}

	methodBodyVar, table := p.parseMethodBodyVar(scope)
	functionBody.AddChild(methodBodyVar)
	functionBody.AddChild(p.parseStatementList())

	if p.thisToken.Type == lexer.RBRACE {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.RBRACE)
	}

	p.printProductions("\n")
	return functionBody, table
}

// Parses the <MethodBodyVar> ::= 'var' '{' <VarDeclRep> '}'  | EPSILON Prod.
func (p *Parser) parseMethodBodyVar(scope string) (*ast.Node, *records.SymbolTable) {
	if !p.skipErrors(ast.METHODBODYVAR) {
		return nil, nil
	}

	methodTable := records.NewTable(scope)
	if scope == "MAIN" {
		methodTable.CreateRow("&__buffer__&", ast.VARIABLE)
		methodTable.AddArrayType("&__buffer__&", ast.New(ast.ARRAYSIZE, ""), ast.New(ast.TYPE, "integer"))
	}

	p.printProductions("<METHODBODYVAR>")
	if p.thisToken.Type == lexer.VAR {
		p.readToken()
	} else {
		return nil, records.NewTable(scope)
	}

	if p.thisToken.Type == lexer.LBRACE {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.LBRACE)
	}

	variableDeclReps := ast.New(ast.VARIABLEDECL, "")

	for {

		p.printProductions("<VARDECLREP>")
		variable := ast.New(ast.VARIABLE, "Variable")

		var var_type *ast.Node

		if !p.tokenIsType() {
			break
		} else {
			var_type = p.parseType()
		}

		variable.AddChild(var_type)
		variable.Line = p.thisToken.Line

		var iden *ast.Node

		if p.thisToken.Type == lexer.IDENT {
			iden = p.parseIdentifier()
		} else {

			p.printMissingMessage(lexer.IDENT)

			if !p.skipErrors(ast.METHODBODYVAR) {
				return nil, methodTable
			}
		}

		variable.AddChild(iden)

		var arraysize *ast.Node = p.parseArraySizeRept()
		variable.AddChild(arraysize)

		variableDeclReps.AddChild(variable)

		// Add To Symbol Table
		correct := methodTable.CreateRow(iden.Lexeme, ast.VARIABLE)

		if !correct {
			p.err_free = correct
		}

		methodTable.AddArrayType(iden.Lexeme, arraysize, var_type)
		methodTable.SetVisibility(iden.Lexeme, ast.FUNCTION)

		if p.thisToken.Type != lexer.SEMICOLON && p.nextToken.Type == lexer.RBRACE {
			break
		} else if p.thisToken.Type != lexer.SEMICOLON {
			p.printMissingMessage(lexer.SEMICOLON)
		} else {
			p.readToken()
		}
	}

	if p.thisToken.Type == lexer.RBRACE {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.RBRACE)
	}

	if !variableDeclReps.HasChildren() {
		return nil, records.NewTable(scope)
	}

	return variableDeclReps, methodTable
}

// Parses the <FuncOrVar> ::= 'id' <FuncOrVarIdnest>.
func (p *Parser) parseFuncOrVar() *ast.Node {
	if !p.skipErrors(ast.FUNCORVAR) {
		return nil
	}

	p.printProductions("<FUNORVAR>")
	var funcOrVar *ast.Node

	if p.thisToken.Type == lexer.IDENT {
		funcOrVar = p.parseIdentifier()
	}

	left, right := p.parseFuncOrVarIdnest()

	funcOrVar.AddChild(left)
	funcOrVar.AddChild(right)

	return funcOrVar
}

// Parses the -->
// <FuncOrVarIdnest> ::= <IndiceRep> <FuncOrVarIdnestTail>
// <FuncOrVarIdnest> ::= '(' <AParams> ')' <FuncOrVarIdnestTail>
func (p *Parser) parseFuncOrVarIdnest() (*ast.Node, *ast.Node) {
	if !p.skipErrors(ast.FUNCORVARIDNEST) {
		return nil, nil
	}

	var left, right *ast.Node

	p.printProductions("<FUNORVARIDNEST")

	if p.thisToken.Type == lexer.LPAREN {

		p.readToken()
		left = p.parseAParams()

		if p.thisToken.Type == lexer.RPAREN {
			p.readToken()
		} else {
			p.printMissingMessage(lexer.RPAREN)
		}

		right = p.parseFuncOrVarIdnestTail()

	} else {
		left = p.parseIndiceRep()
		right = p.parseFuncOrVarIdnestTail()
	}

	return left, right
}

func (p *Parser) parseFuncOrVarIdnestTail() *ast.Node {
	if !p.skipErrors(ast.FUNCORVARIDNESTTAIL) {
		return nil
	}

	p.printProductions("<FUNCORVARIDNESTTAIL>")
	var id *ast.Node

	if p.thisToken.Type == lexer.DOT {

		p.readToken()

		if p.thisToken.Type == lexer.IDENT {
			id = (p.parseIdentifier())
		} else {
			p.printMissingMessage(lexer.IDENT)
		}

		tmp1, tmp2 := (p.parseFuncOrVarIdnest())

		id.AddChild(tmp1)
		id.AddChild(tmp2)

		return id
	} else {
		return nil
	}
}
