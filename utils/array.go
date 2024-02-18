package utils

import (
	"fmt"
	"strings"
)

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

//-------

type SimpleType interface {
	~string | ~bool | ~int8 | ~int16 | ~int32 | ~int64 | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~int
}

type SimpleArray[T SimpleType] []T

func (arr SimpleArray[T]) Join(sep string) string {
	list := []string{}
	for _, v := range arr {
		list = append(list, fmt.Sprintf("%v", v))
	}
	return strings.Join(list, sep)
}

func (arr SimpleArray[T]) Contains(e T) bool {
	for _, v := range arr {
		if v == e {
			return true
		}
	}
	return false
}

func (arr SimpleArray[T]) IndexOf(e T) int {
	for i, v := range arr {
		if v == e {
			return i
		}
	}
	return -1
}

func (arr SimpleArray[T]) LastIndexOf(e T) int {
	for i := len(arr) - 1; i >= 0; i-- {
		if arr[i] == e {
			return i
		}
	}
	return -1
}

func (arr *SimpleArray[T]) Push(e ...T) {
	*arr = append(*arr, e...)
}

func (arr *SimpleArray[T]) Unshift(e ...T) {
	*arr = append(e, *arr...)
}

func (arr *SimpleArray[T]) Insert(i int, e ...T) {
	n := len(*arr)
	if i > n-1 {
		return
	}
	// 构造需要插入元素的切片
	inserts := SimpleArray[T]{}
	inserts = append(inserts, e...)

	// 重新构造切片数组
	result := SimpleArray[T]{}
	result = append(result, (*arr)[:i]...)
	result = append(result, inserts...)
	result = append(result, (*arr)[i:]...)
	*arr = result
}

func (arr *SimpleArray[T]) Concat(next SimpleArray[T]) {
	*arr = append(*arr, next...)
}

func (arr SimpleArray[T]) ForEach(f func(e T)) {
	for _, v := range arr {
		f(v)
	}
}

func (arr SimpleArray[T]) Filter(f func(T) bool) SimpleArray[T] {
	result := SimpleArray[T]{}
	for _, v := range arr {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

func (arr *SimpleArray[T]) Remove(i int) {
	*arr = append((*arr)[:i], (*arr)[i+1:]...)
}

func (arr *SimpleArray[T]) Delete(e T) bool {
	idx := arr.IndexOf(e)
	if idx != -1 {
		arr.Remove(idx)
		return true
	}
	return false
}

func (arr SimpleArray[T]) DeleteMany(e T) int {
	n := 0
	for i := len(arr) - 1; i >= 0; i-- {
		if arr[i] == e {
			arr.Remove(i)
			n++
		}
	}
	return n
}

func (arr SimpleArray[T]) Map(f func(T) T) SimpleArray[T] {
	result := SimpleArray[T]{}
	for _, v := range arr {
		result = append(result, f(v))
	}
	return result
}
