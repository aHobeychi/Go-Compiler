% String tests for MOON

       entry
       addi    r14,r0,topaddr

       % Read two strings.
loop   addi    r1,r0,enta
       sw      -8(r14),r1
       jl      r15,putstr
       addi    r1,r0,stra
       sw      -8(r14),r1
       jl      r15,getstr
       addi    r1,r0,entb
       sw      -8(r14),r1
       jl      r15,putstr
       addi    r1,r0,strb
       sw      -8(r14),r1
       jl      r15,getstr

       % Concatenate the strings.

       addi    r1,r0,catm
       sw      -8(r14),r1
       jl      r15,putstr
       addi    r1,r0,stra
       sw      -8(r14),r1
       addi    r1,r0,strb
       sw      -12(r14),r1
       addi    r1,r0,strc
       sw      -16(r14),r1
       jl      r15,strcat
       addi    r1,r0,strc
       sw      -8(r14),r1
       jl      r15,putstr
       addi    r1,r0,cr
       sw      -8(r14),r1
       jl      r15,putstr

       % Apply each equality test to the strings.

       addi    r1,r0,stra
       sw      -8(r14),r1
       addi    r1,r0,strb
       sw      -12(r14),r1
       jl      r15,streq
       ceqi    r1,r13,0
       bnz     r1,endeq
       addi    r1,r0,repeq
       sw      -8(r14),r1
       jl      r15,putstr
endeq

       addi    r1,r0,stra
       sw      -8(r14),r1
       addi    r1,r0,strb
       sw      -12(r14),r1
       jl      r15,strne
       ceqi    r1,r13,0
       bnz     r1,endne
       addi    r1,r0,repne
       sw      -8(r14),r1
       jl      r15,putstr
endne

       addi    r1,r0,stra
       sw      -8(r14),r1
       addi    r1,r0,strb
       sw      -12(r14),r1
       jl      r15,strlt
       ceqi    r1,r13,0
       bnz     r1,endlt
       addi    r1,r0,replt
       sw      -8(r14),r1
       jl      r15,putstr
endlt

       addi    r1,r0,stra
       sw      -8(r14),r1
       addi    r1,r0,strb
       sw      -12(r14),r1
       jl      r15,strle
       ceqi    r1,r13,0
       bnz     r1,endle
       addi    r1,r0,reple
       sw      -8(r14),r1
       jl      r15,putstr
endle

       addi    r1,r0,stra
       sw      -8(r14),r1
       addi    r1,r0,strb
       sw      -12(r14),r1
       jl      r15,strgt
       ceqi    r1,r13,0
       bnz     r1,endgt
       addi    r1,r0,repgt
       sw      -8(r14),r1
       jl      r15,putstr
endgt

       addi    r1,r0,stra
       sw      -8(r14),r1
       addi    r1,r0,strb
       sw      -12(r14),r1
       jl      r15,strge
       ceqi    r1,r13,0
       bnz     r1,endge
       addi    r1,r0,repge
       sw      -8(r14),r1
       jl      r15,putstr
endge
       jl      r15,more
       j       loop

% Asks user whether to continue testing.
% If yes, returns; otherwise halts.

more  sw     -4(r14),r15
      addi   r1,r0,query
      sw     -8(r14),r1
      jl     r15,putstr
      addi   r1,r0,reply
      sw     -8(r14),r1
      jl     r15,getstr
      addi   r1,r0,0
      lb     r1,reply(r0)
      ceqi   r2,r1,121        % 'y'
      bnz    r2,more1
      hlt
more1 lw     r15,-4(r14)
      jr     r15
query db     "Again? ", 0
reply res    80

enta   db      "Enter string A: ", 0
entb   db      "Enter string B: ", 0
catm   db      "Concatenated: ", 0
cr     db      13, 10, 0
repeq  db      "streq true", 13, 10, 0
repne  db      "strne true", 13, 10, 0
replt  db      "strlt true", 13, 10, 0
reple  db      "strle true", 13, 10, 0
repgt  db      "strgt true", 13, 10, 0
repge  db      "strge true", 13, 10, 0
       align
stra   res      80
strb   res      80
strc   res     160
