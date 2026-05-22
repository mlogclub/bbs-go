package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/permissions"
	"bbs-go/internal/repositories"
	"testing"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
)

func setupPermissionServiceTestDB(t *testing.T) {
	t.Helper()
	db := setupTestDB(t)
	if err := db.AutoMigrate(&models.Role{}, &models.UserRole{}, &models.Permission{}, &models.RolePermission{}); err != nil {
		t.Fatalf("auto migrate permission models: %v", err)
	}
	PermissionService.ClearCache()
}

func mustCreateRole(t *testing.T, code string, status int) *models.Role {
	t.Helper()
	role := &models.Role{
		Type:       constants.RoleTypeCustom,
		Name:       code,
		Code:       code,
		Status:     status,
		CreateTime: dates.NowTimestamp(),
		UpdateTime: dates.NowTimestamp(),
	}
	if err := repositories.RoleRepository.Create(sqls.DB(), role); err != nil {
		t.Fatalf("create role: %v", err)
	}
	return role
}

func mustCreatePermission(t *testing.T, code string, status int) *models.Permission {
	t.Helper()
	permission := &models.Permission{
		Type:       "dashboard",
		Code:       code,
		Name:       code,
		GroupName:  "test",
		Status:     status,
		CreateTime: dates.NowTimestamp(),
		UpdateTime: dates.NowTimestamp(),
	}
	if err := repositories.PermissionRepository.Create(sqls.DB(), permission); err != nil {
		t.Fatalf("create permission: %v", err)
	}
	return permission
}

func mustGrantPermission(t *testing.T, role *models.Role, permission *models.Permission) {
	t.Helper()
	if err := repositories.RolePermissionRepository.Create(sqls.DB(), &models.RolePermission{
		RoleId:       role.Id,
		PermissionId: permission.Id,
		CreateTime:   dates.NowTimestamp(),
	}); err != nil {
		t.Fatalf("create role permission: %v", err)
	}
}

func mustAssignRole(t *testing.T, user *models.User, role *models.Role) {
	t.Helper()
	if err := repositories.UserRoleRepository.Create(sqls.DB(), &models.UserRole{
		UserId:     user.Id,
		RoleId:     role.Id,
		CreateTime: dates.NowTimestamp(),
	}); err != nil {
		t.Fatalf("create user role: %v", err)
	}
}

func TestPermissionService_UserPermissionCodesAggregatesEnabledRoles(t *testing.T) {
	setupPermissionServiceTestDB(t)
	now := dates.NowTimestamp()
	user := mustCreateUser(t, now)
	role := mustCreateRole(t, "moderator", constants.StatusOk)
	permission := mustCreatePermission(t, "dashboard.topic.audit", constants.StatusOk)
	mustAssignRole(t, user, role)
	mustGrantPermission(t, role, permission)

	codes := PermissionService.GetUserPermissionCodes(user)

	if len(codes) != 1 || codes[0] != "dashboard.topic.audit" {
		t.Fatalf("expected dashboard.topic.audit, got %#v", codes)
	}
	if !PermissionService.HasPermission(user, "dashboard.topic.audit") {
		t.Fatalf("expected user to have dashboard.topic.audit")
	}
}

func TestPermissionService_UserPermissionCodesExcludeDisabledPermission(t *testing.T) {
	setupPermissionServiceTestDB(t)
	now := dates.NowTimestamp()
	user := mustCreateUser(t, now)
	role := mustCreateRole(t, "viewer", constants.StatusOk)
	permission := mustCreatePermission(t, "dashboard.setting.update", constants.StatusDeleted)
	mustAssignRole(t, user, role)
	mustGrantPermission(t, role, permission)

	if PermissionService.HasPermission(user, "dashboard.setting.update") {
		t.Fatalf("expected disabled permission to be excluded")
	}
}

func TestPermissionService_HasPermissionAllowsOwner(t *testing.T) {
	setupPermissionServiceTestDB(t)
	user := &models.User{Roles: constants.RoleOwner}

	if !PermissionService.HasPermission(user, "dashboard.setting.update") {
		t.Fatalf("expected owner to have every permission")
	}
}

func TestPermissionService_CanManageOwnResourceWithoutPermission(t *testing.T) {
	setupPermissionServiceTestDB(t)
	user := &models.User{Model: models.Model{Id: 42}}

	if !PermissionService.CanManageOwnedResource(user, 42, "dashboard.topic.delete") {
		t.Fatalf("expected user to manage own resource")
	}
}

func TestPermissionService_CanManageOtherResourceWithPermission(t *testing.T) {
	setupPermissionServiceTestDB(t)
	now := dates.NowTimestamp()
	user := mustCreateUser(t, now)
	role := mustCreateRole(t, "topic-moderator", constants.StatusOk)
	permission := mustCreatePermission(t, "dashboard.topic.delete", constants.StatusOk)
	mustAssignRole(t, user, role)
	mustGrantPermission(t, role, permission)

	if !PermissionService.CanManageOwnedResource(user, 99, "dashboard.topic.delete") {
		t.Fatalf("expected user with permission to manage other resource")
	}
}

func TestPermissionService_CannotManageOtherResourceWithoutPermission(t *testing.T) {
	setupPermissionServiceTestDB(t)
	user := &models.User{Model: models.Model{Id: 42}}

	if PermissionService.CanManageOwnedResource(user, 99, "dashboard.topic.delete") {
		t.Fatalf("expected user without permission to be rejected")
	}
}

func TestPermissionService_CanForbiddenUserUsesSeparatePermanentPermission(t *testing.T) {
	setupPermissionServiceTestDB(t)
	now := dates.NowTimestamp()
	user := mustCreateUser(t, now)
	role := mustCreateRole(t, "user-moderator", constants.StatusOk)
	forbiddenPermission := mustCreatePermission(t, permissions.PermissionUserForbidden.Code, constants.StatusOk)
	mustAssignRole(t, user, role)
	mustGrantPermission(t, role, forbiddenPermission)

	if !PermissionService.CanForbiddenUser(user, 7) {
		t.Fatalf("expected ordinary forbidden permission to allow temporary ban")
	}
	if PermissionService.CanForbiddenUser(user, -1) {
		t.Fatalf("expected ordinary forbidden permission to reject permanent ban")
	}
}

func TestPermissionService_CanForbiddenUserAllowsPermanentPermission(t *testing.T) {
	setupPermissionServiceTestDB(t)
	now := dates.NowTimestamp()
	user := mustCreateUser(t, now)
	role := mustCreateRole(t, "permanent-user-moderator", constants.StatusOk)
	forbiddenForeverPermission := mustCreatePermission(t, permissions.PermissionUserForbiddenForever.Code, constants.StatusOk)
	mustAssignRole(t, user, role)
	mustGrantPermission(t, role, forbiddenForeverPermission)

	if !PermissionService.CanForbiddenUser(user, -1) {
		t.Fatalf("expected permanent forbidden permission to allow permanent ban")
	}
	if PermissionService.CanForbiddenUser(user, 7) {
		t.Fatalf("expected permanent forbidden permission alone to reject temporary ban")
	}
}
