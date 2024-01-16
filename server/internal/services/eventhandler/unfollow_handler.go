package eventhandler

import (
	"bbs-go/internal/pkg/event"
	"bbs-go/internal/services"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.UnFollowEvent{}), handleUnFollowEvent)
}

func handleUnFollowEvent(i interface{}) {
	e := i.(event.UnFollowEvent)

	// 清理该用户下的信息流
	services.UserFeedService.DeleteByUser(e.UserId, e.OtherId)
}
