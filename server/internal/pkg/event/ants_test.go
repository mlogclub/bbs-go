package event

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"sync"
	"testing"
)

func TestType(t *testing.T) {
	//e1 := FollowEvent{}
	//e2 := &FollowEvent{}
	//fmt.Println("type:", reflect.TypeOf(e1))
	//fmt.Println("type:", reflect.TypeOf(e2))
	//
	//reflect.TypeOf(e1)

	//var e1 interface{}
	//
	//e1 = &mq.FollowEvent{}
	//
	//switch e1.(type) {
	//case string:
	//	//...
	//case int:
	//	//...
	//case mq.FollowEvent:
	//case *mq.FollowEvent:
	//	fmt.Println("FollowEvent!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	//}
	//return
}

func TestAnts(t *testing.T) {
	var wg sync.WaitGroup
	runTimes := 1000

	// Use the pool with a function,
	// set 10 to the capacity of goroutine pool and 1 second for expired duration.
	p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
		n := i.(int32)
		fmt.Printf("run with %d\n", n)
		wg.Done()
	})
	defer p.Release()

	// Submit tasks one by one.
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = p.Invoke(int32(i))
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", p.Running())
}
