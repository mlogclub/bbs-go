package eventhandler

import (
	"bbs-go/model/constants"
	"bbs-go/pkg/event"
	"bbs-go/services"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.UserLikeEvent{}), handleUserLike)
}

func handleUserLike(i interface{}) {
	e := i.(event.UserLikeEvent)

	if e.EntityType == constants.EntityTopic {
		services.MessageService.SendTopicLikeMsg(e.EntityId, e.UserId)
	} else if e.EntityType == constants.EntityComment {
		// TODO
	}
}
