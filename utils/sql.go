package utils

import "strconv"

type DeleteStatus int8

const (
	DeleteStatusDefault DeleteStatus = 0
	DeleteStatusDeleted DeleteStatus = 1
	DeleteStatusDroped  DeleteStatus = 2
)

func SqlNotDeleted() string {
	return "deleted_stauts = " + strconv.Itoa(int(DeleteStatusDefault))
}

func SqlDeleted() string {
	return "deleted_stauts = " + strconv.Itoa(int(DeleteStatusDeleted))
}

func SqlWhereGt(field string) string {
	return field + " > ?"
}

func SqlWhereEq(field string) string {
	return field + " = ?"
}

func SqlWhereIsNull(field string) string {
	return field + " is null"
}

func SqlColAs(field string, field2 string) string {
	return field + " as " + field2
}

func SqlWhereLike(field string) string {
	return field + " like ?"
}

func SqlWherePrefix(field string) string {
	return "locate(?, " + field + ") = 1"
}
