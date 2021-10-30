package eventhandler

import (
	"bbs-go/pkg/mq"
	"bbs-go/services"
)

func init() {
	mq.AddEventHandler(mq.EventTypeUnFollow, HandleUnFollow)
}

func HandleUnFollow(e interface{}) error {
	event := e.(*mq.UnFollowEvent)

	// 清理该用户下的信息流
	services.UserFeedService.DeleteByUser(event.UserId, event.OtherId)

	return nil
}
