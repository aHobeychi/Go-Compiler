package records

import (
	"aHobeychi/GoCompiler/ast"
	"strconv"
)

// Class Row.
type ClassRow struct {
	Name          string
	Kind          ast.Production
	Visibility    ast.Production
	Inherits      []string
	InheritedLink map[string]*SymbolTable
	Line          int
	Link          *SymbolTable
	Total_Memory  int
}

// Converts class row to string.
func (vr *ClassRow) ToString() string {
	str_rep := vr.Name + "｜" + ast.TypesStrings[vr.Kind]
	if len(vr.Inherits) != 0 {
		str_rep = str_rep + "｜Inherits: "
		for _, k := range vr.Inherits {
			str_rep = str_rep + k + " "
		}
	}
	if vr.Link != nil {
		str_rep = str_rep + "｜Impl"
	}
	str_rep += "｜Mem_Size: " + strconv.Itoa(vr.Total_Memory)

	return str_rep
}

// Returns kind of class row.
func (vr *ClassRow) GetKind() ast.Production {
	return vr.Kind
}

// Returns kind of class row.
func (vr *ClassRow) GetType() string {
	return vr.Name
}

// Creates new class row.
func NewClassRow(name string, kind ast.Production) *ClassRow {
	return &ClassRow{Name: name, Kind: kind, InheritedLink: make(map[string]*SymbolTable)}
}

// Add Inherited class to class row.
func (vr *ClassRow) AddInherits(class string) {
	vr.Inherits = append(vr.Inherits, class)
}

// Returns If One of the parent classes contains that id.
func (st *SymbolTable) InheritsSearch(class, id string) bool {
	if !st.EntryExists(class) || st.GetRow(class).GetKind() != ast.CLASS {
		return false
	}

	tmp_row := st.GetRow(class)
	classRow := tmp_row.(*ClassRow)
	classTable := classRow.Link

	if classTable.EntryExists(id) && classTable.GetRow(id).GetKind() == ast.VARIABLE {
		return true
	}

	if len(classRow.Inherits) == 0 {
		return false
	}

	for _, inherited := range classRow.Inherits {
		table := classRow.InheritedLink[inherited]
		if table.EntryExists(id) && table.GetRow(id).GetKind() == ast.VARIABLE {
			return true
		}
		if st.InheritsSearch(inherited, id) {
			return true
		}
	}

	return false
}

// Checks if variable id is inherited from a parent class.
func (st *SymbolTable) InheritsFromParen(class, id string) bool {
	if !st.EntryExists(class) || st.GetRow(class).GetKind() != ast.CLASS {
		return false
	}

	tmp_row := st.GetRow(class)
	classRow := tmp_row.(*ClassRow)

	if len(classRow.Inherits) == 0 {
		return false
	}

	for _, inherited := range classRow.Inherits {
		table := classRow.InheritedLink[inherited]

		if table.EntryExists(id) && table.GetRow(id).GetKind() == ast.VARIABLE {
			return true
		}
		if st.InheritsSearch(inherited, id) {
			return true
		}
	}

	return false
}

// Returns Type Of Parent Id
func (st *SymbolTable) GetInheritedVar(class, id string) *VariableRow {
	if !st.EntryExists(class) || st.GetRow(class).GetKind() != ast.CLASS {
		return nil
	}

	tmp_row := st.GetRow(class)
	classRow := tmp_row.(*ClassRow)
	classTable := classRow.Link

	if classTable.EntryExists(id) && classTable.GetRow(id).GetKind() == ast.VARIABLE {
		return classTable.GetRow(id).(*VariableRow)
	}

	if len(classRow.Inherits) == 0 {
		return nil
	}

	for _, inherited := range classRow.Inherits {
		table := classRow.InheritedLink[inherited]
		if table.EntryExists(id) && table.GetRow(id).GetKind() == ast.VARIABLE {
			return table.GetRow(id).(*VariableRow)
		}
		if st.InheritsSearch(inherited, id) {
			return st.GetInheritedVar(inherited, id)
		}
	}

	return nil
}

// Returns the line of the class declaration, useful for error messaging.
func (cr *ClassRow) GetLine() string {
	return strconv.Itoa(cr.Line)
}
