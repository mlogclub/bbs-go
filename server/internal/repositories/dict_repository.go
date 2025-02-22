package repositories

import (
	"bbs-go/internal/models"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var DictRepository = newDictRepository()

func newDictRepository() *dictRepository {
	return &dictRepository{}
}

type dictRepository struct {
}

func (r *dictRepository) Get(db *gorm.DB, id int64) *models.Dict {
	ret := &models.Dict{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *dictRepository) Take(db *gorm.DB, where ...interface{}) *models.Dict {
	ret := &models.Dict{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *dictRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.Dict) {
	cnd.Find(db, &list)
	return
}

func (r *dictRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.Dict {
	ret := &models.Dict{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *dictRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.Dict, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *dictRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.Dict, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.Dict{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *dictRepository) FindBySql(db *gorm.DB, sqlStr string, paramArr... interface{}) (list []models.Dict) {
	db.Raw(sqlStr, paramArr...).Scan(&list)
	return
}

func (r *dictRepository) CountBySql(db *gorm.DB, sqlStr string, paramArr... interface{}) (count int64) {
	db.Raw(sqlStr, paramArr...).Count(&count)
	return
}

func (r *dictRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.Dict{})
}

func (r *dictRepository) Create(db *gorm.DB, t *models.Dict) (err error) {
	err = db.Create(t).Error
	return
}

func (r *dictRepository) Update(db *gorm.DB, t *models.Dict) (err error) {
	err = db.Save(t).Error
	return
}

func (r *dictRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.Dict{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *dictRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.Dict{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *dictRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.Dict{}, "id = ?", id)
}

