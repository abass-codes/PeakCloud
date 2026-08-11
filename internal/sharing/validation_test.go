package sharing

import "testing"

func TestValidPermission(t *testing.T) {
	if !ValidPermission(PermissionViewer) {
		t.Fatal("viewer should be valid")
	}

	if !ValidPermission(PermissionEditor) {
		t.Fatal("editor should be valid")
	}

	if ValidPermission(Permission("owner")) {
		t.Fatal("owner should not be assignable")
	}
}

func TestValidResourceType(t *testing.T) {
	if !ValidResourceType(ResourceFile) {
		t.Fatal("file should be valid")
	}

	if !ValidResourceType(ResourceFolder) {
		t.Fatal("folder should be valid")
	}

	if ValidResourceType(ResourceType("other")) {
		t.Fatal("unexpected valid resource type")
	}
}

func TestNormalizeEmail(t *testing.T) {
	got := NormalizeEmail(" USER@Example.COM ")

	if got != "user@example.com" {
		t.Fatalf("unexpected email: %q", got)
	}
}
