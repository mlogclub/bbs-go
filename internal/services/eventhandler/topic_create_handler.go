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
)

func init() {
	event.RegHandler(reflect.TypeOf(event.TopicCreateEvent{}), handleTopicCreateEvent)
}

func handleTopicCreateEvent(i interface{}) {
	e := i.(event.TopicCreateEvent)

	services.UserFollowService.ScanFans(e.UserId, func(fansId int64) {
		slog.With(slog.Any("topicId", e.TopicId), slog.Any("userId", e.UserId), slog.Any("fansId", fansId)).Info("用户关注，处理帖子")
		if err := services.UserFeedService.Create(&models.UserFeed{
			UserId:     fansId,
			DataId:     e.TopicId,
			DataType:   constants.EntityTopic,
			AuthorId:   e.UserId,
			CreateTime: e.CreateTime,
		}); err != nil {
			slog.Error(err.Error(), slog.Any("err", err))
		}
	})

	// 百度链接推送
	seo.Push(bbsurls.TopicUrl(e.TopicId))
}
