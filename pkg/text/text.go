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
package text

import "strings"

type Pos struct {
	Filename  string
	Line, Col int
}

type Text struct {
	Lines []Line
}

type Line struct {
	Filename   string
	LineNumber int
	Runes      []rune
}

func (t *Text) Append(text Text) {
	t.Lines = append(t.Lines, text.Lines...)
}

func (t *Text) AppendLine(l Line) {
	t.Lines = append(t.Lines, l)
}

func (t *Text) LastLine() Line {
	return t.Lines[len(t.Lines)-1]
}

func (l *Line) Extract(from, to Pos) string {
	return string(l.Runes[from.Col-1 : to.Col-1])
}

func Process(filename string, text string) Text {
	t := Text{}
	text = strings.ReplaceAll(text, "\r\n", "\n")
	curLine := Line{
		Filename:   filename,
		LineNumber: 1,
		Runes:      []rune{},
	}
	for _, r := range []rune(text) {
		curLine.Runes = append(curLine.Runes, r)
		if r == '\n' {
			t.Lines = append(t.Lines, curLine)
			curLine = Line{
				Filename:   filename,
				LineNumber: curLine.LineNumber + 1,
				Runes:      []rune{},
			}
		}
	}
	t.Lines = append(t.Lines, curLine)
	return t
}
