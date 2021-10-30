package eventhandler

import (
	"bbs-go/model"
	"bbs-go/model/constants"
	"bbs-go/pkg/mq"
	"bbs-go/services"
	"github.com/mlogclub/simple/date"
	"github.com/sirupsen/logrus"
)

func init() {
	mq.AddEventHandler(mq.EventTypeTopicCreate, HandleTopicCreate)
}

func HandleTopicCreate(e interface{}) error {
	event := e.(*mq.TopicCreateEvent)

	services.UserFollowService.ScanFans(event.UserId, func(fansId int64) {
		logrus.WithField("topicId", event.TopicId).
			WithField("userId", event.UserId).
			WithField("fansId", fansId).
			Info("用户关注，处理帖子")
		if err := services.UserFeedService.Create(&model.UserFeed{
			UserId:     fansId,
			DataId:     event.TopicId,
			DataType:   constants.EntityTopic,
			AuthorId:   event.UserId,
			CreateTime: date.NowTimestamp(),
		}); err != nil {
			logrus.Error(err)
		}
	})

	return nil
}
