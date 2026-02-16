package repositories

import (
	"bbs-go/internal/models"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var UserRoleRepository = newUserRoleRepository()

func newUserRoleRepository() *userRoleRepository {
	return &userRoleRepository{}
}

type userRoleRepository struct {
}

func (r *userRoleRepository) Get(db *gorm.DB, id int64) *models.UserRole {
	ret := &models.UserRole{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userRoleRepository) Take(db *gorm.DB, where ...interface{}) *models.UserRole {
	ret := &models.UserRole{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userRoleRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.UserRole) {
	cnd.Find(db, &list)
	return
}

func (r *userRoleRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.UserRole {
	ret := &models.UserRole{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *userRoleRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.UserRole, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *userRoleRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.UserRole, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.UserRole{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *userRoleRepository) FindBySql(db *gorm.DB, sqlStr string, paramArr ...interface{}) (list []models.UserRole) {
	db.Raw(sqlStr, paramArr...).Scan(&list)
	return
}

func (r *userRoleRepository) CountBySql(db *gorm.DB, sqlStr string, paramArr ...interface{}) (count int64) {
	db.Raw(sqlStr, paramArr...).Count(&count)
	return
}

func (r *userRoleRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.UserRole{})
}

func (r *userRoleRepository) Create(db *gorm.DB, t *models.UserRole) (err error) {
	err = db.Create(t).Error
	return
}

func (r *userRoleRepository) Update(db *gorm.DB, t *models.UserRole) (err error) {
	err = db.Save(t).Error
	return
}

func (r *userRoleRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.UserRole{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *userRoleRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.UserRole{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *userRoleRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.UserRole{}, "id = ?", id)
}
