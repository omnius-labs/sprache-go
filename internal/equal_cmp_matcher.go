package internal

// ref. https://github.com/KamikazeZirou/equal-cmp/blob/main/equal_cmp.go

import (
	"errors"

	"github.com/google/go-cmp/cmp"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

func EqualCmp(expected interface{}, options ...cmp.Option) types.GomegaMatcher {
	return &equalCmpMatcher{
		expected: expected,
		options:  options,
	}
}

type equalCmpMatcher struct {
	expected interface{}
	options  cmp.Options
}

func (matcher *equalCmpMatcher) Match(actual interface{}) (success bool, err error) {
	if actual == nil && matcher.expected == nil {
		return false, errors.New("refusing to compare <nil> to <nil>.\nbe explicit and use BeNil() instead. this is to avoid mistakes where both sides of an assertion are erroneously uninitialized")
	}
	return cmp.Equal(actual, matcher.expected, matcher.options), nil
}

func (matcher *equalCmpMatcher) FailureMessage(actual interface{}) (message string) {
	diff := cmp.Diff(actual, matcher.expected, matcher.options)
	return format.Message(actual, "to equal", matcher.expected) +
		"\n\nDiff:\n" + format.IndentString(diff, 1)
}

func (matcher *equalCmpMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	diff := cmp.Diff(actual, matcher.expected, matcher.options)
	return format.Message(actual, "not to equal", matcher.expected) +
		"\n\nDiff:\n" + format.IndentString(diff, 1)
}
