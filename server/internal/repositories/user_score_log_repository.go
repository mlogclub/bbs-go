package repositories

import (
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"

	"bbs-go/internal/models"
)

var UserScoreLogRepository = newUserScoreLogRepository()

func newUserScoreLogRepository() *userScoreLogRepository {
	return &userScoreLogRepository{}
}

type userScoreLogRepository struct {
}

func (r *userScoreLogRepository) Get(db *gorm.DB, id int64) *models.UserScoreLog {
	ret := &models.UserScoreLog{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userScoreLogRepository) Take(db *gorm.DB, where ...interface{}) *models.UserScoreLog {
	ret := &models.UserScoreLog{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userScoreLogRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.UserScoreLog) {
	cnd.Find(db, &list)
	return
}

func (r *userScoreLogRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.UserScoreLog {
	ret := &models.UserScoreLog{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *userScoreLogRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.UserScoreLog, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *userScoreLogRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.UserScoreLog, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.UserScoreLog{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *userScoreLogRepository) Create(db *gorm.DB, t *models.UserScoreLog) (err error) {
	err = db.Create(t).Error
	return
}

func (r *userScoreLogRepository) Update(db *gorm.DB, t *models.UserScoreLog) (err error) {
	err = db.Save(t).Error
	return
}

func (r *userScoreLogRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.UserScoreLog{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *userScoreLogRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.UserScoreLog{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *userScoreLogRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.UserScoreLog{}, "id = ?", id)
}
