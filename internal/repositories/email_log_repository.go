package repositories

import (
	"bbs-go/internal/models"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var EmailLogRepository = newEmailLogRepository()

func newEmailLogRepository() *emailLogRepository {
	return &emailLogRepository{}
}

type emailLogRepository struct {
}

func (r *emailLogRepository) Get(db *gorm.DB, id int64) *models.EmailLog {
	ret := &models.EmailLog{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *emailLogRepository) Take(db *gorm.DB, where ...interface{}) *models.EmailLog {
	ret := &models.EmailLog{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *emailLogRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.EmailLog) {
	cnd.Find(db, &list)
	return
}

func (r *emailLogRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.EmailLog {
	ret := &models.EmailLog{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *emailLogRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.EmailLog, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *emailLogRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.EmailLog, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.EmailLog{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *emailLogRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.EmailLog{})
}

func (r *emailLogRepository) Create(db *gorm.DB, t *models.EmailLog) (err error) {
	err = db.Create(t).Error
	return
}

func (r *emailLogRepository) Update(db *gorm.DB, t *models.EmailLog) (err error) {
	err = db.Save(t).Error
	return
}

func (r *emailLogRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.EmailLog{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *emailLogRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.EmailLog{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *emailLogRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.EmailLog{}, "id = ?", id)
}
