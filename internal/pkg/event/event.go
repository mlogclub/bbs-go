package event

import (
	"log/slog"
	"reflect"
	"sync"

	"github.com/panjf2000/ants/v2"
)

var (
	m         sync.RWMutex
	eventPool *ants.PoolWithFunc
	handlers  map[reflect.Type][]func(i any)
	// wg        sync.WaitGroup
)

func init() {
	var err error
	eventPool, err = ants.NewPoolWithFunc(4, dispatch, ants.WithMaxBlockingTasks(1000))
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	}
	handlers = make(map[reflect.Type][]func(i any))
}

func dispatch(i any) {
	handlerList := getHandlerList(i)
	if len(handlerList) == 0 {
		return
	}
	for _, handler := range handlerList {
		handler(i)
		// wg.Done()
	}
}

func Send(e any) {
	if err := eventPool.Invoke(e); err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	} else {
		// wg.Add(len(getHandlerList(e)))
		// wg.Wait()
	}
}

func RegHandler(t reflect.Type, handler func(i any)) {
	m.Lock()
	defer m.Unlock()

	handlerList := handlers[t]
	handlerList = append(handlerList, handler)
	handlers[t] = handlerList
}

func getHandlerList(i any) []func(i any) {
	m.RLock()
	defer m.RUnlock()

	t := reflect.TypeOf(i)
	handlerList, ok := handlers[t]
	if ok {
		return handlerList
	} else {
		slog.Error("没找到任务处理器", slog.String("type", t.String()))
		return nil
	}
}
