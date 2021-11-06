package eventhandler

import (
	"bbs-go/model/constants"
	"bbs-go/pkg/event"
	"bbs-go/services"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.TopicDeleteEvent{}), handleTopicDeleteEvent)
}

func handleTopicDeleteEvent(i interface{}) {
	e := i.(event.TopicDeleteEvent)

	// 处理userFeed
	services.UserFeedService.DeleteByDataId(e.TopicId, constants.EntityTopic)

	// 发送消息
	services.MessageService.SendTopicDeleteMsg(e.TopicId, e.DeleteUserId)

	// 操作日志
	services.OperateLogService.AddOperateLog(e.DeleteUserId, constants.OpTypeDelete, constants.EntityTopic,
		e.TopicId, "", nil)
}
