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

type AddressingMode int

const (
	AM_Register AddressingMode = 1 << iota // A, B, C, ...
	AM_RegisterIndirect                    // (HL), (BC), (DE)
	AM_Indexed                             // (IX + d), (IY + d)
	AM_ExtAddressing                       // (addr)
	AM_Implied                             // I, R
	AM_Immediate                           // nn
)


type ParamPattern struct {
	am AddressingMode
	regs Register       // Valid registers for AM_Register, AM_RegisterIndirect, AM_Implied
}

type Param struct {
	pos text.Pos
	am AddressingMode
	val expr.Node     // Offset if am == AM_Indexed, value if am == AM_Immediate
	r Register       // Valid registers for AM_Register, AM_RegisterIndirect, AM_Implied
}

type CodeGen func(p []Param, errorSink errors.Sink) []expr.Node

type OpCodeEntry struct {
	p1, p2 *ParamPattern
	c CodeGen
}

type OpCodeEntryList []OpCodeEntry

func (l OpCodeEntryList) FindMatch(p []Param) CodeGen {
	return nil
}

var Mnemonics = map[string]OpCodeEntryList{
	"ld": {},
}
