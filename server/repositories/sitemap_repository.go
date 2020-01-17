
package repositories

import (
	"bbs-go/model"
	"github.com/mlogclub/simple"
	"github.com/jinzhu/gorm"
)

var SitemapRepository = newSitemapRepository()

func newSitemapRepository() *sitemapRepository {
	return &sitemapRepository{}
}

type sitemapRepository struct {
}

func (this *sitemapRepository) Get(db *gorm.DB, id int64) *model.Sitemap {
	ret := &model.Sitemap{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *sitemapRepository) Take(db *gorm.DB, where ...interface{}) *model.Sitemap {
	ret := &model.Sitemap{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *sitemapRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Sitemap) {
	cnd.Find(db, &list)
	return
}

func (this *sitemapRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) *model.Sitemap {
	ret := &model.Sitemap{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (this *sitemapRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.Sitemap, paging *simple.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *sitemapRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Sitemap, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Sitemap{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *sitemapRepository) Create(db *gorm.DB, t *model.Sitemap) (err error) {
	err = db.Create(t).Error
	return
}

func (this *sitemapRepository) Update(db *gorm.DB, t *model.Sitemap) (err error) {
	err = db.Save(t).Error
	return
}

func (this *sitemapRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Sitemap{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *sitemapRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Sitemap{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *sitemapRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.Sitemap{}, "id = ?", id)
}

