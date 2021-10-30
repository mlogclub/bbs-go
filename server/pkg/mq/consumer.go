package mq

import (
	"errors"
	"github.com/gookit/event"
)

func handleMsg(msgBytes []byte) error {
	body := string(msgBytes)
	et, e, err := ParseEvent(body)
	if err != nil {
		return err
	}

	err, _ = em.Fire(string(et), map[string]interface{}{
		"type": et,
		"data": e,
	})
	return err
}

func AddEventHandler(t EventType, handle func(e interface{}) error) {
	em.On(string(t), event.ListenerFunc(func(e event.Event) error {
		et := e.Get("type").(EventType)
		if et != t {
			return errors.New("消息类型不匹配...")
		}
		return handle(e.Get("data"))
	}))
}
