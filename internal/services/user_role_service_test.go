package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/repositories"
	"testing"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
)

func setupUserRoleServiceTestDB(t *testing.T) {
	t.Helper()
	db := setupTestDB(t)
	if err := db.AutoMigrate(&models.Role{}, &models.UserRole{}); err != nil {
		t.Fatalf("auto migrate user role models: %v", err)
	}
}

func TestUserRoleService_IsRoleInUse(t *testing.T) {
	setupUserRoleServiceTestDB(t)
	now := dates.NowTimestamp()
	user := mustCreateUser(t, now)
	role := mustCreateRole(t, "role-in-use", constants.StatusOk)
	mustAssignRole(t, user, role)

	if !UserRoleService.IsRoleInUse(role.Id) {
		t.Fatalf("expected assigned role to be in use")
	}

	unusedRole := mustCreateRole(t, "role-unused", constants.StatusOk)
	if UserRoleService.IsRoleInUse(unusedRole.Id) {
		t.Fatalf("expected unassigned role to be unused")
	}
}

func TestUserRoleService_IsRoleInUseIgnoresInvalidRoleId(t *testing.T) {
	setupUserRoleServiceTestDB(t)
	if err := repositories.UserRoleRepository.Create(sqls.DB(), &models.UserRole{
		UserId:     1,
		RoleId:     1,
		CreateTime: dates.NowTimestamp(),
	}); err != nil {
		t.Fatalf("create user role: %v", err)
	}

	if UserRoleService.IsRoleInUse(0) {
		t.Fatalf("expected invalid role id to be unused")
	}
}
