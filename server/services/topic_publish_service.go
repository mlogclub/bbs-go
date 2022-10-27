package services

import (
	"bbs-go/model"
	"bbs-go/model/constants"
	"bbs-go/pkg/es"
	"bbs-go/pkg/event"
	"bbs-go/repositories"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var TopicPublishService = new(topicPublishService)

type topicPublishService struct{}

// Publish 发表
func (s *topicPublishService) Publish(userId int64, form model.CreateTopicForm) (*model.Topic, error) {
	if err := s._CheckParams(userId, form); err != nil {
		return nil, err
	}

	now := dates.NowTimestamp()
	topic := &model.Topic{
		Type:            form.Type,
		UserId:          userId,
		NodeId:          form.NodeId,
		Title:           form.Title,
		Content:         form.Content,
		HideContent:     form.HideContent,
		Status:          constants.StatusOk,
		UserAgent:       form.UserAgent,
		Ip:              form.Ip,
		LastCommentTime: now,
		CreateTime:      now,
	}

	if len(form.ImageList) > 0 {
		imageListStr, err := jsons.ToStr(form.ImageList)
		if err == nil {
			topic.ImageList = imageListStr
		} else {
			logrus.Error(err)
		}
	}

	if err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		tagIds := repositories.TagRepository.GetOrCreates(tx, form.Tags)
		err := repositories.TopicRepository.Create(tx, topic)
		if err != nil {
			return err
		}
		repositories.TopicTagRepository.AddTopicTags(tx, topic.Id, tagIds)
		return nil
	}); err != nil {
		return nil, err
	}
	// 添加索引
	es.UpdateTopicIndex(topic)
	// 用户话题计数
	UserService.IncrTopicCount(userId)
	// 获得积分
	UserService.IncrScoreForPostTopic(topic)
	// 发送事件
	event.Send(event.TopicCreateEvent{
		UserId:     topic.UserId,
		TopicId:    topic.Id,
		CreateTime: topic.CreateTime,
	})
	return topic, nil
}

func (s topicPublishService) _CheckParams(userId int64, form model.CreateTopicForm) (err error) {
	// TODO 帖子内容、标题字数限制可配置
	if form.Type == constants.TopicTypeTweet {
		if strs.IsBlank(form.Content) && len(form.ImageList) == 0 {
			return web.NewErrorMsg("内容或图片不能为空")
		}
	} else {
		if strs.IsBlank(form.Title) {
			return web.NewErrorMsg("标题不能为空")
		}

		if strs.IsBlank(form.Content) {
			return web.NewErrorMsg("内容不能为空")
		}

		if strs.RuneLen(form.Title) > 128 {
			return web.NewErrorMsg("标题长度不能超过128")
		}
	}

	if form.NodeId <= 0 {
		form.NodeId = SysConfigService.GetConfig().DefaultNodeId
		if form.NodeId <= 0 {
			return web.NewErrorMsg("请选择节点")
		}
	}
	node := repositories.TopicNodeRepository.Get(sqls.DB(), form.NodeId)
	if node == nil || node.Status != constants.StatusOk {
		return web.NewErrorMsg("节点不存在")
	}
	return nil
}
