package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlice(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	type user struct {
		age int
	}
	users := AnySlice[int, user](s).Exchange(func(i int) user {
		return user{age: i}
	})
	b := AnySlice[user, int](users).Exchange(func(u user) int {
		return u.age
	})
	user1 := AnySlice[int, user](s).Filter(func(i int) bool {
		return i == 1
	}).Exchange(func(i int) user {
		return user{age: i}
	})
	assert.NotEqual(t, users, user1)
	assert.Equal(t, s, b)
}
