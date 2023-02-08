package utils

import "github.com/pkg/errors"

func WithRecover(fn func() error) (err error) {
	err = nil
	defer func() {
		if re := recover(); re != nil {
			if e, ok := re.(error); ok {
				err = e
			} else {
				err = errors.Errorf("withRecover: %+v", re)
			}
		}
	}()
	err = fn()
	return err
}
