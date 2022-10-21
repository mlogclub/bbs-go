package repositories

import (
	"bbs-go/model"

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

func (r *userReportRepository) Get(db *gorm.DB, id int64) *model.UserReport {
	ret := &model.UserReport{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userReportRepository) Take(db *gorm.DB, where ...interface{}) *model.UserReport {
	ret := &model.UserReport{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userReportRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []model.UserReport) {
	cnd.Find(db, &list)
	return
}

func (r *userReportRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *model.UserReport {
	ret := &model.UserReport{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *userReportRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []model.UserReport, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *userReportRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []model.UserReport, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.UserReport{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *userReportRepository) FindBySql(db *gorm.DB, sqlStr string, paramArr... interface{}) (list []model.UserReport) {
	db.Raw(sqlStr, paramArr...).Scan(&list)
	return
}

func (r *userReportRepository) CountBySql(db *gorm.DB, sqlStr string, paramArr... interface{}) (count int64) {
	db.Raw(sqlStr, paramArr...).Count(&count)
	return
}

func (r *userReportRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &model.UserReport{})
}

func (r *userReportRepository) Create(db *gorm.DB, t *model.UserReport) (err error) {
	err = db.Create(t).Error
	return
}

func (r *userReportRepository) Update(db *gorm.DB, t *model.UserReport) (err error) {
	err = db.Save(t).Error
	return
}

func (r *userReportRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.UserReport{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *userReportRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.UserReport{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *userReportRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.UserReport{}, "id = ?", id)
}

