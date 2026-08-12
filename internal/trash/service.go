package trash

import (
	"context"
	"errors"
	"fmt"

	"github.com/abass-codes/peakcloud/internal/files"
	"github.com/abass-codes/peakcloud/internal/folders"
)

type fileRepository interface {
	Trash(
		ctx context.Context,
		id string,
		ownerID string,
	) (*files.File, error)

	RestoreFromTrash(
		ctx context.Context,
		id string,
		ownerID string,
	) (*files.File, error)

	DeletePermanently(
		ctx context.Context,
		id string,
		ownerID string,
	) (*files.File, error)

	ListTrashed(
		ctx context.Context,
		ownerID string,
	) ([]files.File, error)
}

type folderRepository interface {
	Trash(
		ctx context.Context,
		id string,
		ownerID string,
	) (*folders.Folder, error)

	RestoreFromTrash(
		ctx context.Context,
		id string,
		ownerID string,
	) (*folders.Folder, error)

	DeletePermanently(
		ctx context.Context,
		id string,
		ownerID string,
	) (*folders.Folder, error)

	ListTrashed(
		ctx context.Context,
		ownerID string,
	) ([]folders.Folder, error)
}

type Service struct {
	files   fileRepository
	folders folderRepository
}

func NewService(
	files fileRepository,
	folders folderRepository,
) *Service {
	return &Service{
		files:   files,
		folders: folders,
	}
}

type Contents struct {
	Files   []files.File     `json:"files"`
	Folders []folders.Folder `json:"folders"`
}

func (s *Service) List(
	ctx context.Context,
	ownerID string,
) (*Contents, error) {
	trashedFiles, err := s.files.ListTrashed(
		ctx,
		ownerID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"list trashed files: %w",
			err,
		)
	}

	trashedFolders, err := s.folders.ListTrashed(
		ctx,
		ownerID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"list trashed folders: %w",
			err,
		)
	}

	return &Contents{
		Files:   trashedFiles,
		Folders: trashedFolders,
	}, nil
}

func (s *Service) Trash(
	ctx context.Context,
	ownerID string,
	resourceType ResourceType,
	resourceID string,
) error {
	if err := ValidateResourceType(resourceType); err != nil {
		return err
	}

	switch resourceType {
	case ResourceFile:
		_, err := s.files.Trash(
			ctx,
			resourceID,
			ownerID,
		)
		if err != nil {
			if errors.Is(err, files.ErrNotFound) {
				return ErrNotFound
			}

			return fmt.Errorf(
				"trash file: %w",
				err,
			)
		}

	case ResourceFolder:
		_, err := s.folders.Trash(
			ctx,
			resourceID,
			ownerID,
		)
		if err != nil {
			if errors.Is(err, folders.ErrNotFound) {
				return ErrNotFound
			}

			return fmt.Errorf(
				"trash folder: %w",
				err,
			)
		}
	}

	return nil
}

func (s *Service) Restore(
	ctx context.Context,
	ownerID string,
	resourceType ResourceType,
	resourceID string,
) error {
	if err := ValidateResourceType(resourceType); err != nil {
		return err
	}

	switch resourceType {
	case ResourceFile:
		_, err := s.files.RestoreFromTrash(
			ctx,
			resourceID,
			ownerID,
		)
		if err != nil {
			if errors.Is(err, files.ErrNotFound) {
				return ErrNotFound
			}

			return fmt.Errorf(
				"restore file: %w",
				err,
			)
		}

	case ResourceFolder:
		_, err := s.folders.RestoreFromTrash(
			ctx,
			resourceID,
			ownerID,
		)
		if err != nil {
			if errors.Is(err, folders.ErrNotFound) {
				return ErrNotFound
			}

			return fmt.Errorf(
				"restore folder: %w",
				err,
			)
		}
	}

	return nil
}

func (s *Service) DeletePermanently(
	ctx context.Context,
	ownerID string,
	resourceType ResourceType,
	resourceID string,
) error {
	if err := ValidateResourceType(resourceType); err != nil {
		return err
	}

	switch resourceType {
	case ResourceFile:
		_, err := s.files.DeletePermanently(
			ctx,
			resourceID,
			ownerID,
		)
		if err != nil {
			if errors.Is(err, files.ErrNotFound) {
				return ErrNotFound
			}

			return fmt.Errorf(
				"permanently delete file: %w",
				err,
			)
		}

	case ResourceFolder:
		_, err := s.folders.DeletePermanently(
			ctx,
			resourceID,
			ownerID,
		)
		if err != nil {
			if errors.Is(err, folders.ErrNotFound) {
				return ErrNotFound
			}

			return fmt.Errorf(
				"permanently delete folder: %w",
				err,
			)
		}
	}

	return nil
}
