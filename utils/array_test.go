package utils

import "testing"

func TestSimpleArray(t *testing.T) {
	a := SimpleArray[int]([]int{1, 2, 3})
	if !a.Contains(1) {
		t.Error("SimpleArray.Contains(1) should be true")
	}
	if !a.Delete(1) {
		t.Error("SimpleArray.Delete(1) should be true")
	}
	if a.Contains(1) {
		t.Error("SimpleArray.Contains(1) should be false after delete")
	}
	if a.Join(",") != "2,3" {
		t.Error("SimpleArray.Join(',') should be '2,3'")
	}
	if a.IndexOf(1) != -1 {
		t.Error("SimpleArray.IndexOf(1) should be -1")
	}
	if a.IndexOf(2) != 0 {
		t.Error("SimpleArray.IndexOf(2) should be 0")
	}
	if a.IndexOf(3) != 1 {
		t.Error("SimpleArray.IndexOf(3) should be 1")
	}
	if a.LastIndexOf(1) != -1 {
		t.Error("SimpleArray.LastIndexOf(1) should be -1")
	}
	if a.LastIndexOf(2) != 0 {
		t.Error("SimpleArray.LastIndexOf(2) should be 0")
	}
	if a.LastIndexOf(3) != 1 {
		t.Error("SimpleArray.LastIndexOf(3) should be 1")
	}
	a.Push(4, 4, 4)

	if a.Join(",") != "2,3,4,4,4" {
		t.Error("SimpleArray.Join(',') should be '2,3,4,4,4'")
	}
	if !a.Delete(4) {
		t.Error("SimpleArray.Delete(4) should be true")
	}
	if a.DeleteMany(4) != 2 {
		t.Error("SimpleArray.DeleteMany(4) should be 2")
	}
}
