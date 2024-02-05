package repositories

import (
	"bbs-go/internal/models"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var RoleMenuRepository = newRoleMenuRepository()

func newRoleMenuRepository() *roleMenuRepository {
	return &roleMenuRepository{}
}

type roleMenuRepository struct {
}

func (r *roleMenuRepository) Get(db *gorm.DB, id int64) *models.RoleMenu {
	ret := &models.RoleMenu{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *roleMenuRepository) Take(db *gorm.DB, where ...interface{}) *models.RoleMenu {
	ret := &models.RoleMenu{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *roleMenuRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.RoleMenu) {
	cnd.Find(db, &list)
	return
}

func (r *roleMenuRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.RoleMenu {
	ret := &models.RoleMenu{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *roleMenuRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.RoleMenu, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *roleMenuRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.RoleMenu, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.RoleMenu{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *roleMenuRepository) FindBySql(db *gorm.DB, sqlStr string, paramArr ...interface{}) (list []models.RoleMenu) {
	db.Raw(sqlStr, paramArr...).Scan(&list)
	return
}

func (r *roleMenuRepository) CountBySql(db *gorm.DB, sqlStr string, paramArr ...interface{}) (count int64) {
	db.Raw(sqlStr, paramArr...).Count(&count)
	return
}

func (r *roleMenuRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.RoleMenu{})
}

func (r *roleMenuRepository) Create(db *gorm.DB, t *models.RoleMenu) (err error) {
	err = db.Create(t).Error
	return
}

func (r *roleMenuRepository) Update(db *gorm.DB, t *models.RoleMenu) (err error) {
	err = db.Save(t).Error
	return
}

func (r *roleMenuRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.RoleMenu{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *roleMenuRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.RoleMenu{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *roleMenuRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.RoleMenu{}, "id = ?", id)
}
