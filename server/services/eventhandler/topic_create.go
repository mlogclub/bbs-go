package eventhandler

import (
	"bbs-go/model"
	"bbs-go/model/constants"
	"bbs-go/pkg/event"
	"bbs-go/pkg/seo"
	"bbs-go/pkg/urls"
	"bbs-go/services"
	"github.com/mlogclub/simple/date"
	"github.com/sirupsen/logrus"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.TopicCreateEvent{}), HandleTopicCreate)
}

func HandleTopicCreate(e interface{}) {
	evt := e.(event.TopicCreateEvent)

	// 百度链接推送
	seo.Push(urls.TopicUrl(evt.TopicId))

	services.UserFollowService.ScanFans(evt.UserId, func(fansId int64) {
		logrus.WithField("topicId", evt.TopicId).
			WithField("userId", evt.UserId).
			WithField("fansId", fansId).
			Info("用户关注，处理帖子")
		if err := services.UserFeedService.Create(&model.UserFeed{
			UserId:     fansId,
			DataId:     evt.TopicId,
			DataType:   constants.EntityTopic,
			AuthorId:   evt.UserId,
			CreateTime: date.NowTimestamp(),
		}); err != nil {
			logrus.Error(err)
		}
	})
}
