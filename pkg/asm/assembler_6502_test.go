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
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assembler := New([]string{})
			src := " .cpu \"6502\"\n .org 0\n " + test.text
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
