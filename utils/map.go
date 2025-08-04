package utils

import "iter"

func Map[T, V any](input []T, fn func(T) V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range input {
			if !yield(fn(v)) {
				return
			}
		}
	}
}
