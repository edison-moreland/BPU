package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func exitE(err error) {
	fmt.Printf("Error!\n %v\n", err)
	os.Exit(1)
}

func warn(msg string) {
	fmt.Printf("Warning! %s\n", msg)
}

func main() {
	if len(os.Args) < 2 {
		exitE(fmt.Errorf("usage: assembler <file>"))
	}

	inputFile, err := os.ReadFile(os.Args[1])
	if err != nil {
		exitE(err)
	}

	program := []byte{}
	programOffset := 0

	labelLocation := map[string]int{} // Program location of each label
	labelUses := map[string][]int{}   // Program location where each label is used

	hasLabel := false // next instruction is labeled
	instLabel := ""
	for _, line := range strings.Split(string(inputFile), "\n") {
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		switch line[0] {
		case '-':
			// Comment
			continue
		case ':':
			// Label
			if hasLabel {
				warn(fmt.Sprintf("unused label: %s", instLabel))
			}
			instLabel = line[1:]
			if _, ok := labelLocation[instLabel]; !ok {
				// This label hasn't been seen before
				labelLocation[instLabel] = -1
				labelUses[instLabel] = []int{}
			}

			hasLabel = true
			continue
		case '%':
			// Assembler directive
			directive := line[1:]
			directive, args := parseInstruction(directive)

			switch directive {
			case "OFFSET":
				if len(args) != 1 {
					exitE(errors.New("OFFSET directive requires one argument"))
				}
				offset, _ := immediateArg(args[0])
				if offset > 0xFF {
					exitE(fmt.Errorf("OFFSET too large: %d", offset))
				}

				programOffset = offset
			}

			continue
		}

		// This line must be an instruction
		if hasLabel {
			// a label precedes this instruction
			instructionLocation := len(program)
			labelLocation[instLabel] = instructionLocation
			hasLabel = false
		}

		opcode, args := parseInstruction(line)
		instruction, immByte, immLabel := encodeInstruction(opcode, args)

		program = append(program, instruction)

		if immLabel != "" {
			// the next byte should be filled with the location of a label
			// this will be filled in with the real address later
			if _, ok := labelLocation[immLabel]; !ok {
				// This label hasn't been seen before
				labelLocation[immLabel] = -1
				labelUses[immLabel] = []int{}
			}

			immLocation := len(program)
			labelUses[immLabel] = append(labelUses[immLabel], immLocation)
			immByte = 0
		}

		if immByte != -1 {
			program = append(program, byte(immByte))
		}
	}

	// Fixup labels
	for name, location := range labelLocation {
		if location == -1 {
			exitE(fmt.Errorf("label with no destination: %s", name))
		}

		location += programOffset
		if location >= 0xFF {
			warn(fmt.Sprintf("label destination overflows address space after offset: %s", name))
		}

		usages := labelUses[name]
		if len(usages) == 0 {
			warn(fmt.Sprintf("unused label: %s", name))
			continue
		}

		for _, useLocation := range labelUses[name] {
			program[useLocation] = byte(location)
		}
	}

	// Print program
	fmt.Println("---")
	fmt.Println("Address         | Data")
	for i, b := range program {
		fmt.Printf("0x%02x 0b%08b | 0x%02x 0b%08b\n", i+programOffset, i+programOffset, b, b)
	}
}

// Parse line into an encoded instruction and possibly an immediate byte or label
func parseInstruction(line string) (string, []string) {
	op, rawArgs, hasArgs := strings.Cut(line, " ")
	var args []string
	if hasArgs {
		arg1, arg2, hasSecondArg := strings.Cut(rawArgs, ",")
		args = append(args, strings.TrimSpace(arg1))
		if hasSecondArg {
			args = append(args, strings.TrimSpace(arg2))
		}
	}

	return op, args
}

var matchCondFlags = regexp.MustCompile(`[ZCEG]{1,4}`)

func encodeInstruction(opcode string, args []string) (byte, int, string) {
	instruction := byte(0)
	immediateByte := -1
	immediateLabel := ""

	switch opcode {
	case "ADD":
		instruction = aluInstruction(0b000, args)
	case "SHR":
		instruction = aluInstruction(0b001, args)
	case "SHL":
		instruction = aluInstruction(0b010, args)
	case "NOT":
		instruction = aluInstruction(0b011, args)
	case "AND":
		instruction = aluInstruction(0b100, args)
	case "OR":
		instruction = aluInstruction(0b101, args)
	case "XOR":
		instruction = aluInstruction(0b110, args)
	case "CMP":
		instruction = aluInstruction(0b111, args)
	case "LD":
		instruction = loadStoreInstruction(false, args)
	case "ST":
		instruction = loadStoreInstruction(true, args)
	case "LDI":
		instruction, immediateByte, immediateLabel = loadImmInstruction(args)
	case "JMPR":
		instruction = jumpRegInstruction(args)
	case "JMP":
		instruction, immediateByte, immediateLabel = jumpImmInstruction(args)
	case "CLF":
		instruction = 0b0000_1100
	case "PSEL":
		instruction = peripheralBusInstruction(output, address, args)
	case "PIN":
		instruction = peripheralBusInstruction(input, data, args)
	case "POUT":
		instruction = peripheralBusInstruction(output, data, args)

	default:
		if len(opcode) > 1 && opcode[0] == 'J' && matchCondFlags.MatchString(opcode[1:]) {
			conditionFlags := opcode[1:]
			instruction, immediateByte, immediateLabel = jumpConditionalInstruction(conditionFlags, args)
		} else {
			exitE(fmt.Errorf("invalid instruction: %s", opcode))
		}
	}

	return instruction, immediateByte, immediateLabel
}

func aluInstruction(op byte, args []string) byte {
	regA := byte(0)
	regB := byte(0)

	switch len(args) {
	case 1:
		regA = registerArg(args[0])
		regB = regA
	case 2:
		regA = registerArg(args[0])
		regB = registerArg(args[1])
	default:
		exitE(fmt.Errorf("ALU op should have 1 or 2 arguments"))

	}

	return 0b0000_0001 | (op << 1) | (regA << 4) | (regB << 6)
}

func loadStoreInstruction(isStore bool, args []string) byte {
	if len(args) != 2 {
		exitE(fmt.Errorf("LD/ST op should have 2 arguments"))
	}
	regA := registerArg(args[0])
	regB := registerArg(args[1])

	inst := 0b0000_0000 | (regA << 4) | (regB << 6)
	if isStore {
		inst |= 0b0000_0010
	}

	return inst
}

func loadImmInstruction(args []string) (byte, int, string) {
	if len(args) != 2 {
		exitE(fmt.Errorf("LDI op should have 2 arguments"))
	}

	regB := registerArg(args[0])
	immByte, immLabel := immediateArg(args[1])

	inst := 0b0000_0100 | (regB << 6)
	return inst, immByte, immLabel
}

func jumpRegInstruction(args []string) byte {
	if len(args) != 1 {
		exitE(fmt.Errorf("JMPR op should have 1 argument"))
	}

	regB := registerArg(args[0])

	return 0b0000_0110 | (regB << 6)
}

func jumpImmInstruction(args []string) (byte, int, string) {
	if len(args) != 1 {
		exitE(fmt.Errorf("JMP op should have 1 argument"))
	}

	immByte, immLabel := immediateArg(args[0])

	return 0b0000_1000, immByte, immLabel
}

func jumpConditionalInstruction(flags string, args []string) (byte, int, string) {
	if len(args) != 1 {
		exitE(fmt.Errorf("conditional jump op should have 1 argument"))
	}

	immByte, immLabel := immediateArg(args[0])

	inst := byte(0b0000_1010)
	for _, c := range flags {
		switch c {
		case 'Z':
			inst |= 0b1000_0000
		case 'C':
			inst |= 0b0100_0000
		case 'E':
			inst |= 0b0010_0000
		case 'G':
			inst |= 0b0001_0000
		}
	}

	return inst, immByte, immLabel
}

type pbusDirection bool

const (
	input  pbusDirection = false
	output pbusDirection = true
)

type pbusKind bool

const (
	data    pbusKind = false
	address pbusKind = true
)

func peripheralBusInstruction(direction pbusDirection, kind pbusKind, args []string) byte {
	if len(args) != 1 {
		exitE(fmt.Errorf("peripheral op should have 1 argument"))
	}

	regB := registerArg(args[0])

	inst := 0b0000_1110 | (regB << 6)
	if direction == output {
		inst |= 0b0001_0000
	}
	if kind == address {
		inst |= 0b0010_0000
	}

	return inst
}

func registerArg(arg string) byte {
	if len(arg) != 2 || arg[0] != 'R' {
		exitE(fmt.Errorf("unknown register: %s", arg))
	}

	regI, err := strconv.Atoi(string(arg[1]))
	if err != nil {
		exitE(fmt.Errorf("unknown register: %s", arg))
	}

	return byte(regI)
}

func immediateArg(arg string) (int, string) {
	if arg[0] == ':' {
		// Arg is a label
		return -1, arg[1:]
	}

	// Arg is a literal
	imm, err := strconv.ParseUint(arg, 0, 8)
	if err != nil {
		exitE(fmt.Errorf("immediate must be label or decimal value: %s", arg))
	}

	return int(imm), ""
}
