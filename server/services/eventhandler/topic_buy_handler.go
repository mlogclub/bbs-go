package eventhandler

import (
	"bbs-go/pkg/event"
	"bbs-go/pkg/msg"
	"bbs-go/services"
	"github.com/sirupsen/logrus"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.TopicBuyEvent{}), handleTopicBuyCreate)
}

func handleTopicBuyCreate(i interface{}) {
	e := i.(event.TopicBuyEvent)

	var (
		from = e.UserId
		to   = e.ToUserId
	)
	if from == to {
		return
	}
	if to <= 0 {
		logrus.Errorf("消息发送失败, to = [%s]", to)
		return
	}

	services.MessageService.SendMsg(from, to,
		msg.TypeBuyHidContent,
		"购买了你的隐藏内容",
		"",
		"《"+e.QuoteContent+"》",
		nil,
	)
}
