package sliceutil

import "errors"

func Pop[T any](s []T) (*T, []T, error) {
	if len(s) == 0 {
		return nil, nil, errors.New("invalid operation: given slice is empty")
	}
	return &s[0], s[1:], nil
}

func Push[T any](s []T, e T) []T {
	return append(s, e)
}

func Unpush[T any](s []T) (*T, []T, error) {
	if len(s) == 0 {
		return nil, nil, errors.New("invalid operation: given slice is empty")
	}
	n := len(s) - 1
	return &s[n], s[:n], nil
}

func Map[T any, U any](s []T, mapper func(t T) U) []U {
	res := make([]U, 0, len(s))
	for _, e := range s {
		res = append(res, mapper(e))
	}
	return res
}

func Filter[T any](s []T, filter func(t T) bool) []T {
	res := make([]T, 0, len(s))
	for _, e := range s {
		if filter(e) {
			res = append(res, e)
		}
	}
	return res
}
