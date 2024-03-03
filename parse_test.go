package parser_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	parser "github.com/lyrise/sprache-go"
	"github.com/lyrise/sprache-go/internal"
)

func TestParse(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Parse Spec")
}

var _ = Describe("Parse Test", func() {

	Context("RuneFunc Test", func() {
		tests := []struct {
			name          string
			input         parser.ParserInput
			parser        parser.Parser[rune]
			wantSucceeded bool
			wantValue     rune
			wantRemainder string
		}{
			{
				name:          "simple",
				input:         parser.NewParserInput("a"),
				parser:        parser.RuneFunc(func(r rune) bool { return r == 'a' }, ""),
				wantSucceeded: true,
				wantValue:     'a',
				wantRemainder: "",
			},
			{
				name:          "simple 2",
				input:         parser.NewParserInput("ab"),
				parser:        parser.RuneFunc(func(r rune) bool { return r == 'a' }, ""),
				wantSucceeded: true,
				wantValue:     'a',
				wantRemainder: "b",
			},
			{
				name:          "simple 3",
				input:         parser.NewParserInput("ab"),
				parser:        parser.RuneFunc(func(r rune) bool { return r == 'b' }, ""),
				wantSucceeded: false,
				wantValue:     0,
				wantRemainder: "ab",
			},
		}
		for _, tt := range tests {
			It(tt.name, func() {
				got := tt.parser(tt.input)
				Expect(got.Succeeded).To(Equal(tt.wantSucceeded))
				Expect(got.Value).To(internal.EqualCmp(tt.wantValue))
				Expect(got.Remainder.String()).To(internal.EqualCmp(tt.wantRemainder))
			})
		}
	})
})
