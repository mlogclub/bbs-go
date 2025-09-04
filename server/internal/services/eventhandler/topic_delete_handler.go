package eventhandler

import (
	"reflect"

	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/event"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/pkg/msg"
	"bbs-go/internal/repositories"
	"bbs-go/internal/services"

	"bbs-go/internal/pkg/simple/sqls"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.TopicDeleteEvent{}), handleTopicDeleteEvent)
}

func handleTopicDeleteEvent(i interface{}) {
	e := i.(event.TopicDeleteEvent)

	// 处理userFeed
	services.UserFeedService.DeleteByDataId(e.TopicId, constants.EntityTopic)

	// 发送消息
	sendTopicDeleteMsg(e.TopicId, e.DeleteUserId)

	// 操作日志
	services.OperateLogService.AddOperateLog(e.DeleteUserId, constants.OpTypeDelete, constants.EntityTopic,
		e.TopicId, "", nil)
}

// sendTopicDeleteMsg 话题被删除消息
func sendTopicDeleteMsg(topicId, deleteUserId int64) {
	topic := repositories.TopicRepository.Get(sqls.DB(), topicId)
	if topic == nil {
		return
	}
	if topic.UserId == deleteUserId {
		return
	}
	var (
		from         int64 = 0
		to                 = topic.UserId
		title              = locales.Get("message.topic_delete_msg_title")
		quoteContent       = "《" + topic.GetTitle() + "》"
	)
	services.MessageService.SendMsg(from, to, msg.TypeTopicDelete, title, "", quoteContent,
		&msg.TopicDeleteExtraData{
			TopicId:      topicId,
			DeleteUserId: deleteUserId,
		},
	)
}
