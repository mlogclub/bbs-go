package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/model"
)

var ArticleRepository = newArticleRepository()

func newArticleRepository() *articleRepository {
	return &articleRepository{}
}

type articleRepository struct {
}

func (this *articleRepository) Get(db *gorm.DB, id int64) *model.Article {
	ret := &model.Article{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *articleRepository) Take(db *gorm.DB, where ...interface{}) *model.Article {
	ret := &model.Article{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *articleRepository) QueryCnd(db *gorm.DB, cnd *simple.QueryCnd) (list []model.Article, err error) {
	err = cnd.DoQuery(db).Find(&list).Error
	return
}

func (this *articleRepository) Query(db *gorm.DB, queries *simple.ParamQueries) (list []model.Article, paging *simple.Paging) {
	queries.StartQuery(db).Find(&list)
	queries.StartCount(db).Model(&model.Article{}).Count(&queries.Paging.Total)
	paging = queries.Paging
	return
}

func (this *articleRepository) Create(db *gorm.DB, t *model.Article) (err error) {
	err = db.Create(t).Error
	return
}

func (this *articleRepository) Update(db *gorm.DB, t *model.Article) (err error) {
	err = db.Save(t).Error
	return
}

func (this *articleRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Article{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *articleRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Article{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *articleRepository) Delete(db *gorm.DB, id int64) {
	db.Model(&model.Article{}).Delete("id", id)
}
