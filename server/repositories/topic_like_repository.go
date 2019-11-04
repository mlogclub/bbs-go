package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
)

var TopicLikeRepository = newTopicLikeRepository()

func newTopicLikeRepository() *topicLikeRepository {
	return &topicLikeRepository{}
}

type topicLikeRepository struct {
}

func (this *topicLikeRepository) Get(db *gorm.DB, id int64) *model.TopicLike {
	ret := &model.TopicLike{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *topicLikeRepository) Take(db *gorm.DB, where ...interface{}) *model.TopicLike {
	ret := &model.TopicLike{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *topicLikeRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.TopicLike) {
	cnd.Find(db, &list)
	return
}

func (this *topicLikeRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.TopicLike) {
	cnd.FindOne(db, &ret)
	return
}

func (this *topicLikeRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.TopicLike, paging *simple.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *topicLikeRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.TopicLike, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.TopicLike{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *topicLikeRepository) Create(db *gorm.DB, t *model.TopicLike) (err error) {
	err = db.Create(t).Error
	return
}

func (this *topicLikeRepository) Update(db *gorm.DB, t *model.TopicLike) (err error) {
	err = db.Save(t).Error
	return
}

func (this *topicLikeRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.TopicLike{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *topicLikeRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.TopicLike{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *topicLikeRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.TopicLike{}, "id = ?", id)
}
