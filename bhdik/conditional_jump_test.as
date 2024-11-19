
- Put a one in the carry flag

LDI R0,1
SHR R0
JC :success

:failure
JMP :failure

:success
JMP :success