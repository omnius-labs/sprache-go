package parse

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRune(t *testing.T) {
	tests := []struct {
		input     ParserInput
		parser    Parser[rune]
		wantValue rune
		wantInput string
	}{
		{
			input:     NewParserInput("a"),
			parser:    RuneFunc(func(r rune) bool { return r == 'a' }, ""),
			wantValue: 'a',
			wantInput: "",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.input), func(t *testing.T) {
			got := tt.parser(tt.input)
			if d := cmp.Diff(got.Value, tt.wantValue); len(d) != 0 {
				t.Errorf("unexpected Value: (-got +want)\n%s", d)
			}
			if d := cmp.Diff(got.Remainder.String(), tt.wantInput); len(d) != 0 {
				t.Errorf("unexpected Input: (-got +want)\n%s", d)
			}
		})
	}
}

func TestToken(t *testing.T) {
	v := "     a      	"
	res := Token(Rune('a'))(NewParserInput(v)).Value
	if res != 'a' {
		t.Errorf("unexpected result: got %v, want %v", res, "a")
	}
}
