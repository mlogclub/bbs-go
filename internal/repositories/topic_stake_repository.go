package repositories

import (
	"bbs-go/internal/models"
	"bbs-go/internal/pkg/params"

	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

var TopicStakeRepository = newTopicStakeRepository()

func newTopicStakeRepository() *topicStakeRepository {
	return &topicStakeRepository{}
}

type topicStakeRepository struct {
}

func (r *topicStakeRepository) Get(db *gorm.DB, id int64) *models.TopicStake {
	ret := &models.TopicStake{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *topicStakeRepository) Take(db *gorm.DB, where ...interface{}) *models.TopicStake {
	ret := &models.TopicStake{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *topicStakeRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.TopicStake) {
	cnd.Find(db, &list)
	return
}

func (r *topicStakeRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.TopicStake {
	ret := &models.TopicStake{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *topicStakeRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.TopicStake, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *topicStakeRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.TopicStake, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.TopicStake{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *topicStakeRepository) Create(db *gorm.DB, t *models.TopicStake) (err error) {
	err = db.Create(t).Error
	return
}

func (r *topicStakeRepository) Update(db *gorm.DB, t *models.TopicStake) (err error) {
	err = db.Save(t).Error
	return
}

func (r *topicStakeRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.TopicStake{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *topicStakeRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.TopicStake{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *topicStakeRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.TopicStake{}, "id = ?", id)
}
