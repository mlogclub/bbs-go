package eventhandler

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/event"
	"bbs-go/internal/services"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.TopicUpdateEvent{}), handleTopicUpdateEvent)
}

func handleTopicUpdateEvent(i interface{}) {
	e := i.(event.TopicUpdateEvent)

	// 操作日志
	services.OperateLogService.AddOperateLog(e.UserId, constants.OpTypeUpdate, constants.EntityTopic, e.TopicId, "", nil)
}
