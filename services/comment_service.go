package services

import (
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
)

var CommentService = newCommentService()

func newCommentService() *commentService {
	return &commentService{}
}

type commentService struct {
}

func (this *commentService) Get(id int64) *model.Comment {
	return repositories.CommentRepository.Get(simple.GetDB(), id)
}

func (this *commentService) Take(where ...interface{}) *model.Comment {
	return repositories.CommentRepository.Take(simple.GetDB(), where...)
}

func (this *commentService) QueryCnd(cnd *simple.QueryCnd) (list []model.Comment, err error) {
	return repositories.CommentRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *commentService) Query(queries *simple.ParamQueries) (list []model.Comment, paging *simple.Paging) {
	return repositories.CommentRepository.Query(simple.GetDB(), queries)
}

func (this *commentService) Create(t *model.Comment) error {
	return repositories.CommentRepository.Create(simple.GetDB(), t)
}

func (this *commentService) Update(t *model.Comment) error {
	return repositories.CommentRepository.Update(simple.GetDB(), t)
}

func (this *commentService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.CommentRepository.Updates(simple.GetDB(), id, columns)
}

func (this *commentService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.CommentRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *commentService) Delete(id int64) {
	repositories.CommentRepository.Delete(simple.GetDB(), id)
}

func (this *commentService) List(entityType string, entityId int64, cursor int64) (list []model.Comment, err error) {
	if cursor > 0 {
		err = simple.GetDB().Where("entity_type = ? and entity_id = ? and status = ? and id < ?", entityType,
			entityId, model.CommentStatusOk, cursor).Order("id desc").Limit(20).Find(&list).Error
	} else {
		err = simple.GetDB().Where("entity_type = ? and entity_id = ? and status = ?", entityType, entityId,
			model.CommentStatusOk).Order("id desc").Limit(20).Find(&list).Error
	}
	return
}
