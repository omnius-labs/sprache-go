package parser

import "slices"

type ParserInput struct {
	source   []rune
	position int
	line     int
	column   int
}

func NewParserInput(source string) ParserInput {
	return ParserInput{
		source:   []rune(source),
		position: 0,
		line:     1,
		column:   1,
	}
}

func (p ParserInput) Equal(other ParserInput) bool {
	return slices.Equal(p.source, other.source) && p.position == other.position && p.line == other.line && p.column == other.column
}

func (p ParserInput) Advance() ParserInput {
	line := p.line
	if p.Current() == '\n' {
		line++
	}

	column := p.column
	if p.Current() == '\n' {
		column = 1
	} else {
		column++
	}

	return ParserInput{
		source:   p.source,
		position: p.position + 1,
		line:     line,
		column:   column,
	}
}

func (p ParserInput) Source() string {
	return string(p.source)
}

func (p ParserInput) Current() rune {
	return p.source[p.position]
}

func (p ParserInput) IsEnd() bool {
	return p.position >= len(p.source)
}

func (p ParserInput) Position() int {
	return p.position
}

func (p ParserInput) Line() int {
	return p.line
}

func (p ParserInput) Column() int {
	return p.column
}

func (p ParserInput) String() string {
	return string(p.source[p.position:])
}
