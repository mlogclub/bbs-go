package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/simple"
)

type ArticleShareRepository struct {
}

func NewArticleShareRepository() *ArticleShareRepository {
	return &ArticleShareRepository{}
}

func (this *ArticleShareRepository) Get(db *gorm.DB, id int64) *model.ArticleShare {
	ret := &model.ArticleShare{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *ArticleShareRepository) Take(db *gorm.DB, where ...interface{}) *model.ArticleShare {
	ret := &model.ArticleShare{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *ArticleShareRepository) QueryCnd(db *gorm.DB, cnd *simple.QueryCnd) (list []model.ArticleShare, err error) {
	err = cnd.DoQuery(db).Find(&list).Error
	return
}

func (this *ArticleShareRepository) Query(db *gorm.DB, queries *simple.ParamQueries) (list []model.ArticleShare, paging *simple.Paging) {
	queries.StartQuery(db).Find(&list)
	queries.StartCount(db).Model(&model.ArticleShare{}).Count(&queries.Paging.Total)
	paging = queries.Paging
	return
}

func (this *ArticleShareRepository) Create(db *gorm.DB, t *model.ArticleShare) (err error) {
	err = db.Create(t).Error
	return
}

func (this *ArticleShareRepository) Update(db *gorm.DB, t *model.ArticleShare) (err error) {
	err = db.Save(t).Error
	return
}

func (this *ArticleShareRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.ArticleShare{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *ArticleShareRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.ArticleShare{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *ArticleShareRepository) Delete(db *gorm.DB, id int64) {
	db.Model(&model.ArticleShare{}).Delete("id", id)
}
