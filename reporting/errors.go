package reporting

import "os"

var error_path *os.File

func SetOutputPath(path *os.File) {
	error_path = path
}

func OutputMessage(message string) {
	error_path.WriteString(message)
}
