package utils

import "strconv"

func StrPad(s string, length int, pad string) string {
	if len(s) >= length {
		return s[:length]
	}
	for i := 0; i <= length-len(s); i++ {
		s += pad
	}
	return s
}

func StrPadPrefix(s string, length int, pad string) string {
	if len(s) >= length {
		return s[:length]
	}
	for i := 0; i <= length-len(s); i++ {
		s = pad + s
	}
	return s
}

func IntToStrWithPad(num int, length int, pad string) string {
	s := strconv.Itoa(num)
	return StrPad(s, length, pad)
}
