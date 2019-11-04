package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
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

func (this *topicRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Topic) {
	cnd.Find(db, &list)
	return
}

func (this *topicRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.Topic) {
	cnd.FindOne(db, &ret)
	return
}

func (this *topicRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.Topic, paging *simple.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *topicRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Topic, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Topic{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
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
