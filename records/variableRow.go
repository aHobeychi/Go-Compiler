package records

import (
	"aHobeychi/GoCompiler/ast"
	"strconv"
	"strings"
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
	var sb strings.Builder

	sb.WriteString(vr.Name)
	sb.WriteString("｜")
	sb.WriteString(ast.TypesStrings[vr.Kind])
	sb.WriteString("｜Visibility: ")
	sb.WriteString(ast.TypesStrings[vr.Visibility])
	sb.WriteString("｜Type: ")
	sb.WriteString(vr.VarType)

	for i := 0; i < vr.ArrDim; i++ {
		sb.WriteString("[]")
	}

	sb.WriteString("｜Size: ")
	sb.WriteString(strconv.Itoa(vr.Memory_Size))
	sb.WriteString("｜Offset: ")
	sb.WriteString(strconv.Itoa(vr.Memory_Address))

	return sb.String()
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
