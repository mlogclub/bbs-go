package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/simple"
)

var TopicTagRepository = newTopicTagRepository()

func newTopicTagRepository() *topicTagRepository {
	return &topicTagRepository{}
}

type topicTagRepository struct {
}

func (this *topicTagRepository) Get(db *gorm.DB, id int64) *model.TopicTag {
	ret := &model.TopicTag{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *topicTagRepository) Take(db *gorm.DB, where ...interface{}) *model.TopicTag {
	ret := &model.TopicTag{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *topicTagRepository) QueryCnd(db *gorm.DB, cnd *simple.QueryCnd) (list []model.TopicTag, err error) {
	err = cnd.DoQuery(db).Find(&list).Error
	return
}

func (this *topicTagRepository) Query(db *gorm.DB, queries *simple.ParamQueries) (list []model.TopicTag, paging *simple.Paging) {
	queries.StartQuery(db).Find(&list)
	queries.StartCount(db).Model(&model.TopicTag{}).Count(&queries.Paging.Total)
	paging = queries.Paging
	return
}

func (this *topicTagRepository) Create(db *gorm.DB, t *model.TopicTag) (err error) {
	err = db.Create(t).Error
	return
}

func (this *topicTagRepository) Update(db *gorm.DB, t *model.TopicTag) (err error) {
	err = db.Save(t).Error
	return
}

func (this *topicTagRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.TopicTag{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *topicTagRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.TopicTag{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *topicTagRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.TopicTag{}, "id = ?", id)
}

func (this *topicTagRepository) AddTopicTags(db *gorm.DB, topicId int64, tagIds []int64) {
	if topicId <= 0 || len(tagIds) == 0 {
		return
	}

	for _, tagId := range tagIds {
		_ = this.Create(db, &model.TopicTag{
			TopicId:    topicId,
			TagId:      tagId,
			CreateTime: simple.NowTimestamp(),
		})
	}
}

func (this *topicTagRepository) RemoveTopicTags(db *gorm.DB, topicId int64) {
	if topicId <= 0 {
		return
	}

	db.Where("topic_id = ?", topicId).Delete(model.TopicTag{})
}
