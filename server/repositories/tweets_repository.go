package repositories

import (
	"bbs-go/model"
	"github.com/mlogclub/simple"
	"github.com/jinzhu/gorm"
)

var TweetsRepository = newTweetsRepository()

func newTweetsRepository() *tweetsRepository {
	return &tweetsRepository{}
}

type tweetsRepository struct {
}

func (r *tweetsRepository) Get(db *gorm.DB, id int64) *model.Tweets {
	ret := &model.Tweets{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *tweetsRepository) Take(db *gorm.DB, where ...interface{}) *model.Tweets {
	ret := &model.Tweets{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *tweetsRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Tweets) {
	cnd.Find(db, &list)
	return
}

func (r *tweetsRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) *model.Tweets {
	ret := &model.Tweets{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *tweetsRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.Tweets, paging *simple.Paging) {
	return r.FindPageByCnd(db, &params.SqlCnd)
}

func (r *tweetsRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Tweets, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Tweets{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *tweetsRepository) Count(db *gorm.DB, cnd *simple.SqlCnd) int {
	return cnd.Count(db, &model.Tweets{})
}

func (r *tweetsRepository) Create(db *gorm.DB, t *model.Tweets) (err error) {
	err = db.Create(t).Error
	return
}

func (r *tweetsRepository) Update(db *gorm.DB, t *model.Tweets) (err error) {
	err = db.Save(t).Error
	return
}

func (r *tweetsRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Tweets{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *tweetsRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Tweets{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *tweetsRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.Tweets{}, "id = ?", id)
}

