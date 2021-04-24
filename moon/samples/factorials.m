% MOON simulator: test of recursion

      entry
      addi   r14,r0,topaddr  % Set stack pointer

loop  lw     r1,arg(r0)     % Fetch argument
      cgti   r2,r1,12       % and decide whether
      bnz    r2,stop        % to stop.

      sw     -8(r14),r1     % Display argument.
      addi   r1,r0,buf
      sw     -12(r14),r1
      jl     r15,intstr
      sw     -8(r14),r13
      jl     r15,putstr

      addi   r1,r0,m1       % "! = "
      sw     -8(r14),r1
      jl     r15,putstr

      lw     r1,arg(r0)     % Fetch argument again
      sw     -8(r14),r1
      jl     r15,fac        % and call factorial
      sw     -8(r14),r13
      addi   r1,r0,buf
      sw     -12(r14),r1
      jl     r15,intstr     % Convert and display result
      sw     -8(r14),r13
      jl     r15,putstr

      addi   r1,r0,m2       % CR
      sw     -8(r14),r1
      jl     r15,putstr

      lw     r1,arg(r0)     % Increment argument
      addi   r1,r1,1
      sw     arg(r0),r1
      j      loop

stop  hlt

arg   dw     0
m1    db     "! = ", 0
m2    db     13, 10, 0
buf   res    20

% Recursive factorial function.
% -8(r14) = argument.
% r13 = result.
fac   sw     -4(r14),r15    % Save link
      lw     r1,-8(r14)     % Get N
      ceqi   r2,r1,0
      bnz    r2,fac1        % branch if N = 0
      subi   r1,r1,1        % N := N - 1
      sw     -16(r14),r1    % Store as argument
      subi   r14,r14,8      % Adjust SP
      jl     r15,fac        % Recursive call
      addi   r14,r14,8      % Adjust SP
      lw     r1,-8(r14)     % Get N again
      mul    r13,r13,r1     % R := R * N
      j      fac2
fac1  addi   r13,r0,1       % R := 1
fac2  lw     r15,-4(r14)    % Restore link
      jr     r15
