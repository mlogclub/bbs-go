package mq

import (
	"github.com/sirupsen/logrus"
)

func Send(typ EventType, event interface{}) {
	body, err := NewEventMsg(typ, event)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = sendText(body)
	if err != nil {
		logrus.Error(err)
	}
}

func sendText(body string) error {
	return handleMsg([]byte(body))
}
