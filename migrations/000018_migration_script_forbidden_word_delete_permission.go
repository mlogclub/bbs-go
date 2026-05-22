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

func migrate_forbidden_word_delete_permission() error {
	seed := permissions.PermissionForbiddenWordDelete

	return sqls.WithTransaction(func(ctx *sqls.TxContext) error {
		now := dates.NowTimestamp()
		lang := config.Instance.Language
		if !lang.IsValid() {
			lang = config.DefaultLanguage
		}

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

		owner := repositories.RoleRepository.Take(ctx.Tx, "code = ?", constants.RoleOwner)
		if owner == nil {
			return nil
		}
		existing := repositories.RolePermissionRepository.Take(ctx.Tx, "role_id = ? and permission_id = ?", owner.Id, permission.Id)
		if existing != nil {
			return nil
		}
		return repositories.RolePermissionRepository.Create(ctx.Tx, &models.RolePermission{
			RoleId:       owner.Id,
			PermissionId: permission.Id,
			CreateTime:   now,
		})
	})
}
