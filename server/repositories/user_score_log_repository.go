package repositories

import (
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"

	"bbs-go/model"
)

var UserScoreLogRepository = newUserScoreLogRepository()

func newUserScoreLogRepository() *userScoreLogRepository {
	return &userScoreLogRepository{}
}

type userScoreLogRepository struct {
}

func (r *userScoreLogRepository) Get(db *gorm.DB, id int64) *model.UserScoreLog {
	ret := &model.UserScoreLog{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userScoreLogRepository) Take(db *gorm.DB, where ...interface{}) *model.UserScoreLog {
	ret := &model.UserScoreLog{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userScoreLogRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []model.UserScoreLog) {
	cnd.Find(db, &list)
	return
}

func (r *userScoreLogRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *model.UserScoreLog {
	ret := &model.UserScoreLog{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *userScoreLogRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []model.UserScoreLog, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *userScoreLogRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []model.UserScoreLog, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.UserScoreLog{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *userScoreLogRepository) Create(db *gorm.DB, t *model.UserScoreLog) (err error) {
	err = db.Create(t).Error
	return
}

func (r *userScoreLogRepository) Update(db *gorm.DB, t *model.UserScoreLog) (err error) {
	err = db.Save(t).Error
	return
}

func (r *userScoreLogRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.UserScoreLog{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *userScoreLogRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.UserScoreLog{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *userScoreLogRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.UserScoreLog{}, "id = ?", id)
}
