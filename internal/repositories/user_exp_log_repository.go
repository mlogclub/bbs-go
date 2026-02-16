package repositories

import (
	"bbs-go/internal/models"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var UserExpLogRepository = newUserExpLogRepository()

func newUserExpLogRepository() *userExpLogRepository {
	return &userExpLogRepository{}
}

type userExpLogRepository struct {
}

func (r *userExpLogRepository) Get(db *gorm.DB, id int64) *models.UserExpLog {
	ret := &models.UserExpLog{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userExpLogRepository) Take(db *gorm.DB, where ...interface{}) *models.UserExpLog {
	ret := &models.UserExpLog{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userExpLogRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.UserExpLog) {
	cnd.Find(db, &list)
	return
}

func (r *userExpLogRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.UserExpLog {
	ret := &models.UserExpLog{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *userExpLogRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.UserExpLog, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *userExpLogRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.UserExpLog, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.UserExpLog{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *userExpLogRepository) FindBySql(db *gorm.DB, sqlStr string, paramArr... interface{}) (list []models.UserExpLog) {
	db.Raw(sqlStr, paramArr...).Scan(&list)
	return
}

func (r *userExpLogRepository) CountBySql(db *gorm.DB, sqlStr string, paramArr... interface{}) (count int64) {
	db.Raw(sqlStr, paramArr...).Count(&count)
	return
}

func (r *userExpLogRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.UserExpLog{})
}

func (r *userExpLogRepository) Create(db *gorm.DB, t *models.UserExpLog) (err error) {
	err = db.Create(t).Error
	return
}

func (r *userExpLogRepository) Update(db *gorm.DB, t *models.UserExpLog) (err error) {
	err = db.Save(t).Error
	return
}

func (r *userExpLogRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.UserExpLog{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *userExpLogRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.UserExpLog{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *userExpLogRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.UserExpLog{}, "id = ?", id)
}

