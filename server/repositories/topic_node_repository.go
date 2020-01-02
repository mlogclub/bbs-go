package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
)

var TopicNodeRepository = newTopicNodeRepository()

func newTopicNodeRepository() *topicNodeRepository {
	return &topicNodeRepository{}
}

type topicNodeRepository struct {
}

func (this *topicNodeRepository) Get(db *gorm.DB, id int64) *model.TopicNode {
	ret := &model.TopicNode{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *topicNodeRepository) Take(db *gorm.DB, where ...interface{}) *model.TopicNode {
	ret := &model.TopicNode{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *topicNodeRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.TopicNode) {
	cnd.Find(db, &list)
	return
}

func (this *topicNodeRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) *model.TopicNode {
	ret := &model.TopicNode{}
	cnd.FindOne(db, &ret)
	return ret
}

func (this *topicNodeRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.TopicNode, paging *simple.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *topicNodeRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.TopicNode, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.TopicNode{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *topicNodeRepository) Create(db *gorm.DB, t *model.TopicNode) (err error) {
	err = db.Create(t).Error
	return
}

func (this *topicNodeRepository) Update(db *gorm.DB, t *model.TopicNode) (err error) {
	err = db.Save(t).Error
	return
}

func (this *topicNodeRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.TopicNode{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *topicNodeRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.TopicNode{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *topicNodeRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.TopicNode{}, "id = ?", id)
}
