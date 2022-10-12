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
package scanner

import (
	"testing"

	"github.com/asig/cbmasm/pkg/errors"
	"github.com/asig/cbmasm/pkg/text"
)

type errorSink struct {
	e []errors.Error
}

func (e *errorSink) AddError(pos text.Pos, message string, args ...interface{}) {
	e.e = append(e.e, errors.Error{pos, message})
}

func TestScanner_Scan_integers(t *testing.T) {
	tests := []struct {
		name string
		text text.Line
		want int64
	}{
		{
			name: "Binary",
			text: text.Process("filename", "%1001").Lines[0],
			want: 9,
		},
		{
			name: "Hex",
			text: text.Process("filename", "$12af").Lines[0],
			want: 4783,
		},
		{
			name: "Octal",
			text: text.Process("filename", "&67").Lines[0],
			want: 55,
		},
		{
			name: "Decimal",
			text: text.Process("filename", "12345").Lines[0],
			want: 12345,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errors := errorSink{}
			scanner := New(test.text, &errors)
			got := scanner.Scan()
			if got.Type != Integer {
				t.Errorf("got token type %s, expected %s", got.Type, Integer)
			}
			if got.IntVal != test.want {
				t.Errorf("got %d, expected %d", got.IntVal, test.want)
			}
		})
	}
}

func TestScanner_Scan_floats(t *testing.T) {
	tests := []struct {
		name string
		text text.Line
		want float64
	}{
		{
			name: "Float",
			text: text.Process("filename", "123.456").Lines[0],
			want: 123.456,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errors := errorSink{}
			scanner := New(test.text, &errors)
			got := scanner.Scan()
			if got.Type != Float {
				t.Errorf("got token type %s, expected %s", got.Type, Float)
			}
			if got.FloatVal != test.want {
				t.Errorf("got %f, expected %f", got.FloatVal, test.want)
			}
		})
	}
}

func TestScanner_Scan_strings(t *testing.T) {
	tests := []struct {
		name string
		text text.Line
		want string
	}{
		{
			name: "Plain and simple",
			text: text.Process("filename", `"Whatever!"`).Lines[0],
			want: "Whatever!",
		},
		{
			name: "Escaped quote",
			text: text.Process("filename", `"What\"ever!"`).Lines[0],
			want: `What"ever!`,
		},
		{
			name: "Escaped backslash",
			text: text.Process("filename", `"What\\ever!"`).Lines[0],
			want: `What\ever!`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errors := errorSink{}
			scanner := New(test.text, &errors)
			got := scanner.Scan()
			if got.Type != String {
				t.Errorf("got token type %s, expected %s", got.Type, String)
			}
			if got.StrVal != test.want {
				t.Errorf("got %s, expected %s", got.StrVal, test.want)
			}
		})
	}
}
