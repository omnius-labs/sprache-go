package parse

import (
	"slices"
	"unicode"
)

func RuneFunc(predicate func(rune) bool) Parser[rune] {
	return func(input ParserInput) (ParserResult[rune], error) {
		if input.IsEnd() {
			return ParserResult[rune]{}, ParseError{}
		}

		if predicate(input.Current()) {
			return ParserResult[rune]{
				Value:     input.Current(),
				Remainder: input.Advance(),
			}, nil
		}

		return ParserResult[rune]{}, nil
	}
}

func RuneExceptFunc(predicate func(rune) bool) Parser[rune] {
	return RuneFunc(func(c rune) bool {
		return !predicate(c)
	})
}

func Rune(c rune) Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return r == c
	})
}

func RuneExcept(c rune) Parser[rune] {
	return RuneExceptFunc(func(r rune) bool {
		return r == c
	})
}

func Runes(rs ...rune) Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return slices.Contains(rs, r)
	})
}

func RunesString(s string) Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return slices.Contains([]rune(s), r)
	})
}

func RunesExcept(rs ...rune) Parser[rune] {
	return RuneExceptFunc(func(r rune) bool {
		return slices.Contains(rs, r)
	})
}

func RunesStringExcept(s string) Parser[rune] {
	return RuneExceptFunc(func(r rune) bool {
		return slices.Contains([]rune(s), r)
	})
}

func RuneIgnoreCase(c rune) Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return unicode.ToLower(r) == unicode.ToLower(c)
	})
}

// func String(s string) Parser[rune] {
// }

// func StringIgnoreCase(s string) Parser[rune] {
// }

func AnyRune() Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return true
	})
}

func Space() Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return unicode.IsSpace(r)
	})
}

func Digit() Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return unicode.IsDigit(r)
	})
}

func Letter() Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return unicode.IsLetter(r)
	})
}

func LetterOrDigit() Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r)
	})
}

func Lower() Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return unicode.IsLower(r)
	})
}

func Upper() Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return unicode.IsUpper(r)
	})
}

func Numeric() Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return unicode.IsNumber(r)
	})
}
