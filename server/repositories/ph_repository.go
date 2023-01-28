package repositories

import (
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"

	"server/model"
)

var PhRepository = newPhRepository()

func newPhRepository() *phRepository {
	return &phRepository{}
}

type phRepository struct {
}

func (r *phRepository) Get(db *gorm.DB, id int64) *model.PurchaseHistory {
	ret := &model.PurchaseHistory{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *phRepository) Take(db *gorm.DB, where ...interface{}) *model.PurchaseHistory {
	ret := &model.PurchaseHistory{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *phRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []model.PurchaseHistory) {
	cnd.Find(db, &list)
	return
}

func (r *phRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *model.PurchaseHistory {
	ret := &model.PurchaseHistory{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *phRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []model.PurchaseHistory, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *phRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []model.PurchaseHistory, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.PurchaseHistory{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *phRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &model.PurchaseHistory{})
}

func (r *phRepository) Create(db *gorm.DB, t *model.PurchaseHistory) (err error) {
	err = db.Create(t).Error
	return
}

func (r *phRepository) Update(db *gorm.DB, t *model.PurchaseHistory) (err error) {
	err = db.Save(t).Error
	return
}

func (r *phRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.PurchaseHistory{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *phRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.PurchaseHistory{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *phRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.PurchaseHistory{}, "id = ?", id)
}
