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
	"strconv"
	"strings"
	"unicode"

	"github.com/asig/cbmasm/pkg/errors"
	"github.com/asig/cbmasm/pkg/text"
)

type TokenType int

const (
	Unknown TokenType = iota
	Ident
	Integer
	String
	Char
	LParen
	RParen
	Plus
	Minus
	Slash
	Asterisk
	Percent
	Dollar
	Ampersand
	Bar
	Dot
	Colon
	Semicolon
	Comma
	Lt
	Le
	Gt
	Ge
	Eq
	Ne
	Hash
	Tilde
	Caret

	// directives
	Cpu
	Platform
	Ifdef
	Ifndef
	If
	Else
	Endif
	Fail
	Include
	Incbin
	Reserve
	Byte
	Word
	Float
	Equ
	Org
	Align
	Macro
	Endm
	Encoding
	Output

	Eol
)

var identToTokenType = map[string]TokenType{
	".cpu":      Cpu,
	".platform": Platform,
	".ifdef":    Ifdef,
	".ifndef":   Ifndef,
	".if":       If,
	".else":     Else,
	".endif":    Endif,
	".fail":     Fail,
	".include":  Include,
	".incbin":   Incbin,
	".reserve":  Reserve,
	".byte":     Byte,
	".word":     Word,
	".float":    Float,
	".equ":      Equ,
	".org":      Org,
	".align":    Align,
	".macro":    Macro,
	".endm":     Endm,
	".encoding": Encoding,
	".output":   Output,
}

var tokenTypeToString = map[TokenType]string{
	Unknown:   "<unknown>",
	Ident:     "identifier",
	Integer:   "integer",
	String:    "string",
	Char:      "character",
	LParen:    "'('",
	RParen:    "')'",
	Plus:      "'+'",
	Minus:     "'-'",
	Slash:     "'/'",
	Asterisk:  "'*'",
	Percent:   "'%'",
	Dollar:    "'$'",
	Ampersand: "'&'",
	Bar:       "'|'",
	Dot:       "'.'",
	Colon:     "':'",
	Semicolon: "';'",
	Comma:     "'.'",
	Lt:        "'<'",
	Le:        "'<='",
	Gt:        "'>'",
	Ge:        "'>='",
	Eq:        "'='",
	Ne:        "'!='",
	Hash:      "'#'",
	Tilde:     "'~'",
	Caret:     "'^'",
	Cpu:       ".cpu",
	Platform:  ".platform",
	Ifdef:     ".ifdef",
	Ifndef:    ".ifndef",
	If:        ".if",
	Else:      ".else",
	Endif:     ".endif",
	Fail:      ".fail",
	Include:   ".include",
	Incbin:    ".incbin'",
	Reserve:   ".reserve",
	Byte:      ".byte",
	Word:      ".word",
	Equ:       ".equ",
	Org:       ".org",
	Align:     ".align",
	Macro:     ".macro",
	Endm:      ".endm",
	Encoding:  ".encoding",
	Output:    ".output",
	Eol:       "EOL",
}

func (t TokenType) String() string {
	return tokenTypeToString[t]
}

type Token struct {
	Type     TokenType
	StrVal   string
	IntVal   int64
	FloatVal float64
	Pos      text.Pos
}

type Scanner struct {
	line      text.Line
	curCol    int
	errorSink errors.Sink
}

func New(line text.Line, errorSink errors.Sink) *Scanner {
	return &Scanner{
		line:      line,
		curCol:    0,
		errorSink: errorSink,
	}
}

func isBinaryDigit(r rune) bool {
	return r >= '0' && r <= '1'
}

func isOctalDigit(r rune) bool {
	return r >= '0' && r <= '7'
}

func isHexDigit(r rune) bool {
	return r >= '0' && r <= '9' || r >= 'a' && r <= 'f' || r >= 'A' && r <= 'F'
}

func isIdentStartChar(r rune) bool {
	return r == '@' || r == '.' || r == '_' || unicode.IsLetter(r)
}

func isIdentChar(r rune) bool {
	return r == '@' || r == '.' || r == '_' || unicode.IsLetter(r)
}

func (scanner *Scanner) readIdent(ch rune) string {
	s := ""
	for isIdentChar(ch) || unicode.IsDigit(ch) {
		s = s + string(ch)
		ch = scanner.getch()
	}
	if ch == '\'' {
		// Consider this being part of the ident for things like "AF'"
		s = s + string(ch)
	} else {
		scanner.ungetch()
	}
	return s
}

func (scanner *Scanner) LineStart() text.Pos {
	return text.Pos{
		Filename: scanner.line.Filename,
		Line:     scanner.line.LineNumber,
		Col:      1,
	}
}

func (scanner *Scanner) CurPos() text.Pos {
	return text.Pos{
		Filename: scanner.line.Filename,
		Line:     scanner.line.LineNumber,
		Col:      scanner.curCol,
	}
}

func (scanner *Scanner) Line() *text.Line {
	return &scanner.line
}

func (scanner *Scanner) Scan() Token {
	// Scan over whitespace
	ch := scanner.getch()
	for unicode.IsSpace(ch) {
		ch = scanner.getch()
	}
	t := Token{
		Type: Unknown,
		Pos:  scanner.CurPos(),
	}
	switch {
	case ch == 0:
		t.Type = Eol
		return t
	case unicode.IsDigit(ch):
		// Read number
		i, s, err := scanner.readInteger(ch, 10, unicode.IsDigit)
		ch = scanner.getch()
		if ch == '.' {
			// Potentially floating point
			ch = scanner.getch()
			if unicode.IsDigit(ch) {
				// Floating point!
				_, s2, err := scanner.readInteger(ch, 10, unicode.IsDigit)
				t.StrVal = s + "." + s2
				t.FloatVal, _ = strconv.ParseFloat(t.StrVal, 64)
				t.Type = Float
				if err != nil {
					scanner.errorSink.AddError(t.Pos, "%s is not a valid floating point number", t.StrVal)
				}
				return t

			}
			scanner.ungetch()
		}
		scanner.ungetch()
		t.IntVal = i
		t.StrVal = s
		t.Type = Integer
		if err != nil {
			scanner.errorSink.AddError(t.Pos, "%s is not a valid integer", t.StrVal)
		}
		return t
	case ch == '.':
		// Ident, floating point number or dot
		ch = scanner.getch()
		if unicode.IsDigit(ch) {
			// Floating point
			_, s, err := scanner.readInteger(ch, 10, unicode.IsDigit)
			t.StrVal = "." + s
			t.Type = Float
			t.FloatVal, _ = strconv.ParseFloat(t.StrVal, 64)
			if err != nil {
				scanner.errorSink.AddError(t.Pos, ".%s is not a valid number", t.StrVal)
			}
		} else if isIdentChar(ch) {
			// Indent
			t.StrVal = "." + scanner.readIdent(ch)
			t.Type = Ident
			if tt, found := identToTokenType[strings.ToLower(t.StrVal)]; found {
				t.Type = tt
			}
		} else {
			// Just the dot
			scanner.ungetch()
			t.Type = Dot
		}
	case isIdentStartChar(ch):
		t.StrVal = scanner.readIdent(ch)
		t.Type = Ident
		if tt, found := identToTokenType[strings.ToLower(t.StrVal)]; found {
			t.Type = tt
		}
	case ch == '%':
		t.StrVal = "%"
		ch = scanner.getch()
		if !isBinaryDigit(ch) {
			scanner.ungetch()
			t.Type = Percent
			return t
		}
		i, s, err := scanner.readInteger(ch, 2, isBinaryDigit)
		t.IntVal = i
		t.StrVal = t.StrVal + s
		t.Type = Integer
		if err != nil {
			scanner.errorSink.AddError(t.Pos, "%s is not a valid number", t.StrVal)
		}
	case ch == '"':
		t.StrVal = scanner.readString(ch)
		t.Type = String
	case ch == '\'':
		pos := scanner.CurPos()
		t.StrVal = scanner.readString(ch)
		if len(t.StrVal) != 1 {
			scanner.errorSink.AddError(pos, "invalid character constant")
		}
		t.Type = Char
	case ch == '&':
		t.StrVal = "&"
		ch = scanner.getch()
		if !isOctalDigit(ch) {
			scanner.ungetch()
			t.Type = Ampersand
			return t
		}
		i, s, err := scanner.readInteger(ch, 8, isOctalDigit)
		t.IntVal = i
		t.StrVal = t.StrVal + s
		t.Type = Integer
		if err != nil {
			scanner.errorSink.AddError(t.Pos, "%s is not a valid number", t.StrVal)
		}
	case ch == '$':
		t.StrVal = "$"
		ch = scanner.getch()
		if !isHexDigit(ch) {
			scanner.ungetch()
			t.Type = Dollar
			return t
		}
		i, s, err := scanner.readInteger(ch, 16, isHexDigit)
		t.IntVal = i
		t.StrVal = t.StrVal + s
		t.Type = Integer
		if err != nil {
			scanner.errorSink.AddError(t.Pos, "%s is not a valid number", t.StrVal)
		}
	case ch == '!':
		ch = scanner.getch()
		if ch == '=' {
			t.StrVal = "!="
			t.Type = Ne
			return t
		}
		scanner.ungetch()
	case ch == '(':
		t.Type = LParen
	case ch == ')':
		t.Type = RParen
	case ch == '*':
		t.Type = Asterisk
	case ch == '/':
		t.Type = Slash
	case ch == ';':
		t.Type = Semicolon
	case ch == ',':
		t.Type = Comma
	case ch == '+':
		t.Type = Plus
	case ch == '-':
		t.Type = Minus
	case ch == '|':
		t.Type = Bar
	case ch == ':':
		t.Type = Colon
	case ch == '<':
		ch = scanner.getch()
		if ch == '=' {
			t.StrVal = "<="
			t.Type = Le
			return t
		}
		scanner.ungetch()
		t.Type = Lt
	case ch == '>':
		ch = scanner.getch()
		if ch == '=' {
			t.StrVal = ">="
			t.Type = Ge
			return t
		}
		scanner.ungetch()
		t.Type = Gt
	case ch == '=':
		t.Type = Eq
	case ch == '#':
		t.Type = Hash
	case ch == '~':
		t.Type = Tilde
	case ch == '^':
		t.Type = Caret
	}
	return t
}

func (scanner *Scanner) readString(separator rune) string {
	s := ""
	ch := scanner.getch()
	for ch != separator {
		p := scanner.CurPos()
		switch ch {
		case '\\':
			ch = scanner.getch()
			switch ch {
			case '\\', '"':
				s = s + string(ch)
			case 'n':
				s = s + "\n"
			case 'r':
				s = s + "\r"
			case 't':
				s = s + "\t"
			case 'b':
				s = s + "\b"
			default:
				scanner.ungetch()
				scanner.errorSink.AddError(p, "Unknown escape sequence")
			}
		case 0, '\n':
			scanner.errorSink.AddError(p, "Unterminated string")
			return s
		default:
			s = s + string(ch)
		}
		ch = scanner.getch()
	}
	return s
}

func (scanner *Scanner) readInteger(ch rune, base int, pred func(rune) bool) (int64, string, error) {
	s := ""
	for pred(ch) {
		s = s + string(ch)
		ch = scanner.getch()
	}
	scanner.ungetch()
	i, err := strconv.ParseInt(s, base, 64)
	return i, s, err
}

func (scanner *Scanner) getch() rune {
	var ch rune
	if scanner.curCol >= len(scanner.line.Runes) {
		ch = 0
	} else {
		ch = scanner.line.Runes[scanner.curCol]
	}
	scanner.curCol = scanner.curCol + 1
	return ch
}

func (scanner *Scanner) ungetch() {
	scanner.curCol = scanner.curCol - 1
	if scanner.curCol < 0 {
		panic("can't unget char!")
	}
}
