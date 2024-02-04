package ps

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExiftoolPool(t *testing.T) {
	pool, err := NewPool(NewExiftoolHandler(), "exiftool", 2, "-j")
	assert.NoError(t, err)
	defer pool.Stop()
	list := []string{
		"./ps_exiftool_test.go",
		"./ps_exiftool.go",
	}
	for _, v := range list {
		filename, err := filepath.Abs(v)
		assert.NoError(t, err)
		buf, err := pool.Execute(filename)
		assert.NoError(t, err)
		fmt.Println(string(buf))
	}
}
