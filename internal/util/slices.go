package util

func RotateSliceLeft[T any](s []T) []T {
	return append(s[1:], s[0])
}

func RotateSliceRight[T any](s []T) []T {
	return append([]T{s[len(s)-1]}, s[:len(s)-1]...)
}
