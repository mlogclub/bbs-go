package migrations

import (
	"bbs-go/internal/models"

	"github.com/mlogclub/simple/sqls"
)

func migrate_remove_comment_admin_permissions() error {
	commentPermissionCodes := []string{
		"dashboard.comment.view",
		"dashboard.comment.delete",
	}

	return sqls.WithTransaction(func(ctx *sqls.TxContext) error {
		var permissionIds []int64
		if err := ctx.Tx.Model(&models.Permission{}).
			Where("code in ?", commentPermissionCodes).
			Pluck("id", &permissionIds).Error; err != nil {
			return err
		}

		if len(permissionIds) > 0 {
			if err := ctx.Tx.
				Where("permission_id in ?", permissionIds).
				Delete(&models.RolePermission{}).Error; err != nil {
				return err
			}
		}

		return ctx.Tx.
			Where("code in ?", commentPermissionCodes).
			Delete(&models.Permission{}).Error
	})
}
