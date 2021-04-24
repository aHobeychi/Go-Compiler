%==============================================================%
% File:        util.m                                          %
% Author:      Nagi B. F. MIKHAIL                              %
% Date:        April, 1995                                     %
% Description: This file contains utility routines written in  %
%              MOON's assembly to handle string operations &   %
%              Input/Output.                                   %
%==============================================================%
%
%
%--------------------------------------------------------------%
% getint                                                       %
%--------------------------------------------------------------%
% Read an integer number from stdin. Read until a non digit char
% is entered. Allowes (+) & (-) signs only as first char. The
% digits are read in ASCII and transformed to numbers.
% Entry : none.
% Exit : result -> r1
%
getint	align
	add	r1,r0,r0		% Initialize input register (r1 = 0)
	add	r2,r0,r0		% Initialize buffer's index i
	add	r4,r0,r0		% Initialize sign flag
getint1	getc	r1			% Input ch
	ceqi	r3,r1,43		% ch = '+' ?
	bnz	r3,getint1		% Branch if true (ch = '+')
	ceqi	r3,r1,45		% ch = '-' ?
	bz	r3,getint2		% Branch if false (ch != '-')
	addi	r4,r0,1			% Set sign flag to true
	j	getint1			% Loop to get next ch
getint2	clti	r3,r1,48		% ch < '0' ?
	bnz	r3,getint3		% Branch if true (ch < '0')
	cgti	r3,r1,57		% ch > '9' ?
	bnz	r3,getint3		% Branch if true (ch > '9')
	sb	getint9(r2),r1		% Store ch in buffer
	addi	r2,r2,1			% i++
	j	getint1			% Loop if not finished
getint3	sb	getint9(r2),r0		% Store end of string in buffer (ch = '\0')
	add	r2,r0,r0		% i = 0
	add	r1,r0,r0		% N = 0
	add	r3,r0,r0		% Initialize r3
getint4	lb	r3,getint9(r2)		% Load ch from buffer
	bz	r3,getint5		% Branch if end of string (ch = '\0')
	subi	r3,r3,48		% Convert to number (M)
	muli	r1,r1,10		% N *= 10
	add	r1,r1,r3		% N += M
	addi	r2,r2,1			% i++
	j	getint4			% Loop if not finished
getint5	bz	r4,getint6		% Branch if sign flag false
	sub	r1,r0,r1		% N = -N
getint6	jr	r15			% Return to the caller
getint9	res	12			% Local buffer (12 bytes)
	align
%
%
%--------------------------------------------------------------%
% putint                                                       %
%--------------------------------------------------------------%
% Write an integer number to stdout. Transform the number into
% ASCII string taking the sign into account.
% Entry : integer number -> r1
% Exit : none.
%
putint	align
	add	r2,r0,r0		% Initialize buffer's index i
	cge	r3,r1,r0		% True if N >= 0
	bnz	r3,putint1		% Branch if True (N >= 0)
	sub	r1,r0,r1		% N = -N
putint1	modi	r4,r1,10		% Rightmost digit
	addi	r4,r4,48		% Convert to ch
	divi	r1,r1,10		% Truncate rightmost digit
	sb	putint9(r2),r4		% Store ch in buffer
	addi	r2,r2,1			% i++
	bnz	r1,putint1		% Loop if not finished
	bnz	r3,putint2		% Branch if True (N >= 0)
	addi	r3,r0,45
	sb	putint9(r2),r3		% Store '-' in buffer
	addi	r2,r2,1			% i++
	add	r1,r0,r0		% Initialize output register (r1 = 0)
putint2	subi	r2,r2,1			% i--
	lb	r1,putint9(r2)		% Load ch from buffer
	putc	r1			% Output ch
	bnz	r2,putint2		% Loop if not finished
	jr	r15			% return to the caller
putint9	res	12			% loacl buffer (12 bytes)
	align
%
%
%--------------------------------------------------------------%
% putstr                                                       %
%--------------------------------------------------------------%
% Write a string to stdout. Write char by char until end of 
% string.
% Entry : address of string -> r1
% Exit : none.
% 
putstr	align
	add	r2,r0,r0		% Initialize output register (r2 = 0)
putstr1	lb	r2,0(r1)		% Load ch from buffer
	bz	r2,putstr2		% Branch if end of string (ch = '\0')
	putc	r2			% Output ch
	addi	r1,r1,1			% i++
	j	putstr1			% Loop if not finished
putstr2	jr	r15			% Return to the caller
	align
%
%
%--------------------------------------------------------------%
% putsub                                                       %
%--------------------------------------------------------------% 
% Write a substring to stdout. Write char by char until end of 
% substring.
% Entry : address of substring -> r1
%         address of end of substring -> r3
% Exit : none.
% 
putsub	align
	add	r2,r0,r0		% Initialize output register (r2 = 0)
putsub1	lb	r2,0(r1)		% Load ch from buffer
	putc	r2			% Output ch
	ceq	r2,r1,r3		% Current pos. is end of substring ?
	bnz	r2,putsub2		% Branch if end of substring
	addi	r1,r1,1			% i++
	j	putsub1			% Loop if not finished
putsub2	jr	r15			% Return to the caller
	align
%
%
%--------------------------------------------------------------%
% getstr                                                       %
%--------------------------------------------------------------%
% Read a string from stdin. Read char by char until CR but do 
% not store CR.
% Entry : address of string var. -> r1
% Exit : address of string var. -> r1
%
getstr	align
	add	r2,r0,r0		% Initialize input register (r2 = 0)
	getc	r2			% Input ch
	ceqi	r3,r2,10		% ch = CR ?
	bnz	r3,getstr1		% branch if true (ch = CR)
	sb	0(r1),r2		% Store ch in buffer
	addi	r1,r1,1			% i++
	j	getstr			% Loop if not finished
getstr1	sb	0(r1),r0		% Store end of string (ch = '\0')
	jr	r15			% Return to the caller
	align
%
%
%--------------------------------------------------------------%
% catstr                                                       %
%--------------------------------------------------------------%
% Append one string to another.
% Entry : address of 1st string -> r1
%         address of 2nd string -> r2
% Exit : address of concatenated strings -> r1
%
catstr	align
	add	r3,r0,r0		% r3 = 0
catstr1	lb	r3,0(r1)		% Load ch from 1st string
	bz	r3,catstr2		% Branch if end of string (ch = `\0')
	addi	r1,r1,1			% i++
	j	catstr1			% Loop if not end of 1st string
catstr2	lb	r3,0(r2)		% Load ch from 2nd string
	sb	0(r1),r3		% Store ch at the end of 1st string
	bz	r3,catstr3		% Branch if end of 2nd string
	addi	r2,r2,1			% j++
	addi	r1,r1,1			% i++
	j	catstr2			% Loop if not end of 2nd string
catstr3	jr	r15			% Return to the caller
	align
%
%
%--------------------------------------------------------------%
% catsub                                                       %
%--------------------------------------------------------------%
% Append a substring to a string.
% Entry : address of string -> r1
%         address of substring -> r2
%         address of end of substring -> r4
% Exit : address of concatenated strings -> r1
% 
catsub	align
	add	r3,r0,r0		% r3 = 0
catsub1	lb	r3,0(r1)		% Load ch from 1st string
	bz	r3,catsub2		% Branch if end of string (ch = `\0')
	addi	r1,r1,1			% i++
	j	catsub1			% Loop if not end of 1st string
catsub2	lb	r3,0(r2)		% Load ch from substring
	sb	0(r1),r3		% Store ch at the end of 1st string
	bz	r3,catsub4		% Branch if end of string ch
	addi	r1,r1,1			% i++
	ceq	r3,r2,r4		% Check end of substring
	bnz	r3,catsub3		% Branch if end of substring
	addi	r2,r2,1			% j++
	j	catsub2			% Loop if not end of substring
catsub3	sb	0(r1),r0		% Store end of string
catsub4	jr	r15			% Return to the caller
	align
%
%
%--------------------------------------------------------------%
% eqstr                                                        %
%--------------------------------------------------------------%
% Check if string1 = string2.
% Entry : address of string1 -> r1
%         address of string2 -> r2
% Exit : r1 = 1 (true)
%        r1 = 0 (false)
%
eqstr	align
	add	r3,r0,r0		% Initialize r3 (r3 = 0)
	add	r4,r0,r0		% Initialize r4 (r4 = 0)
eqstr1	lb	r3,0(r1)		% Load ch1 from 1st string
	lb	r4,0(r2)		% Load ch2 from 2nd string
	ceq	r5,r3,r4		% ch1 = ch2 ?
	bz	r5,eqstr2		% Branch if false
	ceq	r5,r3,r0		% ch1 = '\0' ?
	bnz	r5,eqstr3		% Branch if true
	addi	r1,r1,1			% i++
	addi	r2,r2,1			% j++
	j	eqstr1			% Loop if not finished
eqstr2	add	r1,r0,r0		% Return false
	j	eqstr4
eqstr3	addi	r1,r0,1			% Return true
eqstr4	jr	r15			% Return to the caller
	align
%
%
%--------------------------------------------------------------%
% neqstr                                                       %
%--------------------------------------------------------------%
% Check if string1 != string2.
% Entry : address of string1 -> r1
%         address of string2 -> r2
% Exit : r1 = 1 (true)
%        r1 = 0 (false)
% 
neqstr	align
	add	r3,r0,r0		% Initialize r3 (r3 = 0)
	add	r4,r0,r0		% Initialize r4 (r4 = 0)
neqstr1	lb	r3,0(r1)		% Load ch1 from 1st string
	lb	r4,0(r2)		% Load ch2 from 2nd string
	cne	r5,r3,r4		% ch1 != ch2 ?
	bnz	r5,neqstr3		% Branch if true
	ceq	r5,r3,r0		% ch1 = '\0' ?
	bnz	r5,neqstr2		% Branch if true
	addi	r1,r1,1			% i++
	addi	r2,r2,1			% j++
	j	neqstr1			% Loop if not finished
neqstr2	add	r1,r0,r0		% Return false
	j	neqstr4
neqstr3	addi	r1,r0,1			% Return true
neqstr4	jr	r15			% Return to the caller
	align
%
%
%--------------------------------------------------------------%
% gtstr                                                        %
%--------------------------------------------------------------%
% Check if string1 > string2.
% Entry : address of string1 -> r1
%         address of string2 -> r2
% Exit : r1 = 1 (true)
%        r1 = 0 (false)
%
gtstr	align
	add	r3,r0,r0		% Initialize r3 (r3 = 0)
	add	r4,r0,r0		% Initialize r4 (r4 = 0)
gtstr1	lb	r3,0(r1)		% Load ch1 from 1st string
	lb	r4,0(r2)		% Load ch2 from 2nd string
	cgt	r5,r3,r4		% ch1 > ch2 ?
	bnz	r5,gtstr3		% Branch if true
	clt	r5,r3,r4		% ch1 < ch2 ?
	bnz	r5,gtstr2		% Branch if true
	ceq	r5,r3,r0		% ch1 = '\0' ?
	bnz	r5,gtstr2		% Branch if true
	addi	r1,r1,1			% i++
	addi	r2,r2,1			% j++
	j	gtstr1			% Loop if not finished
gtstr2	add	r1,r0,r0		% Return false
	j	gtstr4
gtstr3	addi	r1,r0,1			% Return true
gtstr4	jr	r15			% Return to the caller
	align
%
%
%--------------------------------------------------------------%
% gtestr                                                       %
%--------------------------------------------------------------%
% Check if string1 >= string2.
% Entry : address of string1 -> r1
%         address of string2 -> r2
% Exit : r1 = 1 (true)
%        r1 = 0 (false)
%
gtestr	align
	add	r3,r0,r0		% Initialize r3 (r3 = 0)
	add	r4,r0,r0		% Initialize r4 (r4 = 0)
gtestr1	lb	r3,0(r1)		% Load ch1 from 1st string
	lb	r4,0(r2)		% Load ch2 from 2nd string
	cgt	r5,r3,r4		% ch1 > ch2 ?
	bnz	r5,gtestr3		% Branch if true
	clt	r5,r3,r4		% ch1 < ch2 ?
	bnz	r5,gtestr2		% Branch if true
	ceq	r5,r3,r0		% ch1 = '\0' ?
	bnz	r5,gtestr3		% Branch if true
	addi	r1,r1,1			% i++
	addi	r2,r2,1			% j++
	j	gtestr1			% Loop if not finished
gtestr2	add	r1,r0,r0		% Return false
	j	gtestr4
gtestr3	addi	r1,r0,1			% Return true
gtestr4	jr	r15			% Return to the caller
	align
%
%
%--------------------------------------------------------------%
% ltstr                                                        %
%--------------------------------------------------------------%
% Check if string1 < string2.
% Entry : address of string1 -> r1
%         address of string2 -> r2
% Exit : r1 = 1 (true)
%        r1 = 0 (false)
%
ltstr	align
	add	r3,r0,r0		% Initialize r3 (r3 = 0)
	add	r4,r0,r0		% Initialize r4 (r4 = 0)
ltstr1	lb	r3,0(r1)		% Load ch1 from 1st string
	lb	r4,0(r2)		% Load ch2 from 2nd string
	clt	r5,r3,r4		% ch1 < ch2 ?
	bnz	r5,ltstr3		% Branch if true
	cgt	r5,r3,r4		% ch1 > ch2 ?
	bnz	r5,ltstr2		% Branch if true
	ceq	r5,r3,r0		% ch1 = '\0' ?
	bnz	r5,ltstr2		% Branch if true
	addi	r1,r1,1			% i++
	addi	r2,r2,1			% j++
	j	ltstr1			% Loop if not finished
ltstr2	add	r1,r0,r0		% Return false
	j	ltstr4
ltstr3	addi	r1,r0,1			% Return true
ltstr4	jr	r15			% Return to the caller
	align
%
%
%--------------------------------------------------------------%
% ltestr                                                       %
%--------------------------------------------------------------%
% Check if string1 <= string2.
% Entry : address of string1 -> r1
%         address of string2 -> r2
% Exit : r1 = 1 (true)
%        r1 = 0 (false)
%
ltestr	align
	add	r3,r0,r0		% Initialize r3 (r3 = 0)
	add	r4,r0,r0		% Initialize r4 (r4 = 0)
ltestr1	lb	r3,0(r1)		% Load ch1 from 1st string
	lb	r4,0(r2)		% Load ch2 from 2nd string
	clt	r5,r3,r4		% ch1 < ch2 ?
	bnz	r5,ltestr3		% Branch if true
	cgt	r5,r3,r4		% ch1 > ch2 ?
	bnz	r5,ltestr2		% Branch if true
	ceq	r5,r3,r0		% ch1 = '\0' ?
	bnz	r5,ltestr3		% Branch if true
	addi	r1,r1,1			% i++
	addi	r2,r2,1			% j++
	j	ltestr1			% Loop if not finished
ltestr2	add	r1,r0,r0		% Return false
	j	ltestr4
ltestr3	addi	r1,r0,1			% Return true
ltestr4	jr	r15			% Return to the caller
	align
%
%
%--------------------------------------------------------------%
% lenstr                                                       %
%--------------------------------------------------------------%
% Returns length of string. '\n' is concidered one char.
% Entry : address of string -> -8(r14)
% Exit : result -> -12(r14)
%
lenstr	align
	sw	-4(r14),r15		% Store link in stack
	sw	-16(r14),r1		% Save registers' old values on stack
	sw	-20(r14),r2
	sw	-24(r14),r3
	sw	-28(r14),r4
	lw	r1,-8(r14)		% Retrieve argument from stack
	add	r2,r0,r0		% Initialize length counter (len = 0)
	add	r3,r0,r0		% Initialize ch holder (r3 = 0)
lenstr1	lb	r3,0(r1)		% Load ch from string
	ceq	r4,r3,r0		% ch = '\0' ?
	bnz	r4,lenstr2		% Branch if true
	addi	r1,r1,1			% i++
	ceqi	r4,r3,13		% ch = LF ?
	bnz	r4,lenstr1		% Skip ch and loop
	addi	r2,r2,1			% len++
	j	lenstr1			% Loop if not finished
lenstr2	add	r1,r0,r2		% Return length of the string
	sw	-12(r14),r1		% Save result on stack
	lw	r1,-16(r14)		% Reset registers to their old values
	lw	r2,-20(r14)
	lw	r3,-24(r14)
	lw	r4,-28(r14)
	lw	r15,-4(r14)		% Load link from stack
	jr	r15			% Return to the caller
	align
%
%
%--------------------------------------------------------------%
% ordstr                                                       %
%--------------------------------------------------------------%
% Returns ASCII number of the 1st char of a string.
% Entry : address of string -> -8(r14)
% Exit : result -> -12(r14)
%
ordstr	align
	sw	-4(r14),r15		% Store link in stack
	sw	-16(r14),r1		% Save registers' old values on stack
	sw	-20(r14),r2
	lw	r1,-8(r14)		% Retrieve argument from stack
	add	r2,r0,r0		% Initialize ch holder (r2 = 0)
	lb	r2,0(r1)		% Load ch from string
	add	r1,r0,r2		% Return ascii code of ch
	sw	-12(r14),r1		% Save result on stack
	lw	r1,-16(r14)		% Reset registers to their old values
	lw	r2,-20(r14)
	lw	r15,-4(r14)		% Load link from stack
	jr	r15			% Return to the caller
	align
%
%
%--------------------------------------------------------------%
% substr                                                       %
%--------------------------------------------------------------%
% Returns the start & end pos. of a substring. '\n' is concidered
% one char. Returns empty string in case of illegal start and/or
% substring pos.
% Entry : address of string -> -8(r14)
%         substring start pos. -> -12(r14)
%         substring end pos. -> -16(r14)
% Exit : address of substring start pos. -> -20(r14)
%        address of substring end pos. -> -24(r14)
%
substr	align
	sw	-4(r14),r15		% Store link in stack
	sw	-28(r14),r1		% Save registers' old values on stack
	sw	-32(r14),r2
	sw	-36(r14),r3
	sw	-40(r14),r4
	sw	-44(r14),r5
	sw	-48(r14),r6
	lw	r1,-8(r14)		% Retrieve arguments from stack
	lw	r2,-12(r14)
	lw	r3,-16(r14)
	add	r5,r0,r0		% Initialize copy reg (r5 = 0)
	cle	r4,r2,r0		% N < 0 ?
	bnz	r4,substr4		% Branch if true
	cle	r4,r2,r3		% M <= N ?
	bz	r4,substr4		% Branch if true
	add	r2,r2,r1		% Start pos. of the substr. in memory
	subi	r2,r2,1
	add	r3,r3,r1		% End pos. of substr. in memory
	subi	r3,r3,1
substr1	lb	r5,0(r1)		% Load ch in copy register starting from ch1
	ceq	r6,r5,r0		% ch = '\0' ?
	bnz	r6,substr4		% Branch if true
	ceq	r4,r2,r1		% Current pos. = Start pos. ?
	bnz	r4,substr2		% Branch if true
	ceqi	r6,r5,13		% ch = LF ?
	bz	r6,substr6		% Branch if false
	addi	r2,r2,1			% Move start pos
	addi	r3,r3,1			% Move end pos
substr6	addi	r1,r1,1			% Move current pos.
	j	substr1			% Loop if current pos. != start pos.
substr2	lb	r5,0(r2)		% Load ch in copy register starting from star pos.
	ceq	r6,r5,r0		% ch = '\0' ?
	bnz	r6,substr4		% Branch if true
	ceqi	r6,r5,13		% ch = LF ?
	bz	r6,substr7		% Branch if false
	addi	r3,r3,1			% Move end pos
substr7	ceq	r4,r3,r2		% Current pos. = end pos. ?
	bnz	r4,substr3		% Branch if true
	addi	r2,r2,1			% Move current pos.
	j	substr2			% Loop if current pos. != end pos.
substr3	lb	r5,0(r3)		% Load ch of end pos. in copy register
	ceq	r6,r5,r0		% ch = '\0' ?
	bnz	r6,substr4		% Branch if true
	j	substr5
substr4	addi	r1,r0,EMPTY		% Start & end pos. pointing to ""
	addi	r3,r0,EMPTY
substr5	sw	-20(r14),r1		% Store results on stack
	sw	-24(r14),r3
	lw	r1,-28(r14)		% Reset regiters to old values
	lw	r2,-32(r14)
	lw	r3,-36(r14)
	lw	r4,-40(r14)
	lw	r5,-44(r14)
	lw	r6,-48(r14)
	lw	r15,-4(r14)		% Load link from stack
	jr	r15			% Return to caller
	align
EMPTY	db	0			% Empty string
	align
%
%
%-------------------------  End of file -----------------------%

