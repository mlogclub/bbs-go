package repositories

import (
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"

	"bbs-go/internal/models"
)

var ArticleRepository = newArticleRepository()

func newArticleRepository() *articleRepository {
	return &articleRepository{}
}

type articleRepository struct {
}

func (r *articleRepository) Get(db *gorm.DB, id int64) *models.Article {
	ret := &models.Article{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *articleRepository) Take(db *gorm.DB, where ...interface{}) *models.Article {
	ret := &models.Article{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *articleRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.Article) {
	cnd.Find(db, &list)
	return
}

func (r *articleRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.Article {
	ret := &models.Article{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *articleRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.Article, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *articleRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.Article, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.Article{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *articleRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.Article{})
}

func (r *articleRepository) Create(db *gorm.DB, t *models.Article) (err error) {
	err = db.Create(t).Error
	return
}

func (r *articleRepository) Update(db *gorm.DB, t *models.Article) (err error) {
	err = db.Save(t).Error
	return
}

func (r *articleRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.Article{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *articleRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.Article{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *articleRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.Article{}, "id = ?", id)
}
