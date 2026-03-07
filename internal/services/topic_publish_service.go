package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/req"
	"bbs-go/internal/pkg/event"
	"bbs-go/internal/pkg/iplocator"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/pkg/search"
	"bbs-go/internal/repositories"
	"errors"
	"log/slog"
	"strings"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/spf13/cast"
)

var TopicPublishService = new(topicPublishService)

type topicPublishService struct{}

// Publish 发表
func (s *topicPublishService) Publish(userId int64, form req.CreateTopicForm) (*models.Topic, error) {
	if err := s.checkParams(userId, form); err != nil {
		return nil, err
	}

	// QA 话题不处理隐藏内容和投票，前端即使传入也忽略。
	if form.Type == constants.TopicTypeQA {
		form.HideContent = ""
		form.Vote = nil
	}

	now := dates.NowTimestamp()
	topic := &models.Topic{
		Type:            form.Type,
		QaStatus:        constants.QaStatusUnsolved,
		UserId:          userId,
		NodeId:          form.NodeId,
		Title:           form.Title,
		ContentType:     form.ContentType,
		Content:         form.Content,
		HideContent:     form.HideContent,
		Status:          constants.StatusOk,
		UserAgent:       form.UserAgent,
		Ip:              form.Ip,
		IpLocation:      iplocator.IpLocation(form.Ip),
		LastCommentTime: now,
		CreateTime:      now,
	}

	if form.Type == constants.TopicTypeQA && form.BountyScore > 0 {
		topic.BountyScore = form.BountyScore
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
	if s._IsNeedReview(form) {
		topic.Status = constants.StatusReview
	}

	if err := sqls.WithTransaction(func(ctx *sqls.TxContext) error {
		var (
			tagIds []int64
			err    error
		)
		// 帖子
		if err = repositories.TopicRepository.Create(ctx.Tx, topic); err != nil {
			return err
		}
		// 投票
		if form.Vote != nil {
			vote, voteErr := VoteService.CreateWithOptionsTx(ctx, topic.Id, userId, form.Vote, now)
			if voteErr != nil {
				return voteErr
			}
			if vote != nil {
				topic.VoteId = vote.Id
				if err = repositories.TopicRepository.UpdateColumn(ctx.Tx, topic.Id, "vote_id", vote.Id); err != nil {
					return err
				}
			}
		}

		// 标签
		if tagIds, err = repositories.TagRepository.GetOrCreates(ctx.Tx, form.Tags); err != nil {
			return err
		}
		if err = repositories.TopicTagRepository.AddTopicTags(ctx.Tx, topic.Id, tagIds); err != nil {
			return err
		}

		// 用户计数
		if err = UserService.IncrTopicCount(ctx, userId); err != nil {
			return err
		}

		// 问答悬赏：扣减题主积分
		if topic.Type == constants.TopicTypeQA && topic.BountyScore > 0 {
			if err = UserService.DecrScoreTx(ctx, userId, topic.BountyScore, constants.SourceTypeQaBounty, cast.ToString(topic.Id), locales.Get("topic.bounty_deduct")); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// 添加索引
	search.UpdateTopicIndexAsync(topic)
	// 发送事件
	event.Send(event.TopicCreateEvent{
		UserId:     topic.UserId,
		TopicId:    topic.Id,
		TopicType:  int(topic.Type),
		CreateTime: topic.CreateTime,
	})
	return topic, nil
}

// IsNeedReview 是否需要审核
func (s *topicPublishService) _IsNeedReview(form req.CreateTopicForm) bool {
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

func (s topicPublishService) checkParams(userId int64, form req.CreateTopicForm) (err error) {
	modules := SysConfigService.GetModules()
	if form.Type == constants.TopicTypeTweet {
		if !modules.Tweet {
			return errors.New("未开启动态功能")
		}
		if strs.IsBlank(form.Content) {
			return errors.New("内容不能为空")
		}
		// if strs.IsBlank(form.Content) && len(form.ImageList) == 0 {
		// 	return errors.New("内容或图片不能为空")
		// }
	} else if constants.IsPostTopicType(form.Type) {
		if !modules.Topic {
			return errors.New("未开启帖子功能")
		}
		if strs.IsBlank(form.Title) {
			return errors.New("标题不能为空")
		}

		if strs.IsBlank(form.Content) {
			return errors.New("内容不能为空")
		}

		if strs.RuneLen(form.Title) > 128 {
			return errors.New("标题长度不能超过128")
		}
	} else {
		return errors.New(locales.Get("topic.type_not_supported"))
	}

	if form.NodeId <= 0 {
		form.NodeId = SysConfigService.GetDefaultNodeId()
		if form.NodeId <= 0 {
			return errors.New("请选择节点")
		}
	}

	node := repositories.TopicNodeRepository.Get(sqls.DB(), form.NodeId)
	if node == nil || node.Status != constants.StatusOk {
		return errors.New("节点不存在")
	}
	if !node.Type.Supports(form.Type) {
		return errors.New(locales.Get("topic.node_type_mismatch"))
	}
	if form.Type == constants.TopicTypeQA {
		form.Vote = nil
		if !SysConfigService.IsEnableQaBounty() {
			form.BountyScore = 0
		} else {
			if form.BountyScore < 0 {
				return errors.New(locales.Get("topic.bounty_invalid"))
			}
			if SysConfigService.IsQaBountyRequired() {
				minVal := SysConfigService.GetQaBountyMin()
				if form.BountyScore < minVal {
					return errors.New(locales.Get("topic.bounty_required"))
				}
			}
			if form.BountyScore > 0 {
				minVal := SysConfigService.GetQaBountyMin()
				maxVal := SysConfigService.GetQaBountyMax()
				if minVal > 0 && form.BountyScore < minVal {
					if maxVal > 0 {
						return errors.New(locales.Getf("topic.bounty_out_of_range_range", minVal, maxVal))
					}
					return errors.New(locales.Getf("topic.bounty_out_of_range_min", minVal))
				}
				if maxVal > 0 && form.BountyScore > maxVal {
					if minVal > 0 {
						return errors.New(locales.Getf("topic.bounty_out_of_range_range", minVal, maxVal))
					}
					return errors.New(locales.Getf("topic.bounty_out_of_range_max", maxVal))
				}
				user := repositories.UserRepository.Get(sqls.DB(), userId)
				if user == nil || user.Score < form.BountyScore {
					return errors.New(locales.Get("topic.insufficient_score"))
				}
			}
		}
	}
	if err = VoteService.CheckCreateForm(form.Vote); err != nil {
		return err
	}

	return nil
}
