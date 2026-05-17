package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/permissions"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
)

var RolePermissionService = newRolePermissionService()

func newRolePermissionService() *rolePermissionService {
	return &rolePermissionService{}
}

type rolePermissionService struct {
}

func (s *rolePermissionService) Find(cnd *sqls.Cnd) []models.RolePermission {
	return repositories.RolePermissionRepository.Find(sqls.DB(), cnd)
}

func (s *rolePermissionService) GetRolePermissionIds(roleId int64) (permissionIds []int64) {
	list := s.Find(sqls.NewCnd().Eq("role_id", roleId))
	for _, item := range list {
		permissionIds = append(permissionIds, item.PermissionId)
	}
	return permissionIds
}

func (s *rolePermissionService) GetRolePermissionCodes(roleId int64) (codes []string) {
	if roleId <= 0 {
		return nil
	}
	err := sqls.DB().
		Model(&models.Permission{}).
		Select("t_permission.code").
		Joins("join t_role_permission on t_role_permission.permission_id = t_permission.id").
		Where("t_role_permission.role_id = ?", roleId).
		Order("t_permission.sort_no asc, t_permission.id desc").
		Pluck("t_permission.code", &codes).Error
	if err != nil {
		return nil
	}
	return codes
}

func (s *rolePermissionService) UpdateRolePermissions(roleId int64, permissionIds []int64) error {
	permissionIds = s.normalizePermissionIds(permissionIds)
	err := sqls.WithTransaction(func(ctx *sqls.TxContext) error {
		if err := ctx.Tx.Delete(&models.RolePermission{}, "role_id = ?", roleId).Error; err != nil {
			return err
		}
		for _, permissionId := range permissionIds {
			if permissionId <= 0 {
				continue
			}
			if err := repositories.RolePermissionRepository.Create(ctx.Tx, &models.RolePermission{
				RoleId:       roleId,
				PermissionId: permissionId,
				CreateTime:   dates.NowTimestamp(),
			}); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	PermissionService.ClearCache()
	return nil
}

func (s *rolePermissionService) normalizePermissionIds(permissionIds []int64) []int64 {
	if len(permissionIds) == 0 {
		return nil
	}

	ids := make([]int64, 0, len(permissionIds)+1)
	seen := make(map[int64]struct{}, len(permissionIds)+1)
	for _, permissionId := range permissionIds {
		if permissionId <= 0 {
			continue
		}
		if _, ok := seen[permissionId]; ok {
			continue
		}
		seen[permissionId] = struct{}{}
		ids = append(ids, permissionId)
	}
	if len(ids) == 0 {
		return nil
	}

	// dashboard.view is the entry permission for /dashboard and /api/admin/common/**.
	// Any role with dashboard permissions must include it, otherwise users can own
	// page/action permissions but still be blocked before entering the dashboard.
	if dashboardView := PermissionService.GetByCode(permissions.PermissionDashboardView.Code); dashboardView != nil {
		if _, ok := seen[dashboardView.Id]; !ok {
			ids = append([]int64{dashboardView.Id}, ids...)
		}
	}
	return ids
}
