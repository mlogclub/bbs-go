package migrations

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/dto"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/sqls"
)

func migrate_attachment_config() error {
	return sqls.WithTransaction(func(ctx *sqls.TxContext) error {
		tx := ctx.Tx
		existing := repositories.SysConfigRepository.GetByKey(tx, constants.SysConfigAttachmentConfig)
		if existing != nil {
			return nil
		}
		cfg := dto.AttachmentConfig{
			Enabled:      true,
			AllowedTypes: []string{".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt", ".md", ".csv", ".zip", ".rar", ".7z", ".tar", ".gz"},
			MaxSizeMB:    10,
			MaxCount:     5,
		}
		value, err := jsons.ToStr(cfg)
		if err != nil {
			return err
		}
		name, desc := attachmentConfigMetaByLanguage()
		return repositories.SysConfigRepository.Create(tx, &models.SysConfig{
			Key:         constants.SysConfigAttachmentConfig,
			Value:       value,
			Name:        name,
			Description: desc,
			CreateTime:  dates.NowTimestamp(),
			UpdateTime:  dates.NowTimestamp(),
		})
	})
}

func attachmentConfigMetaByLanguage() (name, description string) {
	if config.Instance.Language == config.LanguageEnUS {
		return "Attachment Config", "Topic attachment upload settings (enable, allowed types, max size, max count per topic)"
	}
	return "附件配置", "帖子附件上传设置（是否开启、允许类型、单文件大小、每帖最多数量）"
}
