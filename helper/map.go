package helper

func Map[T any, U any](a []T, f func(T) U) []U {
	b := make([]U, len(a))
	for i, v := range a {
		b[i] = f(v)
	}
	return b
}

func RuneToString(r rune) string {
	return string(r)
}
