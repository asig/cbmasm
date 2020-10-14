package asm

import (
	"fmt"
	scanner "github.com/asig/cbmasm/pkg/scanner"

	"github.com/asig/cbmasm/pkg/text"
)

type macro struct {
	pos text.Pos

	params []string
	text *text.Text
}

func (m *macro) addParam(name string) error {
	if m.paramIndex(name) > -1 {
		return fmt.Errorf("Parameter %s already exists", name)
	}
	m.params = append(m.params, name)
	return nil
}

func (m *macro) paramIndex(name string) int {
	for idx, _ := range m.params {
		if m.params[idx] == name {
			return idx
		}
	}
	return -1
}

func (m *macro) replaceParams(actuals []param) []text.Line {

	paramMap := make(map[string]string)
	for idx, p := range m.params {
		paramMap[p] = actuals[idx].rawText
	}

	var res []text.Line
	for _, line := range m.text.Lines {
		substituted := substituteParams(line, paramMap)
		res = append(res, text.Line{line.Filename, line.LineNumber, substituted})
	}
	return res
}


type dummyErrorSink struct {}
func (d *dummyErrorSink) AddError(pos text.Pos, message string, args... interface{}) {}

type replacement struct {
	start, len int;
	text string
}

func substituteParams(line text.Line, paramMap map[string]string) []rune {

	// We want to keep the original text of the non-params, so we just find the positions to replace first.
	var repls []replacement

	// We don't care about errors at this point, they'll surface during instantiation
	s := scanner.New(line, &dummyErrorSink{})
	t := s.Scan()
	for t.Type != scanner.Eol {
		if t.Type == scanner.Ident {
			// Potentially a param
			if val, found := paramMap[t.StrVal]; found {
				// YES. Insert replacement at the beginning
				repls = append([]replacement{replacement{t.Pos.Col-1, len(t.StrVal), val}}, repls...)
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
