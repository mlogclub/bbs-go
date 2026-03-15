package migrations

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/dto"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

func migrate_notification_types_defaults() error {
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		now := dates.NowTimestamp()
		existing := repositories.SysConfigRepository.GetByKey(tx, constants.SysConfigNotificationTypes)
		if existing != nil && existing.Value != "" {
			return nil
		}
		defaults := map[string]dto.NoticeTypeConfig{
			"topicComment":     {Site: true, Email: true},
			"commentReply":     {Site: true, Email: true},
			"topicLike":        {Site: true, Email: true},
			"topicFavorite":    {Site: true, Email: true},
			"topicRecommend":   {Site: true, Email: true},
			"topicDelete":      {Site: true, Email: false},
			"articleComment":   {Site: true, Email: true},
			"userLevelUp":      {Site: true, Email: true},
			"userBadgeGrant":   {Site: true, Email: true},
			"qaAnswerAccepted": {Site: true, Email: true},
		}
		value := jsons.ToJsonStr(defaults)
		if existing == nil {
			return repositories.SysConfigRepository.Create(tx, &models.SysConfig{
				Key:         constants.SysConfigNotificationTypes,
				Value:       value,
				Name:        "通知类型配置",
				Description: "各消息类型的站内信与邮件开关",
				CreateTime:  now,
				UpdateTime:  now,
			})
		}
		existing.Value = value
		existing.UpdateTime = now
		return repositories.SysConfigRepository.Update(tx, existing)
	})
}
