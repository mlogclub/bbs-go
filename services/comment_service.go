package services

import (
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
)

type CommentService struct {
	CommentRepository *repositories.CommentRepository
}

func NewCommentService() *CommentService {
	return &CommentService{
		CommentRepository: repositories.NewCommentRepository(),
	}
}

func (this *CommentService) Get(id int64) *model.Comment {
	return this.CommentRepository.Get(simple.GetDB(), id)
}

func (this *CommentService) Take(where ...interface{}) *model.Comment {
	return this.CommentRepository.Take(simple.GetDB(), where...)
}

func (this *CommentService) QueryCnd(cnd *simple.QueryCnd) (list []model.Comment, err error) {
	return this.CommentRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *CommentService) Query(queries *simple.ParamQueries) (list []model.Comment, paging *simple.Paging) {
	return this.CommentRepository.Query(simple.GetDB(), queries)
}

func (this *CommentService) Create(t *model.Comment) error {
	return this.CommentRepository.Create(simple.GetDB(), t)
}

func (this *CommentService) Update(t *model.Comment) error {
	return this.CommentRepository.Update(simple.GetDB(), t)
}

func (this *CommentService) Updates(id int64, columns map[string]interface{}) error {
	return this.CommentRepository.Updates(simple.GetDB(), id, columns)
}

func (this *CommentService) UpdateColumn(id int64, name string, value interface{}) error {
	return this.CommentRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *CommentService) Delete(id int64) {
	this.CommentRepository.Delete(simple.GetDB(), id)
}

func (this *CommentService) List(entityType string, entityId int64, cursor int64) (list []model.Comment, err error) {
	if cursor > 0 {
		err = simple.GetDB().Where("entity_type = ? and entity_id = ? and status = ? and id < ?", entityType,
			entityId, model.CommentStatusOk, cursor).Order("id desc").Limit(20).Find(&list).Error
	} else {
		err = simple.GetDB().Where("entity_type = ? and entity_id = ? and status = ?", entityType, entityId,
			model.CommentStatusOk).Order("id desc").Limit(20).Find(&list).Error
	}
	return
}
