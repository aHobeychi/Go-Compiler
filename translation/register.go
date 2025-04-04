package translation

import (
	"aHobeychi/GoCompiler/utilities"
)

// Returns New registor pool.
func getNewRegisterPool() *utilities.Stack[string] {

	var reg utilities.Stack[string]

	reg.Push("r12")
	reg.Push("r11")
	reg.Push("r10")
	reg.Push("r9")
	reg.Push("r8")
	reg.Push("r7")
	reg.Push("r6")
	reg.Push("r5")
	reg.Push("r4")
	reg.Push("r3")
	reg.Push("r2")
	reg.Push("r1")
	reg.Push("r14") // reserved for addres stack top.
	reg.Push("r0")  // reserved for 0 literal

	return &reg
}
