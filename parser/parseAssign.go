package parser

import (
	"aHobeychi/GoCompiler/ast"
	"aHobeychi/GoCompiler/lexer"
)

// Func for the <FuncOrAssignStat> ::= 'id' <FuncOrAssignStatIdnest> Prod.
func (p *Parser) parseFuncOrAssignStat() *ast.Node {

	if !p.skipErrors(ast.FUNCORASSIGNSTAT) {
		return nil
	}

	p.printProductions("<FUNCORASSIGNSTAT>")

	var leftId *ast.Node
	var assign *ast.Node

	if p.thisToken.Type == lexer.IDENT {
		leftId = p.parseIdentifier()
	} else {
		p.printMissingMessage(lexer.IDENT)
	}

	tmp1, tmp2 := (p.parseFuncOrAssignStatIdnest())

	leftId.AddChild(tmp1)
	leftId.AddChild(tmp2)

	if p.thisToken.Type == lexer.ASSIGN {
		assign = p.parseAssignOp()
		assign.AddChild(leftId)

		if !p.tokenIsExpression() {
			p.printMissingMessage(lexer.ASSIGN)
			p.readTillSymbol(lexer.SEMICOLON)
			return nil
		}

		assign.AddChild(p.parseExpression())
		return assign
	}

	return leftId
}

// Function for the productions:
// <FuncOrAssignStatIdnest> ::= <IndiceRep> <FuncOrAssignStatIdnestVarTail>
// <FuncOrAssignStatIdnest> ::= '(' <AParams> ')' <FuncOrAssignStatIdnestFuncTail>
func (p *Parser) parseFuncOrAssignStatIdnest() (*ast.Node, *ast.Node) {

	if !p.skipErrors(ast.FUNCORASSIGNSTATIDNEST) {
		return nil, nil
	}

	var left *ast.Node
	var right *ast.Node

	p.printProductions("<FUNCORASSIGNSTATIDNEST>")
	if p.thisToken.Type == lexer.LPAREN {

		p.readToken()

		left = p.parseAParams()

		if p.thisToken.Type == lexer.RPAREN {

			p.readToken()

		} else {

			p.printMissingMessage(lexer.RPAREN)
			p.readTillSymbol(lexer.SEMICOLON)

			return left, right
		}

		right = (p.parseFuncOrAssignStatIdnestFuncTail())

	} else if p.thisToken.Type == lexer.LBRACKET {

		left = p.parseIndiceRep()
		right = (p.parseFuncOrAssignStatIdnestVartail())

	} else {

		return p.parseFuncOrAssignStatIdnestVartail(), nil
	}

	return left, right
}

// Function for the productions:
// <FuncOrAssignStatIdnestFuncTail> ::= '.' 'id' <FuncStatTail> | EPSILON
func (p *Parser) parseFuncOrAssignStatIdnestFuncTail() *ast.Node {

	var id *ast.Node

	p.printProductions("<FUNCORASSIGNSTATIDNESTFUNCTAIL")

	if p.thisToken.Type == lexer.DOT {

		p.readToken()

	} else if p.thisToken.Type == lexer.IDENT {

		p.printProductions(".")

	} else {

		return nil
	}

	id = p.parseIdentifier()
	tmp1, tmp2 := p.parseFuncStatTail()

	id.AddChild(tmp1)
	id.AddChild(tmp2)

	return id
}

// Function for the productions:
// <FuncStatTail> ::= <IndiceRep> '.' 'id' <FuncStatTail>
// <FuncStatTail> ::= '(' <AParams> ')' <FuncStatTailIdnest>
func (p *Parser) parseFuncStatTail() (*ast.Node, *ast.Node) {

	var left *ast.Node
	var right *ast.Node

	if !p.skipErrors(ast.FUNCSTATTAIL) {
		return nil, nil
	}

	p.printProductions("PARSEFUNCSTATTAIL")

	if p.thisToken.Type == lexer.LBRACKET {

		left = p.parseIndiceRep()

		if p.thisToken.Type == lexer.DOT {

			p.readToken()

			id := p.parseIdentifier()
			tmp1, tmp2 := p.parseFuncStatTail()

			id.AddChild(tmp1)
			id.AddChild(tmp2)

			return left, id
		}

	} else if p.thisToken.Type == lexer.LPAREN {

		p.readToken()

		left = p.parseAParams()

		if p.thisToken.Type == lexer.RPAREN {

			p.readToken()

		} else {

			p.readTillSymbol(lexer.SEMICOLON)

			return left, left
		}

		right = p.parseFuncStatTailIdnest()
	}

	return left, right
}

// Function for the
// <FuncStatTailIdnest> ::= '.' 'id' <FuncStatTail> | EPSILON Production
func (p *Parser) parseFuncStatTailIdnest() *ast.Node {

	var tail *ast.Node

	p.printProductions("FUNCSTATTAILIDNEST")

	if p.thisToken.Type == lexer.DOT {

		p.readToken()

		tail = p.parseIdentifier()
		tmp1, tmp2 := p.parseFuncStatTail()

		tail.AddChild(tmp1)
		tail.AddChild(tmp2)

	} else {
		return nil
	}

	return tail
}

// Function to parser the
//       <FuncOrAssignStatIdnestVarTail> ::= '.' 'id' <FuncOrAssignStatIdnest>
//       | <AssignStatTail> Production
func (p *Parser) parseFuncOrAssignStatIdnestVartail() *ast.Node {

	var tail *ast.Node

	p.printProductions("<FUNCORASSIGNSTATIDNESTVARTAIL>")

	if p.thisToken.Type == lexer.DOT {

		p.readToken()

		tail = p.parseIdentifier()

		tmp1, tmp2 := (p.parseFuncOrAssignStatIdnest())

		tail.AddChild(tmp1)
		tail.AddChild(tmp2)

	} else {
		return nil
	}
	return tail
}
