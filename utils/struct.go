package utils

import (
	"reflect"
	"strings"
)

func GetXormStructKeys(ins any) (res Array[string]) {
	return getStructKeys(ins, "xorm", func(str string) string {
		return strings.ReplaceAll(str, "'", "")
	})
}

func GetDbStructKeys(ins any) (res Array[string]) {
	return getStructKeys(ins, "db", func(str string) string {
		return strings.ReplaceAll(str, "'", "")
	})
}

func getStructKeys(ins interface{}, tagName string, filter func(str string) string) (res Array[string]) {
	typeOfCat := reflect.TypeOf(ins)
	for i := 0; i < typeOfCat.NumField(); i++ {
		fieldType := typeOfCat.Field(i)
		if tagName != "" && fieldType.Tag != "" {
			if catType, ok := typeOfCat.FieldByName(fieldType.Name); ok {
				if key, v := catType.Tag.Lookup(tagName); v {
					res = append(res, filter(key))
				} else {
					res = append(res, fieldType.Name)
				}
			}
		} else {
			res = append(res, fieldType.Name)
		}
	}
	return
}

func GetStructKeys(ins interface{}, tagName string) (res Array[string]) {
	typeOfCat := reflect.TypeOf(ins)
	for i := 0; i < typeOfCat.NumField(); i++ {
		fieldType := typeOfCat.Field(i)
		if tagName != "" && fieldType.Tag != "" {
			if catType, ok := typeOfCat.FieldByName(fieldType.Name); ok {
				if key, v := catType.Tag.Lookup(tagName); v {
					res = append(res, key)
				} else {
					res = append(res, fieldType.Name)
				}
			}
		} else {
			res = append(res, fieldType.Name)
		}
	}
	return
}
