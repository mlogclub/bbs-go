package repositories

import (
	"github.com/mlogclub/simple"
	"github.com/mlogclub/simple/date"
	"gorm.io/gorm"

	"bbs-go/model"
)

var TopicTagRepository = newTopicTagRepository()

func newTopicTagRepository() *topicTagRepository {
	return &topicTagRepository{}
}

type topicTagRepository struct {
}

func (r *topicTagRepository) Get(db *gorm.DB, id int64) *model.TopicTag {
	ret := &model.TopicTag{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *topicTagRepository) Take(db *gorm.DB, where ...interface{}) *model.TopicTag {
	ret := &model.TopicTag{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *topicTagRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.TopicTag) {
	cnd.Find(db, &list)
	return
}

func (r *topicTagRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) *model.TopicTag {
	ret := &model.TopicTag{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *topicTagRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.TopicTag, paging *simple.Paging) {
	return r.FindPageByCnd(db, &params.SqlCnd)
}

func (r *topicTagRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.TopicTag, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.TopicTag{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *topicTagRepository) Create(db *gorm.DB, t *model.TopicTag) (err error) {
	err = db.Create(t).Error
	return
}

func (r *topicTagRepository) Update(db *gorm.DB, t *model.TopicTag) (err error) {
	err = db.Save(t).Error
	return
}

func (r *topicTagRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.TopicTag{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *topicTagRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.TopicTag{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *topicTagRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.TopicTag{}, "id = ?", id)
}

func (r *topicTagRepository) AddTopicTags(db *gorm.DB, topicId int64, tagIds []int64) {
	if topicId <= 0 || len(tagIds) == 0 {
		return
	}
	for _, tagId := range tagIds {
		_ = r.Create(db, &model.TopicTag{
			TopicId:    topicId,
			TagId:      tagId,
			CreateTime: date.NowTimestamp(),
		})
	}
}

func (r *topicTagRepository) DeleteTopicTags(db *gorm.DB, topicId int64) {
	if topicId <= 0 {
		return
	}
	db.Where("topic_id = ?", topicId).Delete(model.TopicTag{})
}
