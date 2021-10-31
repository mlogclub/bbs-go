package eventhandler

import (
	"bbs-go/model"
	"bbs-go/model/constants"
	"bbs-go/pkg/event"
	"bbs-go/services"
	"github.com/mlogclub/simple/date"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.FollowEvent{}), HandleFollow)
}

func HandleFollow(e interface{}) {
	handleUserFeedOnFollow(e.(event.FollowEvent))
}

// 将该用户下的帖子添加到信息流
func handleUserFeedOnFollow(evt event.FollowEvent) {
	services.TopicService.ScanByUser(evt.OtherId, func(topics []model.Topic) {
		for _, topic := range topics {
			if topic.Status != constants.StatusOk {
				continue
			}
			_ = services.UserFeedService.Create(&model.UserFeed{
				UserId:     evt.UserId,
				DataType:   constants.EntityTopic,
				DataId:     topic.Id,
				AuthorId:   topic.UserId,
				CreateTime: date.NowTimestamp(),
			})
		}
	})
}
