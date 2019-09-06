
package services

import (
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
	"github.com/mlogclub/simple"
)

var SubjectService = newSubjectService()

func newSubjectService() *subjectService {
	return &subjectService {}
}

type subjectService struct {
}

func (this *subjectService) Get(id int64) *model.Subject {
	return repositories.SubjectRepository.Get(simple.GetDB(), id)
}

func (this *subjectService) Take(where ...interface{}) *model.Subject {
	return repositories.SubjectRepository.Take(simple.GetDB(), where...)
}

func (this *subjectService) QueryCnd(cnd *simple.QueryCnd) (list []model.Subject, err error) {
	return repositories.SubjectRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *subjectService) Query(queries *simple.ParamQueries) (list []model.Subject, paging *simple.Paging) {
	return repositories.SubjectRepository.Query(simple.GetDB(), queries)
}

func (this *subjectService) Create(t *model.Subject) error {
	return repositories.SubjectRepository.Create(simple.GetDB(), t)
}

func (this *subjectService) Update(t *model.Subject) error {
	return repositories.SubjectRepository.Update(simple.GetDB(), t)
}

func (this *subjectService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.SubjectRepository.Updates(simple.GetDB(), id, columns)
}

func (this *subjectService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.SubjectRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *subjectService) Delete(id int64) {
	repositories.SubjectRepository.Delete(simple.GetDB(), id)
}

