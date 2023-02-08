package utils

import (
	"errors"
	"testing"
)

func TestRecover(t *testing.T) {
	err := WithRecover(func() error {
		panic(1)
	})
	if err == nil {
		t.Error("must has error")
	}

	err = WithRecover(func() error {
		return nil
	})
	if err != nil {
		t.Error(err)
	}

	err = WithRecover(func() error {
		return errors.New("err")
	})
	if err == nil {
		t.Error("must has error")
	}

}
