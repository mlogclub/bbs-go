package eventhandler

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/event"
	"bbs-go/internal/pkg/msg"
	"bbs-go/internal/repositories"
	"bbs-go/internal/services"
	"reflect"

	"github.com/mlogclub/simple/sqls"
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
		title              = "你的话题被删除"
		quoteContent       = "《" + topic.GetTitle() + "》"
	)
	services.MessageService.SendMsg(from, to, msg.TypeTopicDelete, title, "", quoteContent,
		&msg.TopicDeleteExtraData{
			TopicId:      topicId,
			DeleteUserId: deleteUserId,
		},
	)
}
