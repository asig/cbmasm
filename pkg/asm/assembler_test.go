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
	"fmt"
	"strings"
	"testing"

	"github.com/asig/cbmasm/pkg/errors"
	"github.com/asig/cbmasm/pkg/text"
)

func TestAssembler_RelativeBranchesAreCheckedForOverflow(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		wantErrors   []errors.Error
		wantWarnings []errors.Error
	}{
		{
			name: "backward branch out of bounds",
			text: `   .org 0
l	.reserve 128
	beq l
`,
			wantErrors:   []errors.Error{{text.Pos{Filename: "", Line: 3, Col: 6}, "Branch target too far away."}},
			wantWarnings: []errors.Error{},
		},
		{
			name: "forward branch out of bounds",
			text: `   .org 0
	beq l
	.reserve 128
l:
`,
			wantErrors:   []errors.Error{{text.Pos{Filename: "", Line: 2, Col: 6}, "Branch target too far away."}},
			wantWarnings: []errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assembler := New([]string{}, "6502", "c128", "plain", "petscii", []string{})
			assembler.Assemble(text.Process("", test.text))
			errs := assembler.Errors()
			if len(errs) != len(test.wantErrors) {
				t.Errorf("Got %d, want %d errs", len(errs), len(test.wantErrors))
			}
			for i := range errs {
				got := errs[i]
				want := test.wantErrors[i]
				if got != want {
					t.Errorf("Error %d: got %+v, want %+v", i+1, got, want)
				}
			}
			warnings := assembler.Warnings()
			if len(warnings) != len(test.wantWarnings) {
				t.Errorf("Got %d, want %d warnings", len(errs), len(test.wantWarnings))
			}
			for i := range warnings {
				got := warnings[i]
				want := test.wantWarnings[i]
				if got != want {
					t.Errorf("Warning %d: got %+v, want %+v", i+1, got, want)
				}
			}
		})
	}
}

func TestAssembler_BadFloatConst(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		wantErrors   []errors.Error
		wantWarnings []errors.Error
	}{
		{
			name: "bad float const",
			text: `   .org 0
	.float "foobar"
`,
			wantErrors:   []errors.Error{{text.Pos{Filename: "", Line: 2, Col: 9}, "Strings are not allowed"}},
			wantWarnings: []errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assembler := New([]string{}, "6502", "c128", "plain", "petscii", []string{})
			assembler.Assemble(text.Process("", test.text))
			errs := assembler.Errors()
			if len(errs) != len(test.wantErrors) {
				t.Errorf("Got %d, want %d errs", len(errs), len(test.wantErrors))
			}
			for i := range errs {
				got := errs[i]
				want := test.wantErrors[i]
				if got != want {
					t.Errorf("Error %d: got %+v, want %+v", i+1, got, want)
				}
			}
			warnings := assembler.Warnings()
			if len(warnings) != len(test.wantWarnings) {
				t.Errorf("Got %d, want %d warnings", len(errs), len(test.wantWarnings))
			}
			for i := range warnings {
				got := warnings[i]
				want := test.wantWarnings[i]
				if got != want {
					t.Errorf("Warning %d: got %+v, want %+v", i+1, got, want)
				}
			}
		})
	}
}

func TestAssembler_assemble(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		wantErrors   []errors.Error
		wantWarnings []errors.Error
		want         []byte
	}{
		{
			name: "Recursive symbol definition",
			text: `   .org 0
t1 .equ t2
t2 .equ t3
t3 .equ t4
t4 .equ t5
   jmp t4
t5 .equ $1234
`,
			want: []byte{0x4c, 0x34, 0x12},
		},

		{
			name: "mixed symbols and labels",
			text: `   .org 0
sym .equ label
   nop
label inx
   jmp sym
`,
			want: []byte{0xea, 0xe8, 0x4c, 0x01, 0x00},
		},

		{
			name: "local labels",
			text: `   .org 0
l1    jmp _l1
     nop
_l1   lda #0
     nop
     jmp _l1
l2    nop
     jmp _l1
     nop
_l1   brk
`,
			want: []byte{0x4c, 0x04, 0x00, 0xea, 0xa9, 0x00, 0xea, 0x4c, 0x04, 0x00, 0xea, 0x4c, 0x0f, 0x00, 0xea, 0x00},
		},

		{
			name: "labels - unresolved",
			text: `   .org 0
     jmp l
     jmp l
     nop
`,
			want: []byte{},
			wantErrors: []errors.Error{
				{text.Pos{Filename: "", Line: 2, Col: 10}, "Undefined label \"l\""},
				{text.Pos{Filename: "", Line: 3, Col: 10}, "Undefined label \"l\""},
			},
		},

		// TODO add test for unresolved local labels

		{
			name: "screen codes",
			text: `   .org 0
	.byte scr("hello")
	.byte scr("h"), scr("e"), scr("l"), scr("l"), scr("o")
	.byte scr('h'), scr('e'), scr('l'), scr('l'), scr('o')
`,
			want: []byte{
				0x08, 0x05, 0x0c, 0x0c, 0x0f,
				0x08, 0x05, 0x0c, 0x0c, 0x0f,
				0x08, 0x05, 0x0c, 0x0c, 0x0f,
			},
		},

		{
			name: "byte constants",
			text: `   .org 0
	.byte $01,$02,$03,$04
`,
			want: []byte{0x01, 0x02, 0x03, 0x04},
		},

		{
			name: "word constants, 2 bytes",
			text: `   .org 0
	.word $0102,$0304
`,
			want: []byte{0x02, 0x01, 0x04, 0x03},
		},

		{
			name: "word constants, 1 bytes",
			text: `   .org 0
	.word $01, 02
`,
			want: []byte{0x01, 0x00, 0x02, 0x00},
		},

		{
			name: "float constants",
			text: `   .org 0
	.float 2
	.float 0.0
	.float .25
	.float .26
	.float .27
	.float .5
	.float -13.2681
	.float 2+0.1
`,
			want: []byte{
				0x82, 0x00, 0x00, 0x00, 0x00, // 2.0
				0x00, 0x00, 0x00, 0x00, 0x00, // 0.0
				0x7f, 0x00, 0x00, 0x00, 0x00, // 0.25
				0x7f, 0x05, 0x1e, 0xb8, 0x51, // 0.26; a C128 actually uses 0x52 as the last byte, but with 0x51 it's also printed as ".26"...
				0x7f, 0x0a, 0x3d, 0x70, 0xa3, // 0.27; a C128 actually uses 0xa4 as the last byte, but with 0xa3 it's also printed as ".27"...
				0x80, 0x00, 0x00, 0x00, 0x00, // 0.5
				0x84, 0xd4, 0x4a, 0x23, 0x39, // -13.2681; a C128 actually uses 0x3a as the last byte...
				0x82, 0x06, 0x66, 0x66, 0x66, // 2.1
			},
		},
		{
			name: "conditional assembly - ifdef",
			text: ` .org 0
foo .equ 1
	.ifdef foo
	.byte $01,$02,$03,$04
	.else
	.byte $05,$06,$07,$08
	.endif
	.ifdef bar
	.byte $09,$0a,$0b,$0c
	.else
	.byte $0d,$0e,$0f,$10
	.endif
`,
			want: []byte{0x01, 0x02, 0x03, 0x04, 0x0d, 0x0e, 0x0f, 0x10},
		},

		{
			name: "conditional assembly - ifndef",
			text: ` .org 0
foo .equ 1
	.ifndef foo
	.byte $01,$02,$03,$04
	.else
	.byte $05,$06,$07,$08
	.endif
	.ifndef bar
	.byte $09,$0a,$0b,$0c
	.else
	.byte $0d,$0e,$0f,$10
	.endif
`,
			want: []byte{0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c},
		},

		{
			name: "conditional assembly - if",
			text: ` .org 0
foo .equ 1
	.if foo = 1
	.byte $01,$02,$03,$04
	.else
	.byte $05,$06,$07,$08
	.endif
	.if foo > 0
	.byte $09,$0a,$0b,$0c
	.else
	.byte $0d,$0e,$0f,$10
	.endif
`,
			want: []byte{0x01, 0x02, 0x03, 0x04, 0x09, 0x0a, 0x0b, 0x0c},
		},

		{
			name: "char constants",
			text: `   .org 0
	lda #'a'
`,
			want: []byte{0xa9, 0x41},
		},

		{
			name: "macros - instantiation",
			text: ` .org 0
m	.macro param
	lda param
	.endm

	m #0
	m $123
`,
			want: []byte{0xa9, 0x00, 0xad, 0x23, 0x01},
		},

		{
			name: "macros - local labels",
			text: ` .org 0
m	.macro
_l	nop
	beq _l
	.endm

	m
	m
`,
			want: []byte{0xea, 0xf0, 0xfd, 0xea, 0xf0, 0xfd},
		},
		{
			name: "macros - local labels are passed in",
			text: ` .org 0
m	.macro dest
	jmp dest
	.endm

start:
	m _l0+1
	m _l0+2
_l0 nop
`,
			want: []byte{0x4c, 0x07, 0x00, 0x4c, 0x08, 0x00, 0xea},
		},
		{
			name: "macros - don't complaing about existing patches",
			text: ` .org 0
m	.macro dest
	jmp dest
	.endm

start:
	jmp _l0
	m _l0+1
	m _l0+2
_l0 nop
`,
			want: []byte{0x4c, 0x09, 0x00, 0x4c, 0x0a, 0x00, 0x4c, 0x0b, 0x00, 0xea},
		},
		{
			name: "STx/LDx - use zero page addressing if possible",
			text: ` .org 0
dp	.equ $fb
	lda dp
	sta dp
	ldx dp
	stx dp
	ldy dp
	sty dp
`,
			want: []byte{0xa5, 0xfb, 0x85, 0xfb, 0xa6, 0xfb, 0x86, 0xfb, 0xa4, 0xfb, 0x84, 0xfb},
		},

		{
			name: "Encodings",
			text: ` .org 0
	.encoding "ascii"
	.byte "hello, world!"
	.encoding "petscii"
	.byte "hello, world!"
`,
			want: []byte{
				0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x21,
				0x48, 0x45, 0x4c, 0x4c, 0x4f, 0x2c, 0x20, 0x57, 0x4f, 0x52, 0x4c, 0x44, 0x21},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assembler := New([]string{}, "6502", "c128", "plain", "petscii", []string{})
			assembler.Assemble(text.Process("", test.text))

			errs := assembler.Errors()
			if len(errs) != len(test.wantErrors) {
				t.Errorf("Got %d, want %d errs", len(errs), len(test.wantErrors))
			} else {
				for i := range errs {
					got := errs[i]
					want := test.wantErrors[i]
					if got != want {
						t.Errorf("Error %d: got %+v, want %+v", i+1, got, want)
					}
				}
			}

			warnings := assembler.Warnings()
			if len(warnings) != len(test.wantWarnings) {
				t.Errorf("Got %d, want %d warnings", len(errs), len(test.wantWarnings))
			} else {
				for i := range warnings {
					got := warnings[i]
					want := test.wantWarnings[i]
					if got != want {
						t.Errorf("Warning %d: got %+v, want %+v", i+1, got, want)
					}
				}
			}
			if len(test.wantErrors) == 0 {
				got := assembler.GetBytes()
				if bytes.Compare(got, test.want) != 0 {
					t.Errorf("Got %s, want %s", toString(got), toString(test.want))
				}
			}
		})
	}
}

func toString(slice []byte) string {
	var parts []string
	for _, b := range slice {
		parts = append(parts, fmt.Sprintf("0x%02x", b))
	}
	return "[ " + strings.Join(parts, ", ") + " ]"

}
