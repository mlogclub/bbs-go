package repositories

import (
	"bbs-go/internal/models"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var UserTaskEventRepository = newUserTaskEventRepository()

func newUserTaskEventRepository() *userTaskEventRepository {
	return &userTaskEventRepository{}
}

type userTaskEventRepository struct {
}

func (r *userTaskEventRepository) Get(db *gorm.DB, id int64) *models.UserTaskEvent {
	ret := &models.UserTaskEvent{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userTaskEventRepository) Take(db *gorm.DB, where ...interface{}) *models.UserTaskEvent {
	ret := &models.UserTaskEvent{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userTaskEventRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.UserTaskEvent) {
	cnd.Find(db, &list)
	return
}

func (r *userTaskEventRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.UserTaskEvent {
	ret := &models.UserTaskEvent{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *userTaskEventRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.UserTaskEvent, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *userTaskEventRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.UserTaskEvent, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.UserTaskEvent{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *userTaskEventRepository) FindBySql(db *gorm.DB, sqlStr string, paramArr... interface{}) (list []models.UserTaskEvent) {
	db.Raw(sqlStr, paramArr...).Scan(&list)
	return
}

func (r *userTaskEventRepository) CountBySql(db *gorm.DB, sqlStr string, paramArr... interface{}) (count int64) {
	db.Raw(sqlStr, paramArr...).Count(&count)
	return
}

func (r *userTaskEventRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.UserTaskEvent{})
}

func (r *userTaskEventRepository) Create(db *gorm.DB, t *models.UserTaskEvent) (err error) {
	err = db.Create(t).Error
	return
}

func (r *userTaskEventRepository) Update(db *gorm.DB, t *models.UserTaskEvent) (err error) {
	err = db.Save(t).Error
	return
}

func (r *userTaskEventRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.UserTaskEvent{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *userTaskEventRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.UserTaskEvent{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *userTaskEventRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.UserTaskEvent{}, "id = ?", id)
}

