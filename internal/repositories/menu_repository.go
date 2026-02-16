package repositories

import (
	"bbs-go/internal/models"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var MenuRepository = newMenuRepository()

func newMenuRepository() *menuRepository {
	return &menuRepository{}
}

type menuRepository struct {
}

func (r *menuRepository) Get(db *gorm.DB, id int64) *models.Menu {
	ret := &models.Menu{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *menuRepository) Take(db *gorm.DB, where ...interface{}) *models.Menu {
	ret := &models.Menu{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *menuRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.Menu) {
	cnd.Find(db, &list)
	return
}

func (r *menuRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.Menu {
	ret := &models.Menu{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *menuRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.Menu, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *menuRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.Menu, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.Menu{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *menuRepository) FindBySql(db *gorm.DB, sqlStr string, paramArr ...interface{}) (list []models.Menu) {
	db.Raw(sqlStr, paramArr...).Scan(&list)
	return
}

func (r *menuRepository) CountBySql(db *gorm.DB, sqlStr string, paramArr ...interface{}) (count int64) {
	db.Raw(sqlStr, paramArr...).Count(&count)
	return
}

func (r *menuRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.Menu{})
}

func (r *menuRepository) Create(db *gorm.DB, t *models.Menu) (err error) {
	err = db.Create(t).Error
	return
}

func (r *menuRepository) Update(db *gorm.DB, t *models.Menu) (err error) {
	err = db.Save(t).Error
	return
}

func (r *menuRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.Menu{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *menuRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.Menu{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *menuRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.Menu{}, "id = ?", id)
}
