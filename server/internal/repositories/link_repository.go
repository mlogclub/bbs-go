package repositories

import (
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"

	"bbs-go/internal/models"
)

var LinkRepository = newLinkRepository()

func newLinkRepository() *linkRepository {
	return &linkRepository{}
}

type linkRepository struct {
}

func (r *linkRepository) Get(db *gorm.DB, id int64) *models.Link {
	ret := &models.Link{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *linkRepository) Take(db *gorm.DB, where ...interface{}) *models.Link {
	ret := &models.Link{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *linkRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.Link) {
	cnd.Find(db, &list)
	return
}

func (r *linkRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.Link {
	ret := &models.Link{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *linkRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.Link, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *linkRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.Link, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.Link{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *linkRepository) Create(db *gorm.DB, t *models.Link) (err error) {
	err = db.Create(t).Error
	return
}

func (r *linkRepository) Update(db *gorm.DB, t *models.Link) (err error) {
	err = db.Save(t).Error
	return
}

func (r *linkRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.Link{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *linkRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.Link{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *linkRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.Link{}, "id = ?", id)
}
