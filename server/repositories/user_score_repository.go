package repositories

import (
	"github.com/mlogclub/simple"
	"gorm.io/gorm"

	"bbs-go/model"
)

var UserScoreRepository = newUserScoreRepository()

func newUserScoreRepository() *userScoreRepository {
	return &userScoreRepository{}
}

type userScoreRepository struct {
}

func (r *userScoreRepository) Get(db *gorm.DB, id int64) *model.UserScore {
	ret := &model.UserScore{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userScoreRepository) Take(db *gorm.DB, where ...interface{}) *model.UserScore {
	ret := &model.UserScore{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userScoreRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.UserScore) {
	cnd.Find(db, &list)
	return
}

func (r *userScoreRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) *model.UserScore {
	ret := &model.UserScore{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *userScoreRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.UserScore, paging *simple.Paging) {
	return r.FindPageByCnd(db, &params.SqlCnd)
}

func (r *userScoreRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.UserScore, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.UserScore{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *userScoreRepository) Create(db *gorm.DB, t *model.UserScore) (err error) {
	err = db.Create(t).Error
	return
}

func (r *userScoreRepository) Update(db *gorm.DB, t *model.UserScore) (err error) {
	err = db.Save(t).Error
	return
}

func (r *userScoreRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.UserScore{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *userScoreRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.UserScore{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *userScoreRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.UserScore{}, "id = ?", id)
}
