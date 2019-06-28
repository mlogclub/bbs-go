package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/model"
)

type ArticleRepository struct {
}

func NewArticleRepository() *ArticleRepository {
	return &ArticleRepository{}
}

func (this *ArticleRepository) Get(db *gorm.DB, id int64) *model.Article {
	ret := &model.Article{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *ArticleRepository) Take(db *gorm.DB, where ...interface{}) *model.Article {
	ret := &model.Article{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *ArticleRepository) QueryCnd(db *gorm.DB, cnd *simple.QueryCnd) (list []model.Article, err error) {
	err = cnd.DoQuery(db).Find(&list).Error
	return
}

func (this *ArticleRepository) Query(db *gorm.DB, queries *simple.ParamQueries) (list []model.Article, paging *simple.Paging) {
	queries.StartQuery(db).Find(&list)
	queries.StartCount(db).Model(&model.Article{}).Count(&queries.Paging.Total)
	paging = queries.Paging
	return
}

func (this *ArticleRepository) Create(db *gorm.DB, t *model.Article) (err error) {
	err = db.Create(t).Error
	return
}

func (this *ArticleRepository) Update(db *gorm.DB, t *model.Article) (err error) {
	err = db.Save(t).Error
	return
}

func (this *ArticleRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Article{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *ArticleRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Article{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *ArticleRepository) Delete(db *gorm.DB, id int64) {
	db.Model(&model.Article{}).Delete("id", id)
}
