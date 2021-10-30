package mq

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

func TestProducer(t *testing.T) {
	Send(EventTypeFollow, &FollowEvent{
		UserId:  33,
		OtherId: 44,
	})
}

func TestConsumer(t *testing.T) {
	AddEventHandler(EventTypeFollow, func(e interface{}) error {
		return nil
	})
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func TestParse(t *testing.T) {
	temp := FollowEvent{
		UserId:  1,
		OtherId: 2,
	}
	typ := reflect.TypeOf(temp)
	if typ.Kind() == reflect.Struct {
		fmt.Println("case1", typ)
	} else if typ.Kind() == reflect.Ptr {
		// typ = typ.Elem()
		if typ.Elem().Kind() == reflect.Struct {
			fmt.Println("case2", typ.Elem().Name())
		}
	}
}
