package permissions

import "testing"

func TestPermissionDefinitionsAreValid(t *testing.T) {
	seenCodes := make(map[string]struct{}, len(Permissions))
	for _, permission := range Permissions {
		if !permission.IsValid() {
			t.Fatalf("invalid permission definition: %#v", permission)
		}
		if _, ok := seenCodes[permission.Code]; ok {
			t.Fatalf("duplicate permission code: %s", permission.Code)
		}
		seenCodes[permission.Code] = struct{}{}
	}
}

func TestTopicStickyPermissionIsRegistered(t *testing.T) {
	if _, ok := FindByCode("dashboard.topic.sticky"); !ok {
		t.Fatalf("expected dashboard.topic.sticky permission to be registered")
	}
}
