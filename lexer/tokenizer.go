// Lexer package contains all the necessary methods to conduct
// the lexical analysis based on the definition in the token package.
package lexer

import (
	"aHobeychi/GoCompiler/utilities"
	"regexp"
	"strings"
)

// RETURN COMMENTS MUST BE SET TO FALSE TO
// IGNORE COMMENTS AND NOT RETURN THEIR TOKENS
const (
	RETURN_COMMENTS = true
)

// Tokenizer struct to handle getting new tokens
type Tokenizer struct {
	input        string
	line         int
	column       int
	readPosition int
	ch           byte
}

// Creates new tokenizer and return it, only accesible from within the pakcage
func newTokenizer(input string) *Tokenizer {
	l := &Tokenizer{input: input, line: 1, readPosition: 0, column: 0}
	l.readChar()
	return l
}

// Read Next Char and track it in the lexer
func (l *Tokenizer) readChar() {

	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.column = l.readPosition
	l.readPosition += 1
}

// Return the next token in the input
func (l *Tokenizer) GetToken() Token {

	var tok Token
	tok.Error_Code = 1
	next_input := byte(9)

	if l.readPosition < len(l.input) {
		next_input = l.input[l.readPosition]
	}

	switch l.ch {
	case ';':
		tok = newToken(SEMICOLON, l.ch, l.line, 1)
	case '.':
		tok = newToken(DOT, l.ch, l.line, 1)
	case ',':
		tok = newToken(COMMA, l.ch, l.line, 1)
	case '?':
		tok = newToken(QUESTION, l.ch, l.line, 1)
	case '!':
		tok = newToken(NOT, l.ch, l.line, 1)
	case '|':
		tok = newToken(OR, l.ch, l.line, 1)
	case '&':
		tok = newToken(AND, l.ch, l.line, 1)
	case '*':
		tok = newToken(MUL, l.ch, l.line, 1)
	case '+':
		tok = newToken(PLUS, l.ch, l.line, 1)
	case '-':
		tok = newToken(MINUS, l.ch, l.line, 1)
	case '/':
		switch next_input {
		case '/':
			tok.Lexeme = l.readRestOfLine()
			tok.Type = LICOMMENT
			tok.Line = l.line
			tok.Error_Code = -2
			if !RETURN_COMMENTS {
				return l.GetToken()
			} else {
				tok.Type = LICOMMENT
				tok.Line = l.line
			}
		case '*':
			tok.Lexeme = l.readBlock()
			if !RETURN_COMMENTS {
				return l.GetToken()
			} else {
				tok.Type = BLCOMMENT
				tok.Line = l.line
				tok.Error_Code = -2
			}
		default:
			tok = newToken(DIV, l.ch, l.line, 1)
		}
	case '(':
		tok = newToken(LPAREN, l.ch, l.line, 1)
	case ')':
		tok = newToken(RPAREN, l.ch, l.line, 1)
	case '{':
		tok = newToken(LBRACE, l.ch, l.line, 1)
	case '}':
		tok = newToken(RBRACE, l.ch, l.line, 1)
	case '[':
		tok = newToken(LBRACKET, l.ch, l.line, 1)
	case ']':
		tok = newToken(RBRACKET, l.ch, l.line, 1)
	case '"':
		tok.Lexeme = l.readStringLiteral()
		tok.Type = STRINGLIT
		tok.Line = l.line
		return tok
	case '=':
		switch next_input {
		case '=':
			tok = newTokenFromString(EQUAL, "==", l.line, 1)
			l.readChar()
		default:
			tok = newToken(ASSIGN, l.ch, l.line, 1)
		}
	case '<':
		switch next_input {
		case '=':
			tok = newTokenFromString(LTE, "<=", l.line, 1)
			l.readChar()
		case '>':
			tok = newTokenFromString(DIAMOND, "<>", l.line, 1)
			l.readChar()
		default:
			tok = newToken(LT, l.ch, l.line, 1)
		}
	case '>':
		switch next_input {
		case '=':
			tok = newTokenFromString(GTE, "<", l.line, 1)
			l.readChar()
		default:
			tok = newToken(GT, l.ch, l.line, 1)
		}
	case ':':
		switch next_input {
		case ':':
			tok = newTokenFromString(DCOLON, "::", l.line, 1)
			l.readChar()
		default:
			tok = newToken(COLON, l.ch, l.line, 1)
		}
	case 0:
		tok.Lexeme = "EOF"
		tok.Type = EOF
		tok.Line = l.line
		tok.Error_Code = 1
	case 9: // EMPTY SPACE
		l.readChar()
		return l.GetToken()
	case 32: // EMPTY SPACE
		l.readChar()
		return l.GetToken()
	case 10: // NEW LINE
		l.nextLine()
		l.readChar()
		return l.GetToken()
	case '~', '@', '#', '%', '\\', '\'', '$', '`', '^':
		tok = newToken(INVALIDCHAR, l.ch, l.line, -1)
	case '_':
		if isLetter(next_input) {
			l.readChar()
			lexeme := "_" + l.readAlphaNum()
			tok = newTokenFromString(INVALIDIDEN, lexeme, l.line, -1)
		} else {
			tok = newToken(INVALIDCHAR, l.ch, l.line, -1)
		}
	default:
		if isLetter(l.ch) {
			tok.Lexeme = l.readAlphaNum()
			tok.Line = l.line
			tok.Type = LookupKeyword(tok.Lexeme)
			return tok
		}
		if isDigit(l.ch) {
			return l.processNumber()
		}
	}

	l.readChar()
	return tok
}

// Reads alphanumerical and returns its content
func (l *Tokenizer) readAlphaNum() string {
	column := l.column
	for isAlphaNum(l.ch) {
		l.readChar()
	}
	return l.input[column:l.column]
}

// Reads the rest of a line.
func (l *Tokenizer) readRestOfLine() string {
	column := l.column
	for l.ch != byte(10) {
		l.readChar()
	}
	return l.input[column:l.column]
}

// Read Content of a block comment.
// Accepts nested blocks
func (l *Tokenizer) readBlock() string {

	column := l.column
	previous := byte('p')
	var stack utilities.Stack

	for {

		previous = l.ch
		l.readChar()

		if previous == byte('/') && l.ch == byte('*') {
			stack.Push("/*")
			l.readChar()
		}

		if previous == byte('*') && l.ch == byte('/') {

			l.readChar()
			stack.Pop()

			if stack.IsEmpty() {
				break
			}
		}

		if l.ch == 0 {
			panic("Block Comment was not closed. Program terminating.")
		}

		previous = l.ch
	}

	return strings.ReplaceAll(l.input[column:l.column], " \n ", "\\n")
}

// Read Content of literal String.
func (l *Tokenizer) readStringLiteral() string {

	column := l.column

	for {

		l.readChar()

		if l.ch == byte('"') {
			l.readChar()
			break
		}

		if l.ch == 0 {
			panic("String literal was not closed, program terminating.\n")
		}

	}

	return strings.ReplaceAll(l.input[column:l.column], "\"", "")
}

// Checks if character is a letter.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

// Checks if character is a digit.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// Checks if character is an alphanumerical or _
func isAlphaNum(ch byte) bool {
	return isDigit(ch) || isLetter(ch) || ch == '_'
}

// Processes Numbers
func (l *Tokenizer) processNumber() Token {

	correct := false
	isFloat := false
	matchedSymbol := false
	column := l.column

	for isDigit(l.ch) {
		l.readChar()
	}

	if l.ch == '+' || l.ch == '-' {
		lexem := l.input[column:l.column]
		re := regexp.MustCompile(IntegerRegex)
		correct = re.MatchString(lexem)

		if correct {
			return newTokenFromString(INT_VALUE, lexem, l.line, 1)
		}

		return newTokenFromString(INVALIDINT, lexem, l.line, -1)
	}

	if l.ch == '.' {
		isFloat = true
	}

	for l.ch == '.' || isAlphaNum(l.ch) ||
		(!matchedSymbol && (l.ch == '+' || l.ch == '-')) {

		if l.ch == '+' || l.ch == '-' {
			matchedSymbol = true
		}

		l.readChar()
	}

	lexem := l.input[column:l.column]

	if isFloat {

		re := regexp.MustCompile(FloatRegex)
		correct = re.MatchString(lexem)

		if correct {
			return newTokenFromString(FLOAT_VALUE, lexem, l.line, 1)
		}
		return newTokenFromString(INVALIDFLOAT, lexem, l.line, -1)

	} else {

		re := regexp.MustCompile(IntegerRegex)
		correct = re.MatchString(lexem)

		if correct {
			return newTokenFromString(INT_VALUE, lexem, l.line, 1)
		}

		return newTokenFromString(INVALIDINT, lexem, l.line, -1)
	}
}

// Creates New Token From a Character
func newToken(tokenType TokenType, ch byte, line int, code int) Token {
	return Token{Type: tokenType, Lexeme: string(ch), Line: line, Error_Code: code}
}

// Creates New Token From a String
func newTokenFromString(tokenType TokenType, ch string, line int, code int) Token {
	return Token{Type: tokenType, Lexeme: ch, Line: line, Error_Code: code}
}

// Increments Line Number.
func (l *Tokenizer) nextLine() {
	l.line++
}
