/* Moon Simulator.
 *
 * Author:			Peter Grogono
 * Last modified:	30 January 1995
 *
 * This source file was created using a tab width of 4 characters.
 */

#include <stdio.h>
#include <ctype.h>
#include <string.h>
#include <math.h>

/* Notes on definitions.
 * The memory occupies about 6 * MEMSIZE bytes.  Enlarging it may create
 * problems for architectures with memory restrictions, such as the PC.
 *
 * The simulator is not completely independent of the underlying
 * processor.  Known variations include:
 *    - trace output depends on whether the host is big- or little endian.
 *    - the MOON `sr' op uses the C `>>' operator.
 */

#define MEMSIZE	    4000	/* Memory size in 4-byte words. */
#define MAXREG		  16	/* Number of registers. */
#define MAXINFILES	  20	/* Restricts # input files. */
#define MAXNAMELEN	  60	/* Restricts path/file name length. */
#define BUFLEN		 255	/* Scanner input buffer. */
#define TOKLEN		 255	/* Scanner token buffer. */
#define RANGE		  20	/* Determines size of memory dump. */

#define FALSE		0
#define TRUE		1
#define BYTE 		unsigned char

/* The following function reports syntax errors.  It is declared here so
 *  that it can be used by memory functions; its definition appears in the
 * loader section.
 */
void syntaxerror (char *message);

/* The following function reports run-time errors.  It is declared here so
 *  that it can be used by memory functions; its definition appears in the
 * execution section.
 */
void runtimeerror (char *message);

/************************ MEMORY ********************************************/

/* A word of memory contains an instruction (A or B format), four bytes,
 * or a 32-bit signed integer.  The simulated memory also contains two
 * flags, to tell whether the word contains an instruction and whether
 * it is a breakpoint.
 *
 * The simulated memory is an array of four-byte words.  Simulated
 * addresses are byte addresses.  Consequently, all addresses are
 * shifted right 2 bits before being used.  Only the memory module
 * should know about this.  Conventionally, <addr> is a byte address
 * and <wordaddr> is a word address.
 */

typedef union {
	struct {
		unsigned op:6;
		unsigned ri:4;
		unsigned rj:4;
		unsigned rk:4;
		unsigned :14;
		} fmta;
	struct {
		unsigned op:6;
		unsigned ri:4;
		unsigned rj:4;
		unsigned   :2;
		int       k:16;     /* Unsigned is not portable! */
		} fmtb;
	BYTE byts [4];
	long data;
	} wordtype;

/* Each component of the memory array contains:
 *  -  word:		The contents of simulated memory.
 *  -  cont:        `a' => format A instruction
 *					`b' => format B instruction
 *					'd' => data
 *					`u' => undefined
 *  -  breakpoint:	True if this is a breakpoint.
 */

struct {
	wordtype word;
	char cont;
	short breakpoint;
	} mem [MEMSIZE];

long ic;				/* Instruction counter: contains address of
							instruction that will be executed next. */
wordtype ir; 			/* Instruction register: contains the instruction
							that will be executed next. */
long mar;				/* Memory address register;
							stores word address of last access. */
wordtype mdr;			/* Memory data register;
							stores result of last access. */
long entrypoint = -1;	/* Address of first instruction */
long cycles = 0;		/* Counts memory cycles. */

/* Initialize the memory: value=zero, kind=undef, no breakpoint.
 * Set "hardware" addresses to an illegal value.
 */
void initmem () {
	long wordaddr;
	for (wordaddr = 0; wordaddr < MEMSIZE; wordaddr++) {
		mem[wordaddr].word.data = 0;
		mem[wordaddr].cont = 'u';
		mem[wordaddr].breakpoint = FALSE;
		ic = -1;
		mar = -1;
		}
	}

/* Report a run-time error if an illegal address is used. */
short outofrange (long addr) {
	if (addr < 0 || (addr >> 2) > MEMSIZE) {
		runtimeerror("address error");
		return 1;
		}
	return 0;
	}

/* Report a run-time error if a memory word access is not on a
 * four-byte boundary. */
short misaligned (long addr) {
	if (addr & 3) {
		runtimeerror("alignment error");
		return 1;
		}
	return 0;
	}

/* Fetch the instruction at <ic>, store it in <ir>, and increment <ic>. */
short fetch () {
	short cont;
	if (outofrange(ic)) {
		ir.data = 0;
		return 0;
		}
	cont = mem[ic >> 2].cont;
	if (!(cont == 'a' || cont == 'b')) {
		runtimeerror("illegal instruction");
		ir.data = 0;
		return 0;
		}
	ir = mem[ic >> 2].word;
	ic += 4;
	cycles += 10;
	return cont;
	}

/* Fetch a data word from memory and return it. */
long getmemword (long addr) {
	long wordaddr;
	if (outofrange(addr) || misaligned(addr))
		return 0;
	wordaddr = addr >> 2;
	if (wordaddr == mar)
		cycles += 1;
	else {
		mar = wordaddr;
		mdr = mem[wordaddr].word;
		cycles += 10;
		}
	return mdr.data;
	}

/* Store a word in memory. */
void putmemword (long addr, long data) {
	long wordaddr;
	if (outofrange(addr) || misaligned(addr))
		return;
	wordaddr = addr >> 2;
	if (mem[wordaddr].cont == 'a' || mem[wordaddr].cont == 'b') {
		runtimeerror("overwriting instructions");
		return;
		}
	mdr.data = data;
	mar = wordaddr;
	mem[mar].word = mdr;
	mem[mar].cont = 'd';
	cycles += 10;
	return;
	}

/* Fetch a byte from memory. */
BYTE getmembyte (long addr) {
	long wordaddr = addr >> 2;
	short offset = addr & 3;
	if (outofrange(addr))
		return 0;
	if (wordaddr == mar)
		cycles += 1;
	else {
		mar = wordaddr;
		cycles += 10;
		}
	mdr = mem[mar].word;
	return mdr.byts[offset];
	}

/* Store a byte in memory. */
void putmembyte (long addr, BYTE byt) {
	long wordaddr = addr >> 2;
	short offset = addr & 3;
	if (outofrange(addr))
		return;
	if (mem[wordaddr].cont == 'a' || mem[wordaddr].cont == 'b') {
		runtimeerror("overwriting instructions");
		return;
		}
	mem[wordaddr].word.byts[offset] = byt & 255;
	mem[wordaddr].cont = 'd';
	return;
	}

/* Store an instruction in memory. Used only by loader. */
void putmeminstr (long addr, wordtype word, char cont) {
	if (addr & 3)
		syntaxerror("alignment error");
	else {
		long wordaddr = addr >> 2;
		mem[wordaddr].word = word;
		mem[wordaddr].cont = cont;
		}
	}

/* Store a character in memory. Used only by loader. */
void putmemchar (long addr, short byte, char cont) {
	long wordaddr = addr >> 2;
	mem[wordaddr].word.byts[addr & 3] = byte;
	mem[wordaddr].cont = cont;
	}

/********************** REGISTERS *******************************************/

/* There are sixteen registers, numbered 0 through 15.  Each register
 *  is a 32-bit word.  Functions check the register address although an
 * error would indicate a fault in the loader rather than the simulator.
 * Register 0 is always 0.
 */

long regs [MAXREG];

/* Fetch the value of a register. */
long fetchreg (unsigned short regnum) {
	if (regnum > 15) {
		runtimeerror("simulator error (illegal register code)");
		return 0;
		}
	return regs[regnum];
	}

/* Store a value in a register. */
void storereg (unsigned short regnum, long data) {
	if (regnum > 15)
		runtimeerror("simulator error (illegal register code)");
	if (regnum > 0)
		regs[regnum] = data;
	}

/******************** INSTRUCTION CODES *************************************/

/* Instruction codes and the corresponding strings; obviously, these two
 * declarations should correspond.  The directives have codes although
 * they are not stored in the memory.  Note that 0 is an illegal
 * instruction.
 */

enum optype {
	bad, lw, lb, sw, sb, add, sub, mul, div, mod,
	and, or, not, ceq, cne, clt, cle, cgt, cge,
	addi, subi, muli, divi, modi, andi, ori,
	ceqi, cnei, clti, clei, cgti, cgei, sl, sr,
	gtc, ptc, bz, bnz, j, jr, jl, jlr, nop, hlt,
	entry, align, org, dw, db, res,
	last
	};

char opnames [last] [6] = {
	"", "lw", "lb", "sw", "sb", "add", "sub", "mul", "div", "mod",
	"and", "or", "not", "ceq", "cne", "clt", "cle", "cgt", "cge",
	"addi", "subi", "muli", "divi", "modi", "andi", "ori",
	"ceqi", "cnei", "clti", "clei", "cgti", "cgei", "sl", "sr",
	"getc", "putc", "bz", "bnz", "j", "jr", "jl", "jlr", "nop", "hlt",
	"entry", "align", "org", "dw", "db", "res"
	};

void showfmta (long addr, wordtype word) {
	char *opcode = opnames[word.fmta.op];
	switch (word.fmta.op) {
		/* Operands Ri, Rj, Rk */
		case add:
		case sub:
		case mul:
		case div:
		case mod:
		case and:
		case or:
		case ceq:
		case cne:
		case clt:
		case cle:
		case cgt:
		case cge:
			printf("%5ld %-6s   r%d, r%d, r%d",
				addr, opcode, word.fmta.ri,
				word.fmta.rj, word.fmta.rk);
			break;

		/* Operands Ri, Rj */
		case not:
		case jlr:
			printf("%5ld %-6s   r%d, r%d",
				addr, opcode, word.fmta.ri, word.fmta.rj);
			break;


		/* No operands */
		case nop:
		case hlt:
			printf("%5ld %-6s",
				addr, opcode);
			break;
		}
	}

void showfmtb (long addr, wordtype word) {
	char *opcode = opnames[word.fmtb.op];
	switch (word.fmtb.op) {

		/* Operands Ri, K(Rj) */
		case lw:
		case lb:
			printf("%5ld %-6s   r%d, %d(r%d)",
				addr, opcode, word.fmtb.ri,
				word.fmtb.k, word.fmtb.rj);
			break;

		/* Operands K(Rj), Ri */
		case sw:
		case sb:
			printf("%5ld %-6s   %d(r%d), r%d",
				addr, opcode, word.fmtb.k,
				word.fmtb.rj, word.fmtb.ri);
			break;

		/* Operands Ri, Rj, K */
		case addi:
		case subi:
		case muli:
		case divi:
		case modi:
		case andi:
		case ori:
		case ceqi:
		case cnei:
		case clti:
		case clei:
		case cgti:
		case cgei:
			printf("%5ld %-6s   r%d, r%d, %d",
				addr, opcode, word.fmtb.ri,
				word.fmtb.rj, word.fmtb.k);
			break;

		/* Operands Ri, K */
		case sl:
		case sr:
		case bz:
		case bnz:
		case jl:
			printf("%5ld %-6s   r%d, %d",
				addr, opcode, word.fmtb.ri, word.fmtb.k);
			break;

		/* Operands Ri */
		case gtc:
		case ptc:
		case jr:
			printf("%5ld %-6s   r%d",
				addr, opcode, word.fmtb.ri);
			break;

		/* Operands K */
		case j:
			printf("%5ld %-6s   %d",
				addr, opcode, word.fmtb.k);
			break;
		}
	}

/* Convert a word to a string of 4 characters. Non-graphics to ".".
 * Exact output depends on whether the host is big-endian or
 * little-endian.
 */
char *wordtochars (char *buf, wordtype word) {
	int i;
	for (i = 0; i < 4; i++) {
		char c = word.byts[i];
		if (32 <= c && c <= 126)
			buf[i] = c;
		else buf[i] = '.';
		}
	buf[4] = '\0';
	return buf;
	}

/* Display one word of memory. */
void showword (long addr) {
	char charbuf [5];
	long wordaddr = addr >> 2;
	wordtype word = mem[wordaddr].word;
	if (addr & 3) {
		printf("Internal error: bad address!\n");
		exit(1);
		}
	switch (mem[wordaddr].cont) {
		case 'a':
			showfmta(addr, word);
			break;
		case 'b':
			showfmtb(addr, word);
			break;
		case 'd':
			printf("%5ld  %08lX  %s  %4ld", addr, word.data,
				wordtochars(charbuf, word), word.data);
			break;
		case 'u':
			printf("%5ld  ??", addr);
			break;
		}
	}

/****************************** SYMBOLS *************************************/

/* There is one <symnode> for each symbol.  It records the name and value
 * of the symbol and how often it has been defined (once is correct).
 * The <symnode> also contains a pointer to a list of <usenodes>'s, each
 * of which contains an address where the value of the symbol is used.
 */

struct symnode {
	char *name;
	long val;
	short defs;
	struct usenode *uses;
	struct symnode *next;
	};

struct usenode {
	long addr;
	struct usenode *next;
	};

/* The base of the symbol table. */
struct symnode *symbols = NULL;

/* Return a pointer to a symbol entry.  This always succeeds, because
 * it creates a new entry if it can't find a matching entry.
 */
struct symnode *findsymbol (char *name) {
	struct symnode *p = symbols;
	while (p) {
		if (!strcmp(name, p->name))
			return p;
		p = p->next;
		}
	/* No entry exists, so make one. */
	p = (struct symnode *) malloc (sizeof(struct symnode));
	if (p == NULL) {
		printf("No more memory!\n");
		exit(1);
		}
	p->name = (char *) malloc(strlen(name) + 1);
	strcpy(p->name, name);
	p->val = 0;
	p->defs = 0;
	p->uses = NULL;
	p->next = symbols;
	symbols = p;
	return p;
	}

/* Define a symbol.  That is, associate the value <val> with
 * the symbol <name>.
 */
void defsymbol (char *name, long val) {
	struct symnode *p = findsymbol(name);
	p->val = val;
	p->defs++;
	}

/* Use a symbol. That is, record the fact that the symbol <name>
 * must be stored at <addr>.
 */
void usesymbol (char *name, long addr) {
	struct symnode *p = findsymbol(name);
	struct usenode *u = (struct usenode *) malloc (sizeof (struct usenode));
	if (u == NULL) {
		printf("No more memory!\n");
		exit(1);
		}
	u->addr = addr;
	u->next = p->uses;
	p->uses = u;
	}

/* Return the value of a symbol, or -1 if it doesn't exist. */
long getsymbolval (char *name) {
	struct symnode *p = symbols;
	while (p) {
		if (!strcmp(name, p->name))
			return (p->val);
		p = p->next;
		}
	return -1;
	}

/* Display all symbols and their uses. */
void showsymbols () {
	short count = 0;
	char reply [80];
	struct symnode *p = symbols;
	while (p) {
		struct usenode *u = p->uses;
		printf("%-8s = %4ld  Used at: ", p->name, p->val);
		while (u) {
			printf("%ld ", u->addr);
			u = u->next;
			}
		printf("\n");
		p = p->next;
		if (++count > 20) {
			printf("Press enter to continue");
			gets(reply);
			count = 0;
			}
		}
	}

/* Check symbol list for errors and return error count. */
int checksymbols () {
	int errors = 0;
	struct symnode *p = symbols;
	while (p) {
		if (p->defs == 0) {
			printf("Undefined symbol: %s.\n", p->name);
			errors++;
			}
		else if (p->defs > 1) {
			printf("Redefined symbol: %s.\n", p->name);
			errors++;
			}
		p = p->next;
		}
	return errors;
	}

/* Store symbols at their respective locations. */
void storesymbols () {
	struct symnode *p = symbols;
	while (p) {
		struct usenode *u = p->uses;
		while (u) {
			long wordaddr = (u -> addr) >> 2;
			switch (mem[wordaddr].cont) {
				case 'b':
					mem[wordaddr].word.fmtb.k = (int) p->val;
					break;
				case 'd':
					mem[wordaddr].word.data = p->val;
					break;
				default:
					printf("Symbol storage error!\n");
					break;
				}
			u = u->next;
			}
		p = p->next;
		}
	}

/***************************** EXECUTION ************************************/

short newreg;		/* Address of a register that has changed */
long newmem;		/* Address of a memory location that has changed */
short running;		/* True if the processor is running, false after errors */
long numsteps;		/* Number of instructions executed in trace mode. */

/* Report a run-time error and stop the program. */
void runtimeerror (char *message) {
	printf("\n%5ld Run-time error: %s.\n", ic, message);
	running = FALSE;
	}

/* Execute the instruction at address <ic>. */
void execinstr (short tracing) {
	long addr, w1, w2, k, rk;
	short cont = fetch();		/* Move next instruction to `ir'. */
	int ch;
	if (!running)
		return;
	newreg = -1;
	newmem = -1;
	switch (cont) {

		/* Format A instructions with register operands. */
		case 'a':
			switch (ir.fmta.op) {

				/* add Ri, Rj, Rk */
				case add:
					storereg(ir.fmta.ri,
						fetchreg(ir.fmta.rj) + fetchreg(ir.fmta.rk));
					newreg = ir.fmta.ri;
					break;

				/* sub Ri, Rj, Rk */
				case sub:
					storereg(ir.fmta.ri,
						fetchreg(ir.fmta.rj) - fetchreg(ir.fmta.rk));
					newreg = ir.fmta.ri;
					break;

				/* mul Ri, Rj, Rk */
				case mul:
					storereg(ir.fmta.ri,
						fetchreg(ir.fmta.rj) * fetchreg(ir.fmta.rk));
					newreg = ir.fmta.ri;
					break;

				/* div Ri, Rj, Rk */
				case div:
					rk = fetchreg(ir.fmta.rk);
					if (rk == 0)
						runtimeerror("division by zero");
					else {
						storereg(ir.fmta.ri,
							fetchreg(ir.fmta.rj) / rk);
						newreg = ir.fmta.ri;
						}
					break;

				/* mod Ri, Rj, Rk */
				case mod:
					rk = fetchreg(ir.fmta.rk);
					if (rk == 0)
						runtimeerror("modulus with zero operand");
					else {
						storereg(ir.fmta.ri,
							fetchreg(ir.fmta.rj) % rk);
						newreg = ir.fmta.ri;
						}
					break;

				/* and Ri, Rj, Rk  (32-bit logical AND) */
				case and:
					storereg(ir.fmta.ri,
						fetchreg(ir.fmta.rj) & fetchreg(ir.fmta.rk));
					newreg = ir.fmta.ri;
					break;

				/* or Ri, Rj, Rk   (32-bit logical OR) */
				case or:
					storereg(ir.fmta.ri,
						fetchreg(ir.fmta.rj) | fetchreg(ir.fmta.rk));
					newreg = ir.fmta.ri;
					break;

				/* ceq Ri, Rj, Rk  (Rj = Rk) */
				case ceq:
					storereg(ir.fmta.ri,
						fetchreg(ir.fmta.rj) == fetchreg(ir.fmta.rk));
					newreg = ir.fmta.ri;
					break;

				/* cne Ri, Rj, Rk */
				case cne:
					storereg(ir.fmta.ri,
						fetchreg(ir.fmta.rj) != fetchreg(ir.fmta.rk));
					newreg = ir.fmta.ri;
					break;

				/* clt Ri, Rj, Rk */
				case clt:
					storereg(ir.fmta.ri,
						fetchreg(ir.fmta.rj) < fetchreg(ir.fmta.rk));
					newreg = ir.fmta.ri;
					break;

				/* cle Ri, Rj, Rk */
				case cle:
					storereg(ir.fmta.ri,
						fetchreg(ir.fmta.rj) <= fetchreg(ir.fmta.rk));
					newreg = ir.fmta.ri;
					break;

				/* cgt Ri, Rj, Rk */
				case cgt:
					storereg(ir.fmta.ri,
						fetchreg(ir.fmta.rj) > fetchreg(ir.fmta.rk));
					newreg = ir.fmta.ri;
					break;

				/* cge Ri, Rj, Rk */
				case cge:
					storereg(ir.fmta.ri,
						fetchreg(ir.fmta.rj) >= fetchreg(ir.fmta.rk));
					newreg = ir.fmta.ri;
					break;

				/* not Ri, Rj  (32-bit complement) */
				case not:
					if (fetchreg(ir.fmta.rj) == 0)
						storereg(ir.fmta.ri, 1);
					else
						storereg(ir.fmta.ri, 0);
					newreg = ir.fmta.ri;
					break;

				/* jlr Ri, Rj  (Jump to register and link) */
				case jlr:
					storereg(ir.fmta.ri, ic);
					ic = fetchreg(ir.fmta.rj);
					newreg = ir.fmta.ri;
					break;

				/* nop */
				case nop:
					break;

				/* hlt */
				case hlt:
					running = FALSE;
					break;
				}
			break;

		/* Format B instructions have a 16-bit immediate operand. */
		case 'b':
			switch (ir.fmtb.op) {

				/* lw Ri, K(Rj)  (Load word) */
				case lw:
					storereg(ir.fmtb.ri,
							getmemword(fetchreg(ir.fmtb.rj) + (long) ir.fmtb.k));
					newreg = ir.fmtb.ri;
					break;

				/* lb Ri, K(Rj)  (Load byte) */
				case lb:
					w1 = getmembyte(fetchreg(ir.fmtb.rj) + (long) ir.fmtb.k);
					w2 = fetchreg(ir.fmtb.ri);
					storereg(ir.fmtb.ri, (w1) | (w2 & ~255));
					newreg = ir.fmtb.ri;
					break;

				/* sw K(Rj), Ri  (Store word) */
				case sw:
					newmem = fetchreg(ir.fmtb.rj) + (long) ir.fmtb.k;
					putmemword(newmem, fetchreg(ir.fmtb.ri));
					break;

				/* sb K(Rj), Ri  (Store byte) */
				case sb:
					newmem = fetchreg(ir.fmtb.rj) + (long) ir.fmtb.k;
					putmembyte(newmem, (BYTE) (fetchreg(ir.fmtb.ri) & 255));
					break;

				/* addi Ri, Rj, K  (Add immediate) */
				case addi:
					storereg(ir.fmtb.ri,
						fetchreg(ir.fmtb.rj) + (long) ir.fmtb.k);
					newreg = ir.fmtb.ri;
					break;

				/* subi Ri, Rj, K */
				case subi:
					storereg(ir.fmtb.ri,
						fetchreg(ir.fmtb.rj) - (long) ir.fmtb.k);
					newreg = ir.fmtb.ri;
					break;

				/* muli Ri, Rj, K */
				case muli:
					storereg(ir.fmtb.ri,
						fetchreg(ir.fmtb.rj) * (long) ir.fmtb.k);
					newreg = ir.fmtb.ri;
					break;

				/* divi Ri, Rj, K */
				case divi:
					k = (long) ir.fmtb.k;
					if (k == 0)
						runtimeerror("division by zero");
					else {
						storereg(ir.fmtb.ri,
							fetchreg(ir.fmtb.rj) / k);
						newreg = ir.fmtb.ri;
						}
					break;

				/* modi Ri, Rj, K */
				case modi:
					k = (long) ir.fmtb.k;
					if (k == 0)
						runtimeerror("division by zero");
					else {
						storereg(ir.fmtb.ri,
							fetchreg(ir.fmtb.rj) % k);
						newreg = ir.fmtb.ri;
						}
					break;

				/* andi Ri, Rj, K */
				case andi:
					storereg(ir.fmtb.ri,
						fetchreg(ir.fmtb.rj) & (long) ir.fmtb.k);
					newreg = ir.fmtb.ri;
					break;

				/* ori Ri, Rj, K */
				case ori:
					storereg(ir.fmtb.ri,
						fetchreg(ir.fmtb.rj) | (long) ir.fmtb.k);
					newreg = ir.fmtb.ri;
					break;

				/* ceqi Ri, Rj, K */
				case ceqi:
					storereg(ir.fmtb.ri,
						fetchreg(ir.fmtb.rj) == (long) ir.fmtb.k);
					newreg = ir.fmtb.ri;
					break;

				/* cnei Ri, Rj, K */
				case cnei:
					storereg(ir.fmtb.ri,
						fetchreg(ir.fmtb.rj) != (long) ir.fmtb.k);
					newreg = ir.fmtb.ri;
					break;

				/* clti Ri, Rj, K */
				case clti:
					storereg(ir.fmtb.ri,
						fetchreg(ir.fmtb.rj) < (long) ir.fmtb.k);
					newreg = ir.fmtb.ri;
					break;

				/* clei Ri, Rj, K */
				case clei:
					storereg(ir.fmtb.ri,
						fetchreg(ir.fmtb.rj) <= (long) ir.fmtb.k);
					newreg = ir.fmtb.ri;
					break;

				/* cgti Ri, Rj, K */
				case cgti:
					storereg(ir.fmtb.ri,
						fetchreg(ir.fmtb.rj) > (long) ir.fmtb.k);
					newreg = ir.fmtb.ri;
					break;

				/* cgei Ri, Rj, K */
				case cgei:
					storereg(ir.fmtb.ri,
						fetchreg(ir.fmtb.rj) >= (long) ir.fmtb.k);
					newreg = ir.fmtb.ri;
					break;

				/* sl Ri, K  (Shift left logical) */
				case sl:
					storereg(ir.fmtb.ri,
						fetchreg(ir.fmtb.ri) << (long) ir.fmtb.k);
					newreg = ir.fmtb.ri;
					break;

				/* sr Ri, K  (Shift right logical) */
				case sr:
					storereg(ir.fmtb.ri,
						fetchreg(ir.fmtb.ri) >> (long) ir.fmtb.k);
					newreg = ir.fmtb.ri;
					break;

				/* bz Ri, K  (Branch to K if Ri == 0) */
				case bz:
					if (fetchreg(ir.fmtb.ri) == 0)
						ic = (long) ir.fmtb.k;
					break;

				/* bnz Ri, K  (Branch to K if Ri != 0) */
				case bnz:
					if (fetchreg(ir.fmtb.ri) != 0)
						ic = (long) ir.fmtb.k;
					break;

				/* jl Ri, K  (Branch to K with link in Ri) */
				case jl:
					storereg(ir.fmtb.ri, ic);
					ic = (long) ir.fmtb.k;
					newreg = ir.fmtb.ri;
					break;

				/* getc Ri  (Read one character to Ri) */
				case gtc:
					if (tracing) {
						char buf [80];
						printf("\nEnter data for getc: ");
						gets(buf);
						if (buf[0] == '\0')
							ch = '\n';
						else
							ch = buf[0];
						}
					else
						ch = (BYTE) getchar();
					storereg(ir.fmtb.ri, ch);
					newreg = ir.fmtb.ri;
					break;

				/* putc Ri  (Write the character in Ri) */
				case ptc:
					if (tracing)
						printf("  Output from putc: %c",
							fetchreg(ir.fmtb.ri));

					else
						printf("%c", fetchreg(ir.fmtb.ri));
					break;

				/* jr Ri  (Jump to Ri) */
				case jr:
					ic = fetchreg(ir.fmtb.ri);
					break;

				/* j K  (Jump to K) */
				case j:
					ic = (long) ir.fmtb.k;
					break;
				}
			break;
		}
	}

/* Dump words of memory from addr-10 to addr+10. */
void dump (long addr) {
	long first = (addr - RANGE) & ~3;
	long last = (addr + RANGE) & ~3;
	if (first < 0)
		first = 0;
	if (last > MEMSIZE)
		last = MEMSIZE;
	for (addr = first; addr < last; addr += 4) {
		showword(addr);
		printf("\n");
		}
	}

/* Display a register. */
void showreg (short regnum) {
	wordtype word;
	char charbuf [5];
	word.data = regs[regnum];
	printf("   r%d =  %08lX  %s  %4ld", regnum, word.data,
		wordtochars(charbuf, word), word.data);
	}

/* Execute one instruction in trace mode. */
void traceinstr () {
	wordtype word;
	char charbuf [5];
	showword(ic);
	execinstr(TRUE);
	if (running) {
		if (newreg >= 0)
			showreg(newreg);
		else if (newmem >= 0) {
			long addr = newmem >> 2;
			printf("   M[%ld] =  %08lX  %s  %ld", newmem,
				mem[addr].word.data,
				wordtochars(charbuf, mem[addr].word),
				mem[addr].word.data);
			}
		printf("\n");
		if (mem[ic >> 2].breakpoint) {
			printf("%5ld Breakpoint\n", ic);
			running = FALSE;
			}
		}
	}

/* Execute <steps> instructions in trace mode. */
void runfor (long steps) {
	long cnt;
	running = TRUE;
	for (cnt = 0; cnt < steps; cnt++) {
		traceinstr();
		if (!running)
			break;
		}
	}

/* 	Fetch the operand of a trace instruction. The operand should be
 * either a number or a symbol.
 */
long getoperand (char *cp) {
	char *first;
	long val;
	while (*cp == ' ' || *cp == '\t')
		cp++;
	first = cp;
	if (isdigit(*cp)) {
		if (sscanf(first, "%ld", &val) == 1)
			return val;
		else {
			printf("?\n");
			return -1;
			}
		}
	else if (isalpha(*cp)) {
		val = getsymbolval(cp);
		if (val < 0)
			printf("?\n");
		return val;
		}
	else {
		printf("?\n");
		return - 1;
		}
	}

/* Show tracing instructions. */
void showtraceusage() {
	printf("The tracer prompts with `IC:-\'.  IC is the instruction counter.\n");
	printf("Upper or lower case letters are accepted.  n must be positive.\n");
	printf("\n");
	printf("<cr>     Trace K instructions.\n");
	printf("n        Trace n instructions.\n");
	printf("B        Display breakpoints.\n");
	printf("Bn       Set a breakpoint at n.\n");
	printf("C        Clear all breakpoints.\n");
	printf("Cn       Clear the breakpoint at n.\n");
	printf("D        Dump memory near IC.\n");
	printf("Dn       Dump memory near n.\n");
	printf("I        Set IC to entry point.\n");
	printf("In       Set IC to n.\n");
	printf("K        Set K (# steps executed by <cr>) to 10.\n");
	printf("Kn       Set K to n.\n");
	printf("Q        Quit.\n");
	printf("R        Show registers.\n");
	printf("S        Show symbols.\n");
	printf("X        Run to next break point.\n");
	printf("Xn       Run until IC = n.\n");
	printf("\n");
	}

/* Execute the program with tracing. */
void exectrace () {
	char cmd [BUFLEN];
	long addr;
	short regnum;
	ic = entrypoint;
	numsteps = 10;
	running = TRUE;
	while (1) {
		printf("%5ld:- ", ic);
		gets(cmd);
		if (!strcmp(cmd, "q") || !strcmp(cmd, "Q"))
			break;
		else if (!strcmp(cmd, ""))
			runfor(numsteps);
		else if (isdigit(*cmd)) {
			long steps = getoperand(cmd);
			if (steps > 0)
				runfor(steps);
			}
		else {
			char *cp = cmd;
			switch (*cp++) {

				/* B = show all breakpoints; Bn = set breakpoint. */
				case 'b': case 'B':
					if (*cp == '\0') {
						printf("Breakpoints are at: ");
						for (addr = 0; addr < MEMSIZE; addr++) {
							if (mem[addr].breakpoint)
								printf("%ld ", addr << 2);
							}
						printf("\n");
						}
					else {
						addr = getoperand(cp);
						if (addr >= 0)
							mem[addr >> 2].breakpoint = TRUE;
						}
					break;

				/* C = clear all breakpoints; Cn = clear a breakpoint. */
				case 'c': case 'C':
					if (*cp == '\0') {
						for (addr = 0; addr < MEMSIZE; addr++)
							mem[addr].breakpoint = FALSE;
						}
					else {
						addr = getoperand(cp);
						if (addr >= 0)
							mem[addr >> 2].breakpoint = FALSE;
						}
					break;

				/* D = dump memory near <ic>; Dn = dump memory near n. */
				case 'd': case 'D':
					if (*cp == '\0')
						dump(ic);
					else {
						addr = getoperand(cp);
						if (addr >= 0)
							dump(addr);
						}
					break;

				/* Explain how to use it. */
				case 'h': case 'H': case '?':
					showtraceusage();
					break;

				/* I = set <ic> to entry point; In = set <ic> to n. */
				case 'i': case 'I':
					if (*cp == '\0')
						ic = entrypoint;
					else {
						addr = getoperand(cp);
						if (addr >= 0)
							ic = addr;
						}
					break;

				/* K = set steps to 10; Kn = set steps to n. */
				case 'k': case 'K':
					if (*cp == '\0')
						numsteps = 10;
					else {
						numsteps = getoperand(cp);
						if (numsteps < 0)
							numsteps = 10;
						}
					break;

				/* R = show registers. */
				case 'r': case 'R':
					for (regnum = 0; regnum < MAXREG; regnum++) {
						showreg(regnum);
						printf("\n");
						}
					break;

				/* S = show symbols. */
				case 's': case 'S':
					showsymbols();
					break;

				/* X = run to next breakpoint; Xn = run to n. */
				case 'x': case 'X':
					if (*cp == '\0') {
						running = TRUE;
						while (running) {
							execinstr(FALSE);
							if (mem[ic >> 2].breakpoint) {
								printf("%5ld Breakpoint\n", ic);
								break;
								}
							}
						}
					else {
						addr = getoperand(cp);
						if (addr < 0)
							printf("?\n");
						else {
							running = TRUE;
							while (running) {
								execinstr(FALSE);
								if (ic == addr)
									break;
								}
							}
						}
					break;

				default:
					printf("?\n");
					break;
				}
			}
		}
	printf("\n%ld cycles.\n", cycles);
	}

/* Execute the program without tracing. */
void exec () {
	ic = entrypoint;
	running = TRUE;
	while (running) {
		execinstr(FALSE);
		}
	printf("\n%ld cycles.\n", cycles);
	}

/******************************* PARSING ***********************************/

enum tokentype {
	T_BAD, T_REG, T_OP, T_SYM, T_NUM, T_STR, T_COMMA, T_LP, T_RP, T_NULL
	};

struct {
	char symval [TOKLEN];	/* Characters of the token */
	enum tokentype kind;	/* Chosen from the enumeration */
	char *pos;				/* Pointer to start of token */
	short reg;				/* Register number for T_REG */
	short op;				/* Op code for T_OP */
	long intval;			/* Value for T_NUM */
	} token;

char oldval [TOKLEN];

char errmes [BUFLEN];	/* Error message */
int errorcount = 0;		/* Number of errors detected. */
char *bp;				/* Buffer pointer */

/* Record an error for reporting later; only the first error is recorded. */
void syntaxerror (char *message) {
	errorcount++;
	if (!strcmp(errmes, "")) {
		strcpy(errmes, "Error at `");
		strcat(errmes, oldval);
		strcat(errmes, " ");
		strcat(errmes, token.symval);
		strcat(errmes, "': ");
		strcat(errmes, message);
		}
	}

/* True if character can occur in a symbol */
short issymchar (char c) {
	return (isalnum(c) || c == '_');
	}

/* True if the string is a valid register name. */
short isreg (char *p) {
	long regnum = 0;
	if (!(*p == 'R' || *p == 'r'))
		return FALSE;
	p++;
	while (*p) {
		if (isdigit(*p))
			regnum = 10 * regnum + *p - '0';
		else
			return FALSE;
		p++;
		}
	if (regnum < MAXREG) {
		token.reg = (short) regnum;
		return TRUE;
		}
	else {
		syntaxerror("Illegal symbol");
		return FALSE;
		}
	}

/* Read a token and store appropriate values in the structure <token>.
 * For error reporting, the token repesenting the string must be left
 * in <token.symval>.
 */
void next () {
	char *t = token.symval;
	short op;
	strcpy(oldval, token.symval);
	while (*bp == ' ' || *bp == '\t')
		bp++;
	token.pos = bp;
	if (isalpha(*bp)) {
		/* Read a register, op code, directive, or symbol */
		while (issymchar(*bp))
			*t++ = *bp++;
		*t = '\0';
/*
		printf("<%s>\n", token.symval);
*/
		if (isreg(token.symval)) {
			token.kind = T_REG;
			return;
			}
		for (op = lw; op < last; op++) {
			if (!strcmp(token.symval, opnames[op])) {
				token.op = op;
				token.kind = T_OP;
				return;
				}
			}
		token.kind = T_SYM;
		return;
		}
	else if (*bp == '-' || *bp == '+' || isdigit(*bp)) {
		/* Read a signed decimal integer */
		*t++ = *bp++;
		while (isdigit(*bp))
			*t++ = *bp++;
		*t = '\0';
		sscanf(token.symval, "%ld", &token.intval);
		token.kind = T_NUM;
		return;
		}
	else if (*bp == '"') {
		/* Read a character string enclosed in quotes */
		bp++;
		while (1) {
			if (*bp == '"') {
				*t = '\0';
				token.kind = T_STR;
				bp++;
				break;
				}
			if (*bp == '\0' || *bp == '\n') {
				syntaxerror("unterminated string");
				token.kind = T_BAD;
				*t = '\0';
				break;
				}
			*t++ = *bp++;
			}
		return;
		}
	else if (*bp == ',') {
		bp++;
		token.kind = T_COMMA;
		strcpy(token.symval, ",");
		return;
		}
	else if (*bp == '(') {
		bp++;
		token.kind = T_LP;
		strcpy(token.symval, "(");
		return;
		}
	else if (*bp == ')') {
		bp++;
		token.kind = T_RP;
		strcpy(token.symval, ")");
		return;
		}
	else if (*bp == '%' || *bp == '\n' || *bp == '\0') {
		token.kind = T_NULL;
		strcpy(token.symval, " ");
		return;
		}
	else {
		token.kind = T_BAD;
		strcpy(token.symval, " ");
		}
	}

/* Match a token */
void match (enum tokentype kind) {
	if (token.kind == kind) {
		next();
		return;
		}
	switch (kind) {
		case T_COMMA:
			syntaxerror("',' expected");
			break;
		case T_LP:
			syntaxerror("'(' expected");
			break;
		case T_RP:
			syntaxerror("')' expected");
			break;
		default:
			syntaxerror("Syntax error");
			break;
		}
	}

/* Parse an opcode.  The error should never occur, since this function
 * is called only when the token type is know.
 */
short getop () {
	if (token.kind == T_OP) {
		short res = token.op;
		next();
		return res;
		}
	syntaxerror("Opcode expected");
	return 0;
	}

/* Parse a register and return the register number */
short getreg () {
	if (token.kind == T_REG) {
		short res = token.reg;
		next();
		return res;
		}
	syntaxerror("Register expected");
	return 0;
	}

/* Parse a constant (number or symbol) and return value */
long getlong (long addr) {
	if (token.kind == T_NUM) {
		long res = token.intval;
		next();
		return res;
		}
	else if (token.kind == T_SYM) {
		usesymbol(token.symval, addr);
		next();
		return 0;
		}
	syntaxerror("Constant expected");
	return 0;
	}

/* 	Similar to getlong(), but checks that its argument can be stored
 * in 16 bits.
 */
int getint (long addr) {
	long val = getlong(addr);
	if (labs(val) <= 32767)
		return (int) val;
	syntaxerror("Value cannot be represented with 16 bits");
	return 0;
	}

char buffer [BUFLEN];		/* Input buffer */
long addr = 0;
int linenum = 0;

/* Parse one line of source code from the buffer. */
void readline () {
	int c;
	wordtype word;
	word.data = 0;
	bp = buffer;
	strcpy(errmes, "");
	strcpy(oldval, "");
	next();
	while (token.kind == T_SYM) {
		defsymbol(token.symval, addr);
		next();
		}
	if (token.kind == T_OP) {
		switch (token.op) {

			/* Format A -- registers only */

			/* Operands Ri, Rj, Rk */
			case add:
			case sub:
			case mul:
			case div:
			case mod:
			case and:
			case or:
			case ceq:
			case cne:
			case clt:
			case cle:
			case cgt:
			case cge:
				word.fmta.op = getop();
				word.fmta.ri = getreg();
				match(T_COMMA);
				word.fmta.rj = getreg();
				match(T_COMMA);
				word.fmta.rk = getreg();
				putmeminstr(addr, word, 'a');
				addr += 4;
				break;

			/* Operands Ri, Rj */
			case not:
			case jlr:
				word.fmta.op = getop();
				word.fmta.ri = getreg();
				match(T_COMMA);
				word.fmta.rj = getreg();
				putmeminstr(addr, word, 'a');
				addr += 4;
				break;

			/* No operands */
			case nop:
			case hlt:
				word.fmta.op = getop();
				putmeminstr(addr, word, 'a');
				addr += 4;
				break;

			/* Format B - operands and constant fields */

			/* Operands Ri, K(Rj) */
			case lw:
			case lb:
				word.fmtb.op = getop();
				word.fmtb.ri = getreg();
				match(T_COMMA);
				word.fmtb.k = getint(addr);
				match(T_LP);
				word.fmtb.rj = getreg();
				match(T_RP);
				putmeminstr(addr, word, 'b');
				addr += 4;
				break;

			/* Operands K(Rj), Ri */
			case sw:
			case sb:
				word.fmtb.op = getop();
				word.fmtb.k = getint(addr);
				match(T_LP);
				word.fmtb.rj = getreg();
				match(T_RP);
				match(T_COMMA);
				word.fmtb.ri = getreg();
				putmeminstr(addr, word, 'b');
				addr += 4;
				break;

			/* Operands Ri, Rj, K */
			case addi:
			case subi:
			case muli:
			case divi:
			case modi:
			case andi:
			case ori:
			case ceqi:
			case cnei:
			case clti:
			case clei:
			case cgti:
			case cgei:
				word.fmtb.op = getop();
				word.fmtb.ri = getreg();
				match(T_COMMA);
				word.fmtb.rj = getreg();
				match(T_COMMA);
				word.fmtb.k = getint(addr);
				putmeminstr(addr, word, 'b');
				addr += 4;
				break;

			/* Operands Ri, K */
			case sl:
			case sr:
			case bz:
			case bnz:
			case jl:
				word.fmtb.op = getop();
				word.fmtb.ri = getreg();
				match(T_COMMA);
				word.fmtb.k = getint(addr);
				putmeminstr(addr, word, 'b');
				addr += 4;
				break;

			/* Operands Ri */
			case gtc:
			case ptc:
			case jr:
				word.fmtb.op = getop();
				word.fmtb.ri = getreg();
				putmeminstr(addr, word, 'b');
				addr += 4;
				break;

			/* Operands K */
			case j:
				word.fmtb.op = getop();
				word.fmtb.k = getint(addr);
				putmeminstr(addr, word, 'b');
				addr += 4;
				break;

			/* Set the entry point of the program. */
			case entry:
				next();
				if (entrypoint < 0) {
					entrypoint = addr;
					break;
					}
				syntaxerror("More than one entry point");
				break;

			/* Adjust the address to the next word boundary. */
			case align:
				next();
				if (addr & 3)
					addr = (addr & ~3) + 4;
				break;

			/* Set the address to the given value. */
			case org:
				next();
				addr = getlong(addr);
				break;

			/* Store words. */
			case dw:
				next();
				while (token.kind == T_NUM || token.kind == T_SYM) {
					word.data = getlong(addr);
					putmeminstr(addr, word, 'd');
					addr += 4;
					if (token.kind == T_COMMA)
						next();
					else
						break;
					}
				break;

			/* Store bytes */
			case db:
				next();
				while (1) {
					if (token.kind == T_NUM) {
						if (0 <= token.intval && token.intval <= 255) {
							putmemchar(addr, token.intval, 'd');
							addr++;
							}
						else
							syntaxerror("Value cannot be represented with 8 bits");
						next();
						}
					else if (token.kind == T_STR) {
						char *t = token.symval;
						while (*t) {
							putmemchar(addr, *t++, 'd');
							addr++;
							}
						next();
						}
					if (token.kind == T_COMMA)
						next();
					else if (token.kind == T_NULL)
						break;
					else {
						syntaxerror("Syntax error in byte list");
						break;
						}
					}
				break;

			/* Reserve the given number of words. */
			case res:
				next();
				addr += getlong(addr);
				break;

			/* Should never get here. */
			default:
				syntaxerror("Unrecognized statement");
				break;
			}
		}
	if (token.kind != T_NULL && errorcount < 5) {
		printf("Warning: junk following `%s' on next line.\n", token.symval);
		printf("%4d  %s\n", linenum, buffer);
		errorcount++;
		}
	}

/* Load a source file. */
void load (FILE *inp, FILE *out, short listing) {
	linenum = 0;
	while (fgets(buffer, BUFLEN - 1, inp)) {
		long oldaddr = addr;
		linenum++;
		readline();
		if (listing)
			fprintf(out, "%5d %5ld %s", linenum, oldaddr, buffer);
		if (strcmp(errmes, "")) {
			printf("%5d %5ld %s", linenum, oldaddr, buffer);
			printf("      >>>>> %s\n", errmes);
			if (listing)
				fprintf(out, "      >>>>> %s\n", errmes);
			}
		}
	}

/*************************** USER INSTRUCTION *******************************/

void showusage () {
	printf("Usage:\n");
	printf("         moon { option | filename }\n");
	printf("The command line may contain source file names and options in any order.\n");
	printf("There should be at least one source file.  Source files will be loaded\n");
	printf("in the order in which they are given.\n");
	printf("Options:\n");
	printf("       +p           print listing\n");
	printf("       -p (default) do not print listing\n");
	printf("       +s           display symbol values\n");
	printf("       -s (default) do not display symbol values\n");
	printf("       +t           start in trace mode\n");
	printf("       -t (default) execute without tracing\n");
	printf("       +x (default) execute the program\n");
	printf("       -x           do not execute the program\n");
	printf("Input files:\n");
	printf("       If an input file name does not contain `.', the suffix\n");
	printf("       `.n' will be appended to it.\n");
	printf("Listing files:\n");
	printf("       Source files may be listed selectively.  The command\n");
	printf("            moon -p lib +p appl\n");
	printf("       would create a listing for `appl.m' but not for `lib.m'.\n");
	printf("       The list file is named `moon.prn' by default.\n");
	printf("       Use +o or -o followed by a name to changes the list file name.\n");
	}

/*************************** MAIN PROGRAM ***********************************/

void main (int argc, char *argv[]) {
	struct {
		char name [MAXNAMELEN];
		short list;
		} filedescs [MAXINFILES];
	int numfiles = 0;
	char outname [MAXNAMELEN] = "";
	short dump = FALSE;			/* D Dump memory */
	short listing = FALSE;		/* P Generate a listing of the source code */
	short symbols = FALSE;		/* S Display symbol values */
	short tracing = FALSE;		/* T Execute program in trace mode */
	short execute = TRUE;		/* X Execute the program after loading */
	short listreq = FALSE;		/* A listing is needed */
	long addr;
	int arg, fil;
	FILE *inp, *out;
	regs[0] = 0;   /* Register 0 is always 0. */

	/* If there are no arguments, help the poor user. */

	if (argc <= 1) {
		showusage();
		exit(0);
		}

	/* Process command line arguments.
	 * There should be at least one argument.  Arguments that start with
	 * + (-) turn flags on (off).  Other arguments are input file names.
	 * The listing switch (+-l) is processed in sequence so that files
	 * may be listed selectively.
	 */

	for (arg = 1; arg < argc; arg++) {
		char *p = argv[arg];
		if (*p == '+') {
			p++;
			switch (*p++) {
				case 'd': case 'D':
					dump = TRUE;
					break;
				case 'o': case 'O':
					strcpy(outname, p);
					break;
				case 'p': case 'P':
					listing = TRUE;
					break;
				case 's': case 'S':
					symbols = TRUE;
					break;
				case 't': case 'T':
					tracing = TRUE;
					break;
				case 'x': case 'X':
					execute = TRUE;
					break;
				default:
					printf("Illegal option: +%s\n", --p);
					exit(1);
				}
			}
		else if (*p == '-') {
			p++;
			switch (*p++) {
				case 'd': case 'D':
					dump = FALSE;
					break;
				case 'o': case 'O':
					strcpy(outname, p);
					break;
				case 'p': case 'P':
					listing = FALSE;
					break;
				case 's': case 'S':
					symbols = FALSE;
					break;
				case 't': case 'T':
					tracing = FALSE;
					break;
				case 'x': case 'X':
					execute = FALSE;
					break;
				default:
					printf("Illegal option: -%s\n", --p);
					exit(1);
				}
			}
		else {
			if (numfiles >= MAXINFILES) {
				printf("Too many input files!\n");
				exit(1);
				}
			strcpy(filedescs[numfiles].name, p);
			filedescs[numfiles].list = listing;
			if (listing)
				listreq = TRUE;
			numfiles++;
			}
		}

	/* Nothing to do if there were no files on the command line. */

	if (numfiles == 0) {
		printf("No input files!\n");
		exit(1);
		}

	/* 	Attempt to open an output file if a listing is required.
	 * If no output file was named, use a default name.
	 */

	if (listreq) {
		if (strlen(outname) == 0)
			strcpy(outname,"moon.prn");
		if ((out = fopen(outname, "w")) == NULL) {
			printf("Unable to open listing file %s.\n", outname);
			exit(1);
			}
		printf("Writing listing to %s.\n", outname);
		}

	/* Process each input file. If no extension is given, assume .m.
	 * If the file can be opened, load assembler code from it.
	 */

	initmem();
	defsymbol("topaddr", 4 * MEMSIZE);
	for (fil = 0; fil < numfiles; fil++) {
		if (!strchr(filedescs[fil].name, '.'))
			strcat(filedescs[fil].name, ".m");
		if ((inp = fopen(filedescs[fil].name, "r")) == NULL) {
			printf("Unable to open input file: %s.\n", filedescs[fil].name);
			exit(1);
			}
		else {
			short listing = filedescs[fil].list;
			printf("Loading %s.\n", filedescs[fil].name);
			if (listing) {
				fprintf(out, "MOON listing of %s.\n\n", filedescs[fil].name);
				load(inp, out, TRUE);
				fprintf(out, "\n");
				}
			else
				load(inp, out, FALSE);
			fclose(inp);
			}
		}
	if (listreq)
		fclose(out);

	/* Check symbols and entry point. If there are errors, stop now. */
	errorcount += checksymbols();
	if (entrypoint < 0) {
		printf("There is no `entry' directive.\n");
		errorcount++;
		}
	if (errorcount > 0) {
		printf("Loader errors -- no execution.\n");
		exit(1);
		}

	/* Store values of symbols where they are used in the program.
	 * Display them if requested.
	 */
	storesymbols();
	if (symbols)
		showsymbols();

	/* If a dump was requested, dump the memory.  This option is not
	 * advertised and therefore need not be supported.
	 */

	if (dump) {
		printf("Memory dump:\n");
		for (addr = 0; addr < MEMSIZE; addr += 4) {
			if (mem[addr >> 2].cont != 'u') {
				showword(addr);
				printf("\n");
				}
			}
		printf("\n");
		}

	/* Execute the program in normal or trace mode. */
	if (execute) {
		if (tracing)
			exectrace();
		else exec();
		}
	}
