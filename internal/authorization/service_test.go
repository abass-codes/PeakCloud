package authorization

import (
	"context"
	"testing"
)

type fakeRepository struct {
	fileAccess   Access
	fileErr      error
	folderAccess Access
	folderErr    error
}

func (r *fakeRepository) FileAccess(
	_ context.Context,
	_ string,
	_ string,
) (Access, error) {
	return r.fileAccess, r.fileErr
}

func (r *fakeRepository) FolderAccess(
	_ context.Context,
	_ string,
	_ string,
) (Access, error) {
	return r.folderAccess, r.folderErr
}

func TestOwnerPermissions(t *testing.T) {
	access := Access{
		Owner:         true,
		Permission:    PermissionOwner,
		AllowDownload: true,
	}

	actions := []Action{
		ActionRead,
		ActionDownload,
		ActionEdit,
		ActionDelete,
		ActionShare,
	}

	for _, action := range actions {
		if !access.Allows(action) {
			t.Fatalf("owner should be allowed to %s", action)
		}
	}
}

func TestViewerPermissions(t *testing.T) {
	access := Access{
		Permission:    PermissionViewer,
		AllowDownload: true,
	}

	if !access.Allows(ActionRead) {
		t.Fatal("viewer should be allowed to read")
	}

	if !access.Allows(ActionDownload) {
		t.Fatal("viewer with download permission should be allowed to download")
	}

	if access.Allows(ActionEdit) {
		t.Fatal("viewer should not be allowed to edit")
	}

	if access.Allows(ActionDelete) {
		t.Fatal("viewer should not be allowed to delete")
	}

	if access.Allows(ActionShare) {
		t.Fatal("viewer should not be allowed to share")
	}
}

func TestViewerDownloadDenied(t *testing.T) {
	access := Access{
		Permission:    PermissionViewer,
		AllowDownload: false,
	}

	if !access.Allows(ActionRead) {
		t.Fatal("viewer should still be allowed to read")
	}

	if access.Allows(ActionDownload) {
		t.Fatal("viewer should not download when downloads are disabled")
	}
}

func TestEditorPermissions(t *testing.T) {
	access := Access{
		Permission:    PermissionEditor,
		AllowDownload: true,
	}

	if !access.Allows(ActionRead) {
		t.Fatal("editor should be allowed to read")
	}

	if !access.Allows(ActionDownload) {
		t.Fatal("editor with download permission should be allowed to download")
	}

	if !access.Allows(ActionEdit) {
		t.Fatal("editor should be allowed to edit")
	}

	if access.Allows(ActionDelete) {
		t.Fatal("editor should not be allowed to delete")
	}

	if access.Allows(ActionShare) {
		t.Fatal("editor should not be allowed to share")
	}
}

func TestNoAccessPermissions(t *testing.T) {
	access := Access{}

	actions := []Action{
		ActionRead,
		ActionDownload,
		ActionEdit,
		ActionDelete,
		ActionShare,
	}

	for _, action := range actions {
		if access.Allows(action) {
			t.Fatalf("user without access should not be allowed to %s", action)
		}
	}
}

func TestAuthorizeAllowed(t *testing.T) {
	repository := &fakeRepository{
		fileAccess: Access{
			Permission:    PermissionViewer,
			AllowDownload: true,
		},
	}

	service := NewService(repository)

	access, err := service.Authorize(
		context.Background(),
		"user-1",
		ResourceFile,
		"file-1",
		ActionRead,
	)
	if err != nil {
		t.Fatalf("expected authorization to succeed: %v", err)
	}

	if access.Permission != PermissionViewer {
		t.Fatalf(
			"expected viewer permission, got %q",
			access.Permission,
		)
	}
}

func TestAuthorizeForbidden(t *testing.T) {
	repository := &fakeRepository{
		fileAccess: Access{
			Permission: PermissionViewer,
		},
	}

	service := NewService(repository)

	_, err := service.Authorize(
		context.Background(),
		"user-1",
		ResourceFile,
		"file-1",
		ActionEdit,
	)

	if err != ErrForbidden {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestAuthorizeInvalidResource(t *testing.T) {
	service := NewService(&fakeRepository{})

	_, err := service.Authorize(
		context.Background(),
		"user-1",
		ResourceType("invalid"),
		"resource-1",
		ActionRead,
	)

	if err != ErrInvalidResource {
		t.Fatalf("expected ErrInvalidResource, got %v", err)
	}
}

func TestCanFolder(t *testing.T) {
	repository := &fakeRepository{
		folderAccess: Access{
			Permission: PermissionEditor,
		},
	}

	service := NewService(repository)

	allowed, err := service.Can(
		context.Background(),
		"user-1",
		ResourceFolder,
		"folder-1",
		ActionEdit,
	)
	if err != nil {
		t.Fatalf("expected no error: %v", err)
	}

	if !allowed {
		t.Fatal("editor should be allowed to edit folder")
	}
}
