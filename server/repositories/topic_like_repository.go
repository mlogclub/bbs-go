package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"bbs-go/model"
)

var TopicLikeRepository = newTopicLikeRepository()

func newTopicLikeRepository() *topicLikeRepository {
	return &topicLikeRepository{}
}

type topicLikeRepository struct {
}

func (r *topicLikeRepository) Get(db *gorm.DB, id int64) *model.TopicLike {
	ret := &model.TopicLike{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *topicLikeRepository) Take(db *gorm.DB, where ...interface{}) *model.TopicLike {
	ret := &model.TopicLike{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *topicLikeRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.TopicLike) {
	cnd.Find(db, &list)
	return
}

func (r *topicLikeRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) *model.TopicLike {
	ret := &model.TopicLike{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *topicLikeRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.TopicLike, paging *simple.Paging) {
	return r.FindPageByCnd(db, &params.SqlCnd)
}

func (r *topicLikeRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.TopicLike, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.TopicLike{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *topicLikeRepository) Create(db *gorm.DB, t *model.TopicLike) (err error) {
	err = db.Create(t).Error
	return
}

func (r *topicLikeRepository) Update(db *gorm.DB, t *model.TopicLike) (err error) {
	err = db.Save(t).Error
	return
}

func (r *topicLikeRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.TopicLike{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *topicLikeRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.TopicLike{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *topicLikeRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.TopicLike{}, "id = ?", id)
}
