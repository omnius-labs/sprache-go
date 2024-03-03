package parser

type Parser[T any] func(ParserInput) ParserResult[T]
