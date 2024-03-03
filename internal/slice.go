package internal

func Map[T any, U any](a []T, f func(T) U) []U {
	b := make([]U, len(a))
	for i, v := range a {
		b[i] = f(v)
	}
	return b
}

func Union[T any](x, y []T) []T {
	r := make([]T, len(x)+len(y))
	copy(r, x)
	copy(r[len(x):], y)
	return r
}
