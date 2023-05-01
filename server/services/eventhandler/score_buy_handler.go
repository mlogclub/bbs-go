package eventhandler

import (
	"bbs-go/pkg/event"
	"bbs-go/pkg/msg"
	"bbs-go/services"
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.ScorePayEvent{}), handleScorePayCreate)
}

func handleScorePayCreate(i interface{}) {
	e := i.(event.ScorePayEvent)

	if e.ToUserId <= 0 {
		logrus.Warnf("消息发送失败, to = [%s]", e.ToUserId)
		return
	}
	quoteContent := fmt.Sprintf("积分充值到账: [%d]", e.Score)
	services.MessageService.SendMsg(0, e.ToUserId,
		msg.TypePayScore,
		"积分变更",
		"",
		quoteContent,
		nil,
	)
}
