package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
)

var CommentRepository = newCommentRepository()

func newCommentRepository() *commentRepository {
	return &commentRepository{}
}

type commentRepository struct {
}

func (this *commentRepository) Get(db *gorm.DB, id int64) *model.Comment {
	ret := &model.Comment{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *commentRepository) Take(db *gorm.DB, where ...interface{}) *model.Comment {
	ret := &model.Comment{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *commentRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Comment) {
	cnd.Find(db, &list)
	return
}

func (this *commentRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.Comment) {
	cnd.FindOne(db, &ret)
	return
}

func (this *commentRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.Comment, paging *simple.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *commentRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Comment, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Comment{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *commentRepository) Create(db *gorm.DB, t *model.Comment) (err error) {
	err = db.Create(t).Error
	return
}

func (this *commentRepository) Update(db *gorm.DB, t *model.Comment) (err error) {
	err = db.Save(t).Error
	return
}

func (this *commentRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Comment{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *commentRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Comment{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *commentRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.Comment{}, "id = ?", id)
}
