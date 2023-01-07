package repositories

import (
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"

	"bbs-go/model"
)

var AdvertRepository = newAdvertRepository()

func newAdvertRepository() *advertRepository {
	return &advertRepository{}
}

type advertRepository struct {
}

func (r *advertRepository) Get(db *gorm.DB, id int64) *model.Advert {
	ret := &model.Advert{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *advertRepository) Take(db *gorm.DB, where ...interface{}) *model.Advert {
	ret := &model.Advert{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *advertRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []model.Advert) {
	cnd.Find(db, &list)
	return
}

func (r *advertRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *model.Advert {
	ret := &model.Advert{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *advertRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []model.Advert, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *advertRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []model.Advert, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Advert{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *advertRepository) Create(db *gorm.DB, t *model.Advert) (err error) {
	err = db.Create(t).Error
	return
}

func (r *advertRepository) Update(db *gorm.DB, t *model.Advert) (err error) {
	err = db.Save(t).Error
	return
}

func (r *advertRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Advert{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *advertRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Advert{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *advertRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.Advert{}, "id = ?", id)
}
