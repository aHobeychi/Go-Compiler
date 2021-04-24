% Arithmetic tests for MOON simulator

      entry
      addi   r14,r0,topaddr  % Set stack pointer

start addi   r1,r0,entx    % Ask for X
      sw     -8(r14),r1
      jl     r15,putstr

      addi   r1,r0,buf     % Get X
      sw     -8(r14),r1
      jl     r15,getstr
      jl     r15,strint    % Convert to integer
      sw     x(r0),r13     % Store X

      addi   r1,r0,enty    % Ask for Y
      sw     -8(r14),r1
      jl     r15,putstr

      addi   r1,r0,buf     % Get Y
      sw     -8(r14),r1
      jl     r15,getstr
      jl     r15,strint    % Convert to integer
      sw     y(r0),r13     % Store Y

      addi   r1,r0,0
      lw     r3,x(r0)      % r3 := X
      lw     r4,y(r0)      % r4 := Y

% Perform each arithmetic operation with X and Y

      add    r2,r3,r4
      sw     tab(r1),r2
      addi   r1,r1,4

      sub    r2,r3,r4
      sw     tab(r1),r2
      addi   r1,r1,4

      mul    r2,r3,r4
      sw     tab(r1),r2
      addi   r1,r1,4

      div    r2,r3,r4
      sw     tab(r1),r2
      addi   r1,r1,4

      mod    r2,r3,r4
      sw     tab(r1),r2
      addi   r1,r1,4

      ceq    r2,r3,r4
      sw     tab(r1),r2
      addi   r1,r1,4

      cne    r2,r3,r4
      sw     tab(r1),r2
      addi   r1,r1,4

      clt    r2,r3,r4
      sw     tab(r1),r2
      addi   r1,r1,4

      cle    r2,r3,r4
      sw     tab(r1),r2
      addi   r1,r1,4

      cgt    r2,r3,r4
      sw     tab(r1),r2
      addi   r1,r1,4

      cge    r2,r3,r4
      sw     tab(r1),r2
      addi   r1,r1,4

% Report the results

      addi   r1,r0,buf
      sw     -12(r14),r1
      addi   r9,r0,0
      addi   r10,r0,0

disp  addi   r3,r0,mess
      add    r3,r3,r10
      addi   r4,r0,0
      lb     r4,0(r3)      % Check for null string
      ceqi   r5,r4,0
      bnz    r5,again
      sw     -8(r14),r3
      jl     r15,putstr

      lw     r3,tab(r9)
      sw     -8(r14),r3
      jl     r15,intstr
      sw     -8(r14),r13
      jl     r15,putstr

      addi   r3,r0,cr
      sw     -8(r14),r3
      jl     r15,putstr

      addi   r9,r9,4
      addi   r10,r10,4
      j      disp

again jl     r15,more
      j      start

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

x     dw     0
y     dw     0
tab   res    48            % Store results of operations
entx  db     "Enter X: ", 0
enty  db     "Enter Y: ", 0
cr    db     13, 10, 0
buf   res    20
      align
mess  db     "+  ", 0
      db     "-  ", 0
      db     "*  ", 0
      db     "/  ", 0
      db     "\  ", 0
      db     "=  ", 0
      db     "!= ", 0
      db     "<  ", 0
      db     "<= ", 0
      db     ">  ", 0
      db     ">= ", 0
      db     0
