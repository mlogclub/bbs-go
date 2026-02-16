package repositories

import (
	"bbs-go/internal/models"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var UserReportRepository = newUserReportRepository()

func newUserReportRepository() *userReportRepository {
	return &userReportRepository{}
}

type userReportRepository struct {
}

func (r *userReportRepository) Get(db *gorm.DB, id int64) *models.UserReport {
	ret := &models.UserReport{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userReportRepository) Take(db *gorm.DB, where ...interface{}) *models.UserReport {
	ret := &models.UserReport{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userReportRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.UserReport) {
	cnd.Find(db, &list)
	return
}

func (r *userReportRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.UserReport {
	ret := &models.UserReport{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *userReportRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.UserReport, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *userReportRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.UserReport, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.UserReport{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *userReportRepository) FindBySql(db *gorm.DB, sqlStr string, paramArr ...interface{}) (list []models.UserReport) {
	db.Raw(sqlStr, paramArr...).Scan(&list)
	return
}

func (r *userReportRepository) CountBySql(db *gorm.DB, sqlStr string, paramArr ...interface{}) (count int64) {
	db.Raw(sqlStr, paramArr...).Count(&count)
	return
}

func (r *userReportRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.UserReport{})
}

func (r *userReportRepository) Create(db *gorm.DB, t *models.UserReport) (err error) {
	err = db.Create(t).Error
	return
}

func (r *userReportRepository) Update(db *gorm.DB, t *models.UserReport) (err error) {
	err = db.Save(t).Error
	return
}

func (r *userReportRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.UserReport{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *userReportRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.UserReport{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *userReportRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.UserReport{}, "id = ?", id)
}
