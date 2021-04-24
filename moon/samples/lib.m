% A Simple MOON Library
% Author: Peter Grogono
% Last modified: 27 Jan 1995

% Conventions
%   r14 is the Stack Pointer
%   -4(r14) is used to store the link, if necessary.
%   -8(r14) contains the first argument.
%   -12(r14) contains the second argument, and so on.
%   r15 contains the link.
%   r13 contains the result of a function.

%   Strings are null terminated.  Strings are passed and returned to
%   library functions as pointers.  It is the caller's responsibility to
%   provide storage for strings.

%   "->" is read "points to".

          align               % In case previous file misaligned

% Write a string to stdout.
% Entry: -8(r14) -> string argument.

putstr    lw    r1,-8(r14)    % i := r1
          addi  r2,r0,0
putstr1   lb    r2,0(r1)      % ch := B[i]
          ceqi  r3,r2,0
          bnz   r3,putstr2    % branch if ch = 0
          putc  r2
          addi  r1,r1,1       % i++
          j     putstr1
putstr2   jr    r15

% Read a string from stdin.  Read until CR (ASCII 13) but do not store
% the CR.
% Entry: -8(r14) -> buffer.

getstr    lw    r1,-8(r14)    % i := r1
getstr1   getc  r2            % get ch
          ceqi  r3,r2,10
          bnz   r3,getstr2    % branch if ch = CR
          sb    0(r1),r2      % B[i] := ch
          addi  r1,r1,1       % i++
          j     getstr1
getstr2   sb    0(r1),r0      % B[i] := '\0'
          jr    r15

% Convert string to integer.  Skip leading blanks.  Accept leading sign.
% Entry: -8(r14) -> string.
% Exit:  result in r13

strint    addi  r13,r0,0      % R := 0 (result)
          addi  r4,r0,0       % S := 0 (sign)
          lw    r1,-8(r14)    % i := r1
          addi  r2,r0,0
strint1   lb    r2,0(r1)      % ch := B[i]
          cnei  r3,r2,32
          bnz   r3,strint2    % branch if ch != blank
          addi  r1,r1,1
          j     strint1
strint2   cnei  r3,r2,43
          bnz   r3,strint3    % branch if ch != "+"
          j     strint4
strint3   cnei  r3,r2,45
          bnz   r3,strint5    % branch if ch != "-"
          addi  r4,r4,1       % S := 1
strint4   addi  r1,r1,1       % i++
          lb    r2,0(r1)      % ch := B[i]
strint5   clti  r3,r2,48
          bnz   r3,strint6    % branch if ch < "0"
          cgti  r3,r2,57
          bnz   r3,strint6    % branch if ch > "9"
          subi  r2,r2,48      % ch -= "0"
          muli  r13,r13,10    % R *= 10
          add   r13,r13,r2    % R += ch
          j     strint4
strint6   ceqi  r3,r4,0
          bnz   r3,strint7    % branch if S = 0
          sub   r13,r0,r13    % R := -R
strint7   jr    r15

% Convert signed integer to string.
% Entry: -8(r14) is the integer.
%        -12(r14) -> buffer containing at least 12 bytes.
% Exit:  r13 -> first character of result string.

intstr    lw    r13,-12(r14)
          addi  r13,r13,11    % r13 points to end of buffer
          sb    0(r13),r0     % store terminator
          lw    r1,-8(r14)    % r1 := N (to be converted)
          addi  r2,r0,0       % S := 0 (sign)
          cgei  r3,r1,0
          bnz   r3,intstr1    % branch if N >= 0
          addi  r2,r2,1       % S := 1
          sub   r1,r0,r1      % N := -N
intstr1   addi  r3,r1,0       % D := N (next digit)
          modi  r3,r3,10      % D mod= 10
          addi  r3,r3,48      % D += "0"
          subi  r13,r13,1     % i--
          sb    0(r13),r3     % B[i] := D
          divi  r1,r1,10      % N div= 10
          cnei  r3,r1,0
          bnz   r3,intstr1    % branch if N != 0
          ceqi  r3,r2,0
          bnz   r3,intstr2    % branch if S = 0
          subi  r13,r13,1     % i--
          addi  r3,r0,45
          sb    0(r13),r3     % B[i] := "-"
intstr2   jr    r15

% Return length of string.
% Entry: -8(r14) -> string.
% Exit:  r13 = length of string.

lenstr   lw    r1,-8(r14)     % i -> string
         addi  r13,r0,0       % L := 0
         addi  r2,r0,0
lenstr1  lb    r2,0(r1)       % ch := B[i]
         ceqi  r3,r2,0
         bnz   r3,lenstr2     % branch if ch = 0
         addi  r13,r13,1      % L++
         addi  r1,r1,1        % i++
         j     lenstr1
lenstr2  jr    r15

% Concatenate strings: Z := X + Y.
% Entry:   -8(r14)  ->  X
%         -12(r14)  ->  Y
%         -16(r14)  ->  Z
% The result string is assumed to be large enough to hold the result.

strcat   lw    r1,-16(r14)    % r1 -> Z
         lw    r2,-8(r14)     % r2 -> X
         addi  r3,r0,0        % r3 = current character
strcat1  lb    r3,0(r2)       % char from X
         ceqi  r4,r3,0
         bnz   r4,strcat2     % branch at end of X
         sb    0(r1),r3       % copy char to Z
         addi  r1,r1,1
         addi  r2,r2,1
         j     strcat1
strcat2  lw    r2,-12(r14)    % r2 -> Y
strcat3  lb    r3,0(r2)       % char from Y
         ceqi  r4,r3,0
         bnz   r4,strcat4     % branch at end of Y
         sb    0(r1),r3       % copy char to Z
         addi  r1,r1,1
         addi  r2,r2,1
         j     strcat3
strcat4  sb    0(r1),r0       % Store terminator
         jr    r15

% The string comparison functions all use strcmp, defined below.
% They are all short; a smart compiler could generate this code
% directly, avoiding the overhead of an extra level of function
% call.
% For each function:
% Entry: -8(r14)  -> string A
%        -12(r14) -> string B
% Exit:  r13 = 1 for true and 0 for false.
% Note that r11 is used for the link to avoid saving r15.

streq    jl    r11,strcmp
         lw    r13,eq(r13)
         jr    r15
eq       dw    1,0,0

strne    jl    r11,strcmp
         lw    r13,ne(r13)
         jr    r15
ne       dw    0,1,1

strlt    jl    r11,strcmp
         lw    r13,lt(r13)
         jr    r15
lt       dw    0,1,0

strle    jl    r11,strcmp
         lw    r13,le(r13)
         jr    r15
le       dw    1,1,0

strgt    jl    r11,strcmp
         lw    r13,gt(r13)
         jr    r15
gt       dw    0,0,1

strge    jl    r11,strcmp
         lw    r13,ge(r13)
         jr    r15
ge       dw    1,0,1

% Compare strings.
% Entry: -8(r14)  -> string A
%        -12(r14) -> string B
%             ( 0   if A = B
% Exit: r13 = < 4   if A < B
%             ( 8   if A > B
% Note that r11 is the link.

strcmp   lw    r1,-8(r14)
         lw    r2,-12(r14)
         addi  r3,r0,0
         addi  r4,r0,0
strcmp1  lb    r3,0(r1)       % get A[i]
         lb    r4,0(r2)       % get B[j]
         ceqi  r5,r3,0
         bnz   r5,strcmp2     % branch if end of A
         ceqi  r5,r4,0
         bnz   r5,strcmp4     % branch if end of B
         ceq   r5,r3,r4
         bz    r5,strcmp3     % branch if A[i] != B[i]
         addi  r1,r1,1        % i++
         addi  r2,r2,1        % j++
         j     strcmp1
strcmp2  ceqi  r5,r4,0
         bz    r5,strcmp5     % branch if not end of B
         addi  r13,r0,0       % A = B
         jr    r11
strcmp3  clt   r5,r3,r4
         bnz   r5,strcmp5     % branch if A[i] < B[i]
strcmp4  addi  r13,r0,8       % A > B
         jr    r11
strcmp5  addi  r13,r0,4       % A < B
         jr    r11

% String indexing: return the string S[M].  The value is returned as a
% string rather than as a character for compatibility with other string
% processing functions.  Also, we must distinguish the empty string
% from a character.
% Entry: -8(r14)  -> S
%        -12(r14) -> M
%        -16(r14) -> T, the output string.
% Exit:  the output string contains the selected character, or is null.

stridx   lw    r1,-8(r14)    % i
         lw    r2,-12(r14)
         lw    r3,-16(r14)
         subi  r2,r2,1
         addi  r4,r0,0
stridx1  lb    r4,0(r1)      % ch := S[i]
         ceqi  r5,r4,0
         bnz   r5,stridx3    % branch if ch = 0
         cge   r5,r1,r2
         bnz   r5,stridx2    % branch if i >= M
         addi  r1,r1,1       % i++
         j     stridx1
stridx2  sb    0(r3),r4      % T[0] := S[M]
         addi  r3,r3,1
stridx3  sb    0(r3),r0      % T[k] := 0
         jr    r15

% String indexing: return the string S[M..N].
% Entry: -8(r14)  -> S
%        -12(r14) -> M
%        -16(r14) -> N
%        -20(r14) -> T, the output string.
% Exit:  the output string contains the selected substring, or is null.

strsub   lw    r1,-8(r14)
         lw    r2,-12(r14)
         lw    r3,-16(r14)
         lw    r4,-20(r14)
         subi  r2,r2,1
         subi  r3,r3,1
         addi  r5,r0,0
strsub1  lb    r5,0(r1)      % ch := S[i]
         ceqi  r6,r5,0
         bnz   r6,strsub3    % branch if ch = 0
         clt   r6,r1,r2
         bnz   r6,strsub2    % branch if i < M
         cgt   r6,r1,r3
         bnz   r6,strsub3    % branch if i > N
         sb    0(r3),r5      % T[k] := S[i]
         addi  r3,r3,1       % k++
strsub2  addi  r1,r1,1       % i++
         j     strsub1
strsub3  sb    0(r3),r0      % T[k] := 0
         jr    r15

