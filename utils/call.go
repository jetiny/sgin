package utils

import "github.com/go-errors/errors"

type FuncNoop func()
type FuncType[T any] func() T
type FuncOne func() error
type FuncTwo[T any] func() (T, error)
type FuncThree[T any, V any] func() (T, V, error)

func CallNoop(f FuncNoop) error {
	return Call(func() error {
		f()
		return nil
	})
}

func CallType[T any](f FuncType[T]) (T, error) {
	return Call2(func() (T, error) {
		return f(), nil
	})
}

func Call3[T any, V any](f FuncThree[T, V]) (res T, res2 V, err error) {
	defer func() {
		if r := recover(); r != nil {
			if ref, ok := r.(error); ok {
				err = ref
			} else {
				err = errors.Wrap(r, 1)
			}
		}
	}()
	if f == nil {
		panic(errors.New("function is nil"))
	}
	return f()
}

func Call2[T any](f func() (T, error)) (res T, err error) {
	defer func() {
		if r := recover(); r != nil {
			if ref, ok := r.(error); ok {
				err = ref
			} else {
				err = errors.Wrap(r, 1)
			}
		}
	}()
	if f == nil {
		panic(errors.New("function is nil"))
	}
	return f()
}

func Call(f FuncOne) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if ref, ok := r.(error); ok {
				err = ref
			} else {
				err = errors.Wrap(r, 1)
			}
		}
	}()
	if f == nil {
		panic(errors.New("function is nil"))
	}
	return f()
}

func Wrap3[T any, V any](f FuncThree[T, V]) func() (T, V, error) {
	return func() (T, V, error) {
		return Call3(f)
	}
}

func Wrap2[T any](f FuncTwo[T]) func() (T, error) {
	return func() (T, error) {
		return Call2(f)
	}
}

func Wrap(f FuncOne) func() error {
	return func() error {
		return Call(f)
	}
}
