package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/simple"
)

var TopicRepository = newTopicRepository()

func newTopicRepository() *topicRepository {
	return &topicRepository{}
}

type topicRepository struct {
}

func (this *topicRepository) Get(db *gorm.DB, id int64) *model.Topic {
	ret := &model.Topic{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *topicRepository) Take(db *gorm.DB, where ...interface{}) *model.Topic {
	ret := &model.Topic{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *topicRepository) QueryCnd(db *gorm.DB, cnd *simple.QueryCnd) (list []model.Topic, err error) {
	err = cnd.DoQuery(db).Find(&list).Error
	return
}

func (this *topicRepository) Query(db *gorm.DB, queries *simple.ParamQueries) (list []model.Topic, paging *simple.Paging) {
	queries.StartQuery(db).Find(&list)
	queries.StartCount(db).Model(&model.Topic{}).Count(&queries.Paging.Total)
	paging = queries.Paging
	return
}

func (this *topicRepository) Create(db *gorm.DB, t *model.Topic) (err error) {
	err = db.Create(t).Error
	return
}

func (this *topicRepository) Update(db *gorm.DB, t *model.Topic) (err error) {
	err = db.Save(t).Error
	return
}

func (this *topicRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Topic{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *topicRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Topic{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *topicRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.Topic{}, "id = ?", id)
}
