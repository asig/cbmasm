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

var Mnemonics = map[string]OpCodeEntryList{
	"ld": {},
}
