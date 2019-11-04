package services

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
)

var SubjectService = newSubjectService()

func newSubjectService() *subjectService {
	return &subjectService{}
}

type subjectService struct {
}

func (this *subjectService) Get(id int64) *model.Subject {
	return repositories.SubjectRepository.Get(simple.DB(), id)
}

func (this *subjectService) Take(where ...interface{}) *model.Subject {
	return repositories.SubjectRepository.Take(simple.DB(), where...)
}

func (this *subjectService) Find(cnd *simple.SqlCnd) []model.Subject {
	return repositories.SubjectRepository.Find(simple.DB(), cnd)
}

func (this *subjectService) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.Subject) {
	cnd.FindOne(db, &ret)
	return
}

func (this *subjectService) FindPageByParams(params *simple.QueryParams) (list []model.Subject, paging *simple.Paging) {
	return repositories.SubjectRepository.FindPageByParams(simple.DB(), params)
}

func (this *subjectService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.Subject, paging *simple.Paging) {
	return repositories.SubjectRepository.FindPageByCnd(simple.DB(), cnd)
}

func (this *subjectService) Create(t *model.Subject) error {
	return repositories.SubjectRepository.Create(simple.DB(), t)
}

func (this *subjectService) Update(t *model.Subject) error {
	return repositories.SubjectRepository.Update(simple.DB(), t)
}

func (this *subjectService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.SubjectRepository.Updates(simple.DB(), id, columns)
}

func (this *subjectService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.SubjectRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (this *subjectService) Delete(id int64) {
	repositories.SubjectRepository.Delete(simple.DB(), id)
}
