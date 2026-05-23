package services

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/permissions"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/repositories"
	"fmt"
	"slices"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

var PermissionService = newPermissionService()

func newPermissionService() *permissionService {
	return &permissionService{}
}

type permissionService struct {
}

func (s *permissionService) Get(id int64) *models.Permission {
	return repositories.PermissionRepository.Get(sqls.DB(), id)
}

func (s *permissionService) GetByCode(code string) *models.Permission {
	return repositories.PermissionRepository.FindOne(sqls.DB(), sqls.NewCnd().Eq("code", code))
}

func (s *permissionService) Find(cnd *sqls.Cnd) []models.Permission {
	return repositories.PermissionRepository.Find(sqls.DB(), cnd)
}

func (s *permissionService) SyncDefinitions() error {
	return sqls.WithTransaction(func(ctx *sqls.TxContext) error {
		return s.syncDefinitions(ctx.Tx)
	})
}

func (s *permissionService) syncDefinitions(db *gorm.DB) error {
	now := dates.NowTimestamp()
	lang := config.DefaultLanguage
	if config.Instance != nil {
		lang = config.Instance.Language
	}
	if !lang.IsValid() {
		lang = config.DefaultLanguage
	}

	seen := make(map[string]struct{}, len(permissions.Permissions))
	synced := make([]models.Permission, 0, len(permissions.Permissions))
	for _, seed := range permissions.Permissions {
		if !seed.IsValid() {
			return fmt.Errorf("invalid permission definition: %+v", seed)
		}
		if _, ok := seen[seed.Code]; ok {
			return fmt.Errorf("duplicate permission definition: %s", seed.Code)
		}
		seen[seed.Code] = struct{}{}

		name := seed.NameZh
		if lang == config.LanguageEnUS {
			name = seed.NameEn
		}

		permission := repositories.PermissionRepository.Take(db, "code = ?", seed.Code)
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
			if err := repositories.PermissionRepository.Update(db, permission); err != nil {
				return err
			}
		} else if err := repositories.PermissionRepository.Create(db, permission); err != nil {
			return err
		}
		synced = append(synced, *permission)
	}

	owner := repositories.RoleRepository.Take(db, "code = ?", constants.RoleOwner)
	if owner != nil {
		for _, permission := range synced {
			existing := repositories.RolePermissionRepository.Take(db, "role_id = ? and permission_id = ?", owner.Id, permission.Id)
			if existing != nil {
				continue
			}
			if err := repositories.RolePermissionRepository.Create(db, &models.RolePermission{
				RoleId:       owner.Id,
				PermissionId: permission.Id,
				CreateTime:   now,
			}); err != nil {
				return err
			}
		}
	}

	s.ClearCache()
	return nil
}

func (s *permissionService) GetUserPermissionCodes(user *models.User) []string {
	if user == nil || user.Id <= 0 {
		return nil
	}
	if user.IsOwner() {
		return []string{"*"}
	}
	if codes, ok := cache.PermissionCache.Get(user.Id); ok {
		return codes
	}

	var codes []string
	err := sqls.DB().
		Model(&models.Permission{}).
		Distinct("t_permission.code").
		Joins("join t_role_permission on t_role_permission.permission_id = t_permission.id").
		Joins("join t_user_role on t_user_role.role_id = t_role_permission.role_id").
		Joins("join t_role on t_role.id = t_user_role.role_id").
		Where("t_user_role.user_id = ?", user.Id).
		Where("t_role.status = ?", constants.StatusOk).
		Where("t_permission.status = ?", constants.StatusOk).
		Order("t_permission.sort_no asc, t_permission.id desc").
		Pluck("t_permission.code", &codes).Error
	if err != nil {
		return nil
	}
	cache.PermissionCache.Put(user.Id, codes)
	return codes
}

func (s *permissionService) HasPermission(user *models.User, code string) bool {
	if user == nil || code == "" {
		return false
	}
	if user.IsOwner() {
		return true
	}
	return slices.Contains(s.GetUserPermissionCodes(user), code)
}

func (s *permissionService) HasAnyPermission(user *models.User, codes ...string) bool {
	if user == nil || len(codes) == 0 {
		return false
	}
	if user.IsOwner() {
		return true
	}
	userCodes := s.GetUserPermissionCodes(user)
	for _, code := range codes {
		if slices.Contains(userCodes, code) {
			return true
		}
	}
	return false
}

func (s *permissionService) CanManageOwnedResource(user *models.User, ownerId int64, permissionCode string) bool {
	if user == nil {
		return false
	}
	if user.Id == ownerId {
		return true
	}
	return s.HasPermission(user, permissionCode)
}

func (s *permissionService) CanForbiddenUser(user *models.User, days int) bool {
	if days == -1 {
		return s.HasPermission(user, permissions.PermissionUserForbiddenForever.Code)
	}
	return s.HasPermission(user, permissions.PermissionUserForbidden.Code)
}

func (s *permissionService) InvalidateUser(userId int64) {
	cache.PermissionCache.Invalidate(userId)
}

func (s *permissionService) ClearCache() {
	cache.PermissionCache.Clear()
}
