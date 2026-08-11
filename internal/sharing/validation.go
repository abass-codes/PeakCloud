package sharing

import (
	"errors"
	"strings"
)

var (
	ErrNotFound          = errors.New("share not found")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidResource   = errors.New("invalid resource")
	ErrInvalidPermission = errors.New("invalid permission")
	ErrRecipientNotFound = errors.New("recipient not found")
	ErrCannotShareSelf   = errors.New("cannot share with yourself")
	ErrAlreadyShared     = errors.New("resource already shared with recipient")
	ErrExpired           = errors.New("share link expired")
	ErrRevoked           = errors.New("share link revoked")
	ErrPasswordRequired  = errors.New("share password required")
	ErrInvalidPassword   = errors.New("invalid share password")
	ErrDownloadDenied    = errors.New("download not permitted")
)

func ValidResourceType(value ResourceType) bool {
	return value == ResourceFile || value == ResourceFolder
}

func ValidPermission(value Permission) bool {
	return value == PermissionViewer || value == PermissionEditor
}

func NormalizeEmail(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
