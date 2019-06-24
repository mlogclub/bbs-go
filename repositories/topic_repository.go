
package repositories

import (
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/simple"
	"github.com/jinzhu/gorm"
)

type TopicRepository struct {
}

func NewTopicRepository() *TopicRepository {
	return &TopicRepository{}
}

func (this *TopicRepository) Get(db *gorm.DB, id int64) *model.Topic {
	ret := &model.Topic{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *TopicRepository) Take(db *gorm.DB, where ...interface{}) *model.Topic {
	ret := &model.Topic{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *TopicRepository) QueryCnd(db *gorm.DB, cnd *simple.QueryCnd) (list []model.Topic, err error) {
	err = cnd.DoQuery(db).Find(&list).Error
	return
}

func (this *TopicRepository) Query(db *gorm.DB, queries *simple.ParamQueries) (list []model.Topic, paging *simple.Paging) {
	queries.StartQuery(db).Find(&list)
    queries.StartCount(db).Model(&model.Topic{}).Count(&queries.Paging.Total)
	paging = queries.Paging
	return
}

func (this *TopicRepository) Create(db *gorm.DB, t *model.Topic) (err error) {
	err = db.Create(t).Error
	return
}

func (this *TopicRepository) Update(db *gorm.DB, t *model.Topic) (err error) {
	err = db.Save(t).Error
	return
}

func (this *TopicRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Topic{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *TopicRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Topic{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *TopicRepository) Delete(db *gorm.DB, id int64) {
	db.Model(&model.Topic{}).Delete("id", id)
}

