package repositories

import (
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"

	"bbs-go/internal/models"
)

var TopicTagRepository = newTopicTagRepository()

func newTopicTagRepository() *topicTagRepository {
	return &topicTagRepository{}
}

type topicTagRepository struct {
}

func (r *topicTagRepository) Get(db *gorm.DB, id int64) *models.TopicTag {
	ret := &models.TopicTag{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *topicTagRepository) Take(db *gorm.DB, where ...interface{}) *models.TopicTag {
	ret := &models.TopicTag{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *topicTagRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.TopicTag) {
	cnd.Find(db, &list)
	return
}

func (r *topicTagRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.TopicTag {
	ret := &models.TopicTag{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *topicTagRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.TopicTag, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *topicTagRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.TopicTag, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.TopicTag{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *topicTagRepository) Create(db *gorm.DB, t *models.TopicTag) (err error) {
	err = db.Create(t).Error
	return
}

func (r *topicTagRepository) Update(db *gorm.DB, t *models.TopicTag) (err error) {
	err = db.Save(t).Error
	return
}

func (r *topicTagRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.TopicTag{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *topicTagRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.TopicTag{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *topicTagRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.TopicTag{}, "id = ?", id)
}

func (r *topicTagRepository) AddTopicTags(db *gorm.DB, topicId int64, tagIds []int64) error {
	if topicId <= 0 || len(tagIds) == 0 {
		return nil
	}
	for _, tagId := range tagIds {
		if err := r.Create(db, &models.TopicTag{
			TopicId:    topicId,
			TagId:      tagId,
			CreateTime: dates.NowTimestamp(),
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *topicTagRepository) DeleteTopicTags(db *gorm.DB, topicId int64) {
	if topicId <= 0 {
		return
	}
	db.Where("topic_id = ?", topicId).Delete(models.TopicTag{})
}
