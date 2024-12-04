define(CONSOLE_ADDRESS, 0b1010_1010)
define(NEWLINE_CHAR, 0x0A)
define(CLEAR_CHAR, 0x03)

- Print to console

LDI R0, CONSOLE_ADDRESS
PSEL R0
LDI R0, CLEAR_CHAR
POUT R0

- R0 = Start of string address
LDI R0, :string
LDI R1, 1
XOR R3

:loop
- Load char
LD R0,R2

- Is char null?
CMP R2,R3
JE :done

- Output to console
POUT R2

- Inc place in line
ADD R1,R0
JMP :loop

:done
JMP :done

:string
%DATA 0x48
%DATA 0x65
%DATA 0x6C
%DATA 0x6C
%DATA 0x6F
%DATA 0x20
%DATA 0x77
%DATA 0x6F
%DATA 0x72
%DATA 0x6C
%DATA 0x64
%DATA 0x00
