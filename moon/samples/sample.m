         org   103
 message db    "Hello, world!", 13, 10, 0
         org   217
         align
         entry                  % Start here
         add   r2,r0,r0
 pri     lb    r3,message(r2)   % Get next char
         ceqi  r4,r3,0
         bnz    r4,pr2          % Finished if zero
         putc   r3
         addi   r2,r2,1
         j      pri             % Go for next char
 pr2     addi   r2,r0,name      % Go and get reply
         jl     r15,getname
         hlt                    % All done!

 % Subroutine to read a string
 name    res    59              % Name buffer
         align
 getname getc   r3              % Read from keyboard
         ceqi   r4,r3,10
         bnz    r4,endget       % Finished if CR
         sb     0(r2),r3        % Store char in buffer
         addi   r2,r2,1
         j      getname
 endget  sb     0(r2),r0        % Store terminator
         jr     r15             % Return

 data    dw     1000, -35





