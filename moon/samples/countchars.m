% First test of library

      entry
      addi    r14,r0,topaddr  % Set stack pointer

      addi    r1,r0,mess     % Display "Enter string: "
      sw      -8(r14),r1
      jl      r15,putstr

      addi    r1,r0,reply    % Read reply
      sw      -8(r14),r1
      jl      r15,getstr

      addi    r1,r0,resp     % Display "You said"
      sw      -8(r14),r1
      jl      r15,putstr

      addi    r1,r0,reply    % Display string
      sw      -8(r14),r1
      jl      r15,putstr

      addi    r1,r0,cr       % Display CR
      sw      -8(r14),r1
      jl      r15,putstr

      addi    r1,r0,thwe     % Display "There were "
      sw      -8(r14),r1
      jl      r15,putstr

      addi    r1,r0,reply    % Find length of reply
      sw      -8(r14),r1
      jl      r15,lenstr

      sw      -8(r14),r13    % Convert result to string
      addi    r1,r0,buff
      sw      -12(r14),r1
      jl      r15,intstr

      sw      -8(r14),r13    % Display the string
      jl      r15,putstr

      addi    r1,r0,char     % Display " characters in your string."
      sw      -8(r14),r1
      jl      r15,putstr

      hlt

mess  db   "Enter string: ", 0
resp  db   "You said: ", 0
thwe  db   "There were ", 0
char  db   " characters in your string.", 13, 10, 0
cr    db   13, 10, 0
buff  res  12
reply res  80
