package services

import (
	"errors"
	"strings"

	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
)

var CommentService = newCommentService()

func newCommentService() *commentService {
	return &commentService{}
}

type commentService struct {
}

func (this *commentService) Get(id int64) *model.Comment {
	return repositories.CommentRepository.Get(simple.DB(), id)
}

func (this *commentService) Take(where ...interface{}) *model.Comment {
	return repositories.CommentRepository.Take(simple.DB(), where...)
}

func (this *commentService) Find(cnd *simple.SqlCnd) (list []model.Comment, err error) {
	return repositories.CommentRepository.Find(simple.DB(), cnd)
}

func (this *commentService) FindPageByParams(params *simple.QueryParams) (list []model.Comment, paging *simple.Paging) {
	return repositories.CommentRepository.FindPageByParams(simple.DB(), params)
}

func (this *commentService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.Comment, paging *simple.Paging) {
	return repositories.CommentRepository.FindPageByCnd(simple.DB(), cnd)
}

func (this *commentService) Create(t *model.Comment) error {
	return repositories.CommentRepository.Create(simple.DB(), t)
}

func (this *commentService) Update(t *model.Comment) error {
	return repositories.CommentRepository.Update(simple.DB(), t)
}

func (this *commentService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.CommentRepository.Updates(simple.DB(), id, columns)
}

func (this *commentService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.CommentRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (this *commentService) Delete(id int64) error {
	return repositories.CommentRepository.UpdateColumn(simple.DB(), id, "status", model.CommentStatusDeleted)
}

// 发表评论
func (this *commentService) Publish(userId int64, form *model.CreateCommentForm) (*model.Comment, error) {
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
	comment := &model.Comment{
		UserId:     userId,
		EntityType: form.EntityType,
		EntityId:   form.EntityId,
		Content:    form.Content,
		QuoteId:    form.QuoteId,
		Status:     model.CommentStatusOk,
		CreateTime: simple.NowTimestamp(),
	}
	if err := repositories.CommentRepository.Create(simple.DB(), comment); err != nil {
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
func (this *commentService) Count(entityType string, entityId int64) int64 {
	var count int64 = 0
	simple.DB().Model(&model.Comment{}).Where("entity_type = ? and entity_id = ?", entityType, entityId).Count(&count)
	return count
}

// 列表
func (this *commentService) List(entityType string, entityId int64, cursor int64) (list []model.Comment, err error) {
	if cursor > 0 {
		err = simple.DB().Where("entity_type = ? and entity_id = ? and status = ? and id < ?", entityType,
			entityId, model.CommentStatusOk, cursor).Order("id desc").Limit(20).Find(&list).Error
	} else {
		err = simple.DB().Where("entity_type = ? and entity_id = ? and status = ?", entityType, entityId,
			model.CommentStatusOk).Order("id desc").Limit(20).Find(&list).Error
	}
	return
}
