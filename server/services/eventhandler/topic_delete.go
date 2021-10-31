package eventhandler

import (
	"bbs-go/model/constants"
	"bbs-go/pkg/event"
	"bbs-go/services"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.TopicDeleteEvent{}), HandleTopicDelete)
}

func HandleTopicDelete(e interface{}) {
	evt := e.(event.TopicDeleteEvent)

	// 处理userFeed
	services.UserFeedService.DeleteByDataId(evt.TopicId, constants.EntityTopic)

	// 发送消息
	services.MessageService.SendTopicDeleteMsg(evt.TopicId, evt.DeleteUserId)

	// 操作日志
	services.OperateLogService.AddOperateLog(evt.DeleteUserId, constants.OpTypeDelete, constants.EntityTopic,
		evt.TopicId, "", nil)
}
