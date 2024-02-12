package utils

import (
	"cmp"
	"sort"
)

type SortCompare[T any] func(T, T) bool

func SortSlice[T any](list []T, cmp SortCompare[T]) {
	sort.Slice(list, func(i, j int) bool {
		return cmp(list[i], list[j])
	})
}

func SortSliceStable[T any](list []T, cmp SortCompare[T]) {
	sort.SliceStable(list, func(i, j int) bool {
		return cmp(list[i], list[j])
	})
}

type SortPeek[T any, V cmp.Ordered] func(T) V

type AnySort[T any, V cmp.Ordered] []T

func (s AnySort[T, V]) SortSlice(sc SortCompare[T]) {
	SortSlice(s, sc)
}

func (s AnySort[T, V]) SortSliceStable(sc SortCompare[T]) {
	SortSliceStable(s, sc)
}

func (s AnySort[T, V]) SortSlicePeek(pk SortPeek[T, V]) {
	SortSlice(s, func(t1, t2 T) bool {
		return cmp.Less[V](pk(t1), pk(t2))
	})
}

func (s AnySort[T, V]) SortSliceStablePeek(pk SortPeek[T, V]) {
	SortSliceStable(s, func(t1, t2 T) bool {
		return cmp.Less[V](pk(t1), pk(t2))
	})
}
