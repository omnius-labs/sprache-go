package parse

type ParserResult[T any] struct {
	Value     T
	Remainder ParserInput
	Success   bool
	Message   string
	Errors    []ParseError
}
