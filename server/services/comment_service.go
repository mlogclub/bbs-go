package services

import (
	"bbs-go/model/constants"
	"bbs-go/pkg/event"
	"errors"
	"strings"

	"github.com/mlogclub/simple/date"
	"github.com/mlogclub/simple/json"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/mlogclub/simple"

	"bbs-go/model"
	"bbs-go/repositories"
)

var CommentService = newCommentService()

func newCommentService() *commentService {
	return &commentService{}
}

type commentService struct {
}

func (s *commentService) Get(id int64) *model.Comment {
	return repositories.CommentRepository.Get(simple.DB(), id)
}

func (s *commentService) Take(where ...interface{}) *model.Comment {
	return repositories.CommentRepository.Take(simple.DB(), where...)
}

func (s *commentService) Find(cnd *simple.SqlCnd) []model.Comment {
	return repositories.CommentRepository.Find(simple.DB(), cnd)
}

func (s *commentService) FindOne(cnd *simple.SqlCnd) *model.Comment {
	return repositories.CommentRepository.FindOne(simple.DB(), cnd)
}

func (s *commentService) FindPageByParams(params *simple.QueryParams) (list []model.Comment, paging *simple.Paging) {
	return repositories.CommentRepository.FindPageByParams(simple.DB(), params)
}

func (s *commentService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.Comment, paging *simple.Paging) {
	return repositories.CommentRepository.FindPageByCnd(simple.DB(), cnd)
}

func (s *commentService) Count(cnd *simple.SqlCnd) int64 {
	return repositories.CommentRepository.Count(simple.DB(), cnd)
}

func (s *commentService) Create(t *model.Comment) error {
	return repositories.CommentRepository.Create(simple.DB(), t)
}

func (s *commentService) Update(t *model.Comment) error {
	return repositories.CommentRepository.Update(simple.DB(), t)
}

func (s *commentService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.CommentRepository.Updates(simple.DB(), id, columns)
}

func (s *commentService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.CommentRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (s *commentService) Delete(id int64) error {
	return repositories.CommentRepository.UpdateColumn(simple.DB(), id, "status", constants.StatusDeleted)
}

// Publish 发表评论
func (s *commentService) Publish(userId int64, form model.CreateCommentForm) (*model.Comment, error) {
	form.Content = strings.TrimSpace(form.Content)
	if simple.IsBlank(form.EntityType) {
		return nil, errors.New("参数非法")
	}
	if form.EntityId <= 0 {
		return nil, errors.New("参数非法")
	}
	if simple.IsBlank(form.Content) {
		return nil, errors.New("请输入评论内容")
	}

	comment := &model.Comment{
		UserId:      userId,
		EntityType:  form.EntityType,
		EntityId:    form.EntityId,
		Content:     form.Content,
		ContentType: simple.DefaultIfBlank(form.ContentType, constants.ContentTypeMarkdown),
		QuoteId:     form.QuoteId,
		Status:      constants.StatusOk,
		UserAgent:   form.UserAgent,
		Ip:          form.Ip,
		CreateTime:  date.NowTimestamp(),
	}

	if len(form.ImageList) > 0 {
		imageListStr, err := json.ToStr(form.ImageList)
		if err == nil {
			comment.ImageList = imageListStr
		} else {
			logrus.Error(err)
		}
	}

	err := simple.DB().Transaction(func(tx *gorm.DB) error {
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
func (s *commentService) onComment(tx *gorm.DB, comment *model.Comment) error {
	return repositories.CommentRepository.UpdateColumn(tx, comment.EntityId, "comment_count", gorm.Expr("comment_count + 1"))
}

// // 统计数量
// func (s *commentService) Count(entityType string, entityId int64) int64 {
// 	var count int64 = 0
// 	simple.DB().Model(&model.Comment{}).Where("entity_type = ? and entity_id = ?", entityType, entityId).Count(&count)
// 	return count
// }

// GetComments 列表
func (s *commentService) GetComments(entityType string, entityId int64, cursor int64) (comments []model.Comment, nextCursor int64, hasMore bool) {
	limit := 20
	cnd := simple.NewSqlCnd().Eq("entity_type", entityType).Eq("entity_id", entityId).Eq("status", constants.StatusOk).Desc("id").Limit(limit)
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	comments = repositories.CommentRepository.Find(simple.DB(), cnd)
	if len(comments) > 0 {
		nextCursor = comments[len(comments)-1].Id
		hasMore = len(comments) >= limit
	} else {
		nextCursor = cursor
	}
	return
}

// GetReplies 二级回复列表
func (s *commentService) GetReplies(commentId int64, cursor int64, limit int) (comments []model.Comment, nextCursor int64, hasMore bool) {
	cnd := simple.NewSqlCnd().Eq("entity_type", constants.EntityComment).Eq("entity_id", commentId).Eq("status", constants.StatusOk).Asc("id").Limit(limit)
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
func (s *commentService) ScanByUser(userId int64, callback func(comments []model.Comment)) {
	var cursor int64 = 0
	for {
		list := repositories.CommentRepository.Find(simple.DB(), simple.NewSqlCnd().
			Eq("user_id", userId).Gt("id", cursor).Asc("id").Limit(1000))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}
