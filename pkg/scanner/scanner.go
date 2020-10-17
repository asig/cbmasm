package scanner

import (
	"github.com/asig/cbmasm/pkg/errors"
	"github.com/asig/cbmasm/pkg/text"
	"strconv"
	"strings"
	"unicode"
)

type TokenType int

const (
	Unknown TokenType = iota
	Ident
	Number
	String
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
	Equ
	Org
	Macro
	Endm

	Eol
)

var identToTokenType = map[string]TokenType{
	".ifdef":   Ifdef,
	".ifndef":  Ifndef,
	".if":      If,
	".else":    Else,
	".endif":   Endif,
	".fail":    Fail,
	".include": Include,
	".incbin'": Incbin,
	".reserve": Reserve,
	".byte":    Byte,
	".word":    Word,
	".equ":     Equ,
	".org":     Org,
	".macro":   Macro,
	".endm":    Endm,
}

var tokenTypeToString = map[TokenType]string{
	Unknown:   "<unknown>",
	Ident:     "identifier",
	Number:    "number",
	String:    "string",
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
	Macro:     ".macro",
	Endm:      ".endm",
	Eol:       "EOL",
}

func (t TokenType) String() string {
	return tokenTypeToString[t]
}

type Token struct {
	Type   TokenType
	StrVal string
	IntVal int64
	Pos    text.Pos
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

func isIdentChar(r rune) bool {
	return r == '@' || r == '.' || r == '_' || unicode.IsLetter(r)
}

func (scanner *Scanner) readIdent(ch rune) string {
	s := ""
	for isIdentChar(ch) || unicode.IsDigit(ch) {
		s = s + string(ch)
		ch = scanner.getch()
	}
	scanner.ungetch()
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
		i, s, err := scanner.readNumber(ch, 10, unicode.IsDigit)
		t.IntVal = i
		t.StrVal = s
		t.Type = Number
		if err != nil {
			scanner.errorSink.AddError(t.Pos, "%s is not a valid number", t.StrVal)
		}
		return t
	case isIdentChar(ch):
		t.StrVal = scanner.readIdent(ch)
		if tt, found := identToTokenType[strings.ToLower(t.StrVal)]; found {
			t.Type = tt
		} else {
			t.Type = Ident
		}
	case ch == '%':
		t.StrVal = "%"
		ch = scanner.getch()
		if !isBinaryDigit(ch) {
			scanner.ungetch()
			t.Type = Percent
			return t
		}
		i, s, err := scanner.readNumber(ch, 2, isBinaryDigit)
		t.IntVal = i
		t.StrVal = t.StrVal + s
		t.Type = Number
		if err != nil {
			scanner.errorSink.AddError(t.Pos, "%s is not a valid number", t.StrVal)
		}
	case ch == '"':
		t.StrVal = ""
		ch = scanner.getch()
		for ch != '"' {
			p := scanner.CurPos()
			switch ch {
			case '\\':
				ch = scanner.getch()
				switch ch {
				case '\\', '"':
					t.StrVal = t.StrVal + string(ch)
				case 'n':
					t.StrVal = t.StrVal + "\n"
				case 't':
					t.StrVal = t.StrVal + "\t"
				case 'b':
					t.StrVal = t.StrVal + "\b"
				default:
					scanner.ungetch()
					scanner.errorSink.AddError(p, "Unknown escape sequence")
				}
			case '\n':
				scanner.errorSink.AddError(p, "Unterminated string")
				break
			default:
				t.StrVal = t.StrVal + string(ch)
			}
			ch = scanner.getch()
		}
		t.Type = String
	case ch == '&':
		t.StrVal = "&"
		ch = scanner.getch()
		if !isOctalDigit(ch) {
			scanner.ungetch()
			t.Type = Ampersand
			return t
		}
		i, s, err := scanner.readNumber(ch, 8, isOctalDigit)
		t.IntVal = i
		t.StrVal = t.StrVal + s
		t.Type = Number
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
		i, s, err := scanner.readNumber(ch, 16, isHexDigit)
		t.IntVal = i
		t.StrVal = t.StrVal + s
		t.Type = Number
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
	case ch == '.':
		ch = scanner.getch()
		if isIdentChar(ch) || unicode.IsDigit(ch) {
			t.StrVal = "." + scanner.readIdent(ch)
			t.Type = Ident
		}
		scanner.ungetch()
		t.Type = Dot
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

func (scanner *Scanner) readNumber(ch rune, base int, pred func(rune) bool) (int64, string, error) {
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
