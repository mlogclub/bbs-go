package repositories

import (
	"gorm.io/gorm"

	"bbs-go/internal/models"

	"github.com/mlogclub/simple/sqls"
)

var AttachmentDownloadLogRepository = newAttachmentDownloadLogRepository()

func newAttachmentDownloadLogRepository() *attachmentDownloadLogRepository {
	return &attachmentDownloadLogRepository{}
}

type attachmentDownloadLogRepository struct{}

func (r *attachmentDownloadLogRepository) Get(db *gorm.DB, id int64) *models.AttachmentDownloadLog {
	ret := &models.AttachmentDownloadLog{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *attachmentDownloadLogRepository) GetByUserAndAttachment(db *gorm.DB, userId int64, attachmentId string) *models.AttachmentDownloadLog {
	ret := &models.AttachmentDownloadLog{}
	if err := db.Where("user_id = ? AND attachment_id = ?", userId, attachmentId).First(ret).Error; err != nil {
		return nil
	}
	return ret
}

func (r *attachmentDownloadLogRepository) Exists(db *gorm.DB, userId int64, attachmentId string) bool {
	var count int64
	db.Model(&models.AttachmentDownloadLog{}).Where("user_id = ? AND attachment_id = ?", userId, attachmentId).Count(&count)
	return count > 0
}

func (r *attachmentDownloadLogRepository) FindDownloadedAttachmentIds(db *gorm.DB, userId int64, attachmentIds []string) []string {
	if userId <= 0 || len(attachmentIds) == 0 {
		return nil
	}

	var ids []string
	db.Model(&models.AttachmentDownloadLog{}).
		Where("user_id = ? AND attachment_id IN ?", userId, attachmentIds).
		Pluck("attachment_id", &ids)
	return ids
}

func (r *attachmentDownloadLogRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.AttachmentDownloadLog) {
	cnd.Find(db, &list)
	return
}

func (r *attachmentDownloadLogRepository) Create(db *gorm.DB, t *models.AttachmentDownloadLog) error {
	return db.Create(t).Error
}

func (r *attachmentDownloadLogRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.AttachmentDownloadLog{})
}
