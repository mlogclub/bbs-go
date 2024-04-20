package services

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/event"
	"bbs-go/internal/pkg/iplocator"
	"errors"
	"log/slog"
	"strings"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"gorm.io/gorm"

	"bbs-go/internal/models"
	"bbs-go/internal/repositories"
)

var CommentService = newCommentService()

func newCommentService() *commentService {
	return &commentService{}
}

type commentService struct {
}

func (s *commentService) Get(id int64) *models.Comment {
	return repositories.CommentRepository.Get(sqls.DB(), id)
}

func (s *commentService) Take(where ...interface{}) *models.Comment {
	return repositories.CommentRepository.Take(sqls.DB(), where...)
}

func (s *commentService) Find(cnd *sqls.Cnd) []models.Comment {
	return repositories.CommentRepository.Find(sqls.DB(), cnd)
}

func (s *commentService) FindOne(cnd *sqls.Cnd) *models.Comment {
	return repositories.CommentRepository.FindOne(sqls.DB(), cnd)
}

func (s *commentService) FindPageByParams(params *params.QueryParams) (list []models.Comment, paging *sqls.Paging) {
	return repositories.CommentRepository.FindPageByParams(sqls.DB(), params)
}

func (s *commentService) FindPageByCnd(cnd *sqls.Cnd) (list []models.Comment, paging *sqls.Paging) {
	return repositories.CommentRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *commentService) Count(cnd *sqls.Cnd) int64 {
	return repositories.CommentRepository.Count(sqls.DB(), cnd)
}

func (s *commentService) Create(t *models.Comment) error {
	return repositories.CommentRepository.Create(sqls.DB(), t)
}

func (s *commentService) Update(t *models.Comment) error {
	return repositories.CommentRepository.Update(sqls.DB(), t)
}

func (s *commentService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.CommentRepository.Updates(sqls.DB(), id, columns)
}

func (s *commentService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.CommentRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *commentService) Delete(id int64) error {
	return repositories.CommentRepository.UpdateColumn(sqls.DB(), id, "status", constants.StatusDeleted)
}

// Publish 发表评论
func (s *commentService) Publish(userId int64, form models.CreateCommentForm) (*models.Comment, error) {
	form.Content = strings.TrimSpace(form.Content)
	if strs.IsBlank(form.EntityType) {
		return nil, errors.New("参数非法")
	}
	if form.EntityId <= 0 {
		return nil, errors.New("参数非法")
	}
	if strs.IsBlank(form.Content) {
		return nil, errors.New("请输入评论内容")
	}

	comment := &models.Comment{
		UserId:      userId,
		EntityType:  form.EntityType,
		EntityId:    form.EntityId,
		Content:     form.Content,
		ContentType: constants.ContentTypeText,
		QuoteId:     form.QuoteId,
		Status:      constants.StatusOk,
		UserAgent:   form.UserAgent,
		Ip:          form.Ip,
		IpLocation:  iplocator.IpLocation(form.Ip),
		CreateTime:  dates.NowTimestamp(),
	}

	if len(form.ImageList) > 0 {
		imageListStr, err := jsons.ToStr(form.ImageList)
		if err == nil {
			comment.ImageList = imageListStr
		} else {
			slog.Error(err.Error(), slog.Any("err", err))
		}
	}

	err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := repositories.CommentRepository.Create(tx, comment); err != nil {
			return err
		}

		if form.EntityType == constants.EntityTopic {
			if err := TopicService.onComment(tx, form.EntityId, comment); err != nil {
				return err
			}
		} else if form.EntityType == constants.EntityComment { // 二级评论
			if err := s.onComment(tx, comment); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// 用户跟帖计数
	UserService.IncrCommentCount(userId)
	// 获得积分
	UserService.IncrScoreForPostComment(comment)
	// 发送事件
	event.Send(event.CommentCreateEvent{
		UserId:    userId,
		CommentId: comment.Id,
	})

	return comment, nil
}

// onComment 评论被回复（二级评论）
func (s *commentService) onComment(tx *gorm.DB, comment *models.Comment) error {
	return repositories.CommentRepository.UpdateColumn(tx, comment.EntityId, "comment_count", gorm.Expr("comment_count + 1"))
}

// // 统计数量
// func (s *commentService) Count(entityType string, entityId int64) int64 {
// 	var count int64 = 0
// 	sqls.DB().Model(&model.Comment{}).Where("entity_type = ? and entity_id = ?", entityType, entityId).Count(&count)
// 	return count
// }

// GetComments 列表
func (s *commentService) GetComments(entityType string, entityId int64, cursor int64) (comments []models.Comment, nextCursor int64, hasMore bool) {
	limit := 20
	cnd := sqls.NewCnd().Eq("entity_type", entityType).Eq("entity_id", entityId).Eq("status", constants.StatusOk).Desc("id").Limit(limit)
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	comments = repositories.CommentRepository.Find(sqls.DB(), cnd)
	if len(comments) > 0 {
		nextCursor = comments[len(comments)-1].Id
		hasMore = len(comments) >= limit
	} else {
		nextCursor = cursor
	}
	return
}

// GetReplies 二级回复列表
func (s *commentService) GetReplies(commentId int64, cursor int64, limit int) (comments []models.Comment, nextCursor int64, hasMore bool) {
	cnd := sqls.NewCnd().Eq("entity_type", constants.EntityComment).Eq("entity_id", commentId).Eq("status", constants.StatusOk).Asc("id").Limit(limit)
	if cursor > 0 {
		cnd.Gt("id", cursor)
	}
	comments = s.Find(cnd)
	if len(comments) > 0 {
		nextCursor = comments[len(comments)-1].Id
		hasMore = len(comments) >= limit
	} else {
		nextCursor = cursor
	}
	return
}

// ScanByUser 按照用户扫描数据
func (s *commentService) ScanByUser(userId int64, callback func(comments []models.Comment)) {
	var cursor int64 = 0
	for {
		list := repositories.CommentRepository.Find(sqls.DB(), sqls.NewCnd().
			Eq("user_id", userId).Gt("id", cursor).Asc("id").Limit(1000))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}

// ScanByUser 按照用户扫描数据
func (s *commentService) Scan(callback func(comments []models.Comment)) {
	var cursor int64 = 0
	for {
		logrus.Info("scan comments, cursor:" + cast.ToString(cursor))
		list := repositories.CommentRepository.Find(sqls.DB(), sqls.NewCnd().
			Gt("id", cursor).Asc("id").Limit(1000))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}

func (s *commentService) IsCommented(userId int64, entityType string, entityId int64) bool {
	return s.FindOne(sqls.NewCnd().Where("user_id = ? and entity_id = ? and entity_type = ? and status = ?", userId, entityId, entityType, constants.StatusOk)) != nil
}
