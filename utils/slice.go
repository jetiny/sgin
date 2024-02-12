package utils

import (
	"io/fs"
)

type SliceExchangRunner[T any, V any] func(T) V
type SliceFilterRunner[T any] func(T) bool

func FiName(fi fs.FileInfo) string {
	return fi.Name()
}

func SliceExchange[T any, V any](s []T, exchanger SliceExchangRunner[T, V]) []V {
	result := make([]V, len(s))
	for i, v := range s {
		result[i] = exchanger(v)
	}
	return result
}

func SliceFilter[T any](s []T, filter SliceFilterRunner[T]) []T {
	result := make([]T, 0, len(s))
	for _, v := range s {
		if filter(v) {
			result = append(result, v)
		}
	}
	return result
}

type AnySlice[T any, V any] []T

func (s AnySlice[T, V]) Exchange(exchanger SliceExchangRunner[T, V]) []V {
	return SliceExchange(s, exchanger)
}

func (s AnySlice[T, V]) Filter(filter SliceFilterRunner[T]) AnySlice[T, V] {
	return SliceFilter(s, filter)
}
