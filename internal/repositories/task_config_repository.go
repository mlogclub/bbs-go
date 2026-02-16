package repositories

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"

	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var TaskConfigRepository = newTaskConfigRepository()

func newTaskConfigRepository() *taskConfigRepository {
	return &taskConfigRepository{}
}

type taskConfigRepository struct {
}

func (r *taskConfigRepository) Get(db *gorm.DB, id int64) *models.TaskConfig {
	ret := &models.TaskConfig{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *taskConfigRepository) FindByEventType(eventType string) (list []models.TaskConfig) {
	if strs.IsBlank(eventType) {
		return nil
	}
	return r.Find(sqls.DB(), sqls.NewCnd().
		Eq("event_type", eventType).
		Eq("status", constants.StatusOk).
		Asc("sort_no"))
}

func (r *taskConfigRepository) Take(db *gorm.DB, where ...interface{}) *models.TaskConfig {
	ret := &models.TaskConfig{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *taskConfigRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.TaskConfig) {
	cnd.Find(db, &list)
	return
}

func (r *taskConfigRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.TaskConfig {
	ret := &models.TaskConfig{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *taskConfigRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.TaskConfig, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *taskConfigRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.TaskConfig, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.TaskConfig{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *taskConfigRepository) FindBySql(db *gorm.DB, sqlStr string, paramArr ...interface{}) (list []models.TaskConfig) {
	db.Raw(sqlStr, paramArr...).Scan(&list)
	return
}

func (r *taskConfigRepository) CountBySql(db *gorm.DB, sqlStr string, paramArr ...interface{}) (count int64) {
	db.Raw(sqlStr, paramArr...).Count(&count)
	return
}

func (r *taskConfigRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.TaskConfig{})
}

func (r *taskConfigRepository) Create(db *gorm.DB, t *models.TaskConfig) (err error) {
	err = db.Create(t).Error
	return
}

func (r *taskConfigRepository) Update(db *gorm.DB, t *models.TaskConfig) (err error) {
	err = db.Save(t).Error
	return
}

func (r *taskConfigRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.TaskConfig{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *taskConfigRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.TaskConfig{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *taskConfigRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.TaskConfig{}, "id = ?", id)
}
