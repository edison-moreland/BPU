package main

import "testing"

func Test_parseInstruction(t *testing.T) {
	tests := []struct {
		line            string
		wantInstruction byte
		wantImmByte     int
		wantImmLabel    string
	}{
		{
			line:            "ADD R1, R2",
			wantInstruction: 0b1001_0001,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
		{
			line:            "ADD R1",
			wantInstruction: 0b0101_0001,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
		{
			line:            "SHR R1",
			wantInstruction: 0b0101_0011,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
		{
			line:            "SHL R3",
			wantInstruction: 0b1111_0101,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
		{
			line:            "NOT R0",
			wantInstruction: 0b0000_0111,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
		{
			line:            "AND R3,R0",
			wantInstruction: 0b0011_1001,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
		{
			line:            "OR R3,R0",
			wantInstruction: 0b0011_1011,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
		{
			line:            "XOR R3,R0",
			wantInstruction: 0b0011_1101,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
		{
			line:            "CMP R3,R0",
			wantInstruction: 0b0011_1111,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
		{
			line:            "LD R3, R0",
			wantInstruction: 0b0011_0000,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
		{
			line:            "ST R3, R0",
			wantInstruction: 0b0011_0010,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
		{
			line:            "LDI R3, 1",
			wantInstruction: 0b1100_0100,
			wantImmByte:     1,
			wantImmLabel:    "",
		},
		{
			line:            "LDI R3, :label",
			wantInstruction: 0b1100_0100,
			wantImmByte:     -1,
			wantImmLabel:    "label",
		},
		{
			line:            "JZ :label",
			wantInstruction: 0b1000_1010,
			wantImmByte:     -1,
			wantImmLabel:    "label",
		},
		{
			line:            "JC :label",
			wantInstruction: 0b0100_1010,
			wantImmByte:     -1,
			wantImmLabel:    "label",
		},
		{
			line:            "JE :label",
			wantInstruction: 0b0010_1010,
			wantImmByte:     -1,
			wantImmLabel:    "label",
		},
		{
			line:            "JG :label",
			wantInstruction: 0b0001_1010,
			wantImmByte:     -1,
			wantImmLabel:    "label",
		},
		{
			line:            "CLF",
			wantInstruction: 0b0000_1100,
			wantImmByte:     -1,
			wantImmLabel:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			gotInstruction, gotImmediateByte, gotImmediateLabel := parseInstruction(tt.line)
			if gotInstruction != tt.wantInstruction {
				t.Errorf("parseInstruction() gotInstruction = %08b, want %08b", gotInstruction, tt.wantInstruction)
			}
			if gotImmediateByte != tt.wantImmByte {
				t.Errorf("parseInstruction() gotImmediateByte = %v, want %v", gotImmediateByte, tt.wantImmByte)
			}
			if gotImmediateLabel != tt.wantImmLabel {
				t.Errorf("parseInstruction() gotImmediateLabel = %v, want %v", gotImmediateLabel, tt.wantImmLabel)
			}
		})
	}
}
