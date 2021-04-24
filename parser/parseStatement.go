// File Containing all functions related to parsing statements.
package parser

import (
	"aHobeychi/GoCompiler/ast"
	"aHobeychi/GoCompiler/lexer"
)

// Parses the Statement lists -- >
// <StatementList> ::= <Statement> <StatementList> | EPSILON .
func (p *Parser) parseStatementList() *ast.Node {

	if !p.skipErrors(ast.STATEMENTLIST) {
		return nil
	}

	p.printProductions("\n<STATEMENTLIST>")
	statementList := ast.New(ast.STATEMENT, "")

	for {
		if !p.tokenIsStatement() {
			break
		}
		statementList.AddChild(p.parseStatement())

	}

	if statementList.NumberOfChildren() == 0 {
		return nil
	}

	return statementList
}

// Parses The Stat Block
// <StatBlock> ::= '{' <StatementList> '}'
// <StatBlock> ::= <Statement>
// <StatBlock> ::= EPSILON
func (p *Parser) parseStatBlock() *ast.Node {

	if !p.skipErrors(ast.STATEMENTBLOCK) {
		return nil
	}

	p.printProductions("\n<STATBLOCK>")

	if p.thisToken.Type == lexer.LBRACE {

		p.readToken()
		statList := (p.parseStatementList())

		if p.thisToken.Type == lexer.RBRACE {
			p.readToken()
		} else {
			p.printMissingMessage(lexer.RBRACE)
		}

		p.printProductions("\n")

		return statList

	}

	return p.parseStatementList()
}

// Parses the Statement Productions
// <Statement> ::= <FuncOrAssignStat> ';'
// <Statement> ::= 'if' '(' <Expr> ')' 'then' <StatBlock> 'else' <StatBlock> ';'
// <Statement> ::= 'while' '(' <Expr> ')' <StatBlock> ';'
// <Statement> ::= 'read' '(' <Variable> ')' ';'
// <Statement> ::= 'write' '(' <Expr> ')' ';'
// <Statement> ::= 'return' '(' <Expr> ')' ';'
// <Statement> ::= 'break' ';'
// <Statement> ::= 'continue' ';'
func (p *Parser) parseStatement() *ast.Node {

	if !p.skipErrors(ast.STATEMENT) {
		return nil
	}

	p.printProductions("<STATEMENT>")
	var statement *ast.Node

	switch p.thisToken.Type {
	case lexer.IF:
		statement = p.parseIfStatement()
	case lexer.WHILE:
		statement = p.parseWhileStatement()
	case lexer.READ:
		statement = p.parseReadStatement()
	case lexer.WRITE:
		statement = p.parseWriteStatement()
	case lexer.RETURN:
		statement = p.parseReturnStatement()
	case lexer.BREAK:
		statement = p.parseBreakStatement()
	case lexer.CONTINUE:
		statement = p.parseContinueStatement()
	case lexer.IDENT:
		statement = p.parseFuncOrAssignStat()
	default:
		return nil
	}

	if p.thisToken.Type == lexer.SEMICOLON {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.SEMICOLON)
	}

	return statement
}

// Parses the if statements.
func (p *Parser) parseIfStatement() *ast.Node {

	p.printProductions("<IF>")
	ifStatement := ast.New(ast.IF, "")

	if p.thisToken.Type == lexer.IF {
		p.readToken()
	}

	if p.thisToken.Type == lexer.LPAREN {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.LPAREN)
	}

	ifStatement.AddChild(p.parseExpression())

	if p.thisToken.Type == lexer.RPAREN {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.RPAREN)
	}

	if p.thisToken.Type == lexer.THEN {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.THEN)

		if !p.skipErrors(ast.STATEMENT) {
			return nil
		}
	}

	ifStatement.AddChild(p.parseStatBlock())

	if p.thisToken.Type == lexer.ELSE {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.ELSE)

		if !p.skipErrors(ast.STATEMENT) {
			return nil
		}
	}

	if p.thisToken.Type == lexer.SEMICOLON {
		return ifStatement
	}

	ifStatement.AddChild(p.parseStatBlock())

	return ifStatement
}

// Parses the while statements.
func (p *Parser) parseWhileStatement() *ast.Node {

	p.printProductions("<WHILE>")
	whileStatement := ast.New(ast.WHILE, "while")

	if p.thisToken.Type == lexer.WHILE {
		p.readToken()
	}

	if p.thisToken.Type == lexer.LPAREN {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.LPAREN)
	}

	whileStatement.AddChild(p.parseExpression())

	if p.thisToken.Type == lexer.RPAREN {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.RPAREN)
	}

	whileStatement.AddChild(p.parseStatBlock())

	return whileStatement

}

// Parses the Read statements.
func (p *Parser) parseReadStatement() *ast.Node {

	p.printProductions("<READ>")
	readStatement := ast.New(ast.READ, "read")

	if p.thisToken.Type == lexer.READ {
		p.readToken()
	}

	if p.thisToken.Type == lexer.LPAREN {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.LPAREN)
	}

	readStatement.AddChild(p.parseVariable())

	if p.thisToken.Type == lexer.RPAREN {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.RPAREN)
	}

	return readStatement
}

// Parses the Write Statement.
func (p *Parser) parseWriteStatement() *ast.Node {

	p.printProductions("<WRITE>")
	writeStatement := ast.New(ast.WRITE, "write")

	if p.thisToken.Type == lexer.WRITE {
		p.readToken()
	}

	if p.thisToken.Type == lexer.LPAREN {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.LPAREN)
	}

	writeStatement.AddChild(p.parseExpression())

	if p.thisToken.Type == lexer.RPAREN {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.RPAREN)
	}

	return writeStatement
}

// Parses the Return Statement.
func (p *Parser) parseReturnStatement() *ast.Node {

	returnStatement := ast.New(ast.RETURN, "")
	p.printProductions("<RETURN>")

	if p.thisToken.Type == lexer.RETURN {
		p.readToken()
	}

	if p.thisToken.Type == lexer.LPAREN {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.LPAREN)
	}

	returnStatement.AddChild(p.parseExpression())

	if p.thisToken.Type == lexer.RPAREN {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.RPAREN)
	}

	return returnStatement
}

// Parses the Break Statement.
func (p *Parser) parseBreakStatement() *ast.Node {

	p.printProductions("<BREAk>")

	p.readToken()

	break_node := ast.New(ast.BREAK, "break")
	break_node.Line = p.getLineNumber()

	return break_node
}

// Parses the Continue Statement.
func (p *Parser) parseContinueStatement() *ast.Node {

	p.readToken()
	continue_node := ast.New(ast.CONTINUE, "continue")
	continue_node.Line = p.getLineNumber()

	return continue_node
}

// Parses the <AParams> ::= <Expr> <AParamsTail> | EPSILON Production.
func (p *Parser) parseAParams() *ast.Node {

	if !p.skipErrors(ast.APARAMS) {
		return nil
	}

	aParams := ast.New(ast.APARAMS, "")
	p.printProductions("APARAMS: ")

	if !(p.tokenIsExpression()) {
		return aParams
	}

	for {

		if !p.tokenIsExpression() {
			break
		}

		aParams.AddChild(p.parseExpression())

		if !p.skipErrors(ast.APARAMSTAIL) {
			return nil
		}

		if p.thisToken.Type == lexer.COMMA {
			p.readToken()
		} else {
			break
		}
	}

	return aParams
}
