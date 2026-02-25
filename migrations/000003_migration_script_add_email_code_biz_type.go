package migrations

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"

	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

func migrate_add_email_code_biz_type() error {
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		return tx.Model(&models.EmailCode{}).
			Where("biz_type = '' OR biz_type IS NULL").
			Update("biz_type", constants.EmailCodeBizTypeEmailVerify).
			Error
	})
}
