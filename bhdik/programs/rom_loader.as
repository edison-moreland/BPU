define(ROM_PERIPHERAL_ADDRESS, 0b0000_1010)
define(ROM_SIZE, 64)
define(ENTRYPOINT, 0)

%OFFSET 0xC0

- ROM is loaded into 0x00-0x3F
- This program lives in 0xC0-0xFF
- Programs are free to use 0x40-0xBF

- Select Rom
LDI R0, ROM_PERIPHERAL_ADDRESS
PSEL R0

- R0 = ROM/RAM Address
XOR R0
- R1 = 1
LDI R1, 1
- R2 = ROM_SIZE
LDI R2, ROM_SIZE
- R3 = scratch

:loop
- read next byte from rom
POUT R0
PIN R3
- store in memory
ST R0,R3
- inc memory address
ADD R1,R0
- Are we done?
CMP R0,R2
JE ENTRYPOINT

JMP :loop
