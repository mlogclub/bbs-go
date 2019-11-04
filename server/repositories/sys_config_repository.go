package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
)

var SysConfigRepository = newSysConfigRepository()

func newSysConfigRepository() *sysConfigRepository {
	return &sysConfigRepository{}
}

type sysConfigRepository struct {
}

func (this *sysConfigRepository) Get(db *gorm.DB, id int64) *model.SysConfig {
	ret := &model.SysConfig{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *sysConfigRepository) Take(db *gorm.DB, where ...interface{}) *model.SysConfig {
	ret := &model.SysConfig{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *sysConfigRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.SysConfig) {
	cnd.Find(db, &list)
	return
}

func (this *sysConfigRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.SysConfig) {
	cnd.FindOne(db, &ret)
	return
}

func (this *sysConfigRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.SysConfig, paging *simple.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *sysConfigRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.SysConfig, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.SysConfig{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *sysConfigRepository) Create(db *gorm.DB, t *model.SysConfig) (err error) {
	err = db.Create(t).Error
	return
}

func (this *sysConfigRepository) Update(db *gorm.DB, t *model.SysConfig) (err error) {
	err = db.Save(t).Error
	return
}

func (this *sysConfigRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.SysConfig{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *sysConfigRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.SysConfig{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *sysConfigRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.SysConfig{}, "id = ?", id)
}

func (this *sysConfigRepository) GetByKey(db *gorm.DB, key string) *model.SysConfig {
	if len(key) == 0 {
		return nil
	}
	return this.Take(db, "`key` = ?", key)
}
