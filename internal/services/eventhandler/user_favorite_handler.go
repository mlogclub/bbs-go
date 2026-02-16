package eventhandler

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/event"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/pkg/msg"
	"bbs-go/internal/services"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.UserFavoriteEvent{}), handleUserFavorite)
}

func handleUserFavorite(i interface{}) {
	e := i.(event.UserFavoriteEvent)

	if e.EntityType == constants.EntityTopic {
		sendTopicFavoriteMsg(e.EntityId, e.UserId)
	} else if e.EntityType == constants.EntityArticle {
		// TODO
	}
}

// sendTopicFavoriteMsg 话题被收藏
func sendTopicFavoriteMsg(topicId, favoriteUserId int64) {
	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOk {
		return
	}
	if topic.UserId == favoriteUserId {
		return
	}
	var (
		from         = favoriteUserId
		to           = topic.UserId
		title        = locales.Get("message.topic_favorite_msg_title")
		quoteContent = "《" + topic.GetTitle() + "》"
	)
	services.MessageService.SendMsg(from, to, msg.TypeTopicFavorite, title, "", quoteContent,
		&msg.TopicFavoriteExtraData{
			TopicId:        topicId,
			FavoriteUserId: favoriteUserId,
		})
}
