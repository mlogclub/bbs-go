package eventhandler

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/event"
	"bbs-go/internal/pkg/msg"
	"bbs-go/internal/services"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.UserLikeEvent{}), handleUserLike)
	event.RegHandler(reflect.TypeOf(event.UserUnLikeEvent{}), handleUserUnLike)
}

func handleUserLike(i interface{}) {
	e := i.(event.UserLikeEvent)

	if e.EntityType == constants.EntityTopic {
		sendTopicLikeMsg(e.EntityId, e.UserId)
	} else if e.EntityType == constants.EntityComment {
		// TODO
	}
}

func handleUserUnLike(i interface{}) {
	e := i.(event.UserUnLikeEvent)
	if e.EntityType == constants.EntityTopic {
		// TODO
	}
}

// sendTopicLikeMsg 话题收到点赞
func sendTopicLikeMsg(topicId, likeUserId int64) {
	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOk {
		return
	}
	if topic.UserId == likeUserId {
		return
	}
	var (
		from         = likeUserId
		to           = topic.UserId
		title        = "点赞了你的话题"
		quoteContent = "《" + topic.GetTitle() + "》"
	)
	services.MessageService.SendMsg(from, to, msg.TypeTopicLike, title, "", quoteContent,
		&msg.TopicLikeExtraData{
			TopicId:    topicId,
			LikeUserId: likeUserId,
		})
}
