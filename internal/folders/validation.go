package folders

import (
	"errors"
	"strings"
)

var (
	ErrNotFound      = errors.New("folder not found")
	ErrInvalidName   = errors.New("invalid folder name")
	ErrInvalidParent = errors.New("invalid parent folder")
	ErrDuplicateName = errors.New("folder already exists")
	ErrInvalidMove   = errors.New("invalid folder move")
)

func ValidateName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return ErrInvalidName
	}

	if len(name) > 255 {
		return ErrInvalidName
	}

	if name == "." || name == ".." {
		return ErrInvalidName
	}

	if strings.ContainsAny(name, `/\`) {
		return ErrInvalidName
	}

	return nil
}
