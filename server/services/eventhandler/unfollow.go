package eventhandler

import (
	"bbs-go/pkg/event"
	"bbs-go/services"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.UnFollowEvent{}), HandleUnFollow)
}

func HandleUnFollow(e interface{}) {
	evt := e.(event.UnFollowEvent)

	// 清理该用户下的信息流
	services.UserFeedService.DeleteByUser(evt.UserId, evt.OtherId)
}
