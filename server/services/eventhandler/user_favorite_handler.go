package eventhandler

import (
	"bbs-go/model/constants"
	"bbs-go/pkg/event"
	"bbs-go/services"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.UserFavoriteEvent{}), handleUserFavorite)
}

func handleUserFavorite(i interface{}) {
	e := i.(event.UserFavoriteEvent)

	if e.EntityType == constants.EntityTopic {
		services.MessageService.SendTopicFavoriteMsg(e.EntityId, e.UserId)
	} else if e.EntityType == constants.EntityArticle {
		// TODO
	}
}
