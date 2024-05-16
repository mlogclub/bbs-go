package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/event"
	"bbs-go/internal/pkg/iplocator"
	"bbs-go/internal/pkg/search"
	"bbs-go/internal/repositories"
	"errors"
	"log/slog"
	"strings"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

var TopicPublishService = new(topicPublishService)

type topicPublishService struct{}

// Publish 发表
func (s *topicPublishService) Publish(userId int64, form models.CreateTopicForm) (*models.Topic, error) {
	if err := s._CheckParams(userId, form); err != nil {
		return nil, err
	}

	now := dates.NowTimestamp()
	topic := &models.Topic{
		Type:            form.Type,
		UserId:          userId,
		NodeId:          form.NodeId,
		Title:           form.Title,
		Content:         form.Content,
		HideContent:     form.HideContent,
		Status:          constants.StatusOk,
		UserAgent:       form.UserAgent,
		Ip:              form.Ip,
		IpLocation:      iplocator.IpLocation(form.Ip),
		LastCommentTime: now,
		CreateTime:      now,
	}

	if len(form.ImageList) > 0 {
		imageListStr, err := jsons.ToStr(form.ImageList)
		if err == nil {
			topic.ImageList = imageListStr
		} else {
			slog.Error(err.Error(), slog.Any("err", err))
		}
	}

	// 检查是否需要审核
	if s._IsNeedReview(userId, form) {
		topic.Status = constants.StatusReview
	}

	if err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		var (
			tagIds []int64
			err    error
		)
		// 帖子
		if err := repositories.TopicRepository.Create(tx, topic); err != nil {
			return err
		}

		// 标签
		if tagIds, err = repositories.TagRepository.GetOrCreates(tx, form.Tags); err != nil {
			return err
		}
		if err = repositories.TopicTagRepository.AddTopicTags(tx, topic.Id, tagIds); err != nil {
			return err
		}

		// 用户计数
		if err = UserService.IncrTopicCount(tx, userId); err != nil {
			return err
		}

		// 积分
		UserService.IncrScoreForPostTopic(topic)

		return nil
	}); err != nil {
		return nil, err
	}
	// 添加索引
	search.UpdateTopicIndex(topic)
	// 发送事件
	event.Send(event.TopicCreateEvent{
		UserId:     topic.UserId,
		TopicId:    topic.Id,
		CreateTime: topic.CreateTime,
	})
	return topic, nil
}

// IsNeedReview 是否需要审核
func (s *topicPublishService) _IsNeedReview(userId int64, form models.CreateTopicForm) bool {
	if hits := ForbiddenWordService.Check(form.Title); len(hits) > 0 {
		slog.Info("帖子标题命中违禁词", slog.String("hits", strings.Join(hits, ",")))
		return true
	}

	if hits := ForbiddenWordService.Check(form.Content); len(hits) > 0 {
		slog.Info("帖子内容命中违禁词", slog.String("hits", strings.Join(hits, ",")))
		return true
	}

	return false
}

func (s topicPublishService) _CheckParams(userId int64, form models.CreateTopicForm) (err error) {
	if form.Type == constants.TopicTypeTweet {
		if strs.IsBlank(form.Content) {
			return errors.New("内容不能为空")
		}
		// if strs.IsBlank(form.Content) && len(form.ImageList) == 0 {
		// 	return errors.New("内容或图片不能为空")
		// }
	} else {
		if strs.IsBlank(form.Title) {
			return errors.New("标题不能为空")
		}

		if strs.IsBlank(form.Content) {
			return errors.New("内容不能为空")
		}

		if strs.RuneLen(form.Title) > 128 {
			return errors.New("标题长度不能超过128")
		}
	}

	if form.NodeId <= 0 {
		form.NodeId = SysConfigService.GetConfig().DefaultNodeId
		if form.NodeId <= 0 {
			return errors.New("请选择节点")
		}
	}
	node := repositories.TopicNodeRepository.Get(sqls.DB(), form.NodeId)
	if node == nil || node.Status != constants.StatusOk {
		return errors.New("节点不存在")
	}

	return nil
}
