
package repositories

import (
	"bbs-go/model"
	"github.com/mlogclub/simple"
	"github.com/jinzhu/gorm"
)

var UserScoreLogRepository = newUserScoreLogRepository()

func newUserScoreLogRepository() *userScoreLogRepository {
	return &userScoreLogRepository{}
}

type userScoreLogRepository struct {
}

func (this *userScoreLogRepository) Get(db *gorm.DB, id int64) *model.UserScoreLog {
	ret := &model.UserScoreLog{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *userScoreLogRepository) Take(db *gorm.DB, where ...interface{}) *model.UserScoreLog {
	ret := &model.UserScoreLog{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *userScoreLogRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.UserScoreLog) {
	cnd.Find(db, &list)
	return
}

func (this *userScoreLogRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) *model.UserScoreLog {
	ret := &model.UserScoreLog{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (this *userScoreLogRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.UserScoreLog, paging *simple.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *userScoreLogRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.UserScoreLog, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.UserScoreLog{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *userScoreLogRepository) Create(db *gorm.DB, t *model.UserScoreLog) (err error) {
	err = db.Create(t).Error
	return
}

func (this *userScoreLogRepository) Update(db *gorm.DB, t *model.UserScoreLog) (err error) {
	err = db.Save(t).Error
	return
}

func (this *userScoreLogRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.UserScoreLog{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *userScoreLogRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.UserScoreLog{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *userScoreLogRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.UserScoreLog{}, "id = ?", id)
}

