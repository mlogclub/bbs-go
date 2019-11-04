package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
)

var SubjectContentRepository = newSubjectContentRepository()

func newSubjectContentRepository() *subjectContentRepository {
	return &subjectContentRepository{}
}

type subjectContentRepository struct {
}

func (this *subjectContentRepository) Get(db *gorm.DB, id int64) *model.SubjectContent {
	ret := &model.SubjectContent{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *subjectContentRepository) Take(db *gorm.DB, where ...interface{}) *model.SubjectContent {
	ret := &model.SubjectContent{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *subjectContentRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.SubjectContent) {
	cnd.Find(db, &list)
	return
}

func (this *subjectContentRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.SubjectContent) {
	cnd.FindOne(db, &ret)
	return
}

func (this *subjectContentRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.SubjectContent, paging *simple.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *subjectContentRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.SubjectContent, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.SubjectContent{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *subjectContentRepository) Create(db *gorm.DB, t *model.SubjectContent) (err error) {
	err = db.Create(t).Error
	return
}

func (this *subjectContentRepository) Update(db *gorm.DB, t *model.SubjectContent) (err error) {
	err = db.Save(t).Error
	return
}

func (this *subjectContentRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.SubjectContent{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *subjectContentRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.SubjectContent{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *subjectContentRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.SubjectContent{}, "id = ?", id)
}
