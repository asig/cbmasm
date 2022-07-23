/*
 * Copyright (c) 2020 Andreas Signer <asigner@gmail.com>
 *
 * This file is part of cbmasm.
 *
 * cbmasm is free software: you can redistribute it and/or
 * modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * cbmasm is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with cbmasm.  If not, see <http://www.gnu.org/licenses/>.
 */
package asm

import (
	"bytes"
	"testing"

	"github.com/asig/cbmasm/pkg/text"
)

func TestAssembler_assemble_6502(t *testing.T) {
	tests := []struct {
		name string
		text string
		want []byte
	}{
		{
			name: "Branches",
			text: `
L NOP  
  BVC L
  BVS L
  BMI L
  BPL L
  BCC L
  BCS L
  BNE L
  BEQ L
`,
			want: []byte{0xea, 0x50, 0xfd, 0x70, 0xfb, 0x30, 0xf9, 0x10, 0xf7, 0x90, 0xf5, 0xb0, 0xf3, 0xd0, 0xf1, 0xf0, 0xef},
		},
		{
			name: "Single instruction ADC $0078",
			text: "ADC $0078",
			want: []byte{0x65, 0x78},
		},
		{
			name: "Single instruction ADC $0078,X",
			text: "ADC $0078,X",
			want: []byte{0x75, 0x78},
		},
		{
			name: "Single instruction ADC $1234",
			text: "ADC $1234",
			want: []byte{0x6D, 0x34, 0x12},
		},
		{
			name: "Single instruction ADC $1234,X",
			text: "ADC $1234,X",
			want: []byte{0x7D, 0x34, 0x12},
		},
		{
			name: "Single instruction ADC $1234,Y",
			text: "ADC $1234,Y",
			want: []byte{0x79, 0x34, 0x12},
		},
		{
			name: "Single instruction ADC ($9A,X)",
			text: "ADC ($9A,X)",
			want: []byte{0x61, 0x9A},
		},
		{
			name: "Single instruction ADC ($BC),Y",
			text: "ADC ($BC),Y",
			want: []byte{0x71, 0xBC},
		},
		{
			name: "Single instruction AND $0056",
			text: "AND $0056",
			want: []byte{0x25, 0x56},
		},
		{
			name: "Single instruction AND $0078",
			text: "AND $0078",
			want: []byte{0x25, 0x78},
		},
		{
			name: "Single instruction AND $0078,X",
			text: "AND $0078,X",
			want: []byte{0x35, 0x78},
		},
		{
			name: "Single instruction AND $1234",
			text: "AND $1234",
			want: []byte{0x2D, 0x34, 0x12},
		},
		{
			name: "Single instruction AND $1234,X",
			text: "AND $1234,X",
			want: []byte{0x3D, 0x34, 0x12},
		},
		{
			name: "Single instruction AND $1234,Y",
			text: "AND $1234,Y",
			want: []byte{0x39, 0x34, 0x12},
		},
		{
			name: "Single instruction AND ($9A,X)",
			text: "AND ($9A,X)",
			want: []byte{0x21, 0x9A},
		},
		{
			name: "Single instruction AND ($BC),Y",
			text: "AND ($BC),Y",
			want: []byte{0x31, 0xBC},
		},
		{
			name: "Single instruction ASL $0078",
			text: "ASL $0078",
			want: []byte{0x06, 0x78},
		},
		{
			name: "Single instruction ASL $0078,X",
			text: "ASL $0078,X",
			want: []byte{0x16, 0x78},
		},
		{
			name: "Single instruction ASL $1234",
			text: "ASL $1234",
			want: []byte{0x0E, 0x34, 0x12},
		},
		{
			name: "Single instruction ASL $1234,X",
			text: "ASL $1234,X",
			want: []byte{0x1E, 0x34, 0x12},
		},
		{
			name: "Single instruction ASL A",
			text: "ASL A",
			want: []byte{0x0A},
		},
		{
			name: "Single instruction BIT $0078",
			text: "BIT $0078",
			want: []byte{0x24, 0x78},
		},
		{
			name: "Single instruction BIT $1234",
			text: "BIT $1234",
			want: []byte{0x2C, 0x34, 0x12},
		},
		{
			name: "Single instruction BRK",
			text: "BRK",
			want: []byte{0x00},
		},
		{
			name: "Single instruction CLC",
			text: "CLC",
			want: []byte{0x18},
		},
		{
			name: "Single instruction CLD",
			text: "CLD",
			want: []byte{0xD8},
		},
		{
			name: "Single instruction CLI",
			text: "CLI",
			want: []byte{0x58},
		},
		{
			name: "Single instruction CLV",
			text: "CLV",
			want: []byte{0xB8},
		},
		{
			name: "Single instruction CMP $0056",
			text: "CMP $0056",
			want: []byte{0xC5, 0x56},
		},
		{
			name: "Single instruction CMP $0078",
			text: "CMP $0078",
			want: []byte{0xC5, 0x78},
		},
		{
			name: "Single instruction CMP $0078,X",
			text: "CMP $0078,X",
			want: []byte{0xD5, 0x78},
		},
		{
			name: "Single instruction CMP $1234",
			text: "CMP $1234",
			want: []byte{0xCD, 0x34, 0x12},
		},
		{
			name: "Single instruction CMP $1234,X",
			text: "CMP $1234,X",
			want: []byte{0xDD, 0x34, 0x12},
		},
		{
			name: "Single instruction CMP $1234,Y",
			text: "CMP $1234,Y",
			want: []byte{0xD9, 0x34, 0x12},
		},
		{
			name: "Single instruction CMP ($9A,X)",
			text: "CMP ($9A,X)",
			want: []byte{0xC1, 0x9A},
		},
		{
			name: "Single instruction CMP ($BC),Y",
			text: "CMP ($BC),Y",
			want: []byte{0xD1, 0xBC},
		},
		{
			name: "Single instruction CPX $0056",
			text: "CPX $0056",
			want: []byte{0xE4, 0x56},
		},
		{
			name: "Single instruction CPX $0078",
			text: "CPX $0078",
			want: []byte{0xE4, 0x78},
		},
		{
			name: "Single instruction CPX $1234",
			text: "CPX $1234",
			want: []byte{0xEC, 0x34, 0x12},
		},
		{
			name: "Single instruction CPY $0056",
			text: "CPY $0056",
			want: []byte{0xC4, 0x56},
		},
		{
			name: "Single instruction CPY $0078",
			text: "CPY $0078",
			want: []byte{0xC4, 0x78},
		},
		{
			name: "Single instruction CPY $1234",
			text: "CPY $1234",
			want: []byte{0xCC, 0x34, 0x12},
		},
		{
			name: "Single instruction DEC $0078",
			text: "DEC $0078",
			want: []byte{0xC6, 0x78},
		},
		{
			name: "Single instruction DEC $0078,X",
			text: "DEC $0078,X",
			want: []byte{0xD6, 0x78},
		},
		{
			name: "Single instruction DEC $1234",
			text: "DEC $1234",
			want: []byte{0xCE, 0x34, 0x12},
		},
		{
			name: "Single instruction DEC $1234,X",
			text: "DEC $1234,X",
			want: []byte{0xDE, 0x34, 0x12},
		},
		{
			name: "Single instruction DEX",
			text: "DEX",
			want: []byte{0xCA},
		},
		{
			name: "Single instruction DEY",
			text: "DEY",
			want: []byte{0x88},
		},
		{
			name: "Single instruction EOR $0056",
			text: "EOR $0056",
			want: []byte{0x45, 0x56},
		},
		{
			name: "Single instruction EOR $0078",
			text: "EOR $0078",
			want: []byte{0x45, 0x78},
		},
		{
			name: "Single instruction EOR $0078,X",
			text: "EOR $0078,X",
			want: []byte{0x55, 0x78},
		},
		{
			name: "Single instruction EOR $1234",
			text: "EOR $1234",
			want: []byte{0x4D, 0x34, 0x12},
		},
		{
			name: "Single instruction EOR $1234,X",
			text: "EOR $1234,X",
			want: []byte{0x5D, 0x34, 0x12},
		},
		{
			name: "Single instruction EOR $1234,Y",
			text: "EOR $1234,Y",
			want: []byte{0x59, 0x34, 0x12},
		},
		{
			name: "Single instruction EOR ($9A,X)",
			text: "EOR ($9A,X)",
			want: []byte{0x41, 0x9A},
		},
		{
			name: "Single instruction EOR ($BC),Y",
			text: "EOR ($BC),Y",
			want: []byte{0x51, 0xBC},
		},
		{
			name: "Single instruction INC $0078",
			text: "INC $0078",
			want: []byte{0xE6, 0x78},
		},
		{
			name: "Single instruction INC $0078,X",
			text: "INC $0078,X",
			want: []byte{0xF6, 0x78},
		},
		{
			name: "Single instruction INC $1234",
			text: "INC $1234",
			want: []byte{0xEE, 0x34, 0x12},
		},
		{
			name: "Single instruction INC $1234,X",
			text: "INC $1234,X",
			want: []byte{0xFE, 0x34, 0x12},
		},
		{
			name: "Single instruction INX",
			text: "INX",
			want: []byte{0xE8},
		},
		{
			name: "Single instruction INY",
			text: "INY",
			want: []byte{0xC8},
		},
		{
			name: "Single instruction JMP $1234",
			text: "JMP $1234",
			want: []byte{0x4C, 0x34, 0x12},
		},
		{
			name: "Single instruction JMP ($ABCD)",
			text: "JMP ($ABCD)",
			want: []byte{0x6C, 0xCD, 0xAB},
		},
		{
			name: "Single instruction JSR $1234",
			text: "JSR $1234",
			want: []byte{0x20, 0x34, 0x12},
		},
		{
			name: "Single instruction LDA $0056",
			text: "LDA $0056",
			want: []byte{0xA5, 0x56},
		},
		{
			name: "Single instruction LDA $0078",
			text: "LDA $0078",
			want: []byte{0xA5, 0x78},
		},
		{
			name: "Single instruction LDA $0078,X",
			text: "LDA $0078,X",
			want: []byte{0xB5, 0x78},
		},
		{
			name: "Single instruction LDA $1234",
			text: "LDA $1234",
			want: []byte{0xAD, 0x34, 0x12},
		},
		{
			name: "Single instruction LDA $1234,X",
			text: "LDA $1234,X",
			want: []byte{0xBD, 0x34, 0x12},
		},
		{
			name: "Single instruction LDA $1234,Y",
			text: "LDA $1234,Y",
			want: []byte{0xB9, 0x34, 0x12},
		},
		{
			name: "Single instruction LDA ($9A,X)",
			text: "LDA ($9A,X)",
			want: []byte{0xA1, 0x9A},
		},
		{
			name: "Single instruction LDA ($BC),Y",
			text: "LDA ($BC),Y",
			want: []byte{0xB1, 0xBC},
		},
		{
			name: "Single instruction LDX $0056",
			text: "LDX $0056",
			want: []byte{0xA6, 0x56},
		},
		{
			name: "Single instruction LDX $0078",
			text: "LDX $0078",
			want: []byte{0xA6, 0x78},
		},
		{
			name: "Single instruction LDX $0078,Y",
			text: "LDX $0078,Y",
			want: []byte{0xB6, 0x78},
		},
		{
			name: "Single instruction LDX $1234",
			text: "LDX $1234",
			want: []byte{0xAE, 0x34, 0x12},
		},
		{
			name: "Single instruction LDX $1234,Y",
			text: "LDX $1234,Y",
			want: []byte{0xBE, 0x34, 0x12},
		},
		{
			name: "Single instruction LDY $0056",
			text: "LDY $0056",
			want: []byte{0xA4, 0x56},
		},
		{
			name: "Single instruction LDY $0078",
			text: "LDY $0078",
			want: []byte{0xA4, 0x78},
		},
		{
			name: "Single instruction LDY $0078,X",
			text: "LDY $0078,X",
			want: []byte{0xB4, 0x78},
		},
		{
			name: "Single instruction LDY $1234",
			text: "LDY $1234",
			want: []byte{0xAC, 0x34, 0x12},
		},
		{
			name: "Single instruction LDY $1234,X",
			text: "LDY $1234,X",
			want: []byte{0xBC, 0x34, 0x12},
		},
		{
			name: "Single instruction LSR $0078",
			text: "LSR $0078",
			want: []byte{0x46, 0x78},
		},
		{
			name: "Single instruction LSR $0078,X",
			text: "LSR $0078,X",
			want: []byte{0x56, 0x78},
		},
		{
			name: "Single instruction LSR $1234",
			text: "LSR $1234",
			want: []byte{0x4E, 0x34, 0x12},
		},
		{
			name: "Single instruction LSR $1234,X",
			text: "LSR $1234,X",
			want: []byte{0x5E, 0x34, 0x12},
		},
		{
			name: "Single instruction LSR A",
			text: "LSR A",
			want: []byte{0x4A},
		},
		{
			name: "Single instruction NOP",
			text: "NOP",
			want: []byte{0xEA},
		},
		{
			name: "Single instruction ORA $0056",
			text: "ORA $0056",
			want: []byte{0x05, 0x56},
		},
		{
			name: "Single instruction ORA $0078",
			text: "ORA $0078",
			want: []byte{0x05, 0x78},
		},
		{
			name: "Single instruction ORA $0078,X",
			text: "ORA $0078,X",
			want: []byte{0x15, 0x78},
		},
		{
			name: "Single instruction ORA $1234",
			text: "ORA $1234",
			want: []byte{0x0D, 0x34, 0x12},
		},
		{
			name: "Single instruction ORA $1234,X",
			text: "ORA $1234,X",
			want: []byte{0x1D, 0x34, 0x12},
		},
		{
			name: "Single instruction ORA $1234,Y",
			text: "ORA $1234,Y",
			want: []byte{0x19, 0x34, 0x12},
		},
		{
			name: "Single instruction ORA ($9A,X)",
			text: "ORA ($9A,X)",
			want: []byte{0x01, 0x9A},
		},
		{
			name: "Single instruction ORA ($BC),Y",
			text: "ORA ($BC),Y",
			want: []byte{0x11, 0xBC},
		},
		{
			name: "Single instruction PHA",
			text: "PHA",
			want: []byte{0x48},
		},
		{
			name: "Single instruction PHP",
			text: "PHP",
			want: []byte{0x08},
		},
		{
			name: "Single instruction PLA",
			text: "PLA",
			want: []byte{0x68},
		},
		{
			name: "Single instruction PLP",
			text: "PLP",
			want: []byte{0x28},
		},
		{
			name: "Single instruction ROL $0078",
			text: "ROL $0078",
			want: []byte{0x26, 0x78},
		},
		{
			name: "Single instruction ROL $0078,X",
			text: "ROL $0078,X",
			want: []byte{0x36, 0x78},
		},
		{
			name: "Single instruction ROL $1234",
			text: "ROL $1234",
			want: []byte{0x2E, 0x34, 0x12},
		},
		{
			name: "Single instruction ROL $1234,X",
			text: "ROL $1234,X",
			want: []byte{0x3E, 0x34, 0x12},
		},
		{
			name: "Single instruction ROL A",
			text: "ROL A",
			want: []byte{0x2A},
		},
		{
			name: "Single instruction ROR $0078",
			text: "ROR $0078",
			want: []byte{0x66, 0x78},
		},
		{
			name: "Single instruction ROR $0078,X",
			text: "ROR $0078,X",
			want: []byte{0x76, 0x78},
		},
		{
			name: "Single instruction ROR $1234",
			text: "ROR $1234",
			want: []byte{0x6E, 0x34, 0x12},
		},
		{
			name: "Single instruction ROR $1234,X",
			text: "ROR $1234,X",
			want: []byte{0x7E, 0x34, 0x12},
		},
		{
			name: "Single instruction ROR A",
			text: "ROR A",
			want: []byte{0x6A},
		},
		{
			name: "Single instruction RTI",
			text: "RTI",
			want: []byte{0x40},
		},
		{
			name: "Single instruction RTS",
			text: "RTS",
			want: []byte{0x60},
		},
		{
			name: "Single instruction SBC $0056",
			text: "SBC $0056",
			want: []byte{0xE5, 0x56},
		},
		{
			name: "Single instruction SBC $0078",
			text: "SBC $0078",
			want: []byte{0xE5, 0x78},
		},
		{
			name: "Single instruction SBC $0078,X",
			text: "SBC $0078,X",
			want: []byte{0xF5, 0x78},
		},
		{
			name: "Single instruction SBC $1234",
			text: "SBC $1234",
			want: []byte{0xED, 0x34, 0x12},
		},
		{
			name: "Single instruction SBC $1234,X",
			text: "SBC $1234,X",
			want: []byte{0xFD, 0x34, 0x12},
		},
		{
			name: "Single instruction SBC $1234,Y",
			text: "SBC $1234,Y",
			want: []byte{0xF9, 0x34, 0x12},
		},
		{
			name: "Single instruction SBC ($9A,X)",
			text: "SBC ($9A,X)",
			want: []byte{0xE1, 0x9A},
		},
		{
			name: "Single instruction SBC ($BC),Y",
			text: "SBC ($BC),Y",
			want: []byte{0xF1, 0xBC},
		},
		{
			name: "Single instruction SEC",
			text: "SEC",
			want: []byte{0x38},
		},
		{
			name: "Single instruction SED",
			text: "SED",
			want: []byte{0xF8},
		},
		{
			name: "Single instruction SEI",
			text: "SEI",
			want: []byte{0x78},
		},
		{
			name: "Single instruction STA $0078",
			text: "STA $0078",
			want: []byte{0x85, 0x78},
		},
		{
			name: "Single instruction STA $0078,X",
			text: "STA $0078,X",
			want: []byte{0x95, 0x78},
		},
		{
			name: "Single instruction STA $1234",
			text: "STA $1234",
			want: []byte{0x8D, 0x34, 0x12},
		},
		{
			name: "Single instruction STA $1234,X",
			text: "STA $1234,X",
			want: []byte{0x9D, 0x34, 0x12},
		},
		{
			name: "Single instruction STA $1234,Y",
			text: "STA $1234,Y",
			want: []byte{0x99, 0x34, 0x12},
		},
		{
			name: "Single instruction STA ($9A,X)",
			text: "STA ($9A,X)",
			want: []byte{0x81, 0x9A},
		},
		{
			name: "Single instruction STA ($BC),Y",
			text: "STA ($BC),Y",
			want: []byte{0x91, 0xBC},
		},
		{
			name: "Single instruction STX $0078",
			text: "STX $0078",
			want: []byte{0x86, 0x78},
		},
		{
			name: "Single instruction STX $1234",
			text: "STX $1234",
			want: []byte{0x8E, 0x34, 0x12},
		},
		{
			name: "Single instruction STX $78,Y",
			text: "STX $78,Y",
			want: []byte{0x96, 0x78},
		},
		{
			name: "Single instruction STY $0078",
			text: "STY $0078",
			want: []byte{0x84, 0x78},
		},
		{
			name: "Single instruction STY $1234",
			text: "STY $1234",
			want: []byte{0x8C, 0x34, 0x12},
		},
		{
			name: "Single instruction STY $78,X",
			text: "STY $78,X",
			want: []byte{0x94, 0x78},
		},
		{
			name: "Single instruction TAX",
			text: "TAX",
			want: []byte{0xAA},
		},
		{
			name: "Single instruction TAY",
			text: "TAY",
			want: []byte{0xA8},
		},
		{
			name: "Single instruction TSX",
			text: "TSX",
			want: []byte{0xBA},
		},
		{
			name: "Single instruction TXA",
			text: "TXA",
			want: []byte{0x8A},
		},
		{
			name: "Single instruction TXS",
			text: "TXS",
			want: []byte{0x9A},
		},
		{
			name: "Single instruction TYA",
			text: "TYA",
			want: []byte{0x98},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assembler := New([]string{}, "6502", "c128", []string{})
			src := " .org 0\n " + test.text
			assembler.Assemble(text.Process("", src))
			errs := assembler.Errors()
			if len(errs) != 0 {
				t.Errorf("Got %+v, want 0 errs", errs)
			}
			warnings := assembler.Warnings()
			if len(warnings) != 0 {
				t.Errorf("Got %+v, want 0 warnings", errs)
			}
			got := assembler.GetBytes()
			if bytes.Compare(got, test.want) != 0 {
				t.Errorf("Got %s, want %s", toString(got), toString(test.want))
			}
		})
	}
}

func TestAssembler_expr(t *testing.T) {
	tests := []struct {
		name string
		text string
		want []byte
	}{
		{
			name: "expr",
			text: `
	lda #scr('a')
	lda #'a'
`,
			want: []byte{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assembler := New([]string{}, "6502", "c128", []string{})
			src := " .org 0\n " + test.text
			assembler.Assemble(text.Process("", src))
			errs := assembler.Errors()
			if len(errs) != 0 {
				t.Errorf("Got %+v, want 0 errs", errs)
			}
			warnings := assembler.Warnings()
			if len(warnings) != 0 {
				t.Errorf("Got %+v, want 0 warnings", errs)
			}
			got := assembler.GetBytes()
			if bytes.Compare(got, test.want) != 0 {
				t.Errorf("Got %s, want %s", toString(got), toString(test.want))
			}
		})
	}
}
