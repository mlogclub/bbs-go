package repositories

import (
	"bbs-go/internal/pkg/params"

	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"

	"bbs-go/internal/models"
)

var MentionRepository = newMentionRepository()

func newMentionRepository() *mentionRepository {
	return &mentionRepository{}
}

type mentionRepository struct {
}

func (r *mentionRepository) Get(db *gorm.DB, id int64) *models.Mention {
	ret := &models.Mention{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *mentionRepository) Take(db *gorm.DB, where ...interface{}) *models.Mention {
	ret := &models.Mention{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *mentionRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.Mention) {
	cnd.Find(db, &list)
	return
}

func (r *mentionRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.Mention {
	ret := &models.Mention{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *mentionRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.Mention, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *mentionRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.Mention, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.Mention{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *mentionRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.Mention{})
}

func (r *mentionRepository) Create(db *gorm.DB, t *models.Mention) (err error) {
	err = db.Create(t).Error
	return
}

func (r *mentionRepository) Update(db *gorm.DB, t *models.Mention) (err error) {
	err = db.Save(t).Error
	return
}

func (r *mentionRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.Mention{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *mentionRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.Mention{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *mentionRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.Mention{}, "id = ?", id)
}