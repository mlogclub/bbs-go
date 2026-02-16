package repositories

import (
	"bbs-go/internal/models"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var MigrationRepository = newMigrationRepository()

func newMigrationRepository() *migrationRepository {
	return &migrationRepository{}
}

type migrationRepository struct {
}

func (r *migrationRepository) Get(db *gorm.DB, id int64) *models.Migration {
	ret := &models.Migration{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *migrationRepository) Take(db *gorm.DB, where ...interface{}) *models.Migration {
	ret := &models.Migration{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *migrationRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.Migration) {
	cnd.Find(db, &list)
	return
}

func (r *migrationRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.Migration {
	ret := &models.Migration{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *migrationRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.Migration, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *migrationRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.Migration, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.Migration{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *migrationRepository) FindBySql(db *gorm.DB, sqlStr string, paramArr ...interface{}) (list []models.Migration) {
	db.Raw(sqlStr, paramArr...).Scan(&list)
	return
}

func (r *migrationRepository) CountBySql(db *gorm.DB, sqlStr string, paramArr ...interface{}) (count int64) {
	db.Raw(sqlStr, paramArr...).Count(&count)
	return
}

func (r *migrationRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.Migration{})
}

func (r *migrationRepository) Create(db *gorm.DB, t *models.Migration) (err error) {
	err = db.Create(t).Error
	return
}

func (r *migrationRepository) Update(db *gorm.DB, t *models.Migration) (err error) {
	err = db.Save(t).Error
	return
}

func (r *migrationRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.Migration{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *migrationRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.Migration{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *migrationRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.Migration{}, "id = ?", id)
}
