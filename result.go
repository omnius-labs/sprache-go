package parse

type ParserResult[T any] struct {
	Value        T
	Remainder    ParserInput
	Succeeded    bool
	Message      string
	Expectations []string
}

func NewSuccessResult[T any](value T, remainder ParserInput) ParserResult[T] {
	return ParserResult[T]{
		Value:     value,
		Remainder: remainder,
		Succeeded: true,
	}
}

func NewFailureResult[T any](remainder ParserInput, message string, expectations []string) ParserResult[T] {
	return ParserResult[T]{
		Remainder:    remainder,
		Message:      message,
		Expectations: expectations,
		Succeeded:    false,
	}
}
