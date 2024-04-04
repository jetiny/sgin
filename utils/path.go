package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func GetExecutablePath(useCwd bool) string {
	if useCwd {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		return cwd
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	if strings.Contains(dir, "go-build") {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		return cwd
	}
	return dir
}
