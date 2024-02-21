package parse

type ParserInput interface {
	Advance() ParserInput
	Source() string
	Current() rune
	IsEnd() bool
	Position() int
	Line() int
	Column() int
	String() string
}

type parserInput struct {
	source   []rune
	position int
	line     int
	column   int
}

func NewParserInput(source string) ParserInput {
	return parserInput{
		source:   []rune(source),
		position: 0,
		line:     1,
		column:   1,
	}
}

func (p parserInput) Advance() ParserInput {
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

	return parserInput{
		source:   p.source,
		position: p.position + 1,
		line:     line,
		column:   column,
	}
}

func (p parserInput) Source() string {
	return string(p.source)
}

func (p parserInput) Current() rune {
	return p.source[p.position]
}

func (p parserInput) IsEnd() bool {
	return p.position >= len(p.source)
}

func (p parserInput) Position() int {
	return p.position
}

func (p parserInput) Line() int {
	return p.line
}

func (p parserInput) Column() int {
	return p.column
}

func (p parserInput) String() string {
	return string(p.source[p.position:])
}
