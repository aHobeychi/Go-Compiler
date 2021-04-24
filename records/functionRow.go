package records

import (
	"aHobeychi/GoCompiler/ast"
	"strconv"
)

// Function row.
type FunctionRow struct {
	Name           string
	ReturnType     []string
	ClassMethod    []string
	Kind           ast.Production
	ListParameters []Parameters
	Visibility     []ast.Production
	Link           []*SymbolTable
}

// Returns link at specific index.
func (fr *FunctionRow) getLink(index int) *SymbolTable {
	return fr.Link[index]
}

// Returns a specific link with parameters and return_type.
func (fr *FunctionRow) getSpecificLink(param Parameters, return_type string) *SymbolTable {

	parameters := fr.ListParameters

	for index, params := range parameters {
		if fr.ReturnType[index] == return_type && params.MatchParameter(param) {
			return fr.Link[index]
		}
	}

	return nil
}

// adds link to function row.
func (fr *FunctionRow) addLink(table *SymbolTable) {
	fr.Link = append(fr.Link, table)
}

// add Visibility to row.
func (fr *FunctionRow) addVisibility(visibility ast.Production) {
	fr.Visibility = append(fr.Visibility, visibility)
}

// add class method to table.
func (sym *SymbolTable) AddClassMethod(index string, class *ast.Node) {

	if sym.GetRow(index).GetKind() != ast.FUNCTION {
		return
	}

	row := sym.GetRow(index).(*FunctionRow)

	if class == nil || !class.HasChildren() {
		row.ClassMethod = append(row.ClassMethod, "")
		return
	}

	class_method := class.GetChildren()[0].Lexeme
	row.ClassMethod = append(row.ClassMethod, class_method)
}

// Converts row information to string
func (vr *FunctionRow) ToString() string {

	str_rep := vr.Name + "｜" + ast.TypesStrings[vr.Kind]
	ListParameters := vr.ListParameters

	if vr.Name == "MAIN" {

		str_rep := vr.Name

		if vr.Link != nil {
			str_rep = str_rep + "｜Implemented"
		}

		return str_rep
	}

	str_rep = str_rep + "｜Ret:("

	for index, return_t := range vr.ReturnType {
		str_rep += return_t
		if index == len(vr.ReturnType)-1 {
			break
		}
		str_rep += ","
	}

	str_rep = str_rep + ")｜Vis:("

	for index, vis := range vr.Visibility {

		str_rep += ast.TypesStrings[vis]
		if index == len(vr.Visibility)-1 {
			break
		}

		str_rep += ","
	}

	str_rep += ")｜"

	str_rep = str_rep + "ClassM:("
	for index, class := range vr.ClassMethod {
		if class == "" {
			str_rep += "none"
		}
		str_rep += class
		if index == len(vr.ClassMethod)-1 {
			break
		}
		str_rep += ","
	}

	str_rep = str_rep + ")｜Param: "

	for i := 0; i < len(ListParameters); i++ {
		parameters := ListParameters[i]
		str_rep = str_rep + "("

		for j := 0; j < len(parameters); j++ {

			if parameters[j].BaseType == "" {
				continue
			}

			str_rep = str_rep + parameters[j].BaseType

			for k := 0; k < parameters[j].ArrDim; k++ {
				str_rep = str_rep + "[]"
			}

			if j == len(parameters)-1 {
				break
			}

			str_rep = str_rep + ","
		}

		str_rep = str_rep + ")"
	}

	if vr.Link != nil {
		str_rep = str_rep + "｜Impl[" + strconv.Itoa(len(vr.Link)) + "]"
	}

	return str_rep
}

// Returns the kind of the row
func (vr *FunctionRow) GetKind() ast.Production {
	return vr.Kind
}

// Adds a return type to the funtion row.
func (fr *FunctionRow) AddReturnType(return_t string) {
	fr.ReturnType = append(fr.ReturnType, return_t)
}

// Returns the kind of the row
func (vr *FunctionRow) GetType() string {
	return ""
}

// Creates new function row.
func NewFunctionRow(name string, kind ast.Production) *FunctionRow {
	return &FunctionRow{Name: name, Kind: kind}
}

// Tries to match one list of parameters, true if matched, false otherwise.
func (fr *FunctionRow) MatchOneParam(aparam *ast.Node) bool {

	parameters := formatFparams(aparam)

	for _, fparam := range fr.ListParameters {

		if len(fparam) != len(*parameters) {
			continue
		}
		if fparam.MatchParameter(*parameters) {
			return true
		}

	}

	return false
}

// Returns the link for a function of a class with a specific list of parameters.
func (fr *FunctionRow) GetLinkForClass(class string, fparams *ast.Node) *SymbolTable {

	if fr.Name == "MAIN" {
		return fr.Link[0]
	}

	parameters := formatFparams(fparams)

	classMethods := fr.ClassMethod

	for index, classNames := range classMethods {
		if classNames == class && fr.ListParameters[index].MatchParameter(*parameters) {
			return fr.Link[index]
		}
	}

	return nil
}

func (fr *FunctionRow) GetLinkFromParam(class string, params *Parameters) *SymbolTable {

	for index, params := range fr.ListParameters {

		if params.MatchParameter(params) {
			return fr.Link[index]
		}
	}

	return nil

}

// Returns return type, for class and parameters list.
func (fr *FunctionRow) GetReturn(class string, params *Parameters) string {

	if class == "MAIN" {

		for index, param := range fr.ListParameters {

			if fr.ClassMethod[index] == "" {

				if param.MatchParameter(*params) {
					return fr.ReturnType[index]
				}
			}
		}
	}

	classMethods := fr.ClassMethod

	for index, classes := range classMethods {

		if classes == class {

			if fr.ListParameters[index].MatchParameter(*params) {
				return fr.ReturnType[index]
			}
		}
	}

	return "invalid"
}
