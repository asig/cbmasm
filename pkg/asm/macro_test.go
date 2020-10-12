package asm

import (
	"github.com/asig/cbmasm/pkg/text"
	"testing"
)

func TestMacro_ReplaceParams(t *testing.T) {
	tests := []struct {
		name string
		text string
		paramNames []string
		paramVals []string
		want string
	}{
		{
			name: "Simple parameter substitution",
			text: "12 foo bar (baz+bar) foobar",
			paramNames: []string{"foo", "bar", "foobar"},
			paramVals: []string{"#12", "(A),X", "X"},
			want: "12 #12 (A),X (baz+(A),X) X",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := macro{}
			for _, p := range test.paramNames {
				m.addParam(p)
			}
			m.text = []text.Line{ {1, []rune(test.text)}}

			var actuals []param
			for _, p := range test.paramVals {
				actuals = append(actuals, param{rawText: p})
			}
			got := m.replaceParams(actuals);
			if len(got) != 1 {
				t.Fatalf("Expected 1 line, got %d", len(got))
			}
			if string(got[0].Runes) != test.want {
				t.Fatalf("Expected %q, got %q", string(got[0].Runes), test.want)
			}
		})
	}
}
