package text

import "strings"

type Text struct {
	Lines []Line
}

type Line struct {
	Filename string
	LineNumber int
	Runes      []rune
}

func (t *Text) LastLine() Line {
	return t.Lines[len(t.Lines)-1]
}

func (l *Line) Extract(from, to Pos) string {
	return string(l.Runes[from.Col-1:to.Col-1]);
}

func Process(filename string, text string) Text {
	t := Text{}
	text = strings.ReplaceAll(text, "\r\n", "\n")
	curLine := Line{
		LineNumber: 1,
		Runes:       []rune{},
	}
	for _, r := range []rune(text) {
		curLine.Runes = append(curLine.Runes, r)
		if r == '\n' {
			t.Lines = append(t.Lines, curLine)
			curLine = Line{
				Filename: filename,
				LineNumber: curLine.LineNumber + 1,
				Runes:       []rune{},
			}
		}
	}
	t.Lines = append(t.Lines, curLine)
	return t
}
