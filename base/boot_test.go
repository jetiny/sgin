package base

import "testing"

func TestBootLoadEnv(t *testing.T) {
	r, err := boot(BootWithEnv)
	if err != nil {
		t.Error(err)
	}
	if r == nil {
		t.Fail()
	}
}
