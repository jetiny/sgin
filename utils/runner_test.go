package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestRunner(t *testing.T) {
	r := NewRunner(func(a any) {
		time.Sleep(time.Microsecond * 100)
		fmt.Println(a)
	})
	{
		r.Start()
		for i := 0; i < 10; i++ {
			fmt.Println("push", i)
			r.Push(i)
		}
		fmt.Println("close")
		r.Wait()
		fmt.Println("done")
	}
	time.Sleep(time.Second * 1)
	{
		r.Start()
		fmt.Println("push list")
		r.PushList(1, 2, 3, 4, 5, 6, 7, 8, 9)
		fmt.Println("close")
		r.Wait()
		fmt.Println("done")
	}
}
