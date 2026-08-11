package authorization

type Action string

const (
	ActionRead     Action = "read"
	ActionDownload Action = "download"
	ActionEdit     Action = "edit"
	ActionDelete   Action = "delete"
	ActionShare    Action = "share"
)

type ResourceType string

const (
	ResourceFile   ResourceType = "file"
	ResourceFolder ResourceType = "folder"
)

type Permission string

const (
	PermissionNone   Permission = ""
	PermissionViewer Permission = "viewer"
	PermissionEditor Permission = "editor"
	PermissionOwner  Permission = "owner"
)

type Access struct {
	Owner         bool
	Permission    Permission
	AllowDownload bool
}

func (a Access) CanRead() bool {
	return a.Owner ||
		a.Permission == PermissionEditor ||
		a.Permission == PermissionViewer
}

func (a Access) CanEdit() bool {
	return a.Owner || a.Permission == PermissionEditor
}

func (a Access) CanDownload() bool {
	return a.Owner || (a.CanRead() && a.AllowDownload)
}

func (a Access) CanDelete() bool {
	return a.Owner
}

func (a Access) CanShare() bool {
	return a.Owner
}

func (a Access) Allows(action Action) bool {
	switch action {
	case ActionRead:
		return a.CanRead()
	case ActionDownload:
		return a.CanDownload()
	case ActionEdit:
		return a.CanEdit()
	case ActionDelete:
		return a.CanDelete()
	case ActionShare:
		return a.CanShare()
	default:
		return false
	}
}
