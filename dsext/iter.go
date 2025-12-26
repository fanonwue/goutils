package dsext

import (
	"iter"
	"strings"
)

func MapSeq[T, U any](iterator iter.Seq[T], f func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for v := range iterator {
			mapped := f(v)
			if !yield(mapped) {
				return
			}
		}
	}
}

func FilterSeq[T any](seq iter.Seq[T], test func(T) bool) (ret iter.Seq[T]) {
	return func(yield func(T) bool) {
		for v := range seq {
			if test(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func JoinSeq[T any](values iter.Seq[T], sep string, transform func(T) string) string {
	var stringified []string
	for v := range values {
		stringified = append(stringified, transform(v))
	}
	return strings.Join(stringified, sep)
}
