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

var rbacPermissionSeeds = []permissions.PermissionDefinition{
	permissions.PermissionDashboardView,
	permissions.PermissionTopicView,
	permissions.PermissionTopicRecommend,
	permissions.PermissionTopicSticky,
	permissions.PermissionTopicAudit,
	permissions.PermissionTopicDelete,
	permissions.PermissionTopicSolve,
	permissions.PermissionArticleView,
	permissions.PermissionArticleUpdate,
	permissions.PermissionArticleAudit,
	permissions.PermissionArticleDelete,
	permissions.PermissionArticleTags,
	permissions.PermissionCommentDelete,
	permissions.PermissionCategoryView,
	permissions.PermissionCategoryCreate,
	permissions.PermissionCategoryUpdate,
	permissions.PermissionCategoryDelete,
	permissions.PermissionCategorySort,
	permissions.PermissionLinkView,
	permissions.PermissionLinkCreate,
	permissions.PermissionLinkUpdate,
	permissions.PermissionForbiddenWordView,
	permissions.PermissionForbiddenWordCreate,
	permissions.PermissionForbiddenWordUpdate,
	permissions.PermissionForbiddenWordDelete,
	permissions.PermissionUserView,
	permissions.PermissionUserCreate,
	permissions.PermissionUserUpdate,
	permissions.PermissionUserForbidden,
	permissions.PermissionUserForbiddenForever,
	permissions.PermissionUserUpdatePassword,
	permissions.PermissionUserResetPassword,
	permissions.PermissionBadgeView,
	permissions.PermissionBadgeCreate,
	permissions.PermissionBadgeUpdate,
	permissions.PermissionBadgeDelete,
	permissions.PermissionLevelView,
	permissions.PermissionLevelUpdate,
	permissions.PermissionTaskView,
	permissions.PermissionTaskCreate,
	permissions.PermissionTaskUpdate,
	permissions.PermissionTaskDelete,
	permissions.PermissionRoleView,
	permissions.PermissionRoleCreate,
	permissions.PermissionRoleUpdate,
	permissions.PermissionRoleDelete,
	permissions.PermissionRoleSort,
	permissions.PermissionRolePermissionUpdate,
	permissions.PermissionSettingView,
	permissions.PermissionSettingUpdate,
	permissions.PermissionSearchReindex,
	permissions.PermissionEmailLogView,
	permissions.PermissionUserTaskLogView,
	permissions.PermissionUserExpLogView,
	permissions.PermissionUserBadgeView,
	permissions.PermissionUserReportView,
	permissions.PermissionOperateLogView,
}

func migrate_init_rbac_permissions() error {
	return sqls.WithTransaction(func(ctx *sqls.TxContext) error {
		now := dates.NowTimestamp()
		lang := config.Instance.Language
		if !lang.IsValid() {
			lang = config.DefaultLanguage
		}

		var permissions []models.Permission
		for _, seed := range rbacPermissionSeeds {
			name := seed.NameZh
			if lang == config.LanguageEnUS {
				name = seed.NameEn
			}
			existing := repositories.PermissionRepository.Take(ctx.Tx, "code = ?", seed.Code)
			if existing != nil {
				existing.Type = seed.Type
				existing.Name = name
				existing.GroupName = seed.GroupName
				existing.Description = seed.Description
				existing.SortNo = seed.SortNo
				existing.Status = constants.StatusOk
				existing.UpdateTime = now
				if err := repositories.PermissionRepository.Update(ctx.Tx, existing); err != nil {
					return err
				}
				permissions = append(permissions, *existing)
				continue
			}

			permission := &models.Permission{
				Type:        seed.Type,
				Code:        seed.Code,
				Name:        name,
				GroupName:   seed.GroupName,
				Description: seed.Description,
				SortNo:      seed.SortNo,
				Status:      constants.StatusOk,
				CreateTime:  now,
				UpdateTime:  now,
			}
			if err := repositories.PermissionRepository.Create(ctx.Tx, permission); err != nil {
				return err
			}
			permissions = append(permissions, *permission)
		}

		owner := repositories.RoleRepository.Take(ctx.Tx, "code = ?", constants.RoleOwner)
		if owner == nil {
			return nil
		}
		for _, permission := range permissions {
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
