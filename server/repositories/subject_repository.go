
package repositories

import (
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/simple"
	"github.com/jinzhu/gorm"
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

func (this *subjectRepository) QueryCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Subject, err error) {
	err = cnd.Exec(db).Find(&list).Error
	return
}

func (this *subjectRepository) Query(db *gorm.DB, params *simple.QueryParams) (list []model.Subject, paging *simple.Paging) {
	params.StartQuery(db).Find(&list)
    params.StartCount(db).Model(&model.Subject{}).Count(&params.Paging.Total)
	paging = params.Paging
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

