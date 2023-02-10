package utils

import "strconv"

type DeleteStatus int8

const (
	DeleteStatusDefault DeleteStatus = 0
	DeleteStatusDeleted DeleteStatus = 1
	DeleteStatusDroped  DeleteStatus = 2
)

func SqlNotDeleted() string {
	return "deleted_status = " + strconv.Itoa(int(DeleteStatusDefault))
}

func SqlDeleted() string {
	return "deleted_status = " + strconv.Itoa(int(DeleteStatusDeleted))
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

type SqlFieldType string

func Field(field string) SqlFieldType {
	return SqlFieldType(field)
}

func (s SqlFieldType) Eq(value any) (string, any) {
	return string(s) + " = ? ", value
}

func (s SqlFieldType) Neq(value any) (string, any) {
	return string(s) + " <> ? ", value
}

func (s SqlFieldType) Gt(value any) (string, any) {
	return string(s) + " > ? ", value
}

func (s SqlFieldType) Gte(value any) (string, any) {
	return string(s) + " >= ? ", value
}

func (s SqlFieldType) Lt(value any) (string, any) {
	return string(s) + " < ? ", value
}

func (s SqlFieldType) Lte(value any) (string, any) {
	return string(s) + " <= ? ", value
}

func (s SqlFieldType) NotIsNull() string {
	return string(s) + " not is null "
}

func (s SqlFieldType) IsNull() string {
	return string(s) + " is null "
}

func (s SqlFieldType) Like(value any) (string, any) {
	return string(s) + " like ? ", value
}

func (s SqlFieldType) Prefix(value any) (string, any) {
	return "locate(?, " + string(s) + ") = 1", value
}

// xorm.OrderBy field desc
func (s SqlFieldType) OrderDesc() string {
	return "`" + string(s) + "` desc "
}

// xorm.OrderBy field asc
func (s SqlFieldType) OrderAsc() string {
	return "`" + string(s) + "` asc "
}
