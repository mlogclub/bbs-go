package eventhandler

import (
	"bbs-go/model/constants"
	"bbs-go/pkg/mq"
	"bbs-go/services"
)

func init() {
	mq.AddEventHandler(mq.EventTypeTopicDelete, HandleTopicDelete)
}

func HandleTopicDelete(e interface{}) error {
	event := e.(*mq.TopicDeleteEvent)

	// 处理userFeed
	services.UserFeedService.DeleteByDataId(event.TopicId, constants.EntityTopic)

	// 发送消息
	services.MessageService.SendTopicDeleteMsg(event.TopicId, event.DeleteUserId)

	// 操作日志
	services.OperateLogService.AddOperateLog(event.DeleteUserId, constants.OpTypeDelete, constants.EntityTopic,
		event.TopicId, "", nil)

	return nil
}
