package trash

import "errors"

var (
	ErrNotFound       = errors.New("trash item not found")
	ErrForbidden      = errors.New("trash operation forbidden")
	ErrInvalidType    = errors.New("invalid trash resource type")
	ErrAlreadyDeleted = errors.New("resource already deleted")
)
