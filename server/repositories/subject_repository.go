package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
)

var SubjectRepository = newSubjectRepository()

func newSubjectRepository() *subjectRepository {
	return &subjectRepository{}
}

type subjectRepository struct {
}

func (this *subjectRepository) Get(db *gorm.DB, id int64) *model.Subject {
	ret := &model.Subject{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *subjectRepository) Take(db *gorm.DB, where ...interface{}) *model.Subject {
	ret := &model.Subject{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *subjectRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Subject) {
	cnd.Find(db, &list)
	return
}

func (this *subjectRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.Subject) {
	cnd.FindOne(db, &ret)
	return
}

func (this *subjectRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.Subject, paging *simple.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *subjectRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Subject, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Subject{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *subjectRepository) Create(db *gorm.DB, t *model.Subject) (err error) {
	err = db.Create(t).Error
	return
}

func (this *subjectRepository) Update(db *gorm.DB, t *model.Subject) (err error) {
	err = db.Save(t).Error
	return
}

func (this *subjectRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Subject{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *subjectRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Subject{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *subjectRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.Subject{}, "id = ?", id)
}
