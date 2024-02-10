package utils

import (
	"fmt"
	"testing"
)

func TestRandom(t *testing.T) {
	s := RandomChar(10)
	fmt.Println(s)
	if len(s) != 10 {
		t.Error("RandomChar(10) should return 10 characters")
	}
	s = RandomNumber(10)
	if len(s) != 10 {
		t.Error("RandomNumber(10) should return 10 characters")
	}
	s = RandomString(10)
	if len(s) != 10 {
		t.Error("RandomString(10) should return 10 characters")
	}
	s = RandomCharset(10, "abc")
	if len(s) != 10 {
		t.Error("RandomCharset(10, \"abc\") should return 10 characters")
	}
	s = RandomCharset(1, "我是中文")
	fmt.Println(s)
	if len([]rune(s)) != 1 {
		t.Error("RandomCharset(1, \"我是中文\") should return 1 character")
	}
	s = RandomCharset(10, "我是中文")
	fmt.Println(s)
	if len([]rune(s)) != 10 {
		t.Error("RandomCharset(1, \"我是中文\") should return 10 character")
	}
}
