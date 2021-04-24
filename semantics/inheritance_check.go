// Links Class with their inherited Classes and check for circular inheritance.
package semantics

import (
	"aHobeychi/GoCompiler/ast"
	"aHobeychi/GoCompiler/records"
	"aHobeychi/GoCompiler/utilities"
)

// Links Subclass to Parent class, Checks if inherited class exists.
func Inheritance_check(global *records.SymbolTable) bool {

	gb_table = global

	correct := true

	for k := range gb_table.Rows {

		row := gb_table.GetRow(k)

		if row.GetKind() == ast.FUNCTION || row.GetKind() == ast.VARIABLE {
			continue
		}

		classrow := row.(*records.ClassRow)
		inheritedClasses := classrow.Inherits

		for _, class := range inheritedClasses {

			_, exists := gb_table.Rows[class]

			if exists && gb_table.GetRow(class).GetKind() == ast.CLASS {

				classrow.InheritedLink[class] = gb_table.GetRow(class).(*records.ClassRow).Link

			} else {

				outputMessage("Class " + class + " doesn't exist at line: ")
				outputMessage(classrow.GetLine() + ".\n")
				correct = false

			}
		}
	}

	if correct == true {

		return checkForCircularInherits(gb_table)

	} else {

		checkForCircularInherits(gb_table)
		return false
	}
}

// Checks if file contains inheritance loops
func checkForCircularInherits(gb_table *records.SymbolTable) bool {

	for k := range gb_table.Rows {

		row := gb_table.GetRow(k)

		if row.GetKind() == ast.FUNCTION || row.GetKind() == ast.VARIABLE {
			continue
		}

		classrow := row.(*records.ClassRow)
		inheritedClasses := classrow.Inherits

		if len(inheritedClasses) == 0 {
			continue
		}

		for _, class := range inheritedClasses {

			if k == class {
				outputMessage("Class Cannot import itself. Line: " + classrow.GetLine() + "\n")
			}
		}

		var stack utilities.Stack
		inherits := make(map[string]bool)

		stack.Push(k)

		for !stack.IsEmpty() {

			_, class := stack.Pop()

			if inherits[class] {

				outputMessage("Contains Inheritance loop: " + class)
				outputMessage(". Line: " + classrow.GetLine() + ".\n")
				return false
			}

			inherits[class] = true

			if !rowIsClass(class) {
				continue
			}

			next_row := gb_table.GetRow(class)
			next_class_row := next_row.(*records.ClassRow)
			next_inherits := next_class_row.Inherits

			for _, str := range next_inherits {

				if str == "" {
					continue
				}

				stack.Push(str)
			}
		}
	}

	return true
}

// Checks if the class exists in the gb_table.
func rowIsClass(class string) bool {

	if !gb_table.EntryExists(class) {
		return false
	}

	if gb_table.GetRow(class).GetKind() == ast.FUNCTION ||

		gb_table.GetRow(class).GetKind() == ast.VARIABLE {
		return false
	}

	return true
}

// Add Parent var into the child class's symbol Table
func mergeTable(root *ast.Node) {

	classes := gb_table.GetClassList()

	for i := 0; i < len(*classes); i++ {
		for _, class := range *classes {

			classRow := gb_table.GetRow(class).(*records.ClassRow)
			inheritedClasses := classRow.Inherits
			addParentNodes(classRow.Link, inheritedClasses)
		}
	}

}

func addParentNodes(classTable *records.SymbolTable, parents []string) {

	for _, parent := range parents {

		parentClass := gb_table.GetRow(parent).(*records.ClassRow)
		parentTable := parentClass.Link
		parentVar := parentTable.GetVariableList()

		for _, parentVariable := range *parentVar {

			curr_var := parentTable.GetRow(parentVariable).(*records.VariableRow)
			copy_name, copy_row := curr_var.ReturnCopy()

			if _, contains := classTable.Rows[copy_name]; !contains {

				classTable.AddVariable(copy_name)
				classTable.Rows[copy_name] = copy_row
			}
		}

	}
}
