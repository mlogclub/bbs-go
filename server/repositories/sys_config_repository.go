package repositories

import (
	"github.com/mlogclub/simple"
	"gorm.io/gorm"

	"bbs-go/model"
)

var SysConfigRepository = newSysConfigRepository()

func newSysConfigRepository() *sysConfigRepository {
	return &sysConfigRepository{}
}

type sysConfigRepository struct {
}

func (r *sysConfigRepository) Get(db *gorm.DB, id int64) *model.SysConfig {
	ret := &model.SysConfig{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *sysConfigRepository) Take(db *gorm.DB, where ...interface{}) *model.SysConfig {
	ret := &model.SysConfig{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *sysConfigRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.SysConfig) {
	cnd.Find(db, &list)
	return
}

func (r *sysConfigRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) *model.SysConfig {
	ret := &model.SysConfig{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *sysConfigRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.SysConfig, paging *simple.Paging) {
	return r.FindPageByCnd(db, &params.SqlCnd)
}

func (r *sysConfigRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.SysConfig, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.SysConfig{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *sysConfigRepository) Create(db *gorm.DB, t *model.SysConfig) (err error) {
	err = db.Create(t).Error
	return
}

func (r *sysConfigRepository) Update(db *gorm.DB, t *model.SysConfig) (err error) {
	err = db.Save(t).Error
	return
}

func (r *sysConfigRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.SysConfig{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *sysConfigRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.SysConfig{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *sysConfigRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.SysConfig{}, "id = ?", id)
}

func (r *sysConfigRepository) GetByKey(db *gorm.DB, key string) *model.SysConfig {
	if len(key) == 0 {
		return nil
	}
	return r.Take(db, "`key` = ?", key)
}
