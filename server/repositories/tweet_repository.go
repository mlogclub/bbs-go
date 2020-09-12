package repositories

import (
	"github.com/mlogclub/simple"
	"gorm.io/gorm"

	"bbs-go/model"
)

var TweetRepository = newTweetRepository()

func newTweetRepository() *tweetRepository {
	return &tweetRepository{}
}

type tweetRepository struct {
}

func (r *tweetRepository) Get(db *gorm.DB, id int64) *model.Tweet {
	ret := &model.Tweet{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *tweetRepository) Take(db *gorm.DB, where ...interface{}) *model.Tweet {
	ret := &model.Tweet{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *tweetRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Tweet) {
	cnd.Find(db, &list)
	return
}

func (r *tweetRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) *model.Tweet {
	ret := &model.Tweet{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *tweetRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.Tweet, paging *simple.Paging) {
	return r.FindPageByCnd(db, &params.SqlCnd)
}

func (r *tweetRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Tweet, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Tweet{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *tweetRepository) Count(db *gorm.DB, cnd *simple.SqlCnd) int64 {
	return cnd.Count(db, &model.Tweet{})
}

func (r *tweetRepository) Create(db *gorm.DB, t *model.Tweet) (err error) {
	err = db.Create(t).Error
	return
}

func (r *tweetRepository) Update(db *gorm.DB, t *model.Tweet) (err error) {
	err = db.Save(t).Error
	return
}

func (r *tweetRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Tweet{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *tweetRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Tweet{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *tweetRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.Tweet{}, "id = ?", id)
}
