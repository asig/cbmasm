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

type AddressingMode int

const (
	AM_Implied AddressingMode = 1 << iota
	AM_Immediate
	AM_Accumulator
	AM_ZeroPage         // $aa
	AM_ZeroPageIndexedX // $aa,X
	AM_ZeroPageIndexedY // $aa,Y
	AM_Absolute         // $aaaa
	AM_AbsoluteIndirect // ($aaaa)
	AM_AbsoluteIndexedX // $aaaa,X
	AM_AbsoluteIndexedY // $aaaa,Y
	AM_IndexedIndirect  // ($aa,X)
	AM_IndirectIndexed  // ($aa),Y
	AM_Relative         // $aa
)

func (am AddressingMode) withIndex(register string) AddressingMode {
	switch am {
	case AM_Absolute:
		switch register {
		case "X", "x":
			return AM_AbsoluteIndexedX
		case "Y", "y":
			return AM_AbsoluteIndexedY
		}
	case AM_ZeroPage:
		switch register {
		case "X", "x":
			return AM_ZeroPageIndexedX
		case "Y", "y":
			return AM_ZeroPageIndexedY
		}
	}
	return am
}

func (am AddressingMode) withSize(size int) AddressingMode {
	switch am {
	case AM_ZeroPage:
		if size == 2 {
			return AM_Absolute
		}
	case AM_ZeroPageIndexedX:
		if size == 2 {
			return AM_AbsoluteIndexedX
		}
	case AM_ZeroPageIndexedY:
		if size == 2 {
			return AM_AbsoluteIndexedY
		}
	case AM_Absolute:
		if size == 1 {
			return AM_ZeroPage
		}
	case AM_AbsoluteIndexedX:
		if size == 1 {
			return AM_ZeroPageIndexedX
		}
	case AM_AbsoluteIndexedY:
		if size == 1 {
			return AM_ZeroPageIndexedY
		}
	}
	return am
}

type OpCodes map[AddressingMode]byte

var Mnemonics = map[string]OpCodes{
	"adc": {
		AM_Immediate:        0x69,
		AM_ZeroPage:         0x65,
		AM_ZeroPageIndexedX: 0x75,
		AM_Absolute:         0x6d,
		AM_AbsoluteIndexedX: 0x7d,
		AM_AbsoluteIndexedY: 0x79,
		AM_IndexedIndirect:  0x61,
		AM_IndirectIndexed:  0x71,
	},
	"and": {
		AM_Immediate:        0x29,
		AM_ZeroPage:         0x25,
		AM_ZeroPageIndexedX: 0x35,
		AM_Absolute:         0x2d,
		AM_AbsoluteIndexedX: 0x3d,
		AM_AbsoluteIndexedY: 0x39,
		AM_IndexedIndirect:  0x21,
		AM_IndirectIndexed:  0x31,
	},
	"asl": {
		AM_Implied:          0x0a,
		AM_Accumulator:      0x0a,
		AM_ZeroPage:         0x06,
		AM_ZeroPageIndexedX: 0x16,
		AM_Absolute:         0x0e,
		AM_AbsoluteIndexedX: 0x1e,
	},
	"bit": {
		AM_ZeroPage: 0x24,
		AM_Absolute: 0x2c,
	},
	"bpl": {
		AM_Relative: 0x10,
	},
	"bmi": {
		AM_Relative: 0x30,
	},
	"bvc": {
		AM_Relative: 0x50,
	},
	"bvs": {
		AM_Relative: 0x70,
	},
	"bcc": {
		AM_Relative: 0x90,
	},
	"bcs": {
		AM_Relative: 0xb0,
	},
	"bne": {
		AM_Relative: 0xd0,
	},
	"beq": {
		AM_Relative: 0xf0,
	},
	"brk": {
		AM_Implied: 0x00,
	},
	"cmp": {
		AM_Immediate:        0xc9,
		AM_ZeroPage:         0xc5,
		AM_ZeroPageIndexedX: 0xd5,
		AM_Absolute:         0xcd,
		AM_AbsoluteIndexedX: 0xdd,
		AM_AbsoluteIndexedY: 0xd9,
		AM_IndexedIndirect:  0xc1,
		AM_IndirectIndexed:  0xd1,
	},
	"cpx": {
		AM_Immediate: 0xe0,
		AM_ZeroPage:  0xe4,
		AM_Absolute:  0xec,
	},
	"cpy": {
		AM_Immediate: 0xc0,
		AM_ZeroPage:  0xc4,
		AM_Absolute:  0xcc,
	},
	"dec": {
		AM_ZeroPage:         0xc6,
		AM_ZeroPageIndexedX: 0xd6,
		AM_Absolute:         0xce,
		AM_AbsoluteIndexedX: 0xde,
	},
	"eor": {
		AM_Immediate:        0x49,
		AM_ZeroPage:         0x45,
		AM_ZeroPageIndexedX: 0x55,
		AM_Absolute:         0x4d,
		AM_AbsoluteIndexedX: 0x5d,
		AM_AbsoluteIndexedY: 0x59,
		AM_IndexedIndirect:  0x41,
		AM_IndirectIndexed:  0x51,
	},
	"clc": {
		AM_Implied: 0x18,
	},
	"sec": {
		AM_Implied: 0x38,
	},
	"cli": {
		AM_Implied: 0x58,
	},
	"sei": {
		AM_Implied: 0x78,
	},
	"clv": {
		AM_Implied: 0xb8,
	},
	"cld": {
		AM_Implied: 0xd8,
	},
	"sed": {
		AM_Implied: 0xf8,
	},
	"inc": {
		AM_ZeroPage:         0xe6,
		AM_ZeroPageIndexedX: 0xf6,
		AM_Absolute:         0xee,
		AM_AbsoluteIndexedX: 0xfe,
	},
	"jmp": {
		AM_Absolute:         0x4c,
		AM_AbsoluteIndirect: 0x6c,
	},
	"jsr": {
		AM_Absolute: 0x20,
	},
	"lda": {
		AM_Immediate:        0xa9,
		AM_ZeroPage:         0xa5,
		AM_ZeroPageIndexedX: 0xb5,
		AM_Absolute:         0xad,
		AM_AbsoluteIndexedX: 0xbd,
		AM_AbsoluteIndexedY: 0xb9,
		AM_IndexedIndirect:  0xa1,
		AM_IndirectIndexed:  0xb1,
	},
	"ldx": {
		AM_Immediate:        0xa2,
		AM_ZeroPage:         0xa6,
		AM_ZeroPageIndexedY: 0xb6,
		AM_Absolute:         0xae,
		AM_AbsoluteIndexedY: 0xbe,
	},
	"ldy": {
		AM_Immediate:        0xa0,
		AM_ZeroPage:         0xa4,
		AM_ZeroPageIndexedX: 0xb4,
		AM_Absolute:         0xac,
		AM_AbsoluteIndexedX: 0xbc,
	},
	"lsr": {
		AM_Implied:          0x4a,
		AM_Accumulator:      0x4a,
		AM_ZeroPage:         0x46,
		AM_ZeroPageIndexedX: 0x56,
		AM_Absolute:         0x4e,
		AM_AbsoluteIndexedX: 0x5e,
	},
	"nop": {
		AM_Implied: 0xea,
	},
	"ora": {
		AM_Immediate:        0x09,
		AM_ZeroPage:         0x05,
		AM_ZeroPageIndexedX: 0x15,
		AM_Absolute:         0x0d,
		AM_AbsoluteIndexedX: 0x1d,
		AM_AbsoluteIndexedY: 0x19,
		AM_IndexedIndirect:  0x01,
		AM_IndirectIndexed:  0x11,
	},
	"tax": {
		AM_Implied: 0xaa,
	},
	"txa": {
		AM_Implied: 0x8a,
	},
	"dex": {
		AM_Implied: 0xca,
	},
	"inx": {
		AM_Implied: 0xe8,
	},
	"tay": {
		AM_Implied: 0xa8,
	},
	"tya": {
		AM_Implied: 0x98,
	},
	"dey": {
		AM_Implied: 0x88,
	},
	"iny": {
		AM_Implied: 0xc8,
	},
	"rol": {
		AM_Implied:          0x2a,
		AM_Accumulator:      0x2a,
		AM_ZeroPage:         0x26,
		AM_ZeroPageIndexedX: 0x36,
		AM_Absolute:         0x2e,
		AM_AbsoluteIndexedX: 0x3e,
	},
	"ror": {
		AM_Implied:          0x6a,
		AM_Accumulator:      0x6a,
		AM_ZeroPage:         0x66,
		AM_ZeroPageIndexedX: 0x76,
		AM_Absolute:         0x6e,
		AM_AbsoluteIndexedX: 0x7e,
	},
	"rti": {
		AM_Implied: 0x40,
	},
	"rts": {
		AM_Implied: 0x60,
	},
	"sbc": {
		AM_Immediate:        0xe9,
		AM_ZeroPage:         0xe5,
		AM_ZeroPageIndexedX: 0xf5,
		AM_Absolute:         0xed,
		AM_AbsoluteIndexedX: 0xfd,
		AM_AbsoluteIndexedY: 0xf9,
		AM_IndexedIndirect:  0xe1,
		AM_IndirectIndexed:  0xf1,
	},
	"sta": {
		AM_ZeroPage:         0x85,
		AM_ZeroPageIndexedX: 0x95,
		AM_Absolute:         0x8d,
		AM_AbsoluteIndexedX: 0x9d,
		AM_AbsoluteIndexedY: 0x99,
		AM_IndexedIndirect:  0x81,
		AM_IndirectIndexed:  0x91,
	},
	"txs": {
		AM_Implied: 0x9a,
	},
	"tsx": {
		AM_Implied: 0xba,
	},
	"pha": {
		AM_Implied: 0x48,
	},
	"pla": {
		AM_Implied: 0x68,
	},
	"php": {
		AM_Implied: 0x08,
	},
	"plp": {
		AM_Implied: 0x28,
	},
	"stx": {
		AM_ZeroPage:         0x86,
		AM_ZeroPageIndexedY: 0x96,
		AM_Absolute:         0x8e,
	},
	"sty": {
		AM_ZeroPage:         0x84,
		AM_ZeroPageIndexedX: 0x94,
		AM_Absolute:         0x8c,
	},
}
