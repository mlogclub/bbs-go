package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/simple"
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

func (this *sysConfigRepository) QueryCnd(db *gorm.DB, cnd *simple.QueryCnd) (list []model.SysConfig, err error) {
	err = cnd.DoQuery(db).Find(&list).Error
	return
}

func (this *sysConfigRepository) Query(db *gorm.DB, queries *simple.ParamQueries) (list []model.SysConfig, paging *simple.Paging) {
	queries.StartQuery(db).Find(&list)
	queries.StartCount(db).Model(&model.SysConfig{}).Count(&queries.Paging.Total)
	paging = queries.Paging
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
