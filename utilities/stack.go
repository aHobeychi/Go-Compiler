// Package containing useful functions and data structures.
package utilities

type Stack []string

// Checks if the Stacks is empty.
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

// Pushes string on top of the stack.
func (s *Stack) Push(str string) {
	*s = append(*s, str)
}

// Pops string from the stake.
func (s *Stack) Pop() (bool, string) {
	if s.IsEmpty() {
		return false, ""
	} else {
		index := len(*s) - 1
		element := (*s)[index]
		*s = (*s)[:index]
		return true, element
	}
}

// Returns true if stack contains specific string.
func (s *Stack) Contains(el string) bool {
	for _, str := range *s {
		if str == el {
			return true
		}
	}
	return false
}
