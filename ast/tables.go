package ast

import (
	"aHobeychi/GoCompiler/lexer"
)

// Returns The first Set for specific productions.
func GenerateFirstSet() Set {

	var firstSet = map[Production]map[lexer.TokenType]bool{}

	// PROGRAM
	firstSet[PROGRAM] = map[lexer.TokenType]bool{}
	firstSet[PROGRAM][lexer.MAIN] = true
	firstSet[PROGRAM][lexer.CLASS] = true
	firstSet[PROGRAM][lexer.FUNCTION] = true

	// EXPR
	firstSet[EXPR] = map[lexer.TokenType]bool{}
	firstSet[EXPR][lexer.INT_VALUE] = true
	firstSet[EXPR][lexer.FLOAT_VALUE] = true
	firstSet[EXPR][lexer.STRINGLIT] = true
	firstSet[EXPR][lexer.LPAREN] = true
	firstSet[EXPR][lexer.NOT] = true
	firstSet[EXPR][lexer.QUESTION] = true
	firstSet[EXPR][lexer.IDENT] = true
	firstSet[EXPR][lexer.PLUS] = true
	firstSet[EXPR][lexer.MINUS] = true

	// EXPRTAIL
	firstSet[EXPRTAIL] = map[lexer.TokenType]bool{}
	firstSet[EXPRTAIL][lexer.EQUAL] = true
	firstSet[EXPRTAIL][lexer.NOTEQ] = true
	firstSet[EXPRTAIL][lexer.LT] = true
	firstSet[EXPRTAIL][lexer.GT] = true
	firstSet[EXPRTAIL][lexer.LTE] = true
	firstSet[EXPRTAIL][lexer.GTE] = true
	firstSet[EXPRTAIL][lexer.EMPTY] = true

	// ARITHEXPR
	firstSet[ARITHEXPR] = map[lexer.TokenType]bool{}
	firstSet[ARITHEXPR][lexer.INT_VALUE] = true
	firstSet[ARITHEXPR][lexer.FLOAT_VALUE] = true
	firstSet[ARITHEXPR][lexer.STRINGLIT] = true
	firstSet[ARITHEXPR][lexer.LPAREN] = true
	firstSet[ARITHEXPR][lexer.NOT] = true
	firstSet[ARITHEXPR][lexer.QUESTION] = true
	firstSet[ARITHEXPR][lexer.IDENT] = true
	firstSet[ARITHEXPR][lexer.PLUS] = true
	firstSet[ARITHEXPR][lexer.MINUS] = true

	// ARITHTAIL
	firstSet[ARITHTAIL] = map[lexer.TokenType]bool{}
	firstSet[ARITHTAIL][lexer.PLUS] = true
	firstSet[ARITHTAIL][lexer.MINUS] = true
	firstSet[ARITHTAIL][lexer.OR] = true
	firstSet[ARITHTAIL][lexer.EMPTY] = true

	// FACTOR
	firstSet[FACTOR] = map[lexer.TokenType]bool{}
	firstSet[FACTOR][lexer.INT_VALUE] = true
	firstSet[FACTOR][lexer.FLOAT_VALUE] = true
	firstSet[FACTOR][lexer.STRINGLIT] = true
	firstSet[FACTOR][lexer.LPAREN] = true
	firstSet[FACTOR][lexer.NOT] = true
	firstSet[FACTOR][lexer.QUESTION] = true
	firstSet[FACTOR][lexer.IDENT] = true
	firstSet[FACTOR][lexer.PLUS] = true
	firstSet[FACTOR][lexer.MINUS] = true

	// ASSIGN
	firstSet[ASSIGN] = map[lexer.TokenType]bool{}
	firstSet[ASSIGN][lexer.ASSIGN] = true

	// TERM
	firstSet[TERM] = map[lexer.TokenType]bool{}
	firstSet[TERM][lexer.INT_VALUE] = true
	firstSet[TERM][lexer.FLOAT_VALUE] = true
	firstSet[TERM][lexer.STRINGLIT] = true
	firstSet[TERM][lexer.LPAREN] = true
	firstSet[TERM][lexer.NOT] = true
	firstSet[TERM][lexer.QUESTION] = true
	firstSet[TERM][lexer.IDENT] = true
	firstSet[TERM][lexer.PLUS] = true
	firstSet[TERM][lexer.MINUS] = true

	// TERMTAIL
	firstSet[TERMTAIL] = map[lexer.TokenType]bool{}
	firstSet[TERMTAIL][lexer.MUL] = true
	firstSet[TERMTAIL][lexer.DIV] = true
	firstSet[TERMTAIL][lexer.AND] = true
	firstSet[TERMTAIL][lexer.EMPTY] = true

	// INDICEREP
	firstSet[INDICEREP_PROD] = map[lexer.TokenType]bool{}
	firstSet[INDICEREP_PROD][lexer.LBRACKET] = true
	firstSet[INDICEREP_PROD][lexer.EMPTY] = true

	// Statement
	firstSet[STATEMENT] = map[lexer.TokenType]bool{}
	firstSet[STATEMENT][lexer.IF] = true
	firstSet[STATEMENT][lexer.WHILE] = true
	firstSet[STATEMENT][lexer.READ] = true
	firstSet[STATEMENT][lexer.WRITE] = true
	firstSet[STATEMENT][lexer.RETURN] = true
	firstSet[STATEMENT][lexer.BREAK] = true
	firstSet[STATEMENT][lexer.CONTINUE] = true
	firstSet[STATEMENT][lexer.IDENT] = true

	// StatementList
	firstSet[STATEMENTLIST] = map[lexer.TokenType]bool{}
	firstSet[STATEMENTLIST][lexer.IF] = true
	firstSet[STATEMENTLIST][lexer.WHILE] = true
	firstSet[STATEMENTLIST][lexer.READ] = true
	firstSet[STATEMENTLIST][lexer.WRITE] = true
	firstSet[STATEMENTLIST][lexer.RETURN] = true
	firstSet[STATEMENTLIST][lexer.BREAK] = true
	firstSet[STATEMENTLIST][lexer.CONTINUE] = true
	firstSet[STATEMENTLIST][lexer.IDENT] = true
	firstSet[STATEMENTLIST][lexer.EMPTY] = true

	// StatementBlock
	firstSet[STATEMENTBLOCK] = map[lexer.TokenType]bool{}
	firstSet[STATEMENTBLOCK][lexer.LBRACE] = true
	firstSet[STATEMENTBLOCK][lexer.IF] = true
	firstSet[STATEMENTBLOCK][lexer.WHILE] = true
	firstSet[STATEMENTBLOCK][lexer.READ] = true
	firstSet[STATEMENTBLOCK][lexer.WRITE] = true
	firstSet[STATEMENTBLOCK][lexer.RETURN] = true
	firstSet[STATEMENTBLOCK][lexer.BREAK] = true
	firstSet[STATEMENTBLOCK][lexer.CONTINUE] = true
	firstSet[STATEMENTBLOCK][lexer.IDENT] = true
	firstSet[STATEMENTBLOCK][lexer.EMPTY] = true

	// AParams
	firstSet[APARAMS] = map[lexer.TokenType]bool{}
	firstSet[APARAMS][lexer.INT_VALUE] = true
	firstSet[APARAMS][lexer.FLOAT_VALUE] = true
	firstSet[APARAMS][lexer.STRINGLIT] = true
	firstSet[APARAMS][lexer.LPAREN] = true
	firstSet[APARAMS][lexer.NOT] = true
	firstSet[APARAMS][lexer.QUESTION] = true
	firstSet[APARAMS][lexer.IDENT] = true
	firstSet[APARAMS][lexer.PLUS] = true
	firstSet[APARAMS][lexer.MINUS] = true
	firstSet[APARAMS][lexer.EMPTY] = true

	// AParamsTail
	firstSet[APARAMSTAIL] = map[lexer.TokenType]bool{}
	firstSet[APARAMSTAIL][lexer.COMMA] = true
	firstSet[APARAMSTAIL][lexer.EMPTY] = true

	// FParams
	firstSet[FPARAMS] = map[lexer.TokenType]bool{}
	firstSet[FPARAMS][lexer.INTEGER] = true
	firstSet[FPARAMS][lexer.FLOAT] = true
	firstSet[FPARAMS][lexer.STRING] = true
	firstSet[FPARAMS][lexer.IDENT] = true
	firstSet[FPARAMS][lexer.EMPTY] = true

	// FParamsTail
	firstSet[FPARAMSTAIL] = map[lexer.TokenType]bool{}
	firstSet[FPARAMSTAIL][lexer.COMMA] = true
	firstSet[FPARAMSTAIL][lexer.EMPTY] = true

	// Types
	firstSet[TYPE] = map[lexer.TokenType]bool{}
	firstSet[TYPE][lexer.FLOAT] = true
	firstSet[TYPE][lexer.STRING] = true
	firstSet[TYPE][lexer.INTEGER] = true
	firstSet[TYPE][lexer.IDENT] = true

	// ARRAY SIZE
	firstSet[ARRAYSIZE] = map[lexer.TokenType]bool{}
	firstSet[ARRAYSIZE][lexer.LBRACKET] = true
	firstSet[ARRAYSIZE][lexer.EMPTY] = true

	// MEMBERDECL
	firstSet[MEMBERDECL] = map[lexer.TokenType]bool{}
	firstSet[MEMBERDECL][lexer.FUNCTION] = true
	firstSet[MEMBERDECL][lexer.INTEGER] = true
	firstSet[MEMBERDECL][lexer.FLOAT] = true
	firstSet[MEMBERDECL][lexer.STRING] = true
	firstSet[MEMBERDECL][lexer.IDENT] = true

	// VARIABLEDECL
	firstSet[VARIABLEDECL] = map[lexer.TokenType]bool{}
	firstSet[VARIABLEDECL][lexer.INTEGER] = true
	firstSet[VARIABLEDECL][lexer.FLOAT] = true
	firstSet[VARIABLEDECL][lexer.STRING] = true
	firstSet[VARIABLEDECL][lexer.IDENT] = true

	// VARIABLEDECLREP
	firstSet[VARIABLEDECLREP] = map[lexer.TokenType]bool{}
	firstSet[VARIABLEDECLREP][lexer.INTEGER] = true
	firstSet[VARIABLEDECLREP][lexer.FLOAT] = true
	firstSet[VARIABLEDECLREP][lexer.STRING] = true
	firstSet[VARIABLEDECLREP][lexer.IDENT] = true
	firstSet[VARIABLEDECLREP][lexer.EMPTY] = true

	// VARIABLE
	firstSet[VARIABLE] = map[lexer.TokenType]bool{}
	firstSet[VARIABLE][lexer.IDENT] = true

	// VARIABLEIDNEST
	firstSet[VARIABLEIDNEST] = map[lexer.TokenType]bool{}
	firstSet[VARIABLEIDNEST][lexer.LBRACKET] = true
	firstSet[VARIABLEIDNEST][lexer.DOT] = true
	firstSet[VARIABLEIDNEST][lexer.EMPTY] = true

	// VARIABLEIDNESTTAIL
	firstSet[VARIABLEIDNESTTAIL] = map[lexer.TokenType]bool{}
	firstSet[VARIABLEIDNESTTAIL][lexer.DOT] = true
	firstSet[VARIABLEIDNESTTAIL][lexer.EMPTY] = true

	// FUNCTION
	firstSet[FUNCTION] = map[lexer.TokenType]bool{}
	firstSet[FUNCTION][lexer.FUNCTION] = true

	// FUNCDEF
	firstSet[FUNCDEF] = map[lexer.TokenType]bool{}
	firstSet[FUNCDEF][lexer.FUNCTION] = true
	firstSet[FUNCDEF][lexer.EMPTY] = true

	// FUNCBODY
	firstSet[FUNCTIONBODY] = map[lexer.TokenType]bool{}
	firstSet[FUNCTIONBODY][lexer.LBRACE] = true

	// FUNCDEF
	firstSet[FUNCTIONHEAD] = map[lexer.TokenType]bool{}
	firstSet[FUNCTIONHEAD][lexer.FUNCTION] = true

	// FUNCDECL
	firstSet[FUNCDECL] = map[lexer.TokenType]bool{}
	firstSet[FUNCDECL][lexer.FUNCTION] = true

	// FUNCDECLTAIL
	firstSet[FUNCDECLTAIL] = map[lexer.TokenType]bool{}
	firstSet[FUNCDECLTAIL][lexer.VOID] = true
	firstSet[FUNCDECLTAIL][lexer.INTEGER] = true
	firstSet[FUNCDECLTAIL][lexer.FLOAT] = true
	firstSet[FUNCDECLTAIL][lexer.STRING] = true
	firstSet[FUNCDECLTAIL][lexer.IDENT] = true

	// FUNCORVAR
	firstSet[FUNCORVAR] = map[lexer.TokenType]bool{}
	firstSet[FUNCORVAR][lexer.IDENT] = true

	// FUNCORASSIGNASTAT
	firstSet[FUNCORASSIGNSTAT] = map[lexer.TokenType]bool{}
	firstSet[FUNCORASSIGNSTAT][lexer.IDENT] = true

	// FUNCORASSIGNASTATIDNEST
	firstSet[FUNCORASSIGNSTATIDNEST] = map[lexer.TokenType]bool{}
	firstSet[FUNCORASSIGNSTATIDNEST][lexer.LPAREN] = true
	firstSet[FUNCORASSIGNSTATIDNEST][lexer.LBRACKET] = true
	firstSet[FUNCORASSIGNSTATIDNEST][lexer.DOT] = true
	firstSet[FUNCORASSIGNSTATIDNEST][lexer.ASSIGN] = true

	// FUNORVARIDNEST
	firstSet[FUNCORVARIDNEST] = map[lexer.TokenType]bool{}
	firstSet[FUNCORVARIDNEST][lexer.LPAREN] = true
	firstSet[FUNCORVARIDNEST][lexer.DOT] = true
	firstSet[FUNCORVARIDNEST][lexer.LBRACKET] = true
	firstSet[FUNCORVARIDNEST][lexer.EMPTY] = true

	// FUNORVARIDNESTTAIL
	firstSet[FUNCORVARIDNESTTAIL] = map[lexer.TokenType]bool{}
	firstSet[FUNCORVARIDNESTTAIL][lexer.DOT] = true
	firstSet[FUNCORVARIDNESTTAIL][lexer.EMPTY] = true

	// FuncStatTail
	firstSet[FUNCSTATTAIL] = map[lexer.TokenType]bool{}
	firstSet[FUNCSTATTAIL][lexer.DOT] = true
	firstSet[FUNCSTATTAIL][lexer.LPAREN] = true
	firstSet[FUNCSTATTAIL][lexer.LBRACE] = true

	// INHERITS
	firstSet[INHERITS] = map[lexer.TokenType]bool{}
	firstSet[INHERITS][lexer.INHERITS] = true
	firstSet[INHERITS][lexer.EMPTY] = true

	// NESTEDID
	firstSet[NESTEDID] = map[lexer.TokenType]bool{}
	firstSet[NESTEDID][lexer.COMMA] = true
	firstSet[NESTEDID][lexer.EMPTY] = true

	// CLASSDECL
	firstSet[CLASSDECL] = map[lexer.TokenType]bool{}
	firstSet[CLASSDECL][lexer.CLASS] = true
	firstSet[CLASSDECL][lexer.EMPTY] = true

	// CLASSDECLBODY
	firstSet[CLASSDECLBODY] = map[lexer.TokenType]bool{}
	firstSet[CLASSDECLBODY][lexer.PUBLIC] = true
	firstSet[CLASSDECLBODY][lexer.PRIVATE] = true
	firstSet[CLASSDECLBODY][lexer.FUNCTION] = true
	firstSet[CLASSDECLBODY][lexer.INTEGER] = true
	firstSet[CLASSDECLBODY][lexer.FLOAT] = true
	firstSet[CLASSDECLBODY][lexer.STRING] = true
	firstSet[CLASSDECLBODY][lexer.IDENT] = true
	firstSet[CLASSDECLBODY][lexer.EMPTY] = true

	// METHODBODYVAR
	firstSet[METHODBODYVAR] = map[lexer.TokenType]bool{}
	firstSet[METHODBODYVAR][lexer.VAR] = true
	firstSet[METHODBODYVAR][lexer.EMPTY] = true

	// VISIBILITY
	firstSet[VISIBILITY] = map[lexer.TokenType]bool{}
	firstSet[VISIBILITY][lexer.PUBLIC] = true
	firstSet[VISIBILITY][lexer.PRIVATE] = true
	firstSet[VISIBILITY][lexer.EMPTY] = true

	// ADDOP
	firstSet[PLUS] = map[lexer.TokenType]bool{}
	firstSet[PLUS][lexer.PLUS] = true
	firstSet[PLUS][lexer.MINUS] = true
	firstSet[PLUS][lexer.OR] = true

	// ADDOP
	firstSet[MUL] = map[lexer.TokenType]bool{}
	firstSet[MUL][lexer.MUL] = true
	firstSet[MUL][lexer.DIV] = true
	firstSet[MUL][lexer.AND] = true

	// RELOP
	firstSet[LT] = map[lexer.TokenType]bool{}
	firstSet[LT][lexer.EQUAL] = true
	firstSet[LT][lexer.NOTEQ] = true
	firstSet[LT][lexer.LT] = true
	firstSet[LT][lexer.GT] = true
	firstSet[LT][lexer.LTE] = true
	firstSet[LT][lexer.GTE] = true

	// SIGN
	firstSet[SIGN] = map[lexer.TokenType]bool{}
	firstSet[SIGN][lexer.PLUS] = true
	firstSet[SIGN][lexer.MINUS] = true

	return firstSet
}

// Returns the follow set for specific productions
func GenerateFollowSet() Set {

	var followSet = map[Production]map[lexer.TokenType]bool{}

	// START
	followSet[PROGRAM] = map[lexer.TokenType]bool{}

	// EXPR
	followSet[EXPR] = map[lexer.TokenType]bool{}
	followSet[EXPR][lexer.SEMICOLON] = true
	followSet[EXPR][lexer.COMMA] = true
	followSet[EXPR][lexer.COLON] = true
	followSet[EXPR][lexer.RBRACKET] = true
	followSet[EXPR][lexer.RPAREN] = true

	// EXPRTAIL
	followSet[EXPRTAIL] = map[lexer.TokenType]bool{}
	followSet[EXPRTAIL][lexer.SEMICOLON] = true
	followSet[EXPRTAIL][lexer.COMMA] = true
	followSet[EXPRTAIL][lexer.COLON] = true
	followSet[EXPRTAIL][lexer.RBRACKET] = true
	followSet[EXPRTAIL][lexer.RPAREN] = true

	// ARITHEXPR
	followSet[ARITHEXPR] = map[lexer.TokenType]bool{}
	followSet[ARITHEXPR][lexer.SEMICOLON] = true
	followSet[ARITHEXPR][lexer.EQUAL] = true
	followSet[ARITHEXPR][lexer.NOTEQ] = true
	followSet[ARITHEXPR][lexer.LT] = true
	followSet[ARITHEXPR][lexer.GT] = true
	followSet[ARITHEXPR][lexer.LTE] = true
	followSet[ARITHEXPR][lexer.GTE] = true
	followSet[ARITHEXPR][lexer.COMMA] = true
	followSet[ARITHEXPR][lexer.COLON] = true
	followSet[ARITHEXPR][lexer.RBRACKET] = true
	followSet[ARITHEXPR][lexer.RPAREN] = true

	// ARITHTAIL
	followSet[ARITHTAIL] = map[lexer.TokenType]bool{}
	followSet[ARITHTAIL][lexer.SEMICOLON] = true
	followSet[ARITHTAIL][lexer.EQUAL] = true
	followSet[ARITHTAIL][lexer.NOTEQ] = true
	followSet[ARITHTAIL][lexer.LT] = true
	followSet[ARITHTAIL][lexer.GT] = true
	followSet[ARITHTAIL][lexer.LTE] = true
	followSet[ARITHTAIL][lexer.GTE] = true
	followSet[ARITHTAIL][lexer.COMMA] = true
	followSet[ARITHTAIL][lexer.COLON] = true
	followSet[ARITHTAIL][lexer.RBRACKET] = true
	followSet[ARITHTAIL][lexer.RPAREN] = true

	// FACTOR
	followSet[FACTOR] = map[lexer.TokenType]bool{}
	followSet[FACTOR][lexer.MUL] = true
	followSet[FACTOR][lexer.DIV] = true
	followSet[FACTOR][lexer.AND] = true
	followSet[FACTOR][lexer.SEMICOLON] = true
	followSet[FACTOR][lexer.EQUAL] = true
	followSet[FACTOR][lexer.NOTEQ] = true
	followSet[FACTOR][lexer.LT] = true
	followSet[FACTOR][lexer.GT] = true
	followSet[FACTOR][lexer.LTE] = true
	followSet[FACTOR][lexer.GTE] = true
	followSet[FACTOR][lexer.PLUS] = true
	followSet[FACTOR][lexer.MINUS] = true
	followSet[FACTOR][lexer.OR] = true
	followSet[FACTOR][lexer.COMMA] = true
	followSet[FACTOR][lexer.COLON] = true
	followSet[FACTOR][lexer.RBRACKET] = true
	followSet[FACTOR][lexer.RPAREN] = true

	// ASSIGNOP
	followSet[ASSIGN] = map[lexer.TokenType]bool{}
	followSet[ASSIGN][lexer.INT_VALUE] = true
	followSet[ASSIGN][lexer.FLOAT_VALUE] = true
	followSet[ASSIGN][lexer.STRINGLIT] = true
	followSet[ASSIGN][lexer.LPAREN] = true
	followSet[ASSIGN][lexer.NOT] = true
	followSet[ASSIGN][lexer.QUESTION] = true
	followSet[ASSIGN][lexer.IDENT] = true
	followSet[ASSIGN][lexer.PLUS] = true
	followSet[ASSIGN][lexer.MINUS] = true

	// TERM
	followSet[TERM] = map[lexer.TokenType]bool{}
	followSet[TERM][lexer.SEMICOLON] = true
	followSet[TERM][lexer.EQUAL] = true
	followSet[TERM][lexer.NOTEQ] = true
	followSet[TERM][lexer.LT] = true
	followSet[TERM][lexer.LTE] = true
	followSet[TERM][lexer.GT] = true
	followSet[TERM][lexer.GTE] = true
	followSet[TERM][lexer.PLUS] = true
	followSet[TERM][lexer.MINUS] = true
	followSet[TERM][lexer.OR] = true
	followSet[TERM][lexer.COMMA] = true
	followSet[TERM][lexer.COLON] = true
	followSet[TERM][lexer.RBRACKET] = true
	followSet[TERM][lexer.RPAREN] = true

	// TERMTAIL
	followSet[TERMTAIL] = map[lexer.TokenType]bool{}
	followSet[TERMTAIL][lexer.SEMICOLON] = true
	followSet[TERMTAIL][lexer.EQUAL] = true
	followSet[TERMTAIL][lexer.NOTEQ] = true
	followSet[TERMTAIL][lexer.LT] = true
	followSet[TERMTAIL][lexer.GT] = true
	followSet[TERMTAIL][lexer.LTE] = true
	followSet[TERMTAIL][lexer.GTE] = true
	followSet[TERMTAIL][lexer.PLUS] = true
	followSet[TERMTAIL][lexer.MINUS] = true
	followSet[TERMTAIL][lexer.OR] = true
	followSet[TERMTAIL][lexer.COMMA] = true
	followSet[TERMTAIL][lexer.COLON] = true
	followSet[TERMTAIL][lexer.RBRACKET] = true
	followSet[TERMTAIL][lexer.RPAREN] = true

	// INDICEREP
	followSet[INDICEREP_PROD] = map[lexer.TokenType]bool{}
	followSet[INDICEREP_PROD][lexer.MUL] = true
	followSet[INDICEREP_PROD][lexer.DIV] = true
	followSet[INDICEREP_PROD][lexer.AND] = true
	followSet[INDICEREP_PROD][lexer.SEMICOLON] = true
	followSet[INDICEREP_PROD][lexer.ASSIGN] = true
	followSet[INDICEREP_PROD][lexer.DOT] = true
	followSet[INDICEREP_PROD][lexer.EQUAL] = true
	followSet[INDICEREP_PROD][lexer.NOTEQ] = true
	followSet[INDICEREP_PROD][lexer.LT] = true
	followSet[INDICEREP_PROD][lexer.GT] = true
	followSet[INDICEREP_PROD][lexer.LTE] = true
	followSet[INDICEREP_PROD][lexer.GTE] = true
	followSet[INDICEREP_PROD][lexer.PLUS] = true
	followSet[INDICEREP_PROD][lexer.MINUS] = true
	followSet[INDICEREP_PROD][lexer.OR] = true
	followSet[INDICEREP_PROD][lexer.COMMA] = true
	followSet[INDICEREP_PROD][lexer.COLON] = true
	followSet[INDICEREP_PROD][lexer.RBRACKET] = true
	followSet[INDICEREP_PROD][lexer.RPAREN] = true

	// Statement
	followSet[STATEMENT] = map[lexer.TokenType]bool{}
	followSet[STATEMENT][lexer.IF] = true
	followSet[STATEMENT][lexer.WHILE] = true
	followSet[STATEMENT][lexer.READ] = true
	followSet[STATEMENT][lexer.WRITE] = true
	followSet[STATEMENT][lexer.RETURN] = true
	followSet[STATEMENT][lexer.BREAK] = true
	followSet[STATEMENT][lexer.CONTINUE] = true
	followSet[STATEMENT][lexer.IDENT] = true
	followSet[STATEMENT][lexer.ELSE] = true
	followSet[STATEMENT][lexer.SEMICOLON] = true
	followSet[STATEMENT][lexer.RBRACE] = true

	// StatementList
	followSet[STATEMENTLIST] = map[lexer.TokenType]bool{}
	followSet[STATEMENTLIST][lexer.RBRACE] = true

	// StatementBlock
	followSet[STATEMENTBLOCK] = map[lexer.TokenType]bool{}
	followSet[STATEMENTBLOCK][lexer.ELSE] = true
	followSet[STATEMENTBLOCK][lexer.SEMICOLON] = true

	//AParams
	followSet[APARAMS] = map[lexer.TokenType]bool{}
	followSet[APARAMS][lexer.RPAREN] = true

	//AParamsTail
	followSet[APARAMSTAIL] = map[lexer.TokenType]bool{}
	followSet[APARAMSTAIL][lexer.RPAREN] = true

	// FParams
	followSet[FPARAMS] = map[lexer.TokenType]bool{}
	followSet[FPARAMS][lexer.RPAREN] = true

	// FParamsTail
	followSet[FPARAMSTAIL] = map[lexer.TokenType]bool{}
	followSet[FPARAMSTAIL][lexer.RPAREN] = true

	// Types
	followSet[TYPE] = map[lexer.TokenType]bool{}
	followSet[TYPE][lexer.LBRACE] = true
	followSet[TYPE][lexer.SEMICOLON] = true
	followSet[TYPE][lexer.IDENT] = true

	// ARRAY SIZE
	followSet[ARRAYSIZE] = map[lexer.TokenType]bool{}
	followSet[ARRAYSIZE][lexer.RPAREN] = true
	followSet[ARRAYSIZE][lexer.COMMA] = true
	followSet[ARRAYSIZE][lexer.SEMICOLON] = true

	// MEMBERDECL
	followSet[MEMBERDECL] = map[lexer.TokenType]bool{}
	followSet[MEMBERDECL][lexer.PUBLIC] = true
	followSet[MEMBERDECL][lexer.PRIVATE] = true
	followSet[MEMBERDECL][lexer.FUNCTION] = true
	followSet[MEMBERDECL][lexer.INTEGER] = true
	followSet[MEMBERDECL][lexer.FLOAT] = true
	followSet[MEMBERDECL][lexer.STRING] = true
	followSet[MEMBERDECL][lexer.IDENT] = true
	followSet[MEMBERDECL][lexer.RBRACE] = true

	// VARIABLEDECL
	followSet[VARIABLEDECL] = map[lexer.TokenType]bool{}
	followSet[VARIABLEDECL][lexer.PUBLIC] = true
	followSet[VARIABLEDECL][lexer.PRIVATE] = true
	followSet[VARIABLEDECL][lexer.FUNCTION] = true
	followSet[VARIABLEDECL][lexer.INTEGER] = true
	followSet[VARIABLEDECL][lexer.FLOAT] = true
	followSet[VARIABLEDECL][lexer.STRING] = true
	followSet[VARIABLEDECL][lexer.IDENT] = true
	followSet[VARIABLEDECL][lexer.RBRACE] = true

	// VARIABLEDECLREP
	followSet[VARIABLEDECLREP] = map[lexer.TokenType]bool{}
	followSet[VARIABLEDECLREP][lexer.RBRACE] = true

	// VARIABLE
	followSet[VARIABLE] = map[lexer.TokenType]bool{}
	followSet[VARIABLE][lexer.RPAREN] = true

	// VARIABLEIDNEST
	followSet[VARIABLEIDNEST] = map[lexer.TokenType]bool{}
	followSet[VARIABLEIDNEST][lexer.RPAREN] = true

	// VARIABLEIDNESTTAIL
	followSet[VARIABLEIDNESTTAIL] = map[lexer.TokenType]bool{}
	followSet[VARIABLEIDNESTTAIL][lexer.RPAREN] = true

	// FUNCTIONHEAD
	followSet[FUNCDEF] = map[lexer.TokenType]bool{}
	followSet[FUNCDEF][lexer.LBRACE] = true

	// FUNCTIONDEF
	followSet[FUNCDEF] = map[lexer.TokenType]bool{}
	followSet[FUNCDEF][lexer.MAIN] = true

	// FUNCBODY
	followSet[FUNCTIONBODY] = map[lexer.TokenType]bool{}
	followSet[FUNCTIONBODY][lexer.MAIN] = true
	followSet[FUNCTIONBODY][lexer.FUNCTION] = true

	// FUNCTION
	followSet[FUNCTION] = map[lexer.TokenType]bool{}
	followSet[FUNCTION][lexer.MAIN] = true
	followSet[FUNCTION][lexer.FUNCTION] = true

	// FUNCDECL
	followSet[FUNCDECL] = map[lexer.TokenType]bool{}
	followSet[FUNCDECL][lexer.PUBLIC] = true
	followSet[FUNCDECL][lexer.PRIVATE] = true
	followSet[FUNCDECL][lexer.FUNCTION] = true
	followSet[FUNCDECL][lexer.INTEGER] = true
	followSet[FUNCDECL][lexer.FLOAT] = true
	followSet[FUNCDECL][lexer.STRING] = true
	followSet[FUNCDECL][lexer.IDENT] = true
	followSet[FUNCDECL][lexer.RBRACE] = true

	// FUNCDECLTAIL
	followSet[FUNCDECLTAIL] = map[lexer.TokenType]bool{}
	followSet[FUNCDECLTAIL][lexer.LBRACE] = true
	followSet[FUNCDECLTAIL][lexer.SEMICOLON] = true

	// FUNCORVAR
	followSet[FUNCORVAR] = map[lexer.TokenType]bool{}
	followSet[FUNCORVAR][lexer.MUL] = true
	followSet[FUNCORVAR][lexer.DIV] = true
	followSet[FUNCORVAR][lexer.AND] = true
	followSet[FUNCORVAR][lexer.SEMICOLON] = true
	followSet[FUNCORVAR][lexer.EQUAL] = true
	followSet[FUNCORVAR][lexer.NOTEQ] = true
	followSet[FUNCORVAR][lexer.LT] = true
	followSet[FUNCORVAR][lexer.GT] = true
	followSet[FUNCORVAR][lexer.LTE] = true
	followSet[FUNCORVAR][lexer.GTE] = true
	followSet[FUNCORVAR][lexer.PLUS] = true
	followSet[FUNCORVAR][lexer.MINUS] = true
	followSet[FUNCORVAR][lexer.OR] = true
	followSet[FUNCORVAR][lexer.COMMA] = true
	followSet[FUNCORVAR][lexer.COLON] = true
	followSet[FUNCORVAR][lexer.RBRACKET] = true
	followSet[FUNCORVAR][lexer.RPAREN] = true

	// FUNCORASSIGNASTAT
	followSet[FUNCORASSIGNSTAT] = map[lexer.TokenType]bool{}
	followSet[FUNCORASSIGNSTAT][lexer.SEMICOLON] = true

	// FUNCORASSIGNASTATIDNEST
	followSet[FUNCORASSIGNSTATIDNEST] = map[lexer.TokenType]bool{}
	followSet[FUNCORASSIGNSTATIDNEST][lexer.SEMICOLON] = true

	// FUNCORVAR
	followSet[FUNCORVARIDNEST] = map[lexer.TokenType]bool{}
	followSet[FUNCORVARIDNEST][lexer.MUL] = true
	followSet[FUNCORVARIDNEST][lexer.DIV] = true
	followSet[FUNCORVARIDNEST][lexer.AND] = true
	followSet[FUNCORVARIDNEST][lexer.SEMICOLON] = true
	followSet[FUNCORVARIDNEST][lexer.EQUAL] = true
	followSet[FUNCORVARIDNEST][lexer.NOTEQ] = true
	followSet[FUNCORVARIDNEST][lexer.LT] = true
	followSet[FUNCORVARIDNEST][lexer.GT] = true
	followSet[FUNCORVARIDNEST][lexer.LTE] = true
	followSet[FUNCORVARIDNEST][lexer.GTE] = true
	followSet[FUNCORVARIDNEST][lexer.PLUS] = true
	followSet[FUNCORVARIDNEST][lexer.MINUS] = true
	followSet[FUNCORVARIDNEST][lexer.OR] = true
	followSet[FUNCORVARIDNEST][lexer.COMMA] = true
	followSet[FUNCORVARIDNEST][lexer.COLON] = true
	followSet[FUNCORVARIDNEST][lexer.RBRACKET] = true
	followSet[FUNCORVARIDNEST][lexer.RPAREN] = true

	// FUNCORVAR
	followSet[FUNCORVARIDNESTTAIL] = map[lexer.TokenType]bool{}
	followSet[FUNCORVARIDNESTTAIL][lexer.MUL] = true
	followSet[FUNCORVARIDNESTTAIL][lexer.DIV] = true
	followSet[FUNCORVARIDNESTTAIL][lexer.AND] = true
	followSet[FUNCORVARIDNESTTAIL][lexer.SEMICOLON] = true
	followSet[FUNCORVARIDNESTTAIL][lexer.EQUAL] = true
	followSet[FUNCORVARIDNESTTAIL][lexer.NOTEQ] = true
	followSet[FUNCORVARIDNESTTAIL][lexer.LT] = true
	followSet[FUNCORVARIDNESTTAIL][lexer.GT] = true
	followSet[FUNCORVARIDNESTTAIL][lexer.LTE] = true
	followSet[FUNCORVARIDNESTTAIL][lexer.GTE] = true
	followSet[FUNCORVARIDNESTTAIL][lexer.PLUS] = true
	followSet[FUNCORVARIDNESTTAIL][lexer.MINUS] = true
	followSet[FUNCORVARIDNESTTAIL][lexer.OR] = true
	followSet[FUNCORVARIDNESTTAIL][lexer.COMMA] = true
	followSet[FUNCORVARIDNESTTAIL][lexer.COLON] = true
	followSet[FUNCORVARIDNESTTAIL][lexer.RBRACKET] = true
	followSet[FUNCORVARIDNESTTAIL][lexer.RPAREN] = true

	// FUNCSTATTAIL
	followSet[FUNCSTATTAIL] = map[lexer.TokenType]bool{}
	followSet[FUNCSTATTAIL][lexer.SEMICOLON] = true

	// INHERITS
	followSet[INHERITS] = map[lexer.TokenType]bool{}
	followSet[INHERITS][lexer.LBRACE] = true

	// NESTEDID
	followSet[NESTEDID] = map[lexer.TokenType]bool{}
	followSet[NESTEDID][lexer.LBRACE] = true

	// CLASSDECLBODY
	followSet[CLASSDECL] = map[lexer.TokenType]bool{}
	followSet[CLASSDECL][lexer.FUNCTION] = true
	followSet[CLASSDECL][lexer.MAIN] = true

	// CLASSDECLBODY
	followSet[CLASSDECLBODY] = map[lexer.TokenType]bool{}
	followSet[CLASSDECLBODY][lexer.RBRACE] = true

	// METHODBODYVAR
	followSet[METHODBODYVAR] = map[lexer.TokenType]bool{}
	followSet[METHODBODYVAR][lexer.IF] = true
	followSet[METHODBODYVAR][lexer.WHILE] = true
	followSet[METHODBODYVAR][lexer.READ] = true
	followSet[METHODBODYVAR][lexer.WRITE] = true
	followSet[METHODBODYVAR][lexer.RETURN] = true
	followSet[METHODBODYVAR][lexer.BREAK] = true
	followSet[METHODBODYVAR][lexer.CONTINUE] = true
	followSet[METHODBODYVAR][lexer.IDENT] = true
	followSet[METHODBODYVAR][lexer.RBRACE] = true

	// VISIBILITY
	followSet[VISIBILITY] = map[lexer.TokenType]bool{}
	followSet[VISIBILITY][lexer.FUNCTION] = true
	followSet[VISIBILITY][lexer.INTEGER] = true
	followSet[VISIBILITY][lexer.FLOAT] = true
	followSet[VISIBILITY][lexer.STRING] = true
	followSet[VISIBILITY][lexer.IDENT] = true

	// ADDOP
	followSet[PLUS] = map[lexer.TokenType]bool{}
	followSet[PLUS][lexer.INT_VALUE] = true
	followSet[PLUS][lexer.FLOAT_VALUE] = true
	followSet[PLUS][lexer.STRINGLIT] = true
	followSet[PLUS][lexer.LPAREN] = true
	followSet[PLUS][lexer.NOT] = true
	followSet[PLUS][lexer.QUESTION] = true
	followSet[PLUS][lexer.IDENT] = true
	followSet[PLUS][lexer.PLUS] = true
	followSet[PLUS][lexer.MINUS] = true

	// MULOP
	followSet[MUL] = map[lexer.TokenType]bool{}
	followSet[MUL][lexer.INT_VALUE] = true
	followSet[MUL][lexer.FLOAT_VALUE] = true
	followSet[MUL][lexer.STRINGLIT] = true
	followSet[MUL][lexer.LPAREN] = true
	followSet[MUL][lexer.NOT] = true
	followSet[MUL][lexer.QUESTION] = true
	followSet[MUL][lexer.IDENT] = true
	followSet[MUL][lexer.PLUS] = true
	followSet[MUL][lexer.MINUS] = true

	// MULOP
	followSet[LT] = map[lexer.TokenType]bool{}
	followSet[LT][lexer.INT_VALUE] = true
	followSet[LT][lexer.FLOAT_VALUE] = true
	followSet[LT][lexer.STRINGLIT] = true
	followSet[LT][lexer.LPAREN] = true
	followSet[LT][lexer.NOT] = true
	followSet[LT][lexer.QUESTION] = true
	followSet[LT][lexer.IDENT] = true
	followSet[LT][lexer.PLUS] = true
	followSet[LT][lexer.MINUS] = true

	// SIGN
	followSet[SIGN] = map[lexer.TokenType]bool{}
	followSet[SIGN][lexer.INT_VALUE] = true
	followSet[SIGN][lexer.FLOAT_VALUE] = true
	followSet[SIGN][lexer.STRINGLIT] = true
	followSet[SIGN][lexer.LPAREN] = true
	followSet[SIGN][lexer.NOT] = true
	followSet[SIGN][lexer.QUESTION] = true
	followSet[SIGN][lexer.IDENT] = true
	followSet[SIGN][lexer.PLUS] = true
	followSet[SIGN][lexer.MINUS] = true

	return followSet
}
