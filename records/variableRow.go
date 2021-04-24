package records

import (
	"aHobeychi/GoCompiler/ast"
	"strconv"
)

// Variable row, used to define variables, their types and values
type VariableRow struct {
	Name           string
	VarType        string
	ArrDim         int
	Memory_Size    int
	Memory_Address int
	ArrayIndex     []int
	Kind           ast.Production
	Visibility     ast.Production
}

// Returns string value of
func (vr *VariableRow) ToString() string {

	str_rep := vr.Name + "｜" + ast.TypesStrings[vr.Kind] + "｜Visibility: " + ast.TypesStrings[vr.Visibility]
	str_rep = str_rep + "｜Type: " + vr.VarType

	for i := 0; i < vr.ArrDim; i++ {
		str_rep = str_rep + "[" + "]"
	}

	str_rep += "｜Size: " + strconv.Itoa(vr.Memory_Size)
	str_rep += "｜Offset: " + strconv.Itoa(vr.Memory_Address)

	return str_rep
}

// return row kind.
func (vr *VariableRow) GetKind() ast.Production {
	return vr.Kind
}

// return row kind.
func (vr *VariableRow) GetType() string {
	return vr.VarType
}

// Returns new variable row.
func NewVariableRow(name string, kind ast.Production) *VariableRow {
	return &VariableRow{Name: name, Kind: kind}
}

func (v *VariableRow) ReturnCopy() (string, *VariableRow) {

	newRow := NewVariableRow(v.Name, v.Kind)

	newRow.VarType = v.VarType
	newRow.ArrDim = v.ArrDim
	newRow.ArrayIndex = v.ArrayIndex

	return v.Name, newRow
}
