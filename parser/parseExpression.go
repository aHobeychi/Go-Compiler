// File containing all the functions related to parsing expressions.
package parser

import (
	"aHobeychi/GoCompiler/ast"
	"aHobeychi/GoCompiler/lexer"
)

// Func for production <Expr> ::= <ArithExpr> <ExprTail>
//                     <ExprTail> ::= <RelOp> <ArithExpr> | EPSILON
func (p *Parser) parseExpression() *ast.Node {

	if !p.skipErrors(ast.EXPR) {
		return nil
	}

	var left *ast.Node
	var right *ast.Node

	p.printProductions("<EXPR>")

	left = p.parseArithExpr()
	left.Line = p.getLineNumber()

	if !p.skipErrors(ast.EXPRTAIL) {
		return nil
	}

	if p.tokenIsRelOp() {

		p.printProductions("<EXPRTAIL>")

		relOp := p.parseRelOp()
		relOp.Line = p.getLineNumber()
		right = p.parseArithExpr()

		relOp.AddChild(left)
		relOp.AddChild(right)
		return relOp
	}

	return left
}

// Func for production <ArithExpr> ::= <Term> <ArithExprTail>
// <ArithExprTail> ::= <AddOp> <Term> <ArithExprTail> | EPSILON
func (p *Parser) parseArithExpr() *ast.Node {

	var left *ast.Node

	if !p.skipErrors(ast.ARITHEXPR) {
		return nil
	}

	p.printProductions("ARITHEXPR")

	left = p.parseTerm()

	if !p.skipErrors(ast.ARITHTAIL) {
		return nil
	}

	if p.tokenIsAddOp() {

		p.printProductions("ARITHEXPRTAIL")

		addOp := p.parseAddOp()
		addOp.Line = p.getLineNumber()
		addOp.AddChild(left)
		addOp.AddChild(p.parseArithExpr())

		return addOp
	}

	return left
}

// Func for production <Term> ::= <Factor> <TermTail>
//                     <TermTail> ::= <MultOp> <Factor> <TermTail> | EPSILON
func (p *Parser) parseTerm() *ast.Node {

	var left *ast.Node

	if !p.skipErrors(ast.TERM) {
		return nil
	}

	p.printProductions("TERM")
	left = p.parseFactor()

	if !p.skipErrors(ast.TERMTAIL) {
		return nil
	}

	if p.tokenIsMultOp() {

		p.printProductions("TERMTAIL")

		multOp := p.parseMultOp()
		multOp.AddChild(left)

		multOp.AddChild(p.parseTerm())
		multOp.Line = p.getLineNumber()

		return multOp
	}

	return left
}

// Func for the following productions
//      <Factor> ::= <FuncOrVar>
//      <Factor> ::= 'intnum'
//      <Factor> ::= 'floatnum'
//      <Factor> ::= 'stringlit'
//      <Factor> ::= '(' <Expr> ')'
//      <Factor> ::= 'not' <Factor>
//      <Factor> ::= <Sign> <Factor>
//      <Factor> ::= 'qm' '[' <Expr> ':' <Expr> ':' <Expr> ']'
func (p *Parser) parseFactor() *ast.Node {

	var factor *ast.Node
	p.printProductions("FACTOR")

	if !p.skipErrors(ast.FACTOR) {
		return nil
	}

	switch p.thisToken.Type {
	case lexer.INT_VALUE:
		factor = p.parseIntNum()
	case lexer.FLOAT_VALUE:
		factor = p.parseFloatNum()
	case lexer.STRINGLIT:
		factor = p.parseStringLit()
	case lexer.LPAREN:
		return p.parseFactorExpr()
	case lexer.NOT:
		return p.parseFactorNot()
	case lexer.PLUS, lexer.MINUS:
		return p.parseSignFactor()
	case lexer.QUESTION:
		return p.parseTernary()
	case lexer.IDENT:
		return p.parseFuncOrVar()
	}

	return factor
}

// Parses factor containing an expression
func (p *Parser) parseFactorExpr() *ast.Node {

	if p.thisToken.Type == lexer.LPAREN {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.LPAREN)
	}

	expr := p.parseExpression()
	expr.Line = p.getLineNumber()

	if p.thisToken.Type == lexer.RPAREN {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.RPAREN)
	}

	return expr
}

// Parses factor containing a Not expression
func (p *Parser) parseFactorNot() *ast.Node {

	var sign *ast.Node

	if p.thisToken.Type == lexer.NOT {
		sign = ast.New(ast.NOT, "")
		sign.Line = p.getLineNumber()
		p.readToken()
	}

	sign.AddChild(p.parseFactor())
	sign.Line = p.getLineNumber()
	return sign
}

// Function to parse signed factor
// <Factor> ::= <Sign> <Factor>
func (p *Parser) parseSignFactor() *ast.Node {

	var sign *ast.Node
	if p.tokenIsSign() {
		sign = p.parseSign()
	}

	sign.AddChild(p.parseFactor())
	sign.Line = p.getLineNumber()

	return sign
}

// Function to parse the ternary operator
// <Factor> ::= 'qm' '[' <Expr> ':' <Expr> ':' <Expr> ']'
func (p *Parser) parseTernary() *ast.Node {

	ternary := ast.New(ast.TERNARY, "?")
	ternary.Line = p.getLineNumber()

	if p.thisToken.Type == lexer.QUESTION {
		p.readToken()
	}

	if p.thisToken.Type == lexer.LBRACKET {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.LBRACKET)
	}

	ternary.AddChild(p.parseExpression())

	if p.thisToken.Type == lexer.COLON {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.COLON)
	}

	ternary.AddChild(p.parseExpression())

	if p.thisToken.Type == lexer.COLON {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.COLON)
	}

	ternary.AddChild(p.parseExpression())

	if p.thisToken.Type == lexer.RBRACKET {
		p.readToken()
	} else {
		p.printMissingMessage(lexer.RBRACKET)
	}

	return ternary
}
