package authorization

import "errors"

var (
	ErrForbidden        = errors.New("forbidden")
	ErrResourceNotFound = errors.New("resource not found")
	ErrInvalidResource  = errors.New("invalid resource type")
)
