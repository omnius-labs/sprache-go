package parse

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRuneFunc(t *testing.T) {
	tests := []struct {
		name          string
		input         ParserInput
		parser        Parser[rune]
		succeeded     bool
		wantValue     rune
		wantRemainder string
	}{
		{
			name:          "simple",
			input:         NewParserInput("a"),
			parser:        RuneFunc(func(r rune) bool { return r == 'a' }, ""),
			succeeded:     true,
			wantValue:     'a',
			wantRemainder: "",
		},
		{
			name:          "simple 2",
			input:         NewParserInput("ab"),
			parser:        RuneFunc(func(r rune) bool { return r == 'a' }, ""),
			succeeded:     true,
			wantValue:     'a',
			wantRemainder: "b",
		},
		{
			name:          "simple 3",
			input:         NewParserInput("ab"),
			parser:        RuneFunc(func(r rune) bool { return r == 'b' }, ""),
			succeeded:     false,
			wantValue:     0,
			wantRemainder: "ab",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.name), func(t *testing.T) {
			got := tt.parser(tt.input)
			if d := cmp.Diff(got.Succeeded, tt.succeeded); len(d) != 0 {
				t.Errorf("unexpected Succeeded: (-got +want)\n%s", d)
			}
			if d := cmp.Diff(got.Value, tt.wantValue); len(d) != 0 {
				t.Errorf("unexpected Value: (-got +want)\n%s", d)
			}
			if d := cmp.Diff(got.Remainder.String(), tt.wantRemainder); len(d) != 0 {
				t.Errorf("unexpected Input: (-got +want)\n%s", d)
			}
		})
	}
}
