
package repositories

import (
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/simple"
	"github.com/jinzhu/gorm"
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

func (this *subjectContentRepository) QueryCnd(db *gorm.DB, cnd *simple.QueryCnd) (list []model.SubjectContent, err error) {
	err = cnd.DoQuery(db).Find(&list).Error
	return
}

func (this *subjectContentRepository) Query(db *gorm.DB, params *simple.ParamQueries) (list []model.SubjectContent, paging *simple.Paging) {
	params.StartQuery(db).Find(&list)
    params.StartCount(db).Model(&model.SubjectContent{}).Count(&params.Paging.Total)
	paging = params.Paging
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

