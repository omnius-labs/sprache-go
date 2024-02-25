package parse

import (
	"fmt"
	"slices"
	"strings"
	"unicode"

	"github.com/lyrise/sprache-go/helper"
)

func RuneFunc(predicate func(rune) bool, description string) Parser[rune] {
	return func(input ParserInput) ParserResult[rune] {
		if input.IsEnd() {
			return NewFailureResult[rune](input, "unexpected end of input", []string{description})
		}

		if predicate(input.Current()) {
			return NewSuccessResult[rune](input.Current(), input.Advance())
		}

		return NewFailureResult[rune](input, fmt.Sprintf("unexpected %v", input.Current()), []string{description})
	}
}

func RuneExceptFunc(predicate func(rune) bool, description string) Parser[rune] {
	return RuneFunc(func(c rune) bool {
		return !predicate(c)
	}, fmt.Sprintf("any character expect %v", description))
}

func Rune(c rune) Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return r == c
	}, string(c))
}

func RuneExcept(c rune) Parser[rune] {
	return RuneExceptFunc(func(r rune) bool {
		return r == c
	}, string(c))
}

func Runes(rs ...rune) Parser[rune] {
	description := strings.Join(helper.Map(rs, helper.RuneToString), "|")
	return RuneFunc(func(r rune) bool {
		return slices.Contains(rs, r)
	}, description)
}

func RunesString(s string) Parser[rune] {
	rs := []rune(s)
	description := strings.Join(helper.Map(rs, helper.RuneToString), "|")
	return RuneFunc(func(r rune) bool {
		return slices.Contains(rs, r)
	}, description)
}

func RunesExcept(rs ...rune) Parser[rune] {
	description := strings.Join(helper.Map(rs, helper.RuneToString), "|")
	return RuneExceptFunc(func(r rune) bool {
		return slices.Contains(rs, r)
	}, description)
}

func RunesStringExcept(s string) Parser[rune] {
	rs := []rune(s)
	description := strings.Join(helper.Map(rs, helper.RuneToString), "|")
	return RuneExceptFunc(func(r rune) bool {
		return slices.Contains([]rune(s), r)
	}, description)
}

func RuneIgnoreCase(c rune) Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return unicode.ToLower(r) == unicode.ToLower(c)
	}, string(c))
}

// func String(s string) Parser[rune] {
// }

// func StringIgnoreCase(s string) Parser[rune] {
// }

func AnyRune() Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return true
	}, "any character")
}

func WhiteSpace() Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return unicode.IsSpace(r)
	}, "whitespace")
}

func Digit() Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return unicode.IsDigit(r)
	}, "digit")
}

func Letter() Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return unicode.IsLetter(r)
	}, "letter")
}

func LetterOrDigit() Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r)
	}, "letter or digit")
}

func Lower() Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return unicode.IsLower(r)
	}, "lowercase letter")
}

func Upper() Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return unicode.IsUpper(r)
	}, "uppercase letter")
}

func Numeric() Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return unicode.IsNumber(r)
	}, "numeric character")
}

func Not[T any](parser Parser[T]) Parser[T] {
	return func(input ParserInput) ParserResult[T] {
		r := parser(input)

		if !r.Succeeded {
			msg := fmt.Sprintf("unexpected %v", strings.Join(r.Expectations, ", "))
			return NewFailureResult[T](input, msg, []string{})
		}

		return NewSuccessResult[T](r.Value, r.Remainder)
	}
}

func Then[T any, U any](first Parser[T], second func(T) Parser[U]) Parser[U] {
	return func(input ParserInput) ParserResult[U] {
		r := first(input)

		return IfSuccess[T, U](r, func(r ParserResult[T]) ParserResult[U] {
			return second(r.Value)(r.Remainder)
		})
	}
}

func Many[T any](parser Parser[T]) Parser[[]T] {
	return func(input ParserInput) ParserResult[[]T] {
		var results []T
		remainder := input

		for {
			r := parser(remainder)

			if !r.Succeeded {
				break
			}

			results = append(results, r.Value)
			remainder = r.Remainder
		}

		return NewSuccessResult[[]T](results, remainder)
	}
}

// XMany

// AtLeastOnce

// XAtLeastOnce

func End[T any](parser Parser[T]) Parser[T] {
	return func(input ParserInput) ParserResult[T] {
		r := parser(input)

		return IfSuccess[T](r, func(r ParserResult[T]) ParserResult[T] {
			if r.Remainder.IsEnd() {
				return r
			}

			return NewFailureResult[T](input, fmt.Sprintf("unexpected %v", r.Remainder.Current()), []string{"end of input"})
		})
	}
}

func Select[T any, U any](parser Parser[T], convert func(T) U) Parser[U] {
	return Then(parser, func(v T) Parser[U] {
		return Return[U](convert(v))
	})
}

func Token[T any](parser Parser[T]) Parser[T] {
	return Then(Many(WhiteSpace()), func(_ []rune) Parser[T] {
		return Then(parser, func(v T) Parser[T] {
			return Then(Many(WhiteSpace()), func(_ []rune) Parser[T] {
				return Return[T](v)
			})
		})
	})
}

func Return[T any](value T) Parser[T] {
	return func(input ParserInput) ParserResult[T] {
		return NewSuccessResult[T](value, input)
	}
}

func Text(parser Parser[[]rune]) Parser[string] {
	return Select(parser, func(rs []rune) string {
		return string(rs)
	})
}
