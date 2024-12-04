define(CONSOLE_ADDRESS, 0b1010_1010)
define(ASCII_PRINTABLE_START, 0x20)
define(ASCII_PRINTABLE_END, 0x7E)
define(SCREEN_COLS, 16)
define(SCREEN_ROWS, 8)
define(NEWLINE_CHAR, 0x0A)
define(CLEAR_CHAR, 0x03)

- Select console peripheral and clear screen
LDI R0, CONSOLE_ADDRESS
PSEL R0
LDI R0, CLEAR_CHAR
POUT R0

- R0 = Current char
define(CHAR_REG, R0)
LDI CHAR_REG,ASCII_PRINTABLE_START

- R1 = Current col
define(COL_REG, R1)
XOR COL_REG

- R2 = Current row
define(ROW_REG, R2)
XOR ROW_REG

- R3 = Scratch
XOR R3

:loop
- Write to screen
POUT CHAR_REG

- Increment char and col
LDI R3,1
ADD R3,CHAR_REG
ADD R3,COL_REG

LDI R3,ASCII_PRINTABLE_END
CMP CHAR_REG,R3
JE :zero_char
:after_inc_char

LDI R3,SCREEN_COLS
CMP COL_REG,R3
JE :zero_col

JMP :loop
- end loop

:zero_char
LDI CHAR_REG,ASCII_PRINTABLE_START
JMP :after_inc_char
- end zero_char

:zero_col
XOR COL_REG

- output newline
LDI R3,NEWLINE_CHAR
POUT R3

- Inc row
LDI R3,1
ADD R3,ROW_REG

LDI R3,SCREEN_ROWS
CMP ROW_REG,R3
JE :zero_row
JMP :loop
- end zero_col

:zero_row
XOR ROW_REG

LDI R3,CLEAR_CHAR
POUT R3
JMP :loop
- end zero_row