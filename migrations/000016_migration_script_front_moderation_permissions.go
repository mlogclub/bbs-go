package migrations

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/permissions"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
)

func migrate_front_moderation_permissions() error {
	seeds := []permissions.PermissionDefinition{
		permissions.PermissionTopicSticky,
		permissions.PermissionCommentDelete,
		permissions.PermissionUserForbiddenForever,
	}

	return sqls.WithTransaction(func(ctx *sqls.TxContext) error {
		now := dates.NowTimestamp()
		lang := config.Instance.Language
		if !lang.IsValid() {
			lang = config.DefaultLanguage
		}

		owner := repositories.RoleRepository.Take(ctx.Tx, "code = ?", constants.RoleOwner)
		for _, seed := range seeds {
			name := seed.NameZh
			if lang == config.LanguageEnUS {
				name = seed.NameEn
			}

			permission := repositories.PermissionRepository.Take(ctx.Tx, "code = ?", seed.Code)
			if permission == nil {
				permission = &models.Permission{
					Type:       seed.Type,
					Code:       seed.Code,
					CreateTime: now,
				}
			}
			permission.Name = name
			permission.GroupName = seed.GroupName
			permission.Description = seed.Description
			permission.SortNo = seed.SortNo
			permission.Status = constants.StatusOk
			permission.UpdateTime = now
			if permission.Id > 0 {
				if err := repositories.PermissionRepository.Update(ctx.Tx, permission); err != nil {
					return err
				}
			} else if err := repositories.PermissionRepository.Create(ctx.Tx, permission); err != nil {
				return err
			}

			if owner == nil {
				continue
			}
			existing := repositories.RolePermissionRepository.Take(ctx.Tx, "role_id = ? and permission_id = ?", owner.Id, permission.Id)
			if existing != nil {
				continue
			}
			if err := repositories.RolePermissionRepository.Create(ctx.Tx, &models.RolePermission{
				RoleId:       owner.Id,
				PermissionId: permission.Id,
				CreateTime:   now,
			}); err != nil {
				return err
			}
		}
		return nil
	})
}
