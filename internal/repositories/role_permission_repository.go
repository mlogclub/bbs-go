package repositories

import (
	"bbs-go/internal/models"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var RolePermissionRepository = newRolePermissionRepository()

func newRolePermissionRepository() *rolePermissionRepository {
	return &rolePermissionRepository{}
}

type rolePermissionRepository struct {
}

func (r *rolePermissionRepository) Get(db *gorm.DB, id int64) *models.RolePermission {
	ret := &models.RolePermission{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *rolePermissionRepository) Take(db *gorm.DB, where ...interface{}) *models.RolePermission {
	ret := &models.RolePermission{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *rolePermissionRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.RolePermission) {
	cnd.Find(db, &list)
	return
}

func (r *rolePermissionRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.RolePermission {
	ret := &models.RolePermission{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *rolePermissionRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.RolePermission, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *rolePermissionRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.RolePermission, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.RolePermission{})
	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *rolePermissionRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.RolePermission{})
}

func (r *rolePermissionRepository) Create(db *gorm.DB, t *models.RolePermission) error {
	return db.Create(t).Error
}

func (r *rolePermissionRepository) Update(db *gorm.DB, t *models.RolePermission) error {
	return db.Save(t).Error
}

func (r *rolePermissionRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.RolePermission{}, "id = ?", id)
}
