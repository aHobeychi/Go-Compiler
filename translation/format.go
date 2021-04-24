package translation

import "os"

var outputPath *os.File
var exec_code string

// Sets the ouput path of the translated file.
func SetOutputPath(path *os.File) {
	outputPath = path
}

// Prints translatedcode to file
func OutputMessage(message string) {
	print(message)
	outputPath.WriteString(message)
}

// Appends the command and the args to the translated code.
func appendExec(command, args string) {
	exec_code += "\t" + command + "\t" + args + "\n"
}

// Appends the scope, command and the args to the translated code.
func appendExecScope(scope, command, args string) {
	exec_code += scope + "\t" + command + "\t" + args + "\n"
}

func appendEmptySpace() {
	exec_code += "\n"
}
