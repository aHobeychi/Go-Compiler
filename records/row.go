package records

import (
	"aHobeychi/GoCompiler/ast"
)

// Row interface, We'll have 3 types of Rows.
type Row interface {
	ToString() string
	GetKind() ast.Production
	GetType() string
}

// Prints Row.
func PrintRow(r Row) {

	switch r.GetKind() {
	case ast.FUNCTION:
		OutputMessage(r.ToString())
	case ast.CLASS:
		OutputMessage(r.ToString())
	case ast.VARIABLE:
		OutputMessage(r.ToString())
	}

	OutputMessage("\n")
}

// Contains variable elements.
type Variables struct {
	Iden     string
	BaseType string
	ArrDim   int
}

// Returns a primity variable of type base_type
func NewVariable(base_type string) *Variables {
	return &Variables{BaseType: base_type}
}

// Returns string representation of the variable.
func (vr *Variables) ToString() string {
	str_rep := vr.BaseType
	for i := 0; i < vr.ArrDim; i++ {
		str_rep += "[]"
	}

	return str_rep
}

// Returns if Variable is of the size type and size and the parameter variable.
func (vr *Variables) CompareVar(other *Variables) bool {
	return vr.BaseType == other.BaseType && vr.ArrDim == other.ArrDim
}

// Struct Containing all the variables
type Parameters []Variables

// Return number of variables in the set of parameters.
func (param *Parameters) NumberOfVariables() int {
	return len(*param)
}

// Returns true if two set of parameters match
func (param *Parameters) MatchParameter(other Parameters) bool {

	if param.NumberOfVariables() != other.NumberOfVariables() {
		return false
	}

	for index, variable := range *param {

		if variable.BaseType != other[index].BaseType {
			return false
		}

		if variable.ArrDim != other[index].ArrDim {
			return false
		}
	}

	return true
}
