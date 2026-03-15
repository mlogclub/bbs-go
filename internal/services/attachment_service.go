package services

import (
	"errors"
	"io"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"

	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/pkg/uploader"
	"bbs-go/internal/repositories"
)

var AttachmentService = new(attachmentService)

type attachmentService struct{}

func (s *attachmentService) extAllowed(ext string, allowedTypes []string) bool {
	if len(allowedTypes) == 0 {
		return false
	}
	ext = strings.ToLower(ext)
	for _, a := range allowedTypes {
		if strings.ToLower(strings.TrimSpace(a)) == ext {
			return true
		}
	}
	return false
}

// Upload 流式上传附件；content 为数据流，contentLength 为文件大小（用于存储 FileSize 与上传 ContentLength）。
func (s *attachmentService) Upload(userId int64, filename string, content io.Reader, contentLength int64, contentType string, downloadScore int) (*models.Attachment, error) {
	if downloadScore < 0 {
		downloadScore = 0
	}
	cfg := SysConfigService.GetAttachmentConfig()
	ext := strings.ToLower(filepath.Ext(filename))
	if !s.extAllowed(ext, cfg.AllowedTypes) {
		return nil, errors.New(locales.Get("attachment.ext_not_allowed"))
	}
	var (
		attId       = strs.UUID()
		key         = uploader.GenerateAttachmentKey(attId, ext)
		disposition = "attachment; filename=\"" + url.QueryEscape(filepath.Base(filename)) + "\""
	)
	fileUrl, err := UploadService.PutObject(key, content, &uploader.PutOptions{ContentType: contentType, ContentDisposition: disposition, ContentLength: contentLength})
	if err != nil {
		return nil, err
	}
	att := &models.Attachment{
		Id:            attId,
		TopicId:       0,
		UserId:        userId,
		FileName:      filename,
		FileUrl:       fileUrl,
		FileSize:      contentLength,
		FileType:      contentType,
		DownloadScore: downloadScore,
		Status:        constants.StatusOk,
		CreateTime:    dates.NowTimestamp(),
		UpdateTime:    dates.NowTimestamp(),
	}
	if err := repositories.AttachmentRepository.Create(sqls.DB(), att); err != nil {
		return nil, err
	}
	return att, nil
}

// UpdateDownloadScore 更新附件的下载积分（仅附件所属用户可更新）
func (s *attachmentService) UpdateDownloadScore(attachmentId string, userId int64, downloadScore int) (*models.Attachment, error) {
	if strs.IsBlank(attachmentId) {
		return nil, errors.New(locales.Get("attachment.not_found"))
	}
	att := repositories.AttachmentRepository.Get(sqls.DB(), attachmentId)
	if att == nil || att.Status != constants.StatusOk {
		return nil, errors.New(locales.Get("attachment.not_found"))
	}
	if att.UserId != userId {
		return nil, errors.New(locales.Get("attachment.no_permission"))
	}
	if downloadScore < 0 {
		downloadScore = 0
	}
	att.DownloadScore = downloadScore
	att.UpdateTime = dates.NowTimestamp()
	return att, repositories.AttachmentRepository.Updates(sqls.DB(), attachmentId, map[string]any{
		"download_score": downloadScore,
		"update_time":    dates.NowTimestamp(),
	})
}

// Get 根据 ID 获取附件（仅返回存在且正常的）
func (s *attachmentService) Get(id string) *models.Attachment {
	att := repositories.AttachmentRepository.Get(sqls.DB(), id)
	if att == nil || att.Status != constants.StatusOk {
		return nil
	}
	return att
}

// ListByTopicId 按帖子查询正常状态的附件
func (s *attachmentService) ListByTopicId(topicId int64) []models.Attachment {
	return repositories.AttachmentRepository.ListByTopicId(sqls.DB(), topicId)
}

// HasDownloaded 当前用户是否已购买该附件
func (s *attachmentService) HasDownloaded(userId int64, attachmentId string) bool {
	if userId <= 0 || strs.IsBlank(attachmentId) {
		return false
	}
	return repositories.AttachmentDownloadLogRepository.Exists(sqls.DB(), userId, attachmentId)
}

func (s *attachmentService) FindDownloadedAttachmentIds(userId int64, attachmentIds []string) []string {
	if userId <= 0 || len(attachmentIds) == 0 {
		return nil
	}

	filteredIds := make([]string, 0, len(attachmentIds))
	seen := make(map[string]bool, len(attachmentIds))
	for _, attachmentId := range attachmentIds {
		if strs.IsBlank(attachmentId) || seen[attachmentId] {
			continue
		}
		seen[attachmentId] = true
		filteredIds = append(filteredIds, attachmentId)
	}
	if len(filteredIds) == 0 {
		return nil
	}
	return repositories.AttachmentDownloadLogRepository.FindDownloadedAttachmentIds(sqls.DB(), userId, filteredIds)
}

// GetDownloadRedirectUrl 根据附件访问地址生成 302 目标 URL（Local 需拼 baseURL）
func (s *attachmentService) GetDownloadRedirectUrl(att *models.Attachment) string {
	return att.FileUrl
}

// Download 鉴权并返回下载重定向 URL；如需扣积分则在事务内扣费并写入 download_log
func (s *attachmentService) Download(attachmentId string, userId int64) (redirectURL string, err error) {
	if strs.IsBlank(attachmentId) {
		return "", errors.New(locales.Get("attachment.not_found"))
	}
	att := repositories.AttachmentRepository.Get(sqls.DB(), attachmentId)
	if att == nil || att.Status != constants.StatusOk {
		return "", errors.New(locales.Get("attachment.not_found"))
	}
	if att.TopicId <= 0 {
		return "", errors.New(locales.Get("attachment.not_found"))
	}

	topic := repositories.TopicRepository.Get(sqls.DB(), att.TopicId)
	if topic == nil || topic.Status == constants.StatusDeleted {
		return "", errors.New(locales.Get("attachment.not_found"))
	}

	// 已购买：直接放行
	if repositories.AttachmentDownloadLogRepository.Exists(sqls.DB(), userId, attachmentId) {
		redirectURL = s.GetDownloadRedirectUrl(att)
		if strs.IsNotBlank(redirectURL) {
			repositories.AttachmentRepository.IncrDownloadCount(sqls.DB(), att.Id)
		}
		return redirectURL, nil
	}

	// 帖主本人或 0 积分：免费，写入 download_log 便于统计
	if att.UserId == userId || att.DownloadScore <= 0 {
		_ = repositories.AttachmentDownloadLogRepository.Create(sqls.DB(), &models.AttachmentDownloadLog{
			UserId:       userId,
			AttachmentId: attachmentId,
			CreateTime:   dates.NowTimestamp(),
		})
		redirectURL = s.GetDownloadRedirectUrl(att)
		if strs.IsNotBlank(redirectURL) {
			repositories.AttachmentRepository.IncrDownloadCount(sqls.DB(), att.Id)
		}
		return redirectURL, nil
	}

	// 需扣积分：事务内扣费 + 写 UserScoreLog + 插入 download_log
	err = sqls.WithTransaction(func(ctx *sqls.TxContext) error {
		user := repositories.UserRepository.Get(ctx.Tx, userId)
		if user == nil {
			return errors.New(locales.Get("common.not_found"))
		}
		if user.Score < att.DownloadScore {
			return errors.New(locales.Get("attachment.insufficient_score"))
		}
		if err := UserService.DecrScoreTx(ctx, userId, att.DownloadScore, constants.SourceTypeAttachmentDownload, attachmentId, locales.Get("attachment.download_deduct")); err != nil {
			return err
		}
		if err := repositories.AttachmentDownloadLogRepository.Create(ctx.Tx, &models.AttachmentDownloadLog{
			UserId:       userId,
			AttachmentId: attachmentId,
			CreateTime:   dates.NowTimestamp(),
		}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	redirectURL = s.GetDownloadRedirectUrl(att)
	if strs.IsNotBlank(redirectURL) {
		repositories.AttachmentRepository.IncrDownloadCount(sqls.DB(), att.Id)
	}
	return redirectURL, nil
}

// SoftDeleteByTopicId 帖子删除时软删除其下所有附件
func (s *attachmentService) SoftDeleteByTopicId(ctx *sqls.TxContext, topicId int64) error {
	return repositories.AttachmentRepository.UpdateColumns(ctx.Tx, topicId, map[string]interface{}{
		"status":      constants.StatusDeleted,
		"update_time": dates.NowTimestamp(),
	})
}

// ReplaceTopicAttachments 编辑帖时全量替换附件
func (s *attachmentService) ReplaceTopicAttachments(ctx *sqls.TxContext, topicId, userId int64, attachmentIds []string) error {
	newSet := make(map[string]bool)
	for _, id := range attachmentIds {
		if strs.IsNotBlank(id) {
			newSet[id] = true
		}
	}

	// 从当前中移除的：解绑 + 软删除
	current := repositories.AttachmentRepository.ListByTopicId(ctx.Tx, topicId)
	for _, att := range current {
		if !newSet[att.Id] {
			if err := repositories.AttachmentRepository.Updates(ctx.Tx, att.Id, map[string]interface{}{
				"topic_id": 0, "status": constants.StatusDeleted, "update_time": dates.NowTimestamp(),
			}); err != nil {
				return err
			}
		}
	}

	// 新列表中的：校验归属且未绑其他帖，再绑定
	for _, aid := range attachmentIds {
		if strs.IsBlank(aid) {
			continue
		}
		att := repositories.AttachmentRepository.Get(ctx.Tx, aid)
		if att == nil || att.UserId != userId {
			return errors.New(locales.Get("attachment.no_permission"))
		}
		if att.TopicId != 0 && att.TopicId != topicId {
			return errors.New(locales.Get("attachment.already_bound"))
		}
		if err := repositories.AttachmentRepository.Updates(ctx.Tx, aid, map[string]interface{}{
			"topic_id": topicId, "status": constants.StatusOk, "update_time": dates.NowTimestamp(),
		}); err != nil {
			return err
		}
	}
	return nil
}

// CheckAttachmentsExistAndOwned 检查 attachmentIds 是否存在且均属于 userId，且未绑定其他帖子（或仅绑定 topicId）
func (s *attachmentService) CheckAttachmentsExistAndOwned(ctx *sqls.TxContext, userId int64, attachmentIds []string, topicId int64) error {
	for _, aid := range attachmentIds {
		if strs.IsBlank(aid) {
			continue
		}
		att := repositories.AttachmentRepository.Get(ctx.Tx, aid)
		if att == nil || att.UserId != userId {
			return errors.New(locales.Get("attachment.no_permission"))
		}
		if att.TopicId != 0 && att.TopicId != topicId {
			return errors.New(locales.Get("attachment.already_bound"))
		}
	}
	return nil
}
