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
			assembler := New([]string{})
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
		name string
		text string
		want []byte
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

		// TODO add test for unresolved local labels

		{
			name: "byte constants",
			text: `   .org 0
	.byte $01,$02,$03,$04
`,
			want: []byte{0x01, 0x02, 0x03, 0x04},
		},

		{
			name: "word constants",
			text: `   .org 0
	.word $0102,$0304
`,
			want: []byte{0x02, 0x01, 0x04, 0x03},
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assembler := New([]string{})
			assembler.Assemble(text.Process("", test.text))
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

func toString(slice []byte) string {
	var parts []string
	for _, b := range slice {
		parts = append(parts, fmt.Sprintf("0x%02x", b))
	}
	return "[ " + strings.Join(parts, ", ") + " ]"

}
