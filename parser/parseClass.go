// File containing all the functions related to parsing classes.
package parser

import (
	"aHobeychi/GoCompiler/ast"
	"aHobeychi/GoCompiler/lexer"
	"aHobeychi/GoCompiler/records"
)

// Parses the -->
// <ClassDecl> ::= 'class' 'id' <Inherit> '{' <ClassDeclBody> '}' ';' <ClassDecl>
//
//	| EPSILON
func (p *Parser) parseClassDecl() *ast.Node {
	p.printProductions("CLASSDECL\n")

	classDecls := ast.New(ast.CLASSDECL, "")

	for {

		if !p.skipErrors(ast.CLASSDECL) {
			return nil
		}

		classDecl := ast.New(ast.CLASS, "")

		if p.thisToken.Type == lexer.CLASS {
			p.readToken()
		} else {
			break
		}

		var iden *ast.Node

		if p.thisToken.Type == lexer.IDENT {
			iden = p.parseIdentifier()
			classDecl.AddChild(iden)
		} else {
			p.printMissingMessage(lexer.IDENT)
		}

		inherits := p.parseInherit()

		if inherits == nil {
			continue
		}

		classDecl.AddChild(inherits)

		if p.thisToken.Type == lexer.LBRACE {
			p.readToken()
		} else {
			p.printMissingMessage(lexer.LBRACE)
		}

		classDeclBody, table := p.parseClassDeclBody(iden.Lexeme)
		classDecl.AddChild(classDeclBody)

		if p.thisToken.Type == lexer.RBRACE {
			p.readToken()
		} else {
			p.printMissingMessage(lexer.RBRACE)
		}

		classDecls.AddChild(classDecl)

		if p.thisToken.Type == lexer.SEMICOLON {
			p.readToken()
		} else {
			p.printMissingMessage(lexer.SEMICOLON)
		}

		// Create Table
		p.table.CreateRow(iden.Lexeme, ast.CLASS)
		p.table.AddInherits(iden.Lexeme, inherits)
		p.table.LinkRowToTable(iden.Lexeme, ast.CLASS, table)
	}

	if !classDecls.HasChildren() {
		return nil
	}
	p.printProductions("\n")
	return classDecls
}

// Parses the <ClassDeclBody> ::= <Visibility> <MemberDecl> <ClassDeclBody>
//
//	| EPSILON
func (p *Parser) parseClassDeclBody(classType string) (*ast.Node, *records.SymbolTable) {
	if !p.skipErrors(ast.CLASSDECLBODY) {
		return nil, nil
	}

	classMembers := ast.New(ast.CLASSMEMBERDECL, "")
	p.printProductions("ClassDeclBody")
	classBodyTable := records.NewTable(classType)

	for p.thisToken.Type == lexer.PUBLIC || p.thisToken.Type == lexer.PRIVATE ||
		p.tokenIsType() || p.thisToken.Type == lexer.FUNCTION {

		vis := p.parseVisibility()
		memdecl := p.parseMemberDecl()

		classMember := ast.New(ast.CLASSMEMBER, "")
		classMember.AddChild(vis)
		classMember.AddChild(memdecl)
		classMembers.AddChild(classMember)

		classBodyTable.AddMemberDecl(classMember)
	}

	p.printProductions("\n")
	return classMembers, classBodyTable
}

// Parses the <MemberDecl> ::= <FuncDecl> | <MemberDecl> ::= <VarDecl> Prod.
func (p *Parser) parseMemberDecl() *ast.Node {
	if !p.skipErrors(ast.MEMBERDECL) {
		return nil
	}

	if p.tokenIsType() {
		return p.parseVarDecl()
	} else {
		return p.parseFuncDecl()
	}
}

// Parses the <VarDecl> ::= <Type> 'id' <ArraySizeRept> ';' Prod.
func (p *Parser) parseVarDecl() *ast.Node {
	if !p.skipErrors(ast.VARIABLEDECL) {
		return nil
	}

	variable := ast.New(ast.VARIABLE, "Variable")

	if p.tokenIsType() {
		variable.AddChild(p.parseType())
	} else {
		p.printMissingMessage(lexer.IDENT)
	}

	if p.thisToken.Type == lexer.IDENT {
		variable.AddChild(p.parseIdentifier())
	} else {
		p.printMissingMessage(lexer.IDENT)
	}

	variable.AddChild(p.parseArraySizeRept())

	if p.thisToken.Type == lexer.SEMICOLON {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.SEMICOLON)
	}

	if !variable.HasChildren() {
		return nil
	}

	return variable
}

// Parses the Inherits Productions
// <Inherit> ::= 'inherits' 'id' <NestedId> | EPSILON
// <NestedId> ::= ',' 'id' <NestedId> | EPSILON
func (p *Parser) parseInherit() *ast.Node {
	var inherits *ast.Node
	p.printProductions("Inherits:")

	if !p.skipErrors(ast.INHERITS) {
		return nil
	}

	if p.thisToken.Type == lexer.INHERITS {
		p.readToken()
	} else {
		return ast.New(ast.INHERITS, "")
	}

	if p.thisToken.Type != lexer.IDENT {
		p.printMissingMessage(lexer.CLASS)
	}

	inherits = ast.New(ast.INHERITS, "")
	inherits.Line = p.getLineNumber()

	for p.thisToken.Type == lexer.IDENT {

		id := p.parseIdentifier()

		if !p.table.EntryExists(id.Lexeme) {
			print("Warning inheriting from undeclared class: " + id.Lexeme)
			print(", Line: " + p.getLine() + "\n")
		}

		inherits.AddChild(id)

		if !p.skipErrors(ast.NESTEDID) {
			return nil
		}

		if p.thisToken.Type != lexer.COMMA {
			break
		} else {

			p.readToken()
			p.printProductions(",")
		}

	}
	return inherits
}

// Parses the <ClassMethod> ::= 'sr' 'id' | EPSILON Prod.
func (p *Parser) parseClassMethod() *ast.Node {
	classMethod := ast.New(ast.CLASSMETHOD, "")
	p.printProductions("ClassMethod: ")

	if p.nextToken.Type != lexer.DCOLON {
		return classMethod
	}

	if p.thisToken.Type == lexer.IDENT {
		classMethod.AddChild(p.parseIdentifier())
	} else {
		p.printMissingMessage(lexer.CLASS)
	}

	switch p.thisToken.Type {
	case lexer.DCOLON:
		p.readToken()
	case lexer.IDENT:
		p.printMissingMessage(lexer.DCOLON)
	default:
		return ast.New(ast.CLASSMETHOD, "")
	}

	p.printProductions("\n")
	return classMethod
}
