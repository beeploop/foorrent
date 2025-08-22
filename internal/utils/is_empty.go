package utils

func IsEmpty[T any](values []T) bool {
	return len(values) == 0
}
