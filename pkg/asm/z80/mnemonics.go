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
package z80

import (
	"strings"

	"github.com/asig/cbmasm/pkg/errors"
	"github.com/asig/cbmasm/pkg/expr"
	"github.com/asig/cbmasm/pkg/text"
)

type Register int

const (
	Reg_A Register = 1 << iota
	Reg_B
	Reg_C
	Reg_D
	Reg_E
	Reg_H
	Reg_L
	Reg_AF
	Reg_BC
	Reg_DE
	Reg_HL
	Reg_SP
	Reg_I
	Reg_R
	Reg_IX
	Reg_IY
)

var (
	drVal = map[Register]int{
		Reg_BC: 0,
		Reg_DE: 1,
		Reg_HL: 2,
		Reg_SP: 3,
	}
	rVal = map[Register]int{
		Reg_B: 0b000,
		Reg_C: 0b001,
		Reg_D: 0b010,
		Reg_E: 0b011,
		Reg_H: 0b100,
		Reg_L: 0b101,
		Reg_A: 0b111,
	}
)

func (r Register) IsDouble() bool {
	switch r {
	case Reg_AF, Reg_BC, Reg_DE, Reg_HL, Reg_SP:
		return true
	}
	return false
}

var (
	stringToReg = map[string]Register{
		"a":  Reg_A,
		"b":  Reg_B,
		"c":  Reg_C,
		"d":  Reg_D,
		"e":  Reg_E,
		"h":  Reg_H,
		"l":  Reg_L,
		"af": Reg_AF,
		"bc": Reg_BC,
		"de": Reg_DE,
		"hl": Reg_HL,
		"sp": Reg_SP,
		"i":  Reg_I,
		"r":  Reg_R,
		"ix": Reg_IX,
		"iy": Reg_IY,
	}
)

func RegisterFromString(s string) (Register, bool) {
	reg, found := stringToReg[strings.ToLower(s)]
	return reg, found
}

type Cond int

const (
	Cond_NZ Cond = 1 << iota
	Cond_Z
	Cond_NC
	Cond_C
	Cond_PO
	Cond_PE
	Cond_P
	Cond_M
)

var (
	stringToCond = map[string]Cond{
		"nz": Cond_NZ,
		"z":  Cond_Z,
		"nc": Cond_NC,
		"c":  Cond_C,
		"po": Cond_PO,
		"pe": Cond_PE,
		"p":  Cond_P,
		"m":  Cond_M,
	}
)

func CondFromString(s string) (Cond, bool) {
	reg, found := stringToCond[strings.ToLower(s)]
	return reg, found
}

type AddressingMode int

const (
	AM_Register         AddressingMode = 1 << iota // A, B, C, ...
	AM_RegisterIndirect                            // (HL), (BC), (DE)
	AM_Indexed                                     // (IX + d), (IY + d)
	AM_ExtAddressing                               // (addr)
	AM_Implied                                     // I, R
	AM_Immediate                                   // nn
	AM_Cond                                        // Not really an addressing mode. Indicates a condition
)

type ParamPattern struct {
	mode  AddressingMode
	regs  Register // Valid registers for AM_Register, AM_RegisterIndirect, AM_Implied
	conds Cond     // Valid conditions for AM_Cond
}

type Param struct {
	Pos  text.Pos
	Mode AddressingMode
	Val  expr.Node // Offset if am == AM_Indexed, value if am == AM_Immediate
	R    Register  // Register for AM_Register, AM_RegisterIndirect, AM_Indexed, AM_Implied
	Cond Cond      // Condition for AM_Cond
}

func (p Param) Matches(pattern ParamPattern) bool {
	if pattern.mode != p.Mode {
		return false
	}
	switch pattern.mode {
	case AM_Register, AM_RegisterIndirect, AM_Indexed, AM_Implied:
		// Register needs to be in pattern mask
		return pattern.regs&p.R != 0
	case AM_Cond:
		// Condition needs to be in pattern mask
		return pattern.conds&p.Cond != 0
	default:
		return true
	}
}

type CodeGen func(p []Param, errorSink errors.Sink) []expr.Node

type OpCodeEntry struct {
	p []ParamPattern
	c CodeGen
}

type OpCodeEntryList []OpCodeEntry

func (l OpCodeEntryList) FindMatch(p []Param) CodeGen {
	for _, entry := range l {
		if len(entry.p) != len(p) {
			continue
		}

		match := true
		for i := range entry.p {
			if !p[i].Matches(entry.p[i]) {
				match = false
			}
		}
		if match {
			return entry.c
		}
	}
	return nil
}

func c(v int) expr.Node {
	return expr.NewConst(text.Pos{}, v, 1)
}

func loByte(p Param) expr.Node {
	return expr.NewUnaryOp(p.Pos, p.Val, expr.LoByte)
}

func hiByte(p Param) expr.Node {
	return expr.NewUnaryOp(p.Pos, p.Val, expr.HiByte)
}

func bytes(bs ...int) OpCodeEntryList {
	var nodes []expr.Node
	for _, b := range bs {
		nodes = append(nodes, c(b))
	}
	return OpCodeEntryList{
		OpCodeEntry{
			[]ParamPattern{},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return nodes
			},
		},
	}
}

var Mnemonics = map[string]OpCodeEntryList{
	"ld": {
		OpCodeEntry{ // LD dd, (nn)
			[]ParamPattern{{mode: AM_Register, regs: Reg_BC | Reg_DE | Reg_SP}, {mode: AM_ExtAddressing}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{
					c(0xed), c(0b01001011 | drVal[p[0].R]<<4), loByte(p[1]), hiByte(p[1]),
				}
			},
		},
		OpCodeEntry{ // LD dd, nn
			[]ParamPattern{{mode: AM_Register, regs: Reg_BC | Reg_DE | Reg_HL | Reg_SP}, {mode: AM_Immediate}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{
					c(0b00000011 | drVal[p[0].R]<<4), loByte(p[1]), hiByte(p[1]),
				}
			},
		},
		OpCodeEntry{ // LD r, n
			[]ParamPattern{{mode: AM_Register, regs: Reg_B | Reg_C | Reg_D | Reg_E | Reg_H | Reg_L | Reg_A}, {mode: AM_Immediate}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{
					c(0b00000110 | rVal[p[0].R]<<3), p[1].Val,
				}
			},
		},
		OpCodeEntry{ // LD r, r'
			[]ParamPattern{{mode: AM_Register, regs: Reg_B | Reg_C | Reg_D | Reg_E | Reg_H | Reg_L | Reg_A}, {mode: AM_Register, regs: Reg_B | Reg_C | Reg_D | Reg_E | Reg_H | Reg_L | Reg_A}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{
					c(0b01000000 | rVal[p[0].R]<<3 | rVal[p[1].R]),
				}
			},
		},
		OpCodeEntry{ // LD (BC),A
			[]ParamPattern{{mode: AM_RegisterIndirect, regs: Reg_BC}, {mode: AM_Register, regs: Reg_A}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{c(0x02)}
			},
		},
		OpCodeEntry{ // LD (DE),A
			[]ParamPattern{{mode: AM_RegisterIndirect, regs: Reg_DE}, {mode: AM_Register, regs: Reg_A}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{c(0x12)}
			},
		},
		OpCodeEntry{ // LD (HL), n
			[]ParamPattern{{mode: AM_RegisterIndirect, regs: Reg_HL}, {mode: AM_Immediate}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{c(0x36), p[1].Val}
			},
		},
		OpCodeEntry{ // LD (HL), r
			[]ParamPattern{{mode: AM_RegisterIndirect, regs: Reg_HL}, {mode: AM_Register, regs: Reg_B | Reg_C | Reg_D | Reg_E | Reg_H | Reg_L | Reg_A}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{c(0b01110000 | rVal[p[1].R])}
			},
		},
		OpCodeEntry{ // LD r, (IX + d)
			[]ParamPattern{{mode: AM_Register, regs: Reg_B | Reg_C | Reg_D | Reg_E | Reg_H | Reg_L | Reg_A}, {mode: AM_Indexed, regs: Reg_IX}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{c(0xdd), c(0b01000110 | rVal[p[0].R<<3]), p[1].Val}
			},
		},
		OpCodeEntry{ // LD r, (IY + d)
			[]ParamPattern{{mode: AM_Register, regs: Reg_B | Reg_C | Reg_D | Reg_E | Reg_H | Reg_L | Reg_A}, {mode: AM_Indexed, regs: Reg_IY}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{c(0xfd), c(0b01000110 | rVal[p[0].R<<3]), p[1].Val}
			},
		},
		OpCodeEntry{ // LD (IX + d), n
			[]ParamPattern{{mode: AM_Indexed, regs: Reg_IX}, {mode: AM_Immediate}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{c(0xdd), c(0x36), p[0].Val, p[1].Val}
			},
		},
		OpCodeEntry{ // LD (IY + d), n
			[]ParamPattern{{mode: AM_Indexed, regs: Reg_IY}, {mode: AM_Immediate}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{c(0xfd), c(0x36), p[0].Val, p[1].Val}
			},
		},
		OpCodeEntry{ // LD (IX + d), r
			[]ParamPattern{{mode: AM_Indexed, regs: Reg_IX}, {mode: AM_Register, regs: Reg_B | Reg_C | Reg_D | Reg_E | Reg_H | Reg_L | Reg_A}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{c(0xdd), c(0b01110000 | rVal[p[1].R]), p[0].Val}
			},
		},
		OpCodeEntry{ // LD (IY + d), r
			[]ParamPattern{{mode: AM_Indexed, regs: Reg_IY}, {mode: AM_Register, regs: Reg_B | Reg_C | Reg_D | Reg_E | Reg_H | Reg_L | Reg_A}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{c(0xfd), c(0b01110000 | rVal[p[1].R]), p[0].Val}
			},
		},
		OpCodeEntry{ // LD A, (nn)
			[]ParamPattern{{mode: AM_Register, regs: Reg_A}, {mode: AM_ExtAddressing}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{
					c(0x3a), loByte(p[1]), hiByte(p[1]),
				}
			},
		},
		OpCodeEntry{ // LD (nn), A
			[]ParamPattern{{mode: AM_ExtAddressing}, {mode: AM_Register, regs: Reg_A}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{
					c(0x32), loByte(p[0]), hiByte(p[0]),
				}
			},
		},
		OpCodeEntry{ // LD (nn), dd
			[]ParamPattern{{mode: AM_ExtAddressing}, {mode: AM_Register, regs: Reg_BC | Reg_DE | Reg_SP}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{
					c(0xed), c(0b01000011 | drVal[p[1].R]<<4), loByte(p[0]), hiByte(p[0]),
				}
			},
		},
		OpCodeEntry{ // LD (nn), HL
			[]ParamPattern{{mode: AM_ExtAddressing}, {mode: AM_Register, regs: Reg_HL}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{
					c(0x22), loByte(p[0]), hiByte(p[0]),
				}
			},
		},
		OpCodeEntry{ // LD (nn), IX
			[]ParamPattern{{mode: AM_ExtAddressing}, {mode: AM_Register, regs: Reg_IX}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{
					c(0xdd), c(0x22), loByte(p[0]), hiByte(p[0]),
				}
			},
		},
		OpCodeEntry{ // LD (nn), IY
			[]ParamPattern{{mode: AM_ExtAddressing}, {mode: AM_Register, regs: Reg_IY}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{
					c(0xfd), c(0x22), loByte(p[0]), hiByte(p[0]),
				}
			},
		},
		OpCodeEntry{ // LD A, (BC)
			[]ParamPattern{{mode: AM_Register, regs: Reg_A}, {mode: AM_RegisterIndirect, regs: Reg_BC}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{
					c(0x0a),
				}
			},
		},
		OpCodeEntry{ // LD A, (DE)
			[]ParamPattern{{mode: AM_Register, regs: Reg_A}, {mode: AM_RegisterIndirect, regs: Reg_DE}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{
					c(0x1a),
				}
			},
		},
		OpCodeEntry{ // LD A, I
			[]ParamPattern{{mode: AM_Register, regs: Reg_A}, {mode: AM_Register, regs: Reg_I}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{
					c(0xed), c(0x57),
				}
			},
		},
		OpCodeEntry{ // LD I, A
			[]ParamPattern{{mode: AM_Register, regs: Reg_I}, {mode: AM_Register, regs: Reg_A}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{
					c(0xed), c(0x47),
				}
			},
		},
		OpCodeEntry{ // LD A, R
			[]ParamPattern{{mode: AM_Register, regs: Reg_A}, {mode: AM_Register, regs: Reg_R}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{
					c(0xed), c(0x5f),
				}
			},
		},
		OpCodeEntry{ // LD HL, (nn)
			[]ParamPattern{{mode: AM_Register, regs: Reg_HL}, {mode: AM_ExtAddressing}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{
					c(0x2a), loByte(p[1]), hiByte(p[1]),
				}
			},
		},
		OpCodeEntry{ // LD IX, nn
			[]ParamPattern{{mode: AM_Register, regs: Reg_IX}, {mode: AM_Immediate}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{
					c(0xdd), c(0x21), loByte(p[1]), hiByte(p[1]),
				}
			},
		},
		OpCodeEntry{ // LD IX, (nn)
			[]ParamPattern{{mode: AM_Register, regs: Reg_IX}, {mode: AM_ExtAddressing}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{
					c(0xdd), c(0x2a), loByte(p[1]), hiByte(p[1]),
				}
			},
		},
		OpCodeEntry{ // LD IY, nn
			[]ParamPattern{{mode: AM_Register, regs: Reg_IY}, {mode: AM_Immediate}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{
					c(0xfd), c(0x21), loByte(p[1]), hiByte(p[1]),
				}
			},
		},
		OpCodeEntry{ // LD IY, (nn)
			[]ParamPattern{{mode: AM_Register, regs: Reg_IY}, {mode: AM_ExtAddressing}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{
					c(0xfd), c(0x2a), loByte(p[1]), hiByte(p[1]),
				}
			},
		},
		OpCodeEntry{ // LD R, A
			[]ParamPattern{{mode: AM_Register, regs: Reg_R}, {mode: AM_Register, regs: Reg_A}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{
					c(0xed), c(0x4f),
				}
			},
		},
		OpCodeEntry{ // LD SP, HL
			[]ParamPattern{{mode: AM_Register, regs: Reg_SP}, {mode: AM_Register, regs: Reg_HL}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{
					c(0xf9),
				}
			},
		},
		OpCodeEntry{ // LD SP, IX
			[]ParamPattern{{mode: AM_Register, regs: Reg_SP}, {mode: AM_Register, regs: Reg_IX}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{
					c(0xdd), c(0xf9),
				}
			},
		},
		OpCodeEntry{ // LD SP, IY
			[]ParamPattern{{mode: AM_Register, regs: Reg_SP}, {mode: AM_Register, regs: Reg_IY}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{
					c(0xfd), c(0xf9),
				}
			},
		},
		OpCodeEntry{ // LD r, (HL)
			[]ParamPattern{{mode: AM_Register, regs: Reg_B | Reg_C | Reg_D | Reg_E | Reg_H | Reg_L | Reg_A}, {mode: AM_RegisterIndirect, regs: Reg_HL}},
			func(p []Param, errorSink errors.Sink) []expr.Node {
				return []expr.Node{c(0b01000110 | rVal[p[0].R])}
			},
		},
	},
	"ldd":  bytes(0xed, 0xa2),
	"lddr": bytes(0xed, 0xb8),
	"ldi":  bytes(0xed, 0xa0),
	"ldir": bytes(0xed, 0xb0),
	"neg":  bytes(0xed, 0x44),
	"nop":  bytes(0x00),
}
