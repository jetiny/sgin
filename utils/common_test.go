package utils

import "testing"

func TestFormatPhone(t *testing.T) {
	{
		value := FormatPhone("18911112222")
		if value != "189****2222" {
			t.Errorf(value)
		}
	}
	{
		value := FormatPhone("63582606")
		if value != "63****06" {
			t.Errorf(value)
		}
	}
	{
		value := FormatPhone("12345")
		if value != "1***5" {
			t.Errorf(value)
		}
	}
}

func TestFormatEmail(t *testing.T) {
	{
		value := FormatEmail("1234567@x.com")
		if value != "1*****7@x.com" {
			t.Errorf(value)
		}
	}
	{
		value := FormatEmail("12345678@x.com")
		if value != "12****78@x.com" {
			t.Errorf(value)
		}
	}
}

func TestSnow(t *testing.T) {
	InitSnowflake(1)
	{
		value := SnowWithPrefix("", nil)
		if len(value) != 12 {
			t.Error(value, len(value))
		}
	}
}

func TestUuid(t *testing.T) {
	{
		value := HashUuid()
		if len(value) != 32 {
			t.Error(value)
		}
	}
	{
		value := ShortUuid()
		if len(value) != 8 {
			t.Error(value)
		}
	}
}
