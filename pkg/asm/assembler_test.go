package asm

import (
	"bytes"
	"fmt"
	"github.com/asig/cbmasm/pkg/text"
	"strings"
	"testing"
)

func TestAssembler_symbolResolution(t *testing.T) {
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
want: []byte{ 0x4c, 0x34, 0x12},
		},

		{
			name: "mixed symbols and labels",
			text: `   .org 0
sym .equ label
    nop
label inx
    jmp sym
`,
			want: []byte{0xea, 0xe8, 0x4c, 0x01, 0x00 },
		},

		// TODO: test cases
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assembler := New(text.Process("", test.text), []string{})
			assembler.Assemble();
			errors := assembler.Errors()
			if len(errors) != 0 {
				t.Errorf("Got %+v, want 0 errors", errors)
			}
			warnings := assembler.Warnings()
			if len(warnings) != 0 {
				t.Errorf("Got %+v, want 0 warnings", errors)
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
