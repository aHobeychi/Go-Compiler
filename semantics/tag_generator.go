package semantics

import "strconv"

var tag_counter int

// Gerenates the semantic tags that will be used for translation.
func generateTag(label_base string) string {
	tag_counter++
	return label_base + strconv.Itoa(tag_counter)
}
