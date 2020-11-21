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

func TestAssembler_assembleZ80(t *testing.T) {
	tests := []struct {
		name string
		text string
		want []byte
	}{
		{
			name: "Relative branch backwards",
			text: `
foo: NOP
  NOP
  DJNZ foo
`,
			want: []byte{0x00, 0x00, 0x10, 0xfc},
		},
		{
			name: "Relative branch forwards",
			text: ` DJNZ foo
  NOP
  NOP
foo: NOP
`,
			want: []byte{0x10, 0x02, 0x00, 0x00, 0x00},
		},

		{
			name: "Relative branch backwards with condition",
			text: `l: NOP
  JR NC,l
  JR Z,l
`,
			want: []byte{0x00, 0x30, 0xfd, 0x28, 0xfb},
		},
		{
			name: "Single instruction RES 1,C",
			text: "RES 1,C",
			want: []byte{0xcb, 0x89},
		},
		{
			name: "Single instruction RES 2,(HL)",
			text: "RES 2,(HL)",
			want: []byte{0xcb, 0x96},
		},
		{
			name: "Single instruction RES 3,(IX+63)",
			text: "RES 3,(IX+63)",
			want: []byte{0xdd, 0xcb, 0x3f, 0x9e},
		},
		{
			name: "Single instruction RES 4,(IY-27)",
			text: "RES 4,(IY-27)",
			want: []byte{0xfd, 0xcb, 0xe5, 0xa6},
		},
		{
			name: "Single instruction ADC A,$56",
			text: "ADC A,$56",
			want: []byte{0xce, 0x56},
		},
		{
			name: "Single instruction ADC A,C",
			text: "ADC A,C",
			want: []byte{0x89},
		},
		{
			name: "Single instruction ADC A,(HL)",
			text: "ADC A,(HL)",
			want: []byte{0x8e},
		},
		{
			name: "Single instruction ADC A,(IX+$12)",
			text: "ADC A,(IX+$12)",
			want: []byte{0xdd, 0x8e, 0x12},
		},
		{
			name: "Single instruction ADC A,(IY-$12)",
			text: "ADC A,(IY-$12)",
			want: []byte{0xfd, 0x8e, 0xee},
		},
		{
			name: "Single instruction ADC HL,SP",
			text: "ADC HL,SP",
			want: []byte{0xed, 0x7a},
		},
		{
			name: "Single instruction ADD A,$56",
			text: "ADD A,$56",
			want: []byte{0xc6, 0x56},
		},
		{
			name: "Single instruction ADD A,C",
			text: "ADD A,C",
			want: []byte{0x81},
		},
		{
			name: "Single instruction ADD A,(HL)",
			text: "ADD A,(HL)",
			want: []byte{0x86},
		},
		{
			name: "Single instruction ADD A,(IX+$12)",
			text: "ADD A,(IX+$12)",
			want: []byte{0xdd, 0x86, 0x12},
		},
		{
			name: "Single instruction ADD A,(IY-$12)",
			text: "ADD A,(IY-$12)",
			want: []byte{0xfd, 0x86, 0xee},
		},
		{
			name: "Single instruction ADD HL,SP",
			text: "ADD HL,SP",
			want: []byte{0x39},
		},
		{
			name: "Single instruction ADD IX,DE",
			text: "ADD IX,DE",
			want: []byte{0xdd, 0x19},
		},
		{
			name: "Single instruction ADD IY,DE",
			text: "ADD IY,DE",
			want: []byte{0xfd, 0x19},
		},
		{
			name: "Single instruction AND $56",
			text: "AND $56",
			want: []byte{0xe6, 0x56},
		},
		{
			name: "Single instruction AND C",
			text: "AND C",
			want: []byte{0xa1},
		},
		{
			name: "Single instruction AND (HL)",
			text: "AND (HL)",
			want: []byte{0xa6},
		},
		{
			name: "Single instruction AND (IX+$12)",
			text: "AND (IX+$12)",
			want: []byte{0xdd, 0xa6, 0x12},
		},
		{
			name: "Single instruction AND (IY-$12)",
			text: "AND (IY-$12)",
			want: []byte{0xfd, 0xa6, 0xee},
		},
		{
			name: "Single instruction BIT 0,(HL)",
			text: "BIT 0,(HL)",
			want: []byte{0xcb, 0x46},
		},
		{
			name: "Single instruction BIT 1,(IX+$12)",
			text: "BIT 1,(IX+$12)",
			want: []byte{0xdd, 0xcb, 0x12, 0x4e},
		},
		{
			name: "Single instruction BIT 2,(IY-$12)",
			text: "BIT 2,(IY-$12)",
			want: []byte{0xfd, 0xcb, 0xee, 0x56},
		},
		{
			name: "Single instruction BIT 3,C",
			text: "BIT 3,C",
			want: []byte{0xcb, 0x59},
		},
		{
			name: "Single instruction CALL $5678",
			text: "CALL $5678",
			want: []byte{0xcd, 0x78, 0x56},
		},
		{
			name: "Single instruction CALL NZ,$5678",
			text: "CALL NZ,$5678",
			want: []byte{0xc4, 0x78, 0x56},
		},
		{
			name: "Single instruction CCF",
			text: "CCF",
			want: []byte{0x3f},
		},
		{
			name: "Single instruction CP $56",
			text: "CP $56",
			want: []byte{0xfe, 0x56},
		},
		{
			name: "Single instruction CP C",
			text: "CP C",
			want: []byte{0xb9},
		},
		{
			name: "Single instruction CPD",
			text: "CPD",
			want: []byte{0xed, 0xa9},
		},
		{
			name: "Single instruction CPDR",
			text: "CPDR",
			want: []byte{0xed, 0xb9},
		},
		{
			name: "Single instruction CP (HL)",
			text: "CP (HL)",
			want: []byte{0xbe},
		},
		{
			name: "Single instruction CPI",
			text: "CPI",
			want: []byte{0xed, 0xa1},
		},
		{
			name: "Single instruction CPIR",
			text: "CPIR",
			want: []byte{0xed, 0xb1},
		},
		{
			name: "Single instruction CPL",
			text: "CPL",
			want: []byte{0x2f},
		},
		{
			name: "Single instruction DAA",
			text: "DAA",
			want: []byte{0x27},
		},
		{
			name: "Single instruction DEC C",
			text: "DEC C",
			want: []byte{0x0d},
		},
		{
			name: "Single instruction DEC DE",
			text: "DEC DE",
			want: []byte{0x1b},
		},
		{
			name: "Single instruction DEC (HL)",
			text: "DEC (HL)",
			want: []byte{0x35},
		},
		{
			name: "Single instruction DEC IX",
			text: "DEC IX",
			want: []byte{0xdd, 0x2b},
		},
		{
			name: "Single instruction DEC (IX+$12)",
			text: "DEC (IX+$12)",
			want: []byte{0xdd, 0x35, 0x12},
		},
		{
			name: "Single instruction DEC IY",
			text: "DEC IY",
			want: []byte{0xfd, 0x2b},
		},
		{
			name: "Single instruction DEC (IY-$12)",
			text: "DEC (IY-$12)",
			want: []byte{0xfd, 0x35, 0xee},
		},
		{
			name: "Single instruction DI",
			text: "DI",
			want: []byte{0xf3},
		},
		{
			name: "Single instruction EI",
			text: "EI",
			want: []byte{0xfb},
		},
		{
			name: "Single instruction EX AF, AF'",
			text: "EX AF, AF'",
			want: []byte{0x08},
		},
		{
			name: "Single instruction EX DE, HL",
			text: "EX DE, HL",
			want: []byte{0xeb},
		},
		{
			name: "Single instruction EX (SP), HL",
			text: "EX (SP), HL",
			want: []byte{0xe3},
		},
		{
			name: "Single instruction EX (SP), IX",
			text: "EX (SP), IX",
			want: []byte{0xdd, 0xe3},
		},
		{
			name: "Single instruction EX (SP), IY",
			text: "EX (SP), IY",
			want: []byte{0xfd, 0xe3},
		},
		{
			name: "Single instruction EXX",
			text: "EXX",
			want: []byte{0xd9},
		},
		{
			name: "Single instruction HALT",
			text: "HALT",
			want: []byte{0x76},
		},
		{
			name: "Single instruction IM 0",
			text: "IM 0",
			want: []byte{0xed, 0x46},
		},
		{
			name: "Single instruction IM 1",
			text: "IM 1",
			want: []byte{0xed, 0x56},
		},
		{
			name: "Single instruction IM 2",
			text: "IM 2",
			want: []byte{0xed, 0x5e},
		},
		{
			name: "Single instruction IN A,($78)",
			text: "IN A,($78)",
			want: []byte{0xdb, 0x78},
		},
		{
			name: "Single instruction IN C,(C)",
			text: "IN C,(C)",
			want: []byte{0xed, 0x48},
		},
		{
			name: "Single instruction INC C",
			text: "INC C",
			want: []byte{0x0c},
		},
		{
			name: "Single instruction INC DE",
			text: "INC DE",
			want: []byte{0x13},
		},
		{
			name: "Single instruction INC (HL)",
			text: "INC (HL)",
			want: []byte{0x34},
		},
		{
			name: "Single instruction INC IX",
			text: "INC IX",
			want: []byte{0xdd, 0x23},
		},
		{
			name: "Single instruction INC (IX+$12)",
			text: "INC (IX+$12)",
			want: []byte{0xdd, 0x34, 0x12},
		},
		{
			name: "Single instruction INC IY",
			text: "INC IY",
			want: []byte{0xfd, 0x23},
		},
		{
			name: "Single instruction INC (IY-$12)",
			text: "INC (IY-$12)",
			want: []byte{0xfd, 0x34, 0xee},
		},
		{
			name: "Single instruction IND",
			text: "IND",
			want: []byte{0xed, 0xaa},
		},
		{
			name: "Single instruction INDR",
			text: "INDR",
			want: []byte{0xed, 0xba},
		},
		{
			name: "Single instruction INI",
			text: "INI",
			want: []byte{0xed, 0xa2},
		},
		{
			name: "Single instruction INIR",
			text: "INIR",
			want: []byte{0xed, 0xb2},
		},
		{
			name: "Single instruction JP $5678",
			text: "JP $5678",
			want: []byte{0xc3, 0x78, 0x56},
		},
		{
			name: "Single instruction JP (HL)",
			text: "JP (HL)",
			want: []byte{0xe9},
		},
		{
			name: "Single instruction JP (IX)",
			text: "JP (IX)",
			want: []byte{0xdd, 0xe9},
		},
		{
			name: "Single instruction JP (IY)",
			text: "JP (IY)",
			want: []byte{0xfd, 0xe9},
		},
		{
			name: "Single instruction JP NC,$5678",
			text: "JP NC,$5678",
			want: []byte{0xd2, 0x78, 0x56},
		},
		{
			name: "Single instruction LD ($5678), A",
			text: "LD ($5678), A",
			want: []byte{0x32, 0x78, 0x56},
		},
		{
			name: "Single instruction LD ($5678), BC",
			text: "LD ($5678), BC",
			want: []byte{0xed, 0x43, 0x78, 0x56},
		},
		{
			name: "Single instruction LD ($5678), HL",
			text: "LD ($5678), HL",
			want: []byte{0x22, 0x78, 0x56},
		},
		{
			name: "Single instruction LD ($5678), IX",
			text: "LD ($5678), IX",
			want: []byte{0xdd, 0x22, 0x78, 0x56},
		},
		{
			name: "Single instruction LD ($5678), IY",
			text: "LD ($5678), IY",
			want: []byte{0xfd, 0x22, 0x78, 0x56},
		},
		{
			name: "Single instruction LD A, ($5678)",
			text: "LD A, ($5678)",
			want: []byte{0x3a, 0x78, 0x56},
		},
		{
			name: "Single instruction LD A, (BC)",
			text: "LD A, (BC)",
			want: []byte{0x0a},
		},
		{
			name: "Single instruction LD A,C",
			text: "LD A,C",
			want: []byte{0x79},
		},
		{
			name: "Single instruction LD A, (DE)",
			text: "LD A, (DE)",
			want: []byte{0x1a},
		},
		{
			name: "Single instruction LD A, I",
			text: "LD A, I",
			want: []byte{0xed, 0x57},
		},
		{
			name: "Single instruction LD B,C",
			text: "LD B,C",
			want: []byte{0x41},
		},
		{
			name: "Single instruction LD BC, $5678",
			text: "LD BC, $5678",
			want: []byte{0x01, 0x78, 0x56},
		},
		{
			name: "Single instruction LD (BC),A",
			text: "LD (BC),A",
			want: []byte{0x02},
		},
		{
			name: "Single instruction LD C, $56",
			text: "LD C, $56",
			want: []byte{0x0e, 0x56},
		},
		{
			name: "Single instruction LD C, (HL)",
			text: "LD C, (HL)",
			want: []byte{0x4e},
		},
		{
			name: "Single instruction LD C, (IX+$12)",
			text: "LD C, (IX+$12)",
			want: []byte{0xdd, 0x4e, 0x12},
		},
		{
			name: "Single instruction LD C, (IY-$12)",
			text: "LD C, (IY-$12)",
			want: []byte{0xfd, 0x4e, 0xee},
		},
		{
			name: "Single instruction LDD",
			text: "LDD",
			want: []byte{0xed, 0xa8},
		},
		{
			name: "Single instruction LD DE, ($5678)",
			text: "LD DE, ($5678)",
			want: []byte{0xed, 0x5b, 0x78, 0x56},
		},
		{
			name: "Single instruction LD (DE),A",
			text: "LD (DE),A",
			want: []byte{0x12},
		},
		{
			name: "Single instruction LDDR",
			text: "LDDR",
			want: []byte{0xed, 0xb8},
		},
		{
			name: "Single instruction LD (HL), $56",
			text: "LD (HL), $56",
			want: []byte{0x36, 0x56},
		},
		{
			name: "Single instruction LD HL, ($5678)",
			text: "LD HL, ($5678)",
			want: []byte{0x2a, 0x78, 0x56},
		},
		{
			name: "Single instruction LD (HL),C",
			text: "LD (HL),C",
			want: []byte{0x71},
		},
		{
			name: "Single instruction LDI",
			text: "LDI",
			want: []byte{0xed, 0xa0},
		},
		{
			name: "Single instruction LD I, A",
			text: "LD I, A",
			want: []byte{0xed, 0x47},
		},
		{
			name: "Single instruction LDIR",
			text: "LDIR",
			want: []byte{0xed, 0xb0},
		},
		{
			name: "Single instruction LD (IX+$12), $56",
			text: "LD (IX+$12), $56",
			want: []byte{0xdd, 0x36, 0x12, 0x56},
		},
		{
			name: "Single instruction LD (IX+$12),C",
			text: "LD (IX+$12),C",
			want: []byte{0xdd, 0x71, 0x12},
		},
		{
			name: "Single instruction LD IX, ($5678)",
			text: "LD IX, ($5678)",
			want: []byte{0xdd, 0x2a, 0x78, 0x56},
		},
		{
			name: "Single instruction LD IX, $5678",
			text: "LD IX, $5678",
			want: []byte{0xdd, 0x21, 0x78, 0x56},
		},
		{
			name: "Single instruction LD (IY-$12), $56",
			text: "LD (IY-$12), $56",
			want: []byte{0xfd, 0x36, 0xee, 0x56},
		},
		{
			name: "Single instruction LD (IY-$12),C",
			text: "LD (IY-$12),C",
			want: []byte{0xfd, 0x71, 0xee},
		},
		{
			name: "Single instruction LD IY, ($5678)",
			text: "LD IY, ($5678)",
			want: []byte{0xfd, 0x2a, 0x78, 0x56},
		},
		{
			name: "Single instruction LD IY, $5678",
			text: "LD IY, $5678",
			want: []byte{0xfd, 0x21, 0x78, 0x56},
		},
		{
			name: "Single instruction LD R, A",
			text: "LD R, A",
			want: []byte{0xed, 0x4f},
		},
		{
			name: "Single instruction LD SP, HL",
			text: "LD SP, HL",
			want: []byte{0xf9},
		},
		{
			name: "Single instruction LD SP, IX",
			text: "LD SP, IX",
			want: []byte{0xdd, 0xf9},
		},
		{
			name: "Single instruction LD SP, IY",
			text: "LD SP, IY",
			want: []byte{0xfd, 0xf9},
		},
		{
			name: "Single instruction NEG",
			text: "NEG",
			want: []byte{0xed, 0x44},
		},
		{
			name: "Single instruction NOP",
			text: "NOP",
			want: []byte{0x00},
		},
		{
			name: "Single instruction OR $56",
			text: "OR $56",
			want: []byte{0xf6, 0x56},
		},
		{
			name: "Single instruction OR C",
			text: "OR C",
			want: []byte{0xb1},
		},
		{
			name: "Single instruction OR (HL)",
			text: "OR (HL)",
			want: []byte{0xb6},
		},
		{
			name: "Single instruction OR (IX+$12)",
			text: "OR (IX+$12)",
			want: []byte{0xdd, 0xb6, 0x12},
		},
		{
			name: "Single instruction OR (IX+$12)",
			text: "OR (IX+$12)",
			want: []byte{0xdd, 0xb6, 0x12},
		},
		{
			name: "Single instruction OR (IY-$12)",
			text: "OR (IY-$12)",
			want: []byte{0xfd, 0xb6, 0xee},
		},
		{
			name: "Single instruction OR (IY-$12)",
			text: "OR (IY-$12)",
			want: []byte{0xfd, 0xb6, 0xee},
		},
		{
			name: "Single instruction OTDR",
			text: "OTDR",
			want: []byte{0xed, 0xbb},
		},
		{
			name: "Single instruction OTIR",
			text: "OTIR",
			want: []byte{0xed, 0xb3},
		},
		{
			name: "Single instruction OUT (23), A",
			text: "OUT (23), A",
			want: []byte{0xd3, 0x17},
		},
		{
			name: "Single instruction OUT (C),D",
			text: "OUT (C),D",
			want: []byte{0xed, 0x51},
		},
		{
			name: "Single instruction OUTD",
			text: "OUTD",
			want: []byte{0xed, 0xab},
		},
		{
			name: "Single instruction OUTI",
			text: "OUTI",
			want: []byte{0xed, 0xa3},
		},
		{
			name: "Single instruction POP AF",
			text: "POP AF",
			want: []byte{0xf1},
		},
		{
			name: "Single instruction POP IX",
			text: "POP IX",
			want: []byte{0xdd, 0xe1},
		},
		{
			name: "Single instruction POP IY",
			text: "POP IY",
			want: []byte{0xfd, 0xe1},
		},
		{
			name: "Single instruction PUSH AF",
			text: "PUSH AF",
			want: []byte{0xf5},
		},
		{
			name: "Single instruction PUSH IX",
			text: "PUSH IX",
			want: []byte{0xdd, 0xe5},
		},
		{
			name: "Single instruction PUSH IY",
			text: "PUSH IY",
			want: []byte{0xfd, 0xe5},
		},
		{
			name: "Single instruction RES 4,(IY-$12)",
			text: "RES 4,(IY-$12)",
			want: []byte{0xfd, 0xcb, 0xee, 0xa6},
		},
		{
			name: "Single instruction RES 5,(IX+$12)",
			text: "RES 5,(IX+$12)",
			want: []byte{0xdd, 0xcb, 0x12, 0xae},
		},
		{
			name: "Single instruction RES 6,(HL)",
			text: "RES 6,(HL)",
			want: []byte{0xcb, 0xb6},
		},
		{
			name: "Single instruction RES 7,C",
			text: "RES 7,C",
			want: []byte{0xcb, 0xb9},
		},
		{
			name: "Single instruction RET",
			text: "RET",
			want: []byte{0xc9},
		},
		{
			name: "Single instruction RETI",
			text: "RETI",
			want: []byte{0xed, 0x4d},
		},
		{
			name: "Single instruction RETN",
			text: "RETN",
			want: []byte{0xed, 0x45},
		},
		{
			name: "Single instruction RET PE",
			text: "RET PE",
			want: []byte{0xe8},
		},
		{
			name: "Single instruction RLA",
			text: "RLA",
			want: []byte{0x17},
		},
		{
			name: "Single instruction RL C",
			text: "RL C",
			want: []byte{0xcb, 0x11},
		},
		{
			name: "Single instruction RLCA",
			text: "RLCA",
			want: []byte{0x07},
		},
		{
			name: "Single instruction RLC C",
			text: "RLC C",
			want: []byte{0xcb, 0x01},
		},
		{
			name: "Single instruction RLC (HL)",
			text: "RLC (HL)",
			want: []byte{0xcb, 0x06},
		},
		{
			name: "Single instruction RLC (IX+$12)",
			text: "RLC (IX+$12)",
			want: []byte{0xdd, 0xcb, 0x12, 0x06},
		},
		{
			name: "Single instruction RLC (IY-$12)",
			text: "RLC (IY-$12)",
			want: []byte{0xfd, 0xcb, 0xee, 0x06},
		},
		{
			name: "Single instruction RLD",
			text: "RLD",
			want: []byte{0xed, 0x6f},
		},
		{
			name: "Single instruction RL (HL)",
			text: "RL (HL)",
			want: []byte{0xcb, 0x16},
		},
		{
			name: "Single instruction RL (IX+$12)",
			text: "RL (IX+$12)",
			want: []byte{0xdd, 0xcb, 0x12, 0x16},
		},
		{
			name: "Single instruction RL (IY-$12)",
			text: "RL (IY-$12)",
			want: []byte{0xfd, 0xcb, 0xee, 0x16},
		},
		{
			name: "Single instruction RRA",
			text: "RRA",
			want: []byte{0x1f},
		},
		{
			name: "Single instruction RR C",
			text: "RR C",
			want: []byte{0xcb, 0x19},
		},
		{
			name: "Single instruction RRCA",
			text: "RRCA",
			want: []byte{0x0f},
		},
		{
			name: "Single instruction RRC C",
			text: "RRC C",
			want: []byte{0xcb, 0x09},
		},
		{
			name: "Single instruction RRC (HL)",
			text: "RRC (HL)",
			want: []byte{0xcb, 0x0e},
		},
		{
			name: "Single instruction RRC (IX+$12)",
			text: "RRC (IX+$12)",
			want: []byte{0xdd, 0xcb, 0x12, 0x0e},
		},
		{
			name: "Single instruction RRC (IY-$12)",
			text: "RRC (IY-$12)",
			want: []byte{0xfd, 0xcb, 0xee, 0x0e},
		},
		{
			name: "Single instruction RRD",
			text: "RRD",
			want: []byte{0xed, 0x67},
		},
		{
			name: "Single instruction RR (HL)",
			text: "RR (HL)",
			want: []byte{0xcb, 0x1e},
		},
		{
			name: "Single instruction RR (IX+$12)",
			text: "RR (IX+$12)",
			want: []byte{0xdd, 0xcb, 0x12, 0x1e},
		},
		{
			name: "Single instruction RR (IY-$12)",
			text: "RR (IY-$12)",
			want: []byte{0xfd, 0xcb, 0xee, 0x1e},
		},
		{
			name: "Single instruction RST $30",
			text: "RST $30",
			want: []byte{0xf7},
		},
		{
			name: "Single instruction SBC A,$56",
			text: "SBC A,$56",
			want: []byte{0xde, 0x56},
		},
		{
			name: "Single instruction SBC A,C",
			text: "SBC A,C",
			want: []byte{0x99},
		},
		{
			name: "Single instruction SBC A,(HL)",
			text: "SBC A,(HL)",
			want: []byte{0x9e},
		},
		{
			name: "Single instruction SBC A,(IX+$12)",
			text: "SBC A,(IX+$12)",
			want: []byte{0xdd, 0x9e, 0x12},
		},
		{
			name: "Single instruction SBC A,(IY-$12)",
			text: "SBC A,(IY-$12)",
			want: []byte{0xfd, 0x9e, 0xee},
		},
		{
			name: "Single instruction SBC HL,DE",
			text: "SBC HL,DE",
			want: []byte{0xed, 0x52},
		},
		{
			name: "Single instruction SCF",
			text: "SCF",
			want: []byte{0x37},
		},
		{
			name: "Single instruction SET 0,C",
			text: "SET 0,C",
			want: []byte{0xcb, 0xc1},
		},
		{
			name: "Single instruction SET 1,(HL)",
			text: "SET 1,(HL)",
			want: []byte{0xcb, 0xce},
		},
		{
			name: "Single instruction SET 2,(IX+$12)",
			text: "SET 2,(IX+$12)",
			want: []byte{0xdd, 0xcb, 0x12, 0xd6},
		},
		{
			name: "Single instruction SET 3,(IY-$12)",
			text: "SET 3,(IY-$12)",
			want: []byte{0xfd, 0xcb, 0xee, 0xde},
		},
		{
			name: "Single instruction SLA C",
			text: "SLA C",
			want: []byte{0xcb, 0x21},
		},
		{
			name: "Single instruction SLA (HL)",
			text: "SLA (HL)",
			want: []byte{0xcb, 0x26},
		},
		{
			name: "Single instruction SLA (IX+$12)",
			text: "SLA (IX+$12)",
			want: []byte{0xdd, 0xcb, 0x12, 0x26},
		},
		{
			name: "Single instruction SLA (IY-$12)",
			text: "SLA (IY-$12)",
			want: []byte{0xfd, 0xcb, 0xee, 0x26},
		},
		{
			name: "Single instruction SRA C",
			text: "SRA C",
			want: []byte{0xcb, 0x29},
		},
		{
			name: "Single instruction SRA (HL)",
			text: "SRA (HL)",
			want: []byte{0xcb, 0x2e},
		},
		{
			name: "Single instruction SRA (IX+$12)",
			text: "SRA (IX+$12)",
			want: []byte{0xdd, 0xcb, 0x12, 0x2e},
		},
		{
			name: "Single instruction SRA (IY-$12)",
			text: "SRA (IY-$12)",
			want: []byte{0xfd, 0xcb, 0xee, 0x2e},
		},
		{
			name: "Single instruction SRL C",
			text: "SRL C",
			want: []byte{0xcb, 0x39},
		},
		{
			name: "Single instruction SRL (HL)",
			text: "SRL (HL)",
			want: []byte{0xcb, 0x3e},
		},
		{
			name: "Single instruction SRL (IX+$12)",
			text: "SRL (IX+$12)",
			want: []byte{0xdd, 0xcb, 0x12, 0x3e},
		},
		{
			name: "Single instruction SRL (IY-$12)",
			text: "SRL (IY-$12)",
			want: []byte{0xfd, 0xcb, 0xee, 0x3e},
		},
		{
			name: "Single instruction SUB $56",
			text: "SUB $56",
			want: []byte{0xd6, 0x56},
		},
		{
			name: "Single instruction SUB C",
			text: "SUB C",
			want: []byte{0x91},
		},
		{
			name: "Single instruction SUB (HL)",
			text: "SUB (HL)",
			want: []byte{0x96},
		},
		{
			name: "Single instruction SUB (IX+$12)",
			text: "SUB (IX+$12)",
			want: []byte{0xdd, 0x96, 0x12},
		},
		{
			name: "Single instruction SUB (IY-$12)",
			text: "SUB (IY-$12)",
			want: []byte{0xfd, 0x96, 0xee},
		},
		{
			name: "Single instruction XOR $56",
			text: "XOR $56",
			want: []byte{0xee, 0x56},
		},
		{
			name: "Single instruction XOR C",
			text: "XOR C",
			want: []byte{0xa9},
		},
		{
			name: "Single instruction XOR (HL)",
			text: "XOR (HL)",
			want: []byte{0xae},
		},
		{
			name: "Single instruction XOR (IX+$12)",
			text: "XOR (IX+$12)",
			want: []byte{0xdd, 0xae, 0x12},
		},
		{
			name: "Single instruction XOR (IY-$12)",
			text: "XOR (IY-$12)",
			want: []byte{0xfd, 0xae, 0xee},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			src := " .org 0\n " + test.text
			assembler := New([]string{}, "z80", "c128", []string{})
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
