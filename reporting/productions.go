package reporting

import "os"

var prod_path *os.File

func SetProductPath(path *os.File) {
	prod_path = path
}

func OutputProd(message string) {
	prod_path.WriteString(message)
}
