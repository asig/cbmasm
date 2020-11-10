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
	"fmt"

	"github.com/asig/cbmasm/pkg/scanner"
	"github.com/asig/cbmasm/pkg/text"
)

type macro struct {
	pos text.Pos

	params []string
	text   *text.Text
}

func (m *macro) addParam(name string) error {
	if m.paramIndex(name) > -1 {
		return fmt.Errorf("Parameter %s already exists", name)
	}
	m.params = append(m.params, name)
	return nil
}

func (m *macro) paramIndex(name string) int {
	for idx := range m.params {
		if m.params[idx] == name {
			return idx
		}
	}
	return -1
}

func (m *macro) replaceParams(actuals []string) []text.Line {

	paramMap := make(map[string]string)
	for idx, p := range m.params {
		paramMap[p] = actuals[idx]
	}

	var res []text.Line
	for _, line := range m.text.Lines {
		substituted := substituteParams(line, paramMap)
		res = append(res, text.Line{line.Filename, line.LineNumber, substituted})
	}
	return res
}

type dummyErrorSink struct{}

func (d *dummyErrorSink) AddError(_ text.Pos, _ string, _ ...interface{}) {}

type replacement struct {
	start, len int
	text       string
}

func substituteParams(line text.Line, paramMap map[string]string) []rune {

	// We want to keep the original text of the non-params, so we just find the positions to replace first.
	var repls []replacement

	// We don't care about errors at this point, they'll surface during instantiation
	s := scanner.New(line, &dummyErrorSink{})
	t := s.Scan()
	for t.Type != scanner.Eol {
		if t.Type == scanner.Ident {
			// Potentially a mos6502Param
			if val, found := paramMap[t.StrVal]; found {
				// YES. Insert replacement at the beginning
				repls = append([]replacement{{t.Pos.Col - 1, len(t.StrVal), val}}, repls...)
			}
		}
		t = s.Scan()
	}

	res := string(line.Runes)
	for _, r := range repls {
		res = res[0:r.start] + r.text + res[r.start+r.len:]
	}
	return []rune(res)
}
