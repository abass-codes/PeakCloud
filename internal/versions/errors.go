package versions

import "errors"

var (
	ErrNotFound       = errors.New("version not found")
	ErrInvalidVersion = errors.New("invalid version")
	ErrForbidden      = errors.New("forbidden")
)
