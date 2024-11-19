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

| Encoding          | Operation                         | Desc                              |
|-------------------|-----------------------------------|-----------------------------------|
| bbaaxxx1          | b = alu(op, a, b)                 | ALU operation                     |
| bbaa0000          | b = mem[a]                        | Load byte                         |
| bbaa0010          | mem[a] = b                        | Store byte                        |
| bb__0100 iiiiiiii | b = mem[pc++]                     | Load immediate byte               |
| bb__0110          | pc = b                            | Jump to address in register       | 
| ____1000 iiiiiiii | pc = mem[pc+1]                    | Jump to address in immediate byte |
| zceg1010 iiiiiiii | if (flags == zceg) pc = mem[pc+1] | Jump to immediate if flags set    |
| ____1100          | flags = 0000                      | Clear flags                       |


### Assembly
| Assembly     | Operation                    |
|--------------|------------------------------|
| ADD RA, RB   | B = A + B                    |
| SHR RA, RB   | B = A >> 1                   |
| SHL RA, RB   | B = A << 1                   |
| NOT RA, RB   | B = !A                       |
| AND RA, RB   | B = A && B                   |
| OR RA, RB    | B = A \|\| B                 |
| XOR RA, RB   | B = A ^ B                    |
| CMP RA, RB   | set flags                    |
| LD RA, RB    | B = mem[A]                   |
| ST RA, RB    | mem[A] = B                   |
| LDI RB, imm  | B = imm                      |
| JMPR RB      | PC = B                       |
| JMP addr     | PC = addr                    |
| J{ZCEG} addr | if (flags == zceg) PC = addr |
| CLF          | flags = 0000                 |


