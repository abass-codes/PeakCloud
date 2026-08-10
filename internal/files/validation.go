package files

import (
	"errors"
	"path/filepath"
	"strings"
)

var (
	ErrInvalidFilename = errors.New("invalid filename")
	ErrFileTooLarge    = errors.New("file exceeds maximum upload size")
	ErrNotFound        = errors.New("file not found")
)

func ValidateFilename(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return ErrInvalidFilename
	}

	if len(name) > 255 {
		return ErrInvalidFilename
	}

	if name == "." || name == ".." {
		return ErrInvalidFilename
	}

	if filepath.Base(name) != name {
		return ErrInvalidFilename
	}

	if strings.ContainsAny(name, `/\`) {
		return ErrInvalidFilename
	}

	return nil
}
