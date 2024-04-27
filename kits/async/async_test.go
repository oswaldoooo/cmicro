package async_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/oswaldoooo/cmicro/kits/async"
)

func TestAsync(t *testing.T) {
	task := async.NewTask()
	st := task.Start()
	for i := 1; i < 3; i++ {
		task.AddTask(func(t time.Duration) {
			time.Sleep(t)
			fmt.Println("hello", t)
		}, time.Second*time.Duration(i))
	}
	task.GracefulClose()
	// time.Sleep(7 * time.Second)
	for {
		t, ok := <-st
		if ok {
			fmt.Println(t.Stack, t.Val)
		} else {
			break
		}
	}
}
