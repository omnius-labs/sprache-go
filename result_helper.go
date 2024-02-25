package parse

func IfSuccess[T any, U any](result ParserResult[T], next func(ParserResult[T]) ParserResult[U]) ParserResult[U] {
	if result.Succeeded {
		return next(result)
	}
	return NewFailureResult[U](result.Remainder, result.Message, result.Expectations)
}

func IfFailure[T any](result ParserResult[T], next func(ParserResult[T]) ParserResult[T]) ParserResult[T] {
	if !result.Succeeded {
		return next(result)
	}
	return result
}
