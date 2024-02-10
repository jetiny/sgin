package utils

import (
	"errors"
	"testing"

	gerror "github.com/go-errors/errors"
)

func TestStringError(t *testing.T) {
	_, err := Call2(func() (*any, error) {
		panic("test error")
	})
	if err == nil {
		t.Error("error should not be nil")
	}
}

func TestIntError(t *testing.T) {
	_, err := Call2(func() (*any, error) {
		panic(1)
	})
	if e, ok := err.(*gerror.Error); !ok || e.Error() != "1" {
		t.Error("error should not be nil")
	}
}

func TestErrorError(t *testing.T) {
	_, err := Call2(func() (*any, error) {
		panic(errors.ErrUnsupported)
	})
	if err != errors.ErrUnsupported {
		t.Error("error should be ErrUnsupported")
	}
}

func TestClientError(t *testing.T) {
	_, err := Call2(func() (*any, error) {
		BadRequest.Panic()
		return nil, nil
	})
	if err != BadRequest {
		t.Error("error should be ErrUnsupported")
	}
}

func TestWrap3(t *testing.T) {
	f := Wrap3(func() (*any, *any, error) {
		defer func() {
			panic(errors.New("test error"))
		}()
		return nil, nil, nil
	})
	_, _, err := f()
	if err == nil {
		t.Error("error should not be nil")
	}
}

func TestWrap2(t *testing.T) {
	f := Wrap2(func() (*any, error) {
		defer func() {
			panic(errors.New("test error"))
		}()
		return nil, nil
	})
	_, err := f()
	if err == nil {
		t.Error("error should not be nil")
	}
}

func TestWrap(t *testing.T) {
	f := Wrap(func() error {
		defer func() {
			panic(errors.New("test error"))
		}()
		return nil
	})
	err := f()
	if err == nil {
		t.Error("error should not be nil")
	}
}

func TestCall(t *testing.T) {
	err := Call(func() error {
		defer func() {
			panic(errors.New("test error"))
		}()
		return nil
	})
	if err == nil {
		t.Error("error should not be nil")
	}
}

func TestCall2(t *testing.T) {
	_, err := Call2(func() (*any, error) {
		defer func() {
			panic(errors.New("test error"))
		}()
		return nil, nil
	})
	if err == nil {
		t.Error("error should not be nil")
	}
}

func TestCall3(t *testing.T) {
	{
		_, _, err := Call3(func() (*any, *any, error) {
			defer func() {
				panic(errors.New("test error"))
			}()
			return nil, nil, nil
		})
		if err == nil {
			t.Error("error should not be nil")
		}
	}
	{
		v1, v2, err := Call3(func() (*int, *int, error) {
			n := 10
			return &n, &n, nil
		})
		if err != nil {
			t.Error("error should be nil")
		}
		if *v1 != 10 || *v2 != 10 {
			t.Error("v1 and v2 should be 10")
		}
	}
}

func TestCallNoop(t *testing.T) {
	err := CallNoop(func() {
		panic(errors.New("test error"))
	})
	if err == nil {
		t.Error("error should not be nil")
	}
}

func TestCallType(t *testing.T) {
	{
		_, err := CallType(func() *any {
			defer func() {
				panic(errors.New("test error"))
			}()
			return nil
		})
		if err == nil {
			t.Error("error should not be nil")
		}
	}
	{
		n, err := CallType(func() *int {
			n := 10
			return &n
		})
		if err != nil {
			t.Error("error should not not be nil")
		}
		if *n != 10 {
			t.Error("n should be 10")
		}
	}
}
