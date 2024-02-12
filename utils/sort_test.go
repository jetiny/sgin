package utils

import (
	"cmp"
	"fmt"
	"io/fs"
	"os"
	"testing"
	"time"
)

type fileInfo struct {
	name string
}

func (f *fileInfo) Name() string       { return f.name }
func (f *fileInfo) Size() int64        { return 0 }
func (f *fileInfo) Mode() fs.FileMode  { return 0 }
func (f *fileInfo) ModTime() time.Time { return time.Time{} }
func (f *fileInfo) IsDir() bool        { return false }
func (f *fileInfo) Sys() interface{}   { return nil }

func TestSort(t *testing.T) {
	{
		arr := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}
		SortSlice[int](arr, cmp.Less[int])
		fmt.Println(arr)
	}
	{
		arr := []string{"apple_1", "apple_2", "apple_10", "apple_11", "apple_21"}
		SortSlice[string](arr, cmp.Less[string])
		fmt.Println(arr)
	}
	{
		arr := []string{"apple_001", "apple_002", "apple_010", "apple_011", "apple_021"}
		SortSlice[string](arr, cmp.Less[string])
		fmt.Println(arr)
	}
	{
		arr := []fs.FileInfo{
			&fileInfo{name: "apple_1"},
			&fileInfo{name: "apple_2"},
			&fileInfo{name: "apple_10"},
			&fileInfo{name: "apple_11"},
			&fileInfo{name: "apple_21"},
		}
		SortSliceStable[os.FileInfo](arr, func(fi1, fi2 fs.FileInfo) bool {
			return fi1.Name() < fi2.Name()
		})
		fmt.Println(SliceExchange(arr, FiName))
	}
	{
		arr := []fs.FileInfo{
			&fileInfo{name: "apple_1"},
			&fileInfo{name: "apple_2"},
			&fileInfo{name: "apple_10"},
			&fileInfo{name: "apple_11"},
			&fileInfo{name: "apple_21"},
		}
		AnySort[fs.FileInfo, string](arr).SortSlicePeek(FiName)
		fmt.Println(SliceExchange(arr, FiName))
	}
}
