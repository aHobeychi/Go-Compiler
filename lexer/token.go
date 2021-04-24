package lexer

import (
	"fmt"
	"strconv"
)

type TokenType uint8

type Token struct {
	Type       TokenType
	Error_Code int
	Lexeme     string
	Line       int
}

// Returns New Token.
func (token *Token) NewToken(tokenType TokenType) *Token {
	return &Token{tokenType, 0, "", 0}
}

var FloatRegex = `^((([1-9]\d*)|0)\.(([0-9]*[1-9])|0)(e[+|-]([1-9]+[0-9]*)|0)?)$`
var IntegerRegex = `^(([1-9]\d*)|0)$`

const (
	PROGRAM = iota
	EOF
	// Identifiers and Literals
	IDENT
	INTEGER
	INT_VALUE
	FLOAT
	FLOAT_VALUE
	BOOL
	STRING
	STRINGLIT
	// Logical Operator
	NOTEQ
	EQUAL
	LTE
	GTE
	LT
	GT
	NOT
	QUESTION
	AND
	OR
	DIAMOND
	// Math Operator
	MUL
	DIV
	MINUS
	PLUS
	ASSIGN
	//
	COMMA
	DOT
	SEMICOLON
	LPAREN
	RPAREN
	LBRACE
	RBRACE
	RBRACKET
	LBRACKET
	QUOTE
	COLON
	DCOLON
	//
	IF
	THEN
	ELSE
	VOID
	PUBLIC
	PRIVATE
	FUNCTION
	VAR
	CLASS
	WHILE
	READ
	WRITE
	RETURN
	MAIN
	INHERITS
	BREAK
	CONTINUE
	//
	LICOMMENT
	BLCOMMENT
	//
	INVALIDCHAR
	INVALIDIDEN
	INVALIDINT
	INVALIDFLOAT
	EMPTY
)

var TokenAsString = map[TokenType]string{
	SEMICOLON:   ";",
	LPAREN:      "(",
	RPAREN:      ")",
	LBRACE:      "{",
	RBRACE:      "}",
	LBRACKET:    "[",
	RBRACKET:    "]",
	IDENT:       "Id",
	PLUS:        "+",
	MINUS:       "-",
	VAR:         "VAR",
	INTEGER:     "INT",
	FLOAT:       "FLOAT",
	EQUAL:       "==",
	ASSIGN:      "=",
	DIV:         "/",
	MUL:         "*",
	COLON:       ":",
	WRITE:       "WRITE",
	WHILE:       "WHILE",
	IF:          "IF",
	ELSE:        "ELSE",
	THEN:        "THEN",
	EOF:         "END OF FILE",
	INT_VALUE:   "INT_VALUE",
	FLOAT_VALUE: "FLOAT_VALUE",
	RETURN:      "RETURN",
	LT:          "<",
	LTE:         "<=",
	GT:          ">",
	GTE:         ">=",
	NOT:         "!",
	NOTEQ:       "!=",
	AND:         "&",
	OR:          "|",
	LICOMMENT:   "LINE_COMMENT",
	BLCOMMENT:   "BLOCKCOMMENT",
}

var reservedWords = map[string]TokenType{

	"if":       IF,
	"then":     THEN,
	"else":     ELSE,
	"void":     VOID,
	"public":   PUBLIC,
	"private":  PRIVATE,
	"func":     FUNCTION,
	"var":      VAR,
	"class":    CLASS,
	"while":    WHILE,
	"read":     READ,
	"write":    WRITE,
	"return":   RETURN,
	"main":     MAIN,
	"inherits": INHERITS,
	"break":    BREAK,
	"continue": CONTINUE,
	"bool":     BOOL,
	"integer":  INTEGER,
	"float":    FLOAT,
	"string":   STRING,
}

// All the delimeters
var Delimiters = map[byte]TokenType{
	'+': PLUS,
	'-': MINUS,
	'*': MUL,
	'/': DIV,
	'(': LPAREN,
	')': RPAREN,
	'[': LBRACE,
	']': RBRACE,
	'{': LBRACKET,
	'}': RBRACKET,
	'!': NOT,
	'?': QUESTION,
	'"': QUOTE,
	'&': AND,
	'|': OR,
}

// Lookups keyword, and returns the token type for the given reserved word.
func LookupKeyword(word string) TokenType {
	if token, exists := reservedWords[word]; exists {
		return token
	}
	return IDENT
}

// Returns a string representation of the token, useful for error messaging.
func (t *Token) ToString() string {
	return fmt.Sprint("[", TokenAsString[t.Type], ", ", t.Line, ", ", t.Lexeme, "]\n")
}

// Returns the line of the token as a string, useful for error messaging.
func (t *Token) GetLine() string {
	return strconv.Itoa(t.Line)
}
