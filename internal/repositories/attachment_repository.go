package repositories

import (
	"gorm.io/gorm"

	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
)

var AttachmentRepository = newAttachmentRepository()

func newAttachmentRepository() *attachmentRepository {
	return &attachmentRepository{}
}

type attachmentRepository struct{}

func (r *attachmentRepository) Get(db *gorm.DB, id string) *models.Attachment {
	ret := &models.Attachment{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *attachmentRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.Attachment) {
	cnd.Find(db, &list)
	return
}

func (r *attachmentRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.Attachment {
	ret := &models.Attachment{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *attachmentRepository) FindPageByParams(db *gorm.DB, p *params.QueryParams) (list []models.Attachment, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &p.Cnd)
}

func (r *attachmentRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.Attachment, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.Attachment{})
	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *attachmentRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.Attachment{})
}

func (r *attachmentRepository) Create(db *gorm.DB, t *models.Attachment) error {
	return db.Create(t).Error
}

func (r *attachmentRepository) Update(db *gorm.DB, t *models.Attachment) error {
	return db.Save(t).Error
}

func (r *attachmentRepository) Updates(db *gorm.DB, id string, columns map[string]interface{}) error {
	return db.Model(&models.Attachment{}).Where("id = ?", id).Updates(columns).Error
}

func (r *attachmentRepository) UpdateColumn(db *gorm.DB, id string, name string, value interface{}) error {
	return db.Model(&models.Attachment{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

func (r *attachmentRepository) UpdateColumns(db *gorm.DB, topicId int64, columns map[string]interface{}) error {
	return db.Model(&models.Attachment{}).Where("topic_id = ?", topicId).Updates(columns).Error
}

func (r *attachmentRepository) ListByTopicId(db *gorm.DB, topicId int64) []models.Attachment {
	return r.Find(db, sqls.NewCnd().Eq("topic_id", topicId).Eq("status", constants.StatusOk).Asc("create_time"))
}

func (r *attachmentRepository) IncrDownloadCount(db *gorm.DB, id string) error {
	return r.UpdateColumn(db, id, "download_count", gorm.Expr("download_count + 1"))
}
