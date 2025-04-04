// Package containing useful functions and data structures.
package utilities

// Stack is a generic implementation of a stack data structure
// The type parameter T must be comparable to enable the Contains method
type Stack[T comparable] []T

// IsEmpty checks if the Stack is empty.
func (s *Stack[T]) IsEmpty() bool {
	return len(*s) == 0
}

// Push adds an item on top of the stack.
func (s *Stack[T]) Push(item T) {
	*s = append(*s, item)
}

// Pop removes and returns the top item from the stack.
func (s *Stack[T]) Pop() (bool, T) {
	var zero T
	if s.IsEmpty() {
		return false, zero
	} else {
		index := len(*s) - 1
		element := (*s)[index]
		*s = (*s)[:index]
		return true, element
	}
}

// Contains returns true if stack contains specific item.
func (s *Stack[T]) Contains(item T) bool {
	for _, val := range *s {
		if val == item {
			return true
		}
	}
	return false
}
