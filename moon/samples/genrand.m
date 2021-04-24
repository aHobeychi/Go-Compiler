	entry			% Start here
	addi	r14,r0,topaddr	% Initialize stack pointer
	addi	r1,r0,L11		% Store string Lit. addr. in register
	jl	r15,putstr		% Call output string function
	jl	r15,getint		% Call input integer function
	sw	L4(r0),r1		% Store value in memory
L12
	addi	r1,r0,L14		% Store string Lit. addr. in register
	jl	r15,putstr		% Call output string function
	subi	r14,r14,0	% Move stack pointer to calle frame
	jl	r15,L8	% Call the procedure/function
	addi	r14,r14,0	% Move stack pointer to caller frame
	addi	r1,r14,-260	% Store addr. of string result in register
	jl	r15,putstr		% Call output string function
	subi	r14,r14,0	% Move stack pointer to calle frame
	jl	r15,L6	% Call the procedure/function
	addi	r14,r14,0	% Move stack pointer to caller frame
	addi	r1,r14,-260	% Store addr. of string result in register
	jl	r15,putstr		% Call output string function
	addi	r1,r0,L15		% Store string Lit. addr. in register
	jl	r15,putstr		% Call output string function
	subi	r14,r14,0	% Move stack pointer to calle frame
	jl	r15,L9	% Call the procedure/function
	addi	r14,r14,0	% Move stack pointer to caller frame
	addi	r1,r14,-260	% Store addr. of string result in register
	jl	r15,putstr		% Call output string function
	addi	r1,r0,L16		% Store string Lit. addr. in register
	jl	r15,putstr		% Call output string function
	addi	r1,r0,L17		% Store string Lit. addr. in register
	jl	r15,putstr		% Call output string function
	subi	r14,r14,0	% Move stack pointer to calle frame
	jl	r15,L8	% Call the procedure/function
	addi	r14,r14,0	% Move stack pointer to caller frame
	addi	r1,r14,-260	% Store addr. of string result in register
	jl	r15,putstr		% Call output string function
	subi	r14,r14,0	% Move stack pointer to calle frame
	jl	r15,L6	% Call the procedure/function
	addi	r14,r14,0	% Move stack pointer to caller frame
	addi	r1,r14,-260	% Store addr. of string result in register
	jl	r15,putstr		% Call output string function
	addi	r1,r0,L18		% Store string Lit. addr. in register
	jl	r15,putstr		% Call output string function
	addi	r1,r0,L19		% Store string Lit. addr. in register
	jl	r15,putstr		% Call output string function
	addi	r1,r0,L10		% Address of var string in memory
	jl	r15,getstr		% Call input string function
	sb	BUF1(r0),r0		% Initialize buffer #1
	addi	r1,r0,BUF1		% Store BUF addr. into register
	addi	r3,r0,1		% Small numbers store immediate
	addi	r2,r0,L10	% Store substring addr. in register
	sw	-8(r14),r2	% Pass arguments to the stack
	sw	-12(r14),r3
	sw	-16(r14),r3
	subi	r14,r14,0	% Move stack pointer to calle frame
	jl	r15,substr	% Call substring function
	addi	r14,r14,0	% Move stack pointer to caller frame
	lw	r2,-20(r14)	% Retrieve results from stack
	lw	r4,-24(r14)
	jl	r15,catsub		% Call substring concatenate function
	addi	r1,r0,BUF1		% Restore BUF addr. into register
	add	r13,r0,r1		% Save r1 in temp register
	addi	r1,r0,L22		% Store string Lit. addr. in register
	add	r2,r0,r1		% Copy r1 in r2
	add	r1,r0,r13		% Restore temp register into r1
	jl	r15,neqstr		% Call not equal string function
	bz	r1,L21		% If cond. false jump to next pair
	j	L13		% Exit loop
	j	L20		% If cond. true jump to end if
L21
L20
	j	L12		% Go to start of loop
L13
	hlt			% End of main program
L5
	sw	-4(r14),r15	% Store link
	sw	-16(r14),r1
	sw	-20(r14),r2
	addi	r1,r0,25173		% Small numbers store immediate
	lw	r2,L4(r0)	% Load value from memory
	mul	r1,r1,r2	% Mul right expr. by left expr.
	addi	r2,r0,13849		% Small numbers store immediate
	add	r1,r1,r2	% Add left expr. to right expr.
	lw	r2,L23(r0)	% Load large num from mem. into register
	mod	r1,r1,r2	% Mod left expr. & right expr.
	sw	L4(r0),r1	% Store value in memory
	lw	r1,L24(r0)	% Load large num from mem. into register
	lw	r2,-8(r14)	% Load value from stack frame
	div	r1,r1,r2	% Div left expr. by right expr.
	lw	r2,L4(r0)	% Load value from memory
	div	r1,r2,r1	% Div left expr. by right expr.
	addi	r2,r0,1		% Small numbers store immediate
	add	r1,r1,r2	% Add left expr. to right expr.
	sw	-12(r14),r1	% Store value in stack frame
	lw	r1,-16(r14)
	lw	r2,-20(r14)
	lw	r15,-4(r14)	% Load the link into r15
	jr	r15		% Jump back to the stmt. after the call
L6
	sw	-4(r14),r15	% Store link
	sw	-268(r14),r1
	addi	r1,r0,10		% Small numbers store immediate
	sw	-276(r14),r1	% Store value in stack frame
	subi	r14,r14,268	% Move stack pointer to calle frame
	jl	r15,L5	% Call the procedure/function
	addi	r14,r14,268	% Move stack pointer to caller frame
	lw	r1,-280(r14)	% Store result in register
	sw	-264(r14),r1	% Store value in stack frame
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,1		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L26		% If cond. false jump to next pair
	addi	r1,r0,L27		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L28
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L29		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L28		% Loop if not finished
L29
	j	L25		% If cond. true jump to end if
L26
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,2		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L30		% If cond. false jump to next pair
	addi	r1,r0,L31		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L32
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L33		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L32		% Loop if not finished
L33
	j	L25		% If cond. true jump to end if
L30
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,3		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L34		% If cond. false jump to next pair
	addi	r1,r0,L35		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L36
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L37		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L36		% Loop if not finished
L37
	j	L25		% If cond. true jump to end if
L34
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,4		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L38		% If cond. false jump to next pair
	addi	r1,r0,L39		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L40
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L41		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L40		% Loop if not finished
L41
	j	L25		% If cond. true jump to end if
L38
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,5		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L42		% If cond. false jump to next pair
	addi	r1,r0,L43		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L44
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L45		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L44		% Loop if not finished
L45
	j	L25		% If cond. true jump to end if
L42
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,6		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L46		% If cond. false jump to next pair
	addi	r1,r0,L47		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L48
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L49		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L48		% Loop if not finished
L49
	j	L25		% If cond. true jump to end if
L46
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,7		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L50		% If cond. false jump to next pair
	addi	r1,r0,L51		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L52
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L53		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L52		% Loop if not finished
L53
	j	L25		% If cond. true jump to end if
L50
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,7		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L54		% If cond. false jump to next pair
	addi	r1,r0,L55		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L56
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L57		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L56		% Loop if not finished
L57
	j	L25		% If cond. true jump to end if
L54
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,8		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L58		% If cond. false jump to next pair
	addi	r1,r0,L59		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L60
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L61		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L60		% Loop if not finished
L61
	j	L25		% If cond. true jump to end if
L58
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,9		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L62		% If cond. false jump to next pair
	addi	r1,r0,L63		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L64
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L65		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L64		% Loop if not finished
L65
	j	L25		% If cond. true jump to end if
L62
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,10		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L66		% If cond. false jump to next pair
	addi	r1,r0,L67		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L68
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L69		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L68		% Loop if not finished
L69
	j	L25		% If cond. true jump to end if
L66
L25
	lw	r1,-268(r14)
	lw	r15,-4(r14)	% Load the link into r15
	jr	r15		% Jump back to the stmt. after the call
L7
	sw	-4(r14),r15	% Store link
	sw	-268(r14),r1
	addi	r1,r0,10		% Small numbers store immediate
	sw	-276(r14),r1	% Store value in stack frame
	subi	r14,r14,268	% Move stack pointer to calle frame
	jl	r15,L5	% Call the procedure/function
	addi	r14,r14,268	% Move stack pointer to caller frame
	lw	r1,-280(r14)	% Store result in register
	sw	-264(r14),r1	% Store value in stack frame
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,1		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L71		% If cond. false jump to next pair
	addi	r1,r0,L72		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L73
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L74		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L73		% Loop if not finished
L74
	j	L70		% If cond. true jump to end if
L71
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,2		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L75		% If cond. false jump to next pair
	addi	r1,r0,L76		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L77
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L78		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L77		% Loop if not finished
L78
	j	L70		% If cond. true jump to end if
L75
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,3		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L79		% If cond. false jump to next pair
	addi	r1,r0,L80		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L81
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L82		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L81		% Loop if not finished
L82
	j	L70		% If cond. true jump to end if
L79
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,4		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L83		% If cond. false jump to next pair
	addi	r1,r0,L84		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L85
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L86		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L85		% Loop if not finished
L86
	j	L70		% If cond. true jump to end if
L83
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,5		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L87		% If cond. false jump to next pair
	addi	r1,r0,L88		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L89
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L90		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L89		% Loop if not finished
L90
	j	L70		% If cond. true jump to end if
L87
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,6		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L91		% If cond. false jump to next pair
	addi	r1,r0,L92		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L93
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L94		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L93		% Loop if not finished
L94
	j	L70		% If cond. true jump to end if
L91
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,7		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L95		% If cond. false jump to next pair
	addi	r1,r0,L96		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L97
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L98		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L97		% Loop if not finished
L98
	j	L70		% If cond. true jump to end if
L95
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,8		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L99		% If cond. false jump to next pair
	addi	r1,r0,L100		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L101
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L102		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L101		% Loop if not finished
L102
	j	L70		% If cond. true jump to end if
L99
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,9		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L103		% If cond. false jump to next pair
	addi	r1,r0,L104		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L105
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L106		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L105		% Loop if not finished
L106
	j	L70		% If cond. true jump to end if
L103
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,10		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L107		% If cond. false jump to next pair
	addi	r1,r0,L108		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L109
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L110		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L109		% Loop if not finished
L110
	j	L70		% If cond. true jump to end if
L107
L70
	lw	r1,-268(r14)
	lw	r15,-4(r14)	% Load the link into r15
	jr	r15		% Jump back to the stmt. after the call
L8
	sw	-4(r14),r15	% Store link
	sw	-268(r14),r1
	sw	-272(r14),r2
	sw	-276(r14),r3
	addi	r1,r0,L111		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L112
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L113		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L112		% Loop if not finished
L113
	addi	r1,r0,0		% Small numbers store immediate
	sw	-264(r14),r1	% Store value in stack frame
L114
	addi	r1,r0,10		% Small numbers store immediate
	sw	-284(r14),r1	% Store value in stack frame
	subi	r14,r14,276	% Move stack pointer to calle frame
	jl	r15,L5	% Call the procedure/function
	addi	r14,r14,276	% Move stack pointer to caller frame
	lw	r1,-288(r14)	% Store result in register
	addi	r2,r0,6		% Small numbers store immediate
	clt	r1,r1,r2		% Is left expr. < right expr
	bz	r1,L117		% If cond. false jump to next pair
	j	L115		% Exit loop
	j	L116		% If cond. true jump to end if
L117
L116
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,0		% Small numbers store immediate
	cgt	r1,r1,r2		% Is left expr. > right expr
	bz	r1,L119		% If cond. false jump to next pair
	sb	BUF1(r0),r0		% Initialize buffer #1
	sb	BUF2(r0),r0		% Initialize buffer #2
	addi	r1,r0,BUF1		% Store BUF addr. into register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	jl	r15,catstr		% Call string concatenate function
	addi	r1,r0,BUF1		% Store BUF addr. into register
	addi	r2,r0,L120		% Store string Lit. addr. in register
	jl	r15,catstr		% Call string concatenate function
	addi	r1,r0,BUF1		% Restore BUF addr. into register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L121
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L122		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L121		% Loop if not finished
L122
	j	L118		% If cond. true jump to end if
L119
L118
	sb	BUF1(r0),r0		% Initialize buffer #1
	sb	BUF2(r0),r0		% Initialize buffer #2
	addi	r1,r0,BUF1		% Store BUF addr. into register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	jl	r15,catstr		% Call string concatenate function
	addi	r1,r0,BUF1		% Store BUF addr. into register
	subi	r14,r14,276	% Move stack pointer to calle frame
	jl	r15,L7	% Call the procedure/function
	addi	r14,r14,276	% Move stack pointer to caller frame
	addi	r2,r14,-536	% Store addr. of string result in register
	jl	r15,catstr		% Call string concatenate function
	addi	r1,r0,BUF1		% Restore BUF addr. into register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L123
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L124		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L123		% Loop if not finished
L124
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,1		% Small numbers store immediate
	add	r1,r1,r2	% Add left expr. to right expr.
	sw	-264(r14),r1	% Store value in stack frame
	j	L114		% Go to start of loop
L115
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,0		% Small numbers store immediate
	cgt	r1,r1,r2		% Is left expr. > right expr
	bz	r1,L126		% If cond. false jump to next pair
	sb	BUF1(r0),r0		% Initialize buffer #1
	sb	BUF2(r0),r0		% Initialize buffer #2
	addi	r1,r0,BUF1		% Store BUF addr. into register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	jl	r15,catstr		% Call string concatenate function
	addi	r1,r0,BUF1		% Store BUF addr. into register
	addi	r2,r0,L127		% Store string Lit. addr. in register
	jl	r15,catstr		% Call string concatenate function
	addi	r1,r0,BUF1		% Restore BUF addr. into register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L128
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L129		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L128		% Loop if not finished
L129
	j	L125		% If cond. true jump to end if
L126
L125
	lw	r1,-268(r14)
	lw	r2,-272(r14)
	lw	r3,-276(r14)
	lw	r15,-4(r14)	% Load the link into r15
	jr	r15		% Jump back to the stmt. after the call
L9
	sw	-4(r14),r15	% Store link
	sw	-268(r14),r1
	addi	r1,r0,10		% Small numbers store immediate
	sw	-276(r14),r1	% Store value in stack frame
	subi	r14,r14,268	% Move stack pointer to calle frame
	jl	r15,L5	% Call the procedure/function
	addi	r14,r14,268	% Move stack pointer to caller frame
	lw	r1,-280(r14)	% Store result in register
	sw	-264(r14),r1	% Store value in stack frame
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,1		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L131		% If cond. false jump to next pair
	addi	r1,r0,L132		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L133
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L134		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L133		% Loop if not finished
L134
	j	L130		% If cond. true jump to end if
L131
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,2		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L135		% If cond. false jump to next pair
	addi	r1,r0,L136		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L137
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L138		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L137		% Loop if not finished
L138
	j	L130		% If cond. true jump to end if
L135
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,3		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L139		% If cond. false jump to next pair
	addi	r1,r0,L140		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L141
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L142		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L141		% Loop if not finished
L142
	j	L130		% If cond. true jump to end if
L139
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,4		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L143		% If cond. false jump to next pair
	addi	r1,r0,L144		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L145
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L146		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L145		% Loop if not finished
L146
	j	L130		% If cond. true jump to end if
L143
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,5		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L147		% If cond. false jump to next pair
	addi	r1,r0,L148		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L149
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L150		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L149		% Loop if not finished
L150
	j	L130		% If cond. true jump to end if
L147
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,6		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L151		% If cond. false jump to next pair
	addi	r1,r0,L152		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L153
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L154		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L153		% Loop if not finished
L154
	j	L130		% If cond. true jump to end if
L151
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,7		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L155		% If cond. false jump to next pair
	addi	r1,r0,L156		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L157
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L158		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L157		% Loop if not finished
L158
	j	L130		% If cond. true jump to end if
L155
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,8		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L159		% If cond. false jump to next pair
	addi	r1,r0,L160		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L161
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L162		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L161		% Loop if not finished
L162
	j	L130		% If cond. true jump to end if
L159
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,9		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L163		% If cond. false jump to next pair
	addi	r1,r0,L164		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L165
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L166		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L165		% Loop if not finished
L166
	j	L130		% If cond. true jump to end if
L163
	lw	r1,-264(r14)	% Load value from stack frame
	addi	r2,r0,10		% Small numbers store immediate
	ceq	r1,r1,r2		% Is left expr. = right expr
	bz	r1,L167		% If cond. false jump to next pair
	addi	r1,r0,L168		% Store string Lit. addr. in register
	addi	r2,r0,-260		% Store offset addr. of string Var. in register
	add	r2,r2,r14		% Add the stack pointer to the register
	add	r4,r0,r0		% Initialize copy register
L169
	lb	r4,0(r1)		% Get ch from source string
	sb	0(r2),r4		% Store ch in target string
	bz	r4,L170		% Branch if end of source string
	addi	r1,r1,1		% Move to next ch in source string
	addi	r2,r2,1		% Move to next ch in target string
	j	L169		% Loop if not finished
L170
	j	L130		% If cond. true jump to end if
L167
L130
	lw	r1,-268(r14)
	lw	r15,-4(r14)	% Load the link into r15
	jr	r15		% Jump back to the stmt. after the call
L4
	res	4			% randomseed
L10
	db	0			% string reply
	res	255
	align
L11	db	"Please enter a number: ",0
	align
L14	db	13,10,"The",0
	align
L15	db	" ",0
	align
L16	db	" ",0
	align
L17	db	"the",0
	align
L18	db	".",13,10,0
	align
L19	db	13,10,"Another? ",0
	align
L22	db	"y",0
	align
L23	dw	65536		%% Store large numbers in memory
L24	dw	65536		%% Store large numbers in memory
L27	db	"dog",0
	align
L31	db	"computer",0
	align
L35	db	"student",0
	align
L39	db	"professor",0
	align
L43	db	"hippopotamus",0
	align
L47	db	"university",0
	align
L51	db	"administrator",0
	align
L55	db	"lollipop",0
	align
L59	db	"compiler",0
	align
L63	db	"assignment",0
	align
L67	db	"banana",0
	align
L72	db	"red",0
	align
L76	db	"big",0
	align
L80	db	"silly",0
	align
L84	db	"tiny",0
	align
L88	db	"slimy",0
	align
L92	db	"awful",0
	align
L96	db	"wonderful",0
	align
L100	db	"dreaded",0
	align
L104	db	"dangerous",0
	align
L108	db	"anxious",0
	align
L111	db	" ",0
	align
L120	db	", ",0
	align
L127	db	" ",0
	align
L132	db	"insults",0
	align
L136	db	"kicks",0
	align
L140	db	"licks",0
	align
L144	db	"rubs",0
	align
L148	db	"loses",0
	align
L152	db	"eats",0
	align
L156	db	"likes",0
	align
L160	db	"dismisses",0
	align
L164	db	"fails",0
	align
L168	db	"boils",0
	align
BUF1	db	0		% String buffer
	res	255
	align
BUF2	db	0		% String buffer
	res	255
	align
