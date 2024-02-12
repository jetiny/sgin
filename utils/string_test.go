package utils

import (
	"fmt"
	"testing"
)

func TestStringPad(t *testing.T) {
	fmt.Println(StrPad("abc", 5, "0"))
	fmt.Println(StrPad("abcdefg", 5, "0"))
	fmt.Println(StrPadPrefix("abc", 5, "0"))
	fmt.Println(StrPadPrefix("abcdefg", 5, "0"))
}
