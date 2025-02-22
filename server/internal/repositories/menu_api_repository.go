package repositories

import (
	"bbs-go/internal/models"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var MenuApiRepository = newMenuApiRepository()

func newMenuApiRepository() *menuApiRepository {
	return &menuApiRepository{}
}

type menuApiRepository struct {
}

func (r *menuApiRepository) Get(db *gorm.DB, id int64) *models.MenuApi {
	ret := &models.MenuApi{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *menuApiRepository) Take(db *gorm.DB, where ...interface{}) *models.MenuApi {
	ret := &models.MenuApi{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *menuApiRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.MenuApi) {
	cnd.Find(db, &list)
	return
}

func (r *menuApiRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.MenuApi {
	ret := &models.MenuApi{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *menuApiRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.MenuApi, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *menuApiRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.MenuApi, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.MenuApi{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *menuApiRepository) FindBySql(db *gorm.DB, sqlStr string, paramArr... interface{}) (list []models.MenuApi) {
	db.Raw(sqlStr, paramArr...).Scan(&list)
	return
}

func (r *menuApiRepository) CountBySql(db *gorm.DB, sqlStr string, paramArr... interface{}) (count int64) {
	db.Raw(sqlStr, paramArr...).Count(&count)
	return
}

func (r *menuApiRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.MenuApi{})
}

func (r *menuApiRepository) Create(db *gorm.DB, t *models.MenuApi) (err error) {
	err = db.Create(t).Error
	return
}

func (r *menuApiRepository) Update(db *gorm.DB, t *models.MenuApi) (err error) {
	err = db.Save(t).Error
	return
}

func (r *menuApiRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.MenuApi{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *menuApiRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.MenuApi{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *menuApiRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.MenuApi{}, "id = ?", id)
}

