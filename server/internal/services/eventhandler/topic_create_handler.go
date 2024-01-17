package eventhandler

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/bbsurls"
	"bbs-go/internal/pkg/event"
	"bbs-go/internal/pkg/seo"
	"bbs-go/internal/services"
	"log/slog"
	"reflect"

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
		if err := services.UserFeedService.Create(&models.UserFeed{
			UserId:     fansId,
			DataId:     e.TopicId,
			DataType:   constants.EntityTopic,
			AuthorId:   e.UserId,
			CreateTime: e.CreateTime,
		}); err != nil {
			slog.Error(err.Error())
		}
	})

	// 百度链接推送
	seo.Push(bbsurls.TopicUrl(e.TopicId))
}
