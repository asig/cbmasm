package scanner

import (
	"github.com/asig/cbmasm/pkg/errors"
	"github.com/asig/cbmasm/pkg/text"
	"testing"
)

type errorSink struct {
	e []errors.Error
}

func (e *errorSink) AddError(pos text.Pos, message string) {
	e.e = append(e.e, errors.Error{pos, message})
}



func TestScanner_Scan_numbers(t *testing.T) {
	tests := []struct {
		name string
		text text.Line
		want int64
	}{
		{
			name: "Binary",
			text: text.Process("filename","%1001").Lines[0],
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
			scanner := New("filename", test.text, &errors);
			got := scanner.Scan()
			if got.Type != Number {
				t.Errorf("got token type %s, expected %s", got.Type, Number)
			}
			if got.IntVal != test.want {
				t.Errorf("got %d, expected %d", got.IntVal, test.want)
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
			text: text.Process("filename",`"What\"ever!"`).Lines[0],
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
			scanner := New("filename", test.text, &errors);
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
