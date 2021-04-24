// Lexer package contains all the necessary methods to conduct
// the lexical analysis based on the definition in the token package.
package lexer

import (
	"aHobeychi/GoCompiler/reporting"
	"bufio"
	"fmt"
	"os"
)

// Lexer Struct is a Wrapper for the Tokenizer Struct.
type Lexer struct {
	input_path  string
	output_path string
	tokenizer   Tokenizer
}

// Formats the String before assigning it to the struct.
func formatInput(path string) string {

	code := ""
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		code = code + scanner.Text() + " \n "
	}

	return code
}

// Creates New Lexer.
func New(path string, output_path string) *Lexer {

	input := formatInput(path)
	tokenizer := newTokenizer(input)

	l := &Lexer{input_path: path, output_path: output_path,
		tokenizer: *tokenizer}

	return l
}

// Returns The Next
func (l *Lexer) NextToken() Token {

	token := l.tokenizer.GetToken()
	reporting.OutputToken(token.ToString())

	if token.Error_Code == 1 {
		return token
	}

	if token.Error_Code == -1 {
		reporting.OutputMessage(token.ToString())
	}

	return l.NextToken()
}

// Prints the entire input text.
func (l *Lexer) PrintInput() {
	fmt.Println(l.tokenizer.input)
}

// Checks for errors, panics if true
func check(e error) {
	if e != nil {
		panic(e)
	}
}
