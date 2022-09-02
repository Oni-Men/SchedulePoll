package sliceutil

func Pop[T any](s []T) (T, []T) {
	return s[0], s[1:]
}

func Push[T any](s []T, e T) []T {
	return append(s, e)
}
