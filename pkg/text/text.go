package text

import "strings"

type Text struct {
	Filename string
	Lines []Line
}

type Line struct {
	LineNumber int
	Runes      []rune
}

func (t *Text) LastLine() Line {
	return t.Lines[len(t.Lines)-1]
}

func Process(filename string, text string) Text {
	t := Text{
		Filename: filename,
	}
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
				LineNumber: curLine.LineNumber + 1,
				Runes:       []rune{},
			}
		}
	}
	t.Lines = append(t.Lines, curLine)
	return t
}
