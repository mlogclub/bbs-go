package repositories

import (
	"bbs-go/internal/models"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var DictTypeRepository = newDictTypeRepository()

func newDictTypeRepository() *dictTypeRepository {
	return &dictTypeRepository{}
}

type dictTypeRepository struct {
}

func (r *dictTypeRepository) Get(db *gorm.DB, id int64) *models.DictType {
	ret := &models.DictType{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *dictTypeRepository) Take(db *gorm.DB, where ...interface{}) *models.DictType {
	ret := &models.DictType{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *dictTypeRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.DictType) {
	cnd.Find(db, &list)
	return
}

func (r *dictTypeRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.DictType {
	ret := &models.DictType{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *dictTypeRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.DictType, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *dictTypeRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.DictType, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.DictType{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *dictTypeRepository) FindBySql(db *gorm.DB, sqlStr string, paramArr... interface{}) (list []models.DictType) {
	db.Raw(sqlStr, paramArr...).Scan(&list)
	return
}

func (r *dictTypeRepository) CountBySql(db *gorm.DB, sqlStr string, paramArr... interface{}) (count int64) {
	db.Raw(sqlStr, paramArr...).Count(&count)
	return
}

func (r *dictTypeRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.DictType{})
}

func (r *dictTypeRepository) Create(db *gorm.DB, t *models.DictType) (err error) {
	err = db.Create(t).Error
	return
}

func (r *dictTypeRepository) Update(db *gorm.DB, t *models.DictType) (err error) {
	err = db.Save(t).Error
	return
}

func (r *dictTypeRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.DictType{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *dictTypeRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.DictType{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *dictTypeRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.DictType{}, "id = ?", id)
}

