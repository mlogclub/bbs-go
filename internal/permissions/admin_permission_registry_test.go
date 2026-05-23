package permissions

import "testing"

func TestAdminPermissionRegistryMatchesMethodAndPath(t *testing.T) {
	code, ok := GetAdminPermissionCode("POST", "/api/admin/topic/list")
	if !ok {
		t.Fatalf("expected topic list permission to be registered")
	}
	if code != PermissionTopicView.Code {
		t.Fatalf("expected %s, got %s", PermissionTopicView.Code, code)
	}
}

func TestAdminPermissionRegistryAllowsRoleOptionsFromUserUpdate(t *testing.T) {
	codes, ok := GetAdminPermissionCodes("GET", "/api/admin/role/roles")
	if !ok {
		t.Fatalf("expected role options permission to be registered")
	}
	expected := []string{PermissionRoleView.Code, PermissionUserUpdate.Code}
	if len(codes) != len(expected) {
		t.Fatalf("expected %#v, got %#v", expected, codes)
	}
	for i, expectedCode := range expected {
		if codes[i] != expectedCode {
			t.Fatalf("expected %#v, got %#v", expected, codes)
		}
	}
}

func TestAdminPermissionRegistryAllowsEitherUserForbiddenPermission(t *testing.T) {
	codes, ok := GetAdminPermissionCodes("POST", "/api/admin/user/forbidden")
	if !ok {
		t.Fatalf("expected user forbidden permission to be registered")
	}
	expected := []string{PermissionUserForbidden.Code, PermissionUserForbiddenForever.Code}
	if len(codes) != len(expected) {
		t.Fatalf("expected %#v, got %#v", expected, codes)
	}
	for i, expectedCode := range expected {
		if codes[i] != expectedCode {
			t.Fatalf("expected %#v, got %#v", expected, codes)
		}
	}
}

func TestAdminPermissionRegistryProtectsForbiddenWordDelete(t *testing.T) {
	code, ok := GetAdminPermissionCode("POST", "/api/admin/forbidden-word/delete")
	if !ok {
		t.Fatalf("expected forbidden word delete permission to be registered")
	}
	if code != PermissionForbiddenWordDelete.Code {
		t.Fatalf("expected %s, got %s", PermissionForbiddenWordDelete.Code, code)
	}
}

func TestAdminPermissionRegistryProtectsBadgeUpdateSort(t *testing.T) {
	code, ok := GetAdminPermissionCode("POST", "/api/admin/badge/update_sort")
	if !ok {
		t.Fatalf("expected badge sort permission to be registered")
	}
	if code != PermissionBadgeUpdate.Code {
		t.Fatalf("expected %s, got %s", PermissionBadgeUpdate.Code, code)
	}
}

func TestAdminPermissionRegistryProtectsLinkDelete(t *testing.T) {
	code, ok := GetAdminPermissionCode("POST", "/api/admin/link/delete")
	if !ok {
		t.Fatalf("expected link delete permission to be registered")
	}
	if code != PermissionLinkDelete.Code {
		t.Fatalf("expected %s, got %s", PermissionLinkDelete.Code, code)
	}
}

func TestAdminPermissionRegistryRejectsUnknownAdminPath(t *testing.T) {
	if code, ok := GetAdminPermissionCode("POST", "/api/admin/unknown/action"); ok {
		t.Fatalf("expected unknown admin path to be rejected, got %s", code)
	}
}

func TestAdminPermissionRegistryRejectsCommentManagementPaths(t *testing.T) {
	paths := []struct {
		method string
		path   string
	}{
		{method: "GET", path: "/api/admin/comment/1"},
		{method: "POST", path: "/api/admin/comment/list"},
		{method: "DELETE", path: "/api/admin/comment/1"},
	}

	for _, path := range paths {
		if code, ok := GetAdminPermissionCode(path.method, path.path); ok {
			t.Fatalf("expected %s %s to be rejected, got %s", path.method, path.path, code)
		}
	}
}
