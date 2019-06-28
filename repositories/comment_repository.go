package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/model"
)

type CommentRepository struct {
}

func NewCommentRepository() *CommentRepository {
	return &CommentRepository{}
}

func (this *CommentRepository) Get(db *gorm.DB, id int64) *model.Comment {
	ret := &model.Comment{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *CommentRepository) Take(db *gorm.DB, where ...interface{}) *model.Comment {
	ret := &model.Comment{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *CommentRepository) QueryCnd(db *gorm.DB, cnd *simple.QueryCnd) (list []model.Comment, err error) {
	err = cnd.DoQuery(db).Find(&list).Error
	return
}

func (this *CommentRepository) Query(db *gorm.DB, queries *simple.ParamQueries) (list []model.Comment, paging *simple.Paging) {
	queries.StartQuery(db).Find(&list)
	queries.StartCount(db).Model(&model.Comment{}).Count(&queries.Paging.Total)
	paging = queries.Paging
	return
}

func (this *CommentRepository) Create(db *gorm.DB, t *model.Comment) (err error) {
	err = db.Create(t).Error
	return
}

func (this *CommentRepository) Update(db *gorm.DB, t *model.Comment) (err error) {
	err = db.Save(t).Error
	return
}

func (this *CommentRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Comment{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *CommentRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Comment{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *CommentRepository) Delete(db *gorm.DB, id int64) {
	db.Model(&model.Comment{}).Delete("id", id)
}
