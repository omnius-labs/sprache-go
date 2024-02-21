package parse

type Parser[T any] func(ParserInput) (ParserResult[T], error)
