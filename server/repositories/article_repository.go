package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
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

func (this *articleRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Article) {
	cnd.Find(db, &list)
	return
}

func (this *articleRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.Article) {
	cnd.FindOne(db, &ret)
	return
}

func (this *articleRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.Article, paging *simple.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *articleRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Article, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Article{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
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
	db.Delete(&model.Article{}, "id = ?", id)
}
