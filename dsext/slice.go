package dsext

import (
	"strings"
)

func Map[T, U any](arr []T, f func(T) U) []U {
	us := make([]U, len(arr))
	for i := range arr {
		us[i] = f(arr[i])
	}
	return us
}

func Filter[T any](arr []T, test func(T) bool) (ret []T) {
	for _, e := range arr {
		if test(e) {
			ret = append(ret, e)
		}
	}
	return
}

func Join[T any](values []T, sep string, transform func(T) string) string {
	var stringified []string
	for _, v := range values {
		stringified = append(stringified, transform(v))
	}
	return strings.Join(stringified, sep)
}
