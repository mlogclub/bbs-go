package services

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/permissions"
	"testing"
)

func TestRolePermissionService_UpdateRolePermissionsIncludesDashboardAccess(t *testing.T) {
	setupPermissionServiceTestDB(t)
	role := mustCreateRole(t, "admin", constants.StatusOk)
	dashboardView := mustCreatePermission(t, permissions.PermissionDashboardView.Code, constants.StatusOk)
	userView := mustCreatePermission(t, permissions.PermissionUserView.Code, constants.StatusOk)

	if err := RolePermissionService.UpdateRolePermissions(role.Id, []int64{userView.Id}); err != nil {
		t.Fatalf("update role permissions: %v", err)
	}

	codes := RolePermissionService.GetRolePermissionCodes(role.Id)
	if !containsString(codes, permissions.PermissionDashboardView.Code) {
		t.Fatalf("expected %s to be added automatically, got %#v", permissions.PermissionDashboardView.Code, codes)
	}
	if !containsString(codes, permissions.PermissionUserView.Code) {
		t.Fatalf("expected selected permission to remain, got %#v", codes)
	}
	if len(codes) != 2 {
		t.Fatalf("expected exactly selected permission plus dashboard access, got %#v", codes)
	}
	_ = dashboardView
}

func containsString(items []string, value string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}
	return false
}
