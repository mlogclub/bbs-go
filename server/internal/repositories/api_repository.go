package repositories

import (
	"bbs-go/internal/models"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var ApiRepository = newApiRepository()

func newApiRepository() *apiRepository {
	return &apiRepository{}
}

type apiRepository struct {
}

func (r *apiRepository) Get(db *gorm.DB, id int64) *models.Api {
	ret := &models.Api{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *apiRepository) Take(db *gorm.DB, where ...interface{}) *models.Api {
	ret := &models.Api{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *apiRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.Api) {
	cnd.Find(db, &list)
	return
}

func (r *apiRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.Api {
	ret := &models.Api{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *apiRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.Api, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *apiRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.Api, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.Api{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *apiRepository) FindBySql(db *gorm.DB, sqlStr string, paramArr... interface{}) (list []models.Api) {
	db.Raw(sqlStr, paramArr...).Scan(&list)
	return
}

func (r *apiRepository) CountBySql(db *gorm.DB, sqlStr string, paramArr... interface{}) (count int64) {
	db.Raw(sqlStr, paramArr...).Count(&count)
	return
}

func (r *apiRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.Api{})
}

func (r *apiRepository) Create(db *gorm.DB, t *models.Api) (err error) {
	err = db.Create(t).Error
	return
}

func (r *apiRepository) Update(db *gorm.DB, t *models.Api) (err error) {
	err = db.Save(t).Error
	return
}

func (r *apiRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.Api{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *apiRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.Api{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *apiRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.Api{}, "id = ?", id)
}

