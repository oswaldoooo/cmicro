package delaymap_test

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/oswaldoooo/cmicro/api/delaymap"
)

func TestDelayMap(t *testing.T) {
	dm := delaymap.NewRDelayMap[string, int](func(s1, s2 *string) int {
		status := bytes.Compare([]byte(*s1), []byte(*s2))
		// fmt.Printf("compare %s %s status %d\n", *s1, *s2, status)
		return status
	})
	go dm.Run()
	dm.Set("action", 10, 0)
	dm.Set("action2", 20, time.Second*3)
	dm.Set("action3", 20, 5*time.Second)
	dm.Delete("action3")
	fmt.Printf("status: action %v action2 %v action3 %v\n", dm.Get("action"), dm.Get("action2"), dm.Get("action3"))
	// dm.SetCallBackWhenExpire("action2", func() {
	// 	fmt.Println("hello action2 called")
	// })
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println(dm.Get("action") == nil, dm.Get("action2") == nil)
}
