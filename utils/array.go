package utils

import "strings"

// 数组定义
type Array[T any] []T

// 向末尾添加元素
func (arr *Array[T]) Push(e ...T) {
	*arr = append(*arr, e...)
}

// 向头部添加元素
func (arr *Array[T]) Unshift(e ...T) {
	*arr = append(e, *arr...)
}

// 删除末尾元素
func (arr *Array[T]) Pop() any {
	n := len(*arr)
	if n == 0 {
		return nil
	}
	// 获取末尾元素
	v := (*arr)[n-1]
	// 切片修改
	*arr = (*arr)[:n-1]
	return v
}

// 删除数组开头元素
func (arr *Array[T]) Shift() *T {
	n := len(*arr)
	if n == 0 {
		return nil
	}
	// 获取第一个元素
	v := (*arr)[0]
	// 修改切片
	*arr = (*arr)[1:]
	return &v
}

// 任意位置插入元素
func (arr *Array[T]) Insert(i int, e ...T) {
	n := len(*arr)
	if i > n-1 {
		return
	}
	// 构造需要插入元素的切片
	inserts := Array[T]{}
	inserts = append(inserts, e...)

	// 重新构造切片数组
	result := Array[T]{}
	result = append(result, (*arr)[:i]...)
	result = append(result, inserts...)
	result = append(result, (*arr)[i:]...)
	*arr = result
}

// 删除任意位置元素
func (arr *Array[T]) Remove(i int) {
	n := len(*arr)
	if i > n-1 {
		return
	}
	*arr = append((*arr)[:i], (*arr)[i+1:]...)
}

// 数组合并
func (arr *Array[T]) Concat(next Array[T]) {
	*arr = append(*arr, next...)
}

// 迭代器
func (arr *Array[T]) ForEach(f func(e any)) {
	for _, v := range *arr {
		f(v)
	}
}

// 过滤器
func (arr *Array[T]) Filter(f func(any any) bool) Array[T] {
	result := Array[T]{}
	for _, v := range *arr {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

type ArrayType interface {
	string | int
}

func ArrayContains[T ArrayType](arr *Array[T], e T) bool {
	return ArrayIndexOf(arr, e) == -1
}

func ArrayIndexOf[T ArrayType](arr *Array[T], e T) int {
	for i, v := range *arr {
		if v == e {
			return i
		}
	}
	return -1
}

func ArrayDelete[T ArrayType](arr *Array[T], e T) bool {
	idx := ArrayIndexOf(arr, e)
	if idx != -1 {
		arr.Remove(idx)
		return true
	}
	return false
}

func ArrayDeleteMany[T ArrayType](arr *Array[T], e ...T) {
	for _, v := range e {
		ArrayDelete(arr, v)
	}
}

func ArrayJoin(arr *Array[string], sep string) string {
	return strings.Join([]string(*arr), sep)
}
