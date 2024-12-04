# But How Do It Know

Collection of stuff relevant to the computer described in "But How Do It Know".

## Instructions

### ALU Opcodes
| Code | Operation |
|------|-----------|
| 000  | ADD       |
| 001  | SHR       |
| 010  | SHL       |
| 011  | NOT       |
| 100  | AND       |
| 101  | OR        |
| 110  | XOR       |
| 111  | CMP       |

### Encoding
XXX = ALU opcode
aa = Register a
bb = Register b
iiiiiiii = Immediate byte
zceg = Flags nibble (zero, carry, a == b, a > b)


| Encoding          | Operation                               | Desc                                     |
|-------------------|-----------------------------------------|------------------------------------------|
| bbaaxxx1          | b = alu(op, a, b)                       | ALU operation                            |
| bbaa0000          | b = mem[a]                              | Load byte                                |
| bbaa0010          | mem[a] = b                              | Store byte                               |
| bb__0100 iiiiiiii | b = mem[pc++]                           | Load immediate byte                      |
| bb__0110          | pc = b                                  | Jump to address in register              |
| ____1000 iiiiiiii | pc = mem[pc+1]                          | Jump to address in immediate byte        |
| zceg1010 iiiiiiii | if (flags == zceg) pc = mem[pc+1]       | Jump to immediate if flags set           |
| ____1100          | flags = 0000                            | Clear flags                              |
| bbdi1110          | if (i) b = per(i, d) else per(i, d) = b | Send or receive data from peripheral bus |


# Console screen

## characters
All printable chars match the ascii standard but control chars are redefined.

| Hex  | Bin      | Symbol | Description  |
|------|----------|--------|--------------|
| 0x0A | 00001010 | NL     | New line     |
| 0x03 | 00000011 | CL     | Clear screen |