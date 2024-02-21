package parse

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRune(t *testing.T) {
	tests := []struct {
		input      ParserInput
		parser     Parser[rune]
		wantResult rune
		wantInput  string
		wantErr    error
	}{
		{
			input:      NewParserInput("a"),
			parser:     RuneFunc(func(r rune) bool { return r == 'a' }),
			wantResult: 'a',
			wantInput:  "",
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.input), func(t *testing.T) {
			got, err := tt.parser(tt.input)
			if d := cmp.Diff(got.Value, tt.wantResult); len(d) != 0 {
				t.Errorf("unexpected Result: (-got +want)\n%s", d)
			}
			if d := cmp.Diff(got.Remainder.String(), tt.wantInput); len(d) != 0 {
				t.Errorf("unexpected Input: (-got +want)\n%s", d)
			}
			if d := cmp.Diff(err, tt.wantErr); len(d) != 0 {
				t.Errorf("unexpected Error: (-got +want)\n%s", d)
			}
		})
	}
}
