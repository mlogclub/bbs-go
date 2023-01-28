package eventhandler

import (
	"reflect"
	"server/model"
	"server/model/constants"
	"server/pkg/bbsurls"
	"server/pkg/event"
	"server/pkg/seo"
	"server/services"

	"github.com/sirupsen/logrus"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.TopicCreateEvent{}), handleTopicCreateEvent)
}

func handleTopicCreateEvent(i interface{}) {
	e := i.(event.TopicCreateEvent)

	services.UserFollowService.ScanFans(e.UserId, func(fansId int64) {
		logrus.WithField("topicId", e.TopicId).
			WithField("userId", e.UserId).
			WithField("fansId", fansId).
			Info("用户关注，处理帖子")
		if err := services.UserFeedService.Create(&model.UserFeed{
			UserId:     fansId,
			DataId:     e.TopicId,
			DataType:   constants.EntityTopic,
			AuthorId:   e.UserId,
			CreateTime: e.CreateTime,
		}); err != nil {
			logrus.Error(err)
		}
	})

	// 百度链接推送
	seo.Push(bbsurls.TopicUrl(e.TopicId))
}
