
package repositories

import (
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/simple"
	"github.com/jinzhu/gorm"
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

func (this *topicLikeRepository) QueryCnd(db *gorm.DB, cnd *simple.QueryCnd) (list []model.TopicLike, err error) {
	err = cnd.DoQuery(db).Find(&list).Error
	return
}

func (this *topicLikeRepository) Query(db *gorm.DB, params *simple.ParamQueries) (list []model.TopicLike, paging *simple.Paging) {
	params.StartQuery(db).Find(&list)
    params.StartCount(db).Model(&model.TopicLike{}).Count(&params.Paging.Total)
	paging = params.Paging
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

