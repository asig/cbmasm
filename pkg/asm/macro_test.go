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
	"testing"

	"github.com/asig/cbmasm/pkg/text"
)

func TestMacro_ReplaceParams(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		paramNames []string
		paramVals  []string
		want       string
	}{
		{
			name:       "Simple parameter substitution",
			text:       "12 foo bar (baz+bar) foobar",
			paramNames: []string{"foo", "bar", "foobar"},
			paramVals:  []string{"#12", "(A),X", "X"},
			want:       "12 #12 (A),X (baz+(A),X) X",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := macro{}
			for _, p := range test.paramNames {
				err := m.addParam(p)
				if err != nil {
					t.Errorf("Unexpected error while adding param: %s", err)
				}
			}
			m.text = &text.Text{
				Lines: []text.Line{{"dummy", 1, []rune(test.text)}},
			}

			var actuals []param
			for _, p := range test.paramVals {
				actuals = append(actuals, param{rawText: p})
			}
			got := m.replaceParams(actuals)
			if len(got) != 1 {
				t.Fatalf("Expected 1 line, got %d", len(got))
			}
			if string(got[0].Runes) != test.want {
				t.Fatalf("Expected %q, got %q", string(got[0].Runes), test.want)
			}
		})
	}
}
