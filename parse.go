package parser

import (
	"fmt"
	"slices"
	"strings"
	"unicode"

	"github.com/lyrise/sprache-go/internal"
)

// TryParse a single character matching 'predicate'
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

// Parse a single character except those matching 'predicate'
func RuneExceptFunc(predicate func(rune) bool, description string) Parser[rune] {
	return RuneFunc(func(c rune) bool {
		return !predicate(c)
	}, fmt.Sprintf("any character expect %v", description))
}

// Parse a single character c.
func Rune(c rune) Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return r == c
	}, string(c))
}

// Parse a single character except c.
func RuneExcept(c rune) Parser[rune] {
	return RuneExceptFunc(func(r rune) bool {
		return r == c
	}, string(c))
}

// Parse a single character of any in rs
func Runes(rs ...rune) Parser[rune] {
	description := strings.Join(internal.Map(rs, internal.RuneToString), "|")
	return RuneFunc(func(r rune) bool {
		return slices.Contains(rs, r)
	}, description)
}

// Parse a single character of any in s
func RunesString(s string) Parser[rune] {
	rs := []rune(s)
	description := strings.Join(internal.Map(rs, internal.RuneToString), "|")
	return RuneFunc(func(r rune) bool {
		return slices.Contains(rs, r)
	}, description)
}

// Parses a single character except for those in rs
func RunesExcept(rs ...rune) Parser[rune] {
	description := strings.Join(internal.Map(rs, internal.RuneToString), "|")
	return RuneExceptFunc(func(r rune) bool {
		return slices.Contains(rs, r)
	}, description)
}

// Parses a single character except for those in s
func RunesStringExcept(s string) Parser[rune] {
	rs := []rune(s)
	description := strings.Join(internal.Map(rs, internal.RuneToString), "|")
	return RuneExceptFunc(func(r rune) bool {
		return slices.Contains([]rune(s), r)
	}, description)
}

// Parse a single character in a case-insensitive fashion.
func IgnoreCase(c rune) Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return unicode.ToLower(r) == unicode.ToLower(c)
	}, string(c))
}

// Parse a string in a case-insensitive fashion.
func IgnoreCaseString(s string) Parser[[]rune] {
	res := Return([]rune{})
	for _, r := range s {
		res = Concat(res, Once(IgnoreCase(r)))
	}
	return SetExpectationIfError(res, s)
}

// Parse a string of characters.
func String(s string) Parser[[]rune] {
	res := Return([]rune{})
	for _, r := range s {
		res = Concat(res, Once(Rune(r)))
	}
	return SetExpectationIfError(res, s)
}

// Parse any character.
func AnyRune() Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return true
	}, "any character")
}

// Parse a whitespace.
func WhiteSpace() Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return unicode.IsSpace(r)
	}, "whitespace")
}

// Parse a digit.
func Digit() Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return unicode.IsDigit(r)
	}, "digit")
}

// Parse a letter.
func Letter() Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return unicode.IsLetter(r)
	}, "letter")
}

// Parse a letter or digit.
func LetterOrDigit() Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r)
	}, "letter or digit")
}

// Parse a lowercase letter.
func Lower() Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return unicode.IsLower(r)
	}, "lowercase letter")
}

// Parse an uppercase letter.
func Upper() Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return unicode.IsUpper(r)
	}, "uppercase letter")
}

// Parse a numeric character.
func Numeric() Parser[rune] {
	return RuneFunc(func(r rune) bool {
		return unicode.IsNumber(r)
	}, "numeric character")
}

// Constructs a parser that will fail if the given parser succeeds,
// and will succeed if the given parser fails. In any case, it won't
// consume any input. It's like a negative look-ahead in regex.
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

// Parse first, and if successful, then parse second.
func Then[T, U any](first Parser[T], second func(T) Parser[U]) Parser[U] {
	return func(input ParserInput) ParserResult[U] {
		r := first(input)
		return IfSuccess(r, func(r ParserResult[T]) ParserResult[U] {
			return second(r.Value)(r.Remainder)
		})
	}
}

// Parse a stream of elements.
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

// Parse a stream of elements, failing if any element is only partially parsed.
func XMany[T any](parser Parser[T]) Parser[[]T] {
	return Then(Many(parser), func(m []T) Parser[[]T] {
		return XOr(Once(parser), Return(m))
	})
}

// TryParse a stream of elements with at least one item.
func AtLeastOnce[T any](parser Parser[T]) Parser[[]T] {
	return Then(Once(parser), func(t1 []T) Parser[[]T] {
		return Select(Many(parser), func(ts []T) []T {
			return internal.Union(t1, ts)
		})
	})
}

// TryParse a stream of elements with at least one item. Except the first
func XAtLeastOnce[T any](parser Parser[T]) Parser[[]T] {
	return Then(Once(parser), func(t1 []T) Parser[[]T] {
		return Select(XMany(parser), func(ts []T) []T {
			return internal.Union(t1, ts)
		})
	})
}

// Parse end-of-input.
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

// Take the result of parsing, and project it onto a different domain.
func Select[T any, U any](parser Parser[T], convert func(T) U) Parser[U] {
	return Then(parser, func(v T) Parser[U] {
		return Return[U](convert(v))
	})
}

// Parse the token, embedded in any amount of whitespace characters.
func Token[T any](parser Parser[T]) Parser[T] {
	return Then(Many(WhiteSpace()), func(_ []rune) Parser[T] {
		return Then(parser, func(v T) Parser[T] {
			return Then(Many(WhiteSpace()), func(_ []rune) Parser[T] {
				return Return[T](v)
			})
		})
	})
}

// Succeed immediately and return value.
func Return[T any](value T) Parser[T] {
	return func(input ParserInput) ParserResult[T] {
		return NewSuccessResult[T](value, input)
	}
}

// Version of Return with simpler inline syntax.
func ReturnValue[T, U any](parser Parser[T], value U) Parser[U] {
	return Select(parser, func(T) U {
		return value
	})
}

// Convert a stream of characters to a string.
func Text(parser Parser[[]rune]) Parser[string] {
	return Select(parser, func(rs []rune) string {
		return string(rs)
	})
}

// Parse first, if it succeeds, return first, otherwise try second.
func Or[T any](first Parser[T], second Parser[T]) Parser[T] {
	return func(input ParserInput) ParserResult[T] {
		var fr = first(input)
		if !fr.Succeeded {
			return IfFailure(second(input), func(sf ParserResult[T]) ParserResult[T] {
				return determineBestError(fr, sf)
			})
		}

		if fr.Remainder.Equal(input) {
			return IfFailure(second(input), func(sf ParserResult[T]) ParserResult[T] {
				return fr
			})
		}

		return fr
	}
}

// Parse first, if it succeeds, return first, otherwise try second.
func XOr[T any](first Parser[T], second Parser[T]) Parser[T] {
	return func(input ParserInput) ParserResult[T] {
		var fr = first(input)
		if !fr.Succeeded {

			// The 'X' part
			if !fr.Remainder.Equal(input) {
				return fr
			}

			return IfFailure(second(input), func(sf ParserResult[T]) ParserResult[T] {
				return determineBestError(fr, sf)
			})
		}

		if fr.Remainder.Equal(input) {
			return IfFailure(second(input), func(sf ParserResult[T]) ParserResult[T] {
				return fr
			})
		}

		return fr
	}
}

func determineBestError[T any](firstFailure ParserResult[T], secondFailure ParserResult[T]) ParserResult[T] {
	if secondFailure.Remainder.Position() > firstFailure.Remainder.Position() {
		return secondFailure
	}

	if secondFailure.Remainder.Position() == firstFailure.Remainder.Position() {
		unionFailure := NewFailureResult[T](
			firstFailure.Remainder,
			firstFailure.Message,
			internal.Union(firstFailure.Expectations, secondFailure.Expectations))
		return unionFailure
	}

	return firstFailure
}

// Names part of the grammar for help with error messages.
func SetExpectationIfError[T any](parser Parser[T], expectation string) Parser[T] {
	return func(input ParserInput) ParserResult[T] {
		return IfFailure(parser(input), func(f ParserResult[T]) ParserResult[T] {
			if f.Remainder.Equal(input) {
				return NewFailureResult[T](f.Remainder, f.Message, []string{expectation})
			}
			return f
		})
	}
}

// Parse a stream of elements containing only one item.
func Once[T any](parser Parser[T]) Parser[[]T] {
	return Select(parser, func(r T) []T {
		return []T{r}
	})
}

// Concatenate two streams of elements.
func Concat[T any](first, second Parser[[]T]) Parser[[]T] {
	return Then(first, func(fr []T) Parser[[]T] {
		return Select(second, func(sr []T) []T {
			return internal.Union(fr, sr)
		})
	})
}

// Attempt parsing only if the except parser fails.
func Except[T, U any](parser Parser[T], except Parser[U]) Parser[T] {
	return func(input ParserInput) ParserResult[T] {
		r := except(input)
		if r.Succeeded {
			return NewFailureResult[T](input, "Excepted parser succeeded.", []string{"other than the excepted input"})
		}
		return parser(input)
	}
}

// Parse a sequence of items until a terminator is reached.
// Returns the sequence, discarding the terminator.
func Until[T, U any](parser Parser[T], until Parser[U]) Parser[[]T] {
	return Then(Many(Except(parser, until)), func(v []T) Parser[[]T] {
		return ReturnValue(until, v)
	})
}

// Succeed if the parsed value matches predicate.
func Where[T any](parser Parser[T], predicate func(T) bool) Parser[T] {
	return func(input ParserInput) ParserResult[T] {
		return IfSuccess(parser(input), func(s ParserResult[T]) ParserResult[T] {
			if predicate(s.Value) {
				return s
			}

			return NewFailureResult[T](input, fmt.Sprintf("Unexpected %v", s.Value), []string{})
		})
	}
}

// Monadic combinator Then, adapted for Linq comprehension syntax.
func SelectMany[T, U, V any](parser Parser[T], selector func(T) Parser[U], projector func(T, U) V) Parser[V] {
	return Then(parser, func(t T) Parser[V] {
		return Select(selector(t), func(u U) V {
			return projector(t, u)
		})
	})
}

// Parse a number.
func Number() Parser[string] {
	return Text(AtLeastOnce(Numeric()))
}
