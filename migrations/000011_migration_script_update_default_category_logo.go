package migrations

import (
	"github.com/mlogclub/simple/sqls"

	"bbs-go/internal/models"
)

// migrate_update_default_category_logo 重置默认分类图标，默认图标不设置，渲染的时候动态处理
func migrate_update_default_category_logo() error {
	return sqls.WithTransaction(func(ctx *sqls.TxContext) error {
		return ctx.Tx.Model(&models.Category{}).
			Where("logo = ?", "/res/images/category_default.png").
			Updates(map[string]any{"logo": ""}).Error
	})
}
