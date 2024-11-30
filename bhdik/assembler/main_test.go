package main

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

func Test_parseInstruction(t *testing.T) {
	tests := []struct {
		line       string
		wantOpcode string
		wantArgs   []string
	}{
		{
			line:       "ADD R1, R2",
			wantOpcode: "ADD",
			wantArgs:   []string{"R1", "R2"},
		},
		{
			line:       "ADD R1",
			wantOpcode: "ADD",
			wantArgs:   []string{"R1"},
		},
		{
			line:       "SHR R1",
			wantOpcode: "SHR",
			wantArgs:   []string{"R1"},
		},
		{
			line:       "SHL R3",
			wantOpcode: "SHL",
			wantArgs:   []string{"R3"},
		},
		{
			line:       "NOT R0",
			wantOpcode: "NOT",
			wantArgs:   []string{"R0"},
		},
		{
			line:       "AND R3,R0",
			wantOpcode: "AND",
			wantArgs:   []string{"R3", "R0"},
		},
		{
			line:       "OR R3,R0",
			wantOpcode: "OR",
			wantArgs:   []string{"R3", "R0"},
		},
		{
			line:       "XOR R3,R0",
			wantOpcode: "XOR",
			wantArgs:   []string{"R3", "R0"},
		},
		{
			line:       "CMP R3,R0",
			wantOpcode: "CMP",
			wantArgs:   []string{"R3", "R0"},
		},
		{
			line:       "LD R3, R0",
			wantOpcode: "LD",
			wantArgs:   []string{"R3", "R0"},
		},
		{
			line:       "ST R3, R0",
			wantOpcode: "ST",
			wantArgs:   []string{"R3", "R0"},
		},
		{
			line:       "LDI R3, 1",
			wantOpcode: "LDI",
			wantArgs:   []string{"R3", "1"},
		},
		{
			line:       "LDI R3, :label",
			wantOpcode: "LDI",
			wantArgs:   []string{"R3", ":label"},
		},
		{
			line:       "JZ :label",
			wantOpcode: "JZ",
			wantArgs:   []string{":label"},
		},
		{
			line:       "JC :label",
			wantOpcode: "JC",
			wantArgs:   []string{":label"},
		},
		{
			line:       "JE :label",
			wantOpcode: "JE",
			wantArgs:   []string{":label"},
		},
		{
			line:       "JG :label",
			wantOpcode: "JG",
			wantArgs:   []string{":label"},
		},
		{
			line:       "CLF",
			wantOpcode: "CLF",
			wantArgs:   []string{},
		},
		{
			line:       "PSEL R3",
			wantOpcode: "PSEL",
			wantArgs:   []string{"R3"},
		},
		{
			line:       "PIN R2",
			wantOpcode: "PIN",
			wantArgs:   []string{"R2"},
		},
		{
			line:       "POUT R1",
			wantOpcode: "POUT",
			wantArgs:   []string{"R1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			gotOpcode, gotArgs := parseInstruction(tt.line)
			if gotOpcode != tt.wantOpcode {
				t.Errorf("parseInstruction() gotOpcode = %v, wantOpcode %v", gotOpcode, tt.wantOpcode)
			}
			if slices.Compare(gotArgs, tt.wantArgs) != 0 {
				t.Errorf("parseInstruction() gotArgs = %v, wantArgs %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func Test_immediateArg(t *testing.T) {
	tests := []struct {
		arg       string
		wantImm   int
		wantLabel string
	}{
		{
			arg:       "100",
			wantImm:   100,
			wantLabel: "",
		},
		{
			arg:       "0b1001_0110",
			wantImm:   0b1001_0110,
			wantLabel: "",
		},
		{
			arg:       "0x84",
			wantImm:   0x84,
			wantLabel: "",
		},
		{
			arg:       ":label",
			wantImm:   -1,
			wantLabel: "label",
		},
	}
	for _, tt := range tests {
		t.Run(tt.arg, func(t *testing.T) {
			gotImm, gotLabel := immediateArg(tt.arg)
			if gotImm != tt.wantImm {
				t.Errorf("immediateArg() gotImm = %v, wantImm %v", gotImm, tt.wantImm)
			}
			if gotLabel != tt.wantLabel {
				t.Errorf("immediateArg() gotLabel = %v, wantImm %v", gotLabel, tt.wantLabel)
			}
		})
	}
}

func Test_encodeInstruction(t *testing.T) {
	tests := []struct {
		opcode          string
		args            []string
		wantInstruction byte
		wantImmByte     int
		wantImmLabel    string
	}{
		{
			opcode:          "ADD",
			args:            []string{"R1", "R2"},
			wantInstruction: 0b1001_0001,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
		{
			opcode:          "ADD",
			args:            []string{"R1"},
			wantInstruction: 0b0101_0001,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
		{
			opcode:          "SHR",
			args:            []string{"R1"},
			wantInstruction: 0b0101_0011,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
		{
			opcode:          "SHL",
			args:            []string{"R3"},
			wantInstruction: 0b1111_0101,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
		{
			opcode:          "NOT",
			args:            []string{"R0"},
			wantInstruction: 0b0000_0111,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
		{
			opcode:          "AND",
			args:            []string{"R3", "R0"},
			wantInstruction: 0b0011_1001,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
		{
			opcode:          "OR",
			args:            []string{"R3", "R0"},
			wantInstruction: 0b0011_1011,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
		{
			opcode:          "XOR",
			args:            []string{"R3", "R0"},
			wantInstruction: 0b0011_1101,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
		{
			opcode:          "CMP",
			args:            []string{"R3", "R0"},
			wantInstruction: 0b0011_1111,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
		{
			opcode:          "LD",
			args:            []string{"R3", "R0"},
			wantInstruction: 0b0011_0000,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
		{
			opcode:          "ST",
			args:            []string{"R3", "R0"},
			wantInstruction: 0b0011_0010,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
		{
			opcode:          "LDI",
			args:            []string{"R3", "1"},
			wantInstruction: 0b1100_0100,
			wantImmByte:     1,
			wantImmLabel:    "",
		},
		{
			opcode:          "LDI",
			args:            []string{"R3", ":label"},
			wantInstruction: 0b1100_0100,
			wantImmByte:     -1,
			wantImmLabel:    "label",
		},
		{
			opcode:          "JZ",
			args:            []string{":label"},
			wantInstruction: 0b1000_1010,
			wantImmByte:     -1,
			wantImmLabel:    "label",
		},
		{
			opcode:          "JC",
			args:            []string{":label"},
			wantInstruction: 0b0100_1010,
			wantImmByte:     -1,
			wantImmLabel:    "label",
		},
		{
			opcode:          "JE",
			args:            []string{":label"},
			wantInstruction: 0b0010_1010,
			wantImmByte:     -1,
			wantImmLabel:    "label",
		},
		{
			opcode:          "JG",
			args:            []string{":label"},
			wantInstruction: 0b0001_1010,
			wantImmByte:     -1,
			wantImmLabel:    "label",
		},
		{
			opcode:          "CLF",
			args:            []string{},
			wantInstruction: 0b0000_1100,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
		{
			opcode:          "PSEL",
			args:            []string{"R3"},
			wantInstruction: 0b1111_1110,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
		{
			opcode:          "PIN",
			args:            []string{"R2"},
			wantInstruction: 0b1000_1110,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
		{
			opcode:          "POUT",
			args:            []string{"R1"},
			wantInstruction: 0b0101_1110,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s %s", tt.opcode, strings.Join(tt.args, ",")), func(t *testing.T) {
			gotInstruction, gotImmByte, gotImmLabel := encodeInstruction(tt.opcode, tt.args)
			if gotInstruction != tt.wantInstruction {
				t.Errorf("encodeInstruction() gotInstruction = %08b, want %08b", gotInstruction, tt.wantInstruction)
			}
			if gotImmByte != tt.wantImmByte {
				t.Errorf("encodeInstruction() gotImmByte = %v, want %v", gotImmByte, tt.wantImmByte)
			}
			if gotImmLabel != tt.wantImmLabel {
				t.Errorf("encodeInstruction() gotImmLabel = %v, want %v", gotImmLabel, tt.wantImmLabel)
			}
		})
	}
}
