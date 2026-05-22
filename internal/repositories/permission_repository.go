package repositories

import (
	"bbs-go/internal/models"

	"bbs-go/internal/pkg/params"

	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

var PermissionRepository = newPermissionRepository()

func newPermissionRepository() *permissionRepository {
	return &permissionRepository{}
}

type permissionRepository struct {
}

func (r *permissionRepository) Get(db *gorm.DB, id int64) *models.Permission {
	ret := &models.Permission{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *permissionRepository) Take(db *gorm.DB, where ...interface{}) *models.Permission {
	ret := &models.Permission{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *permissionRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.Permission) {
	cnd.Find(db, &list)
	return
}

func (r *permissionRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.Permission {
	ret := &models.Permission{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *permissionRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.Permission, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *permissionRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.Permission, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.Permission{})
	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *permissionRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.Permission{})
}

func (r *permissionRepository) Create(db *gorm.DB, t *models.Permission) error {
	return db.Create(t).Error
}

func (r *permissionRepository) Update(db *gorm.DB, t *models.Permission) error {
	return db.Save(t).Error
}

func (r *permissionRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) error {
	return db.Model(&models.Permission{}).Where("id = ?", id).Updates(columns).Error
}

func (r *permissionRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.Permission{}, "id = ?", id)
}
