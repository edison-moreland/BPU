define(CONSOLE_ADDRESS, 0b1010_1010)
define(ASCII_PRINTABLE_START, 0x20)
define(ASCII_PRINTABLE_END, 0x7E)
define(SCREEN_WIDTH, 16)
define(NEWLINE_CHAR, 0x0A)

LDI R0, CONSOLE_ADDRESS
PSEL R0

- R0 = Current char
LDI R0,ASCII_PRINTABLE_START
- R1 = Place in line
XOR R1
- R2 = 1
LDI R2, 1
- R3 = Scratch

:loop
- Write to screen
POUT R0

- Increment char
ADD R2,R0
LDI R3,ASCII_PRINTABLE_END
CMP R0,R3
JE :zero_char
:after_inc_char

- Increment place in line
ADD R2,R1
LDI R3,SCREEN_WIDTH
CMP R1,R3
JE :zero_line

JMP :loop

:zero_char
LDI R0,ASCII_PRINTABLE_START
JMP :after_inc_char

:zero_line
XOR R1

- output newline
LDI R3,NEWLINE_CHAR
POUT R3

JMP :loop