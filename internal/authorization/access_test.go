package authorization

import "testing"

func TestStrongerAccessEditorBeatsViewer(t *testing.T) {
	viewer := Access{
		Permission:    PermissionViewer,
		AllowDownload: true,
	}

	editor := Access{
		Permission:    PermissionEditor,
		AllowDownload: false,
	}

	got := strongerAccess(viewer, editor)

	if got.Permission != PermissionEditor {
		t.Fatalf(
			"expected editor, got %q",
			got.Permission,
		)
	}
}

func TestStrongerAccessOwnerBeatsEditor(t *testing.T) {
	editor := Access{
		Permission: PermissionEditor,
	}

	owner := ownerAccess()

	got := strongerAccess(editor, owner)

	if !got.Owner {
		t.Fatal("expected owner access")
	}

	if got.Permission != PermissionOwner {
		t.Fatalf(
			"expected owner permission, got %q",
			got.Permission,
		)
	}
}

func TestEqualPermissionsMergeDownload(t *testing.T) {
	first := Access{
		Permission:    PermissionViewer,
		AllowDownload: false,
	}

	second := Access{
		Permission:    PermissionViewer,
		AllowDownload: true,
	}

	got := strongerAccess(first, second)

	if got.Permission != PermissionViewer {
		t.Fatalf(
			"expected viewer, got %q",
			got.Permission,
		)
	}

	if !got.AllowDownload {
		t.Fatal("expected download permission to be preserved")
	}
}

func TestNoAccessDoesNotOverrideViewer(t *testing.T) {
	viewer := Access{
		Permission: PermissionViewer,
	}

	got := strongerAccess(viewer, Access{})

	if got.Permission != PermissionViewer {
		t.Fatalf(
			"expected viewer, got %q",
			got.Permission,
		)
	}
}

func TestPermissionRank(t *testing.T) {
	tests := []struct {
		name   string
		access Access
		want   int
	}{
		{
			name:   "none",
			access: Access{},
			want:   0,
		},
		{
			name: "viewer",
			access: Access{
				Permission: PermissionViewer,
			},
			want: 1,
		},
		{
			name: "editor",
			access: Access{
				Permission: PermissionEditor,
			},
			want: 2,
		},
		{
			name:   "owner",
			access: ownerAccess(),
			want:   3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := permissionRank(test.access)

			if got != test.want {
				t.Fatalf(
					"expected %d, got %d",
					test.want,
					got,
				)
			}
		})
	}
}
