package repositories

import (
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"

	"bbs-go/internal/models"
)

var TopicRepository = newTopicRepository()

func newTopicRepository() *topicRepository {
	return &topicRepository{}
}

type topicRepository struct {
}

func (r *topicRepository) Get(db *gorm.DB, id int64) *models.Topic {
	ret := &models.Topic{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *topicRepository) Take(db *gorm.DB, where ...interface{}) *models.Topic {
	ret := &models.Topic{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *topicRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.Topic) {
	cnd.Find(db, &list)
	return
}

func (r *topicRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.Topic {
	ret := &models.Topic{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *topicRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.Topic, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *topicRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.Topic, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.Topic{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *topicRepository) FindBySql(db *gorm.DB, sqlStr string, paramArr ...interface{}) (list []models.Topic) {
	db.Raw(sqlStr, paramArr...).Scan(&list)
	return
}

func (r *topicRepository) CountBySql(db *gorm.DB, sqlStr string, paramArr ...interface{}) (count int64) {
	db.Raw(sqlStr, paramArr...).Count(&count)
	return
}

func (r *topicRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.Topic{})
}

func (r *topicRepository) Create(db *gorm.DB, t *models.Topic) (err error) {
	err = db.Create(t).Error
	return
}

func (r *topicRepository) Update(db *gorm.DB, t *models.Topic) (err error) {
	err = db.Save(t).Error
	return
}

func (r *topicRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.Topic{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *topicRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.Topic{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *topicRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.Topic{}, "id = ?", id)
}
