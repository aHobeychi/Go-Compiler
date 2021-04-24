package translation

import (
	"aHobeychi/GoCompiler/ast"
	"strconv"
)

// constant containing all the assembly commands.
const (
	add       = "add"
	addi      = "addi"      // add integers ex: addi r1, r0, 8
	sub       = "sub"       //
	subi      = "subi"      //
	mul       = "mul"       //
	muli      = "muli"      // mul integers ex: muli r1, r1, 8
	div       = "div"       //
	divi      = "divi"      //
	load      = "lw"        // Load Word
	store     = "sw"        // Store Word
	reserve   = "res"       // reserve bytes
	entry     = "entry"     // start execution, main function call
	halt      = "halt"      //  stops the execution
	equal     = "ceq"       // compare registers ceq r3, r2, r1
	eqi       = "ceqi"      // compare int ceq1 r2, r1, 8
	gte       = "cge"       // compare int ceq1 r2, r1, 8
	gt        = "cgt"
	lte       = "cle"       // compare int ceq1 r2, r1, 8
	lt        = "clt"
	noteq     = "cne"
	branch    = "bz"        // create branch if zero, will go to jump point
	branch_nz = "bnz"       // create branch if zero, will go to jump point
	jump      = "j"         // unconditional jump
	j_link    = "jl"        // jumps and puts address into register r15
	j_back    = "jr"        // jumps back to addres jr r15
	r0        = "r0"        // r0 register contains zero value
	top       = "topaddr"   // top address
	out       = "putc"      // outputs to the console
	in        = "getc"      // outputs to the console
	res       = "res"       // reserve memory
)

// RelOp map used to get correct instructions during translation.
var relOpInstruction = map[ast.Production]string{

	ast.LT:    lt,
	ast.LTE:   lte,
	ast.GT:    gt,
	ast.GTE:   gte,
	ast.EQ:    equal,
	ast.NOTEQ: noteq,
}

// Arith map used to get correct instructions during translation.
var arithOpInstruction = map[ast.Production]string{

	ast.PLUS:  add,
	ast.MINUS: sub,
	ast.DIV:   div,
	ast.MUL:   mul,
}

// Simplified arith map used to get correct instructions during translation.
var simplifiedArith = map[ast.Production]string{

	ast.PLUS:  addi,
	ast.MINUS: subi,
	ast.DIV:   divi,
	ast.MUL:   muli,
}

// Label counter used to create unique labels.
var label_counter int

// Generates unique label for given id.
func generateLabel(label_base string) string {
	label_counter++
	return label_base + strconv.Itoa(label_counter)
}
