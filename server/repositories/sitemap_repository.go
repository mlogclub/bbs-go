package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"bbs-go/model"
)

var SitemapRepository = newSitemapRepository()

func newSitemapRepository() *sitemapRepository {
	return &sitemapRepository{}
}

type sitemapRepository struct {
}

func (r *sitemapRepository) Get(db *gorm.DB, id int64) *model.Sitemap {
	ret := &model.Sitemap{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *sitemapRepository) Take(db *gorm.DB, where ...interface{}) *model.Sitemap {
	ret := &model.Sitemap{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *sitemapRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Sitemap) {
	cnd.Find(db, &list)
	return
}

func (r *sitemapRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) *model.Sitemap {
	ret := &model.Sitemap{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *sitemapRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.Sitemap, paging *simple.Paging) {
	return r.FindPageByCnd(db, &params.SqlCnd)
}

func (r *sitemapRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Sitemap, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Sitemap{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *sitemapRepository) Create(db *gorm.DB, t *model.Sitemap) (err error) {
	err = db.Create(t).Error
	return
}

func (r *sitemapRepository) Update(db *gorm.DB, t *model.Sitemap) (err error) {
	err = db.Save(t).Error
	return
}

func (r *sitemapRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Sitemap{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *sitemapRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Sitemap{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *sitemapRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.Sitemap{}, "id = ?", id)
}
