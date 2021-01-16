package services

import (
	"bbs-go/model/constants"
	"errors"
	"github.com/mlogclub/simple/date"
	"strings"

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

// 发表评论
func (s *commentService) Publish(userId int64, form *model.CreateCommentForm) (*model.Comment, error) {
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
		CreateTime:  date.NowTimestamp(),
	}
	if err := s.Create(comment); err != nil {
		return nil, err
	}

	if form.EntityType == constants.EntityTopic {
		TopicService.OnComment(form.EntityId, comment)
	}

	UserService.IncrCommentCount(userId)         // 用户跟帖计数
	UserService.IncrScoreForPostComment(comment) // 获得积分
	MessageService.SendCommentMsg(comment)       // 发送消息

	return comment, nil
}

// // 统计数量
// func (s *commentService) Count(entityType string, entityId int64) int64 {
// 	var count int64 = 0
// 	simple.DB().Model(&model.Comment{}).Where("entity_type = ? and entity_id = ?", entityType, entityId).Count(&count)
// 	return count
// }

// 列表
func (s *commentService) GetComments(entityType string, entityId int64, cursor int64) (comments []model.Comment, nextCursor int64) {
	cnd := simple.NewSqlCnd().Eq("entity_type", entityType).Eq("entity_id", entityId).Eq("status", constants.StatusOk).Desc("id").Limit(50)
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	comments = repositories.CommentRepository.Find(simple.DB(), cnd)
	if len(comments) > 0 {
		nextCursor = comments[len(comments)-1].Id
	} else {
		nextCursor = cursor
	}
	return
}
