package repositories

import (
	"github.com/mlogclub/simple"
	"gorm.io/gorm"

	"bbs-go/model"
)

var ArticleRepository = newArticleRepository()

func newArticleRepository() *articleRepository {
	return &articleRepository{}
}

type articleRepository struct {
}

func (r *articleRepository) Get(db *gorm.DB, id int64) *model.Article {
	ret := &model.Article{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *articleRepository) Take(db *gorm.DB, where ...interface{}) *model.Article {
	ret := &model.Article{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *articleRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Article) {
	cnd.Find(db, &list)
	return
}

func (r *articleRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) *model.Article {
	ret := &model.Article{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *articleRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.Article, paging *simple.Paging) {
	return r.FindPageByCnd(db, &params.SqlCnd)
}

func (r *articleRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Article, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Article{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *articleRepository) Create(db *gorm.DB, t *model.Article) (err error) {
	err = db.Create(t).Error
	return
}

func (r *articleRepository) Update(db *gorm.DB, t *model.Article) (err error) {
	err = db.Save(t).Error
	return
}

func (r *articleRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Article{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *articleRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Article{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *articleRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.Article{}, "id = ?", id)
}
