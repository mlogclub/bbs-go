package services

import (
	"errors"
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
	return repositories.CommentRepository.UpdateColumn(simple.DB(), id, "status", model.StatusDeleted)
}

// 发表评论
func (s *commentService) Publish(userId int64, form *model.CreateCommentForm) (*model.Comment, error) {
	form.Content = strings.TrimSpace(form.Content)

	if len(form.EntityType) == 0 {
		return nil, errors.New("参数非法")
	}
	if form.EntityId <= 0 {
		return nil, errors.New("参数非法")
	}
	if len(form.Content) == 0 {
		return nil, errors.New("请输入评论内容")
	}

	contentType := form.ContentType
	if contentType == "" {
		contentType = model.ContentTypeMarkdown
	}

	comment := &model.Comment{
		UserId:      userId,
		EntityType:  form.EntityType,
		EntityId:    form.EntityId,
		Content:     form.Content,
		ContentType: contentType,
		QuoteId:     form.QuoteId,
		Status:      model.StatusOk,
		CreateTime:  simple.NowTimestamp(),
	}
	if err := s.Create(comment); err != nil {
		return nil, err
	}

	// 更新帖子最后回复时间
	if form.EntityType == model.EntityTypeTopic {
		TopicService.OnComment(form.EntityId, simple.NowTimestamp())
	}

	// 发送消息
	MessageService.SendCommentMsg(comment)

	return comment, nil
}

// 统计数量
func (s *commentService) Count(entityType string, entityId int64) int64 {
	var count int64 = 0
	simple.DB().Model(&model.Comment{}).Where("entity_type = ? and entity_id = ?", entityType, entityId).Count(&count)
	return count
}

// 列表
func (s *commentService) GetComments(entityType string, entityId int64, cursor int64) (comments []model.Comment, nextCursor int64) {
	cnd := simple.NewSqlCnd().Eq("entity_type", entityType).Eq("entity_id", entityId).Eq("status", model.StatusOk).Desc("id").Limit(50)
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
