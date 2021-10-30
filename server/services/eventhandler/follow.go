package eventhandler

import (
	"bbs-go/model"
	"bbs-go/model/constants"
	"bbs-go/pkg/mq"
	"bbs-go/services"
	"github.com/mlogclub/simple/date"
)

func init() {
	mq.AddEventHandler(mq.EventTypeFollow, HandleFollow)
}

func HandleFollow(e interface{}) error {
	event := e.(*mq.FollowEvent)
	handleUserFeedOnFollow(event)
	return nil
}

// 将该用户下的帖子添加到信息流
func handleUserFeedOnFollow(event *mq.FollowEvent) {
	services.TopicService.ScanByUser(event.OtherId, func(topics []model.Topic) {
		for _, topic := range topics {
			if topic.Status != constants.StatusOk {
				continue
			}
			_ = services.UserFeedService.Create(&model.UserFeed{
				UserId:     event.UserId,
				DataType:   constants.EntityTopic,
				DataId:     topic.Id,
				AuthorId:   topic.UserId,
				CreateTime: date.NowTimestamp(),
			})
		}
	})
}
