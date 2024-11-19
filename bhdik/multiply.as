- Multiply: R2 = R0 * R1

- R3 = 1, R2 = 0
LDI R3, 1
XOR R2,R2

:start_loop
CLF
SHR R0
JC :do_add
JMP :do_shift

:do_add
CLF
ADD R1,R2

:do_shift
CLF
SHL R1
SHL R3
JC :done
JMP :start_loop

:done
JMP :done
