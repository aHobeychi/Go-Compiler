package reporting

import "os"

var token_path *os.File

func SetTokenOutPath(path *os.File) {
	token_path = path
}

func OutputToken(message string) {
	token_path.WriteString(message)
}
