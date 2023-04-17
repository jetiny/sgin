package utils

import (
	"testing"
)

func TestEnv(t *testing.T) {
	{
		c := GetterDefault("a", "s")
		n := c.Value().(string)
		if n != c.String() {
			t.Fail()
		}
	}
	{
		c := GetterDefault("b", 10)
		n := c.Value().(int)
		if n != c.Int() {
			t.Fail()
		}
	}
}
