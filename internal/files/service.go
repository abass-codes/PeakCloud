package files

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"

	"github.com/abass-codes/peakcloud/internal/authorization"
	"github.com/abass-codes/peakcloud/internal/folders"
	"github.com/abass-codes/peakcloud/internal/storage"
)

type Service struct {
	repository       *Repository
	folderRepository *folders.Repository
	authorization    *authorization.Service
	objectStore      *storage.ObjectStore
	maxUploadSize    int64
}

func NewService(
	repository *Repository,
	folderRepository *folders.Repository,
	authorizationService *authorization.Service,
	objectStore *storage.ObjectStore,
	maxUploadSize int64,
) *Service {
	return &Service{
		repository:       repository,
		folderRepository: folderRepository,
		authorization:    authorizationService,
		objectStore:      objectStore,
		maxUploadSize:    maxUploadSize,
	}
}

func (s *Service) Upload(
	ctx context.Context,
	ownerID string,
	folderID *string,
	filename string,
	contentType string,
	reader io.Reader,
	size int64,
) (*File, error) {
	if err := ValidateFilename(filename); err != nil {
		return nil, err
	}

	if size < 0 || size > s.maxUploadSize {
		return nil, ErrFileTooLarge
	}

	resourceOwnerID := ownerID

	if folderID != nil {
		if err := s.authorizeFolder(
			ctx,
			ownerID,
			*folderID,
			authorization.ActionEdit,
		); err != nil {
			return nil, err
		}

		destination, err := s.folderRepository.GetByID(
			ctx,
			*folderID,
		)
		if err != nil {
			return nil, err
		}

		resourceOwnerID = destination.OwnerID
	}

	contentType = normalizeContentType(filename, contentType)
	objectKey := NewObjectKey(resourceOwnerID)

	objectInfo, err := s.objectStore.Put(
		ctx,
		objectKey,
		reader,
		size,
		contentType,
	)
	if err != nil {
		return nil, fmt.Errorf("store file object: %w", err)
	}

	file := &File{
		OwnerID:      resourceOwnerID,
		FolderID:     folderID,
		ObjectKey:    objectKey,
		OriginalName: filename,
		ContentType:  contentType,
		SizeBytes:    objectInfo.Size,
		ETag:         objectInfo.ETag,
	}

	if err := s.repository.Create(ctx, file); err != nil {
		if deleteErr := s.objectStore.Delete(ctx, objectKey); deleteErr != nil {
			return nil, fmt.Errorf(
				"persist metadata: %w; rollback object: %v",
				err,
				deleteErr,
			)
		}

		return nil, fmt.Errorf("persist metadata: %w", err)
	}

	return file, nil
}

func (s *Service) List(
	ctx context.Context,
	userID string,
	folderID *string,
) ([]File, error) {
	if folderID == nil {
		return s.repository.ListByFolder(
			ctx,
			userID,
			nil,
		)
	}

	if err := s.authorizeFolder(
		ctx,
		userID,
		*folderID,
		authorization.ActionRead,
	); err != nil {
		return nil, err
	}

	folder, err := s.folderRepository.GetByID(
		ctx,
		*folderID,
	)
	if err != nil {
		if errors.Is(err, folders.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return s.repository.ListByFolder(
		ctx,
		folder.OwnerID,
		folderID,
	)
}

func (s *Service) Get(
	ctx context.Context,
	id string,
	userID string,
) (*File, error) {
	if err := s.authorize(
		ctx,
		userID,
		id,
		authorization.ActionRead,
	); err != nil {
		if errors.Is(err, authorization.ErrForbidden) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return s.repository.GetByID(ctx, id)
}

func (s *Service) Content(
	ctx context.Context,
	id string,
	userID string,
) (*File, *minio.Object, error) {
	if err := s.authorize(
		ctx,
		userID,
		id,
		authorization.ActionRead,
	); err != nil {
		return nil, nil, err
	}

	file, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	object, err := s.objectStore.Get(
		ctx,
		file.ObjectKey,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"get stored object: %w",
			err,
		)
	}

	return file, object, nil
}

func (s *Service) Download(
	ctx context.Context,
	id string,
	userID string,
) (*File, *minio.Object, error) {
	if err := s.authorize(
		ctx,
		userID,
		id,
		authorization.ActionDownload,
	); err != nil {
		return nil, nil, err
	}

	file, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	object, err := s.objectStore.Get(ctx, file.ObjectKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get stored object: %w", err)
	}

	return file, object, nil
}

func (s *Service) Rename(
	ctx context.Context,
	id string,
	userID string,
	name string,
) (*File, error) {
	name = strings.TrimSpace(name)

	if err := ValidateFilename(name); err != nil {
		return nil, err
	}

	if err := s.authorize(
		ctx,
		userID,
		id,
		authorization.ActionEdit,
	); err != nil {
		return nil, err
	}

	return s.repository.RenameByID(ctx, id, name)
}

func (s *Service) Move(
	ctx context.Context,
	id string,
	userID string,
	folderID *string,
) (*File, error) {
	source, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.authorize(
		ctx,
		userID,
		id,
		authorization.ActionEdit,
	); err != nil {
		return nil, err
	}

	if folderID == nil {
		if source.OwnerID != userID {
			return nil, authorization.ErrForbidden
		}

		return s.repository.MoveByID(ctx, id, nil)
	}

	if err := s.authorizeFolder(
		ctx,
		userID,
		*folderID,
		authorization.ActionEdit,
	); err != nil {
		return nil, err
	}

	destination, err := s.folderRepository.GetByIDAndOwner(
		ctx,
		*folderID,
		source.OwnerID,
	)
	if err != nil {
		if errors.Is(err, folders.ErrNotFound) {
			return nil, authorization.ErrForbidden
		}

		return nil, err
	}

	return s.repository.MoveByID(
		ctx,
		id,
		&destination.ID,
	)
}

func (s *Service) Copy(
	ctx context.Context,
	id string,
	userID string,
	folderID *string,
) (*File, error) {
	if err := s.authorize(
		ctx,
		userID,
		id,
		authorization.ActionEdit,
	); err != nil {
		return nil, err
	}

	source, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if folderID == nil {
		if source.OwnerID != userID {
			return nil, authorization.ErrForbidden
		}
	} else {
		if err := s.authorizeFolder(
			ctx,
			userID,
			*folderID,
			authorization.ActionEdit,
		); err != nil {
			return nil, err
		}

		if _, err := s.folderRepository.GetByIDAndOwner(
			ctx,
			*folderID,
			source.OwnerID,
		); err != nil {
			if errors.Is(err, folders.ErrNotFound) {
				return nil, authorization.ErrForbidden
			}

			return nil, err
		}
	}

	objectKey := NewObjectKey(source.OwnerID)

	info, err := s.objectStore.Copy(
		ctx,
		source.ObjectKey,
		objectKey,
	)
	if err != nil {
		return nil, err
	}

	copyFile := &File{
		OwnerID:      source.OwnerID,
		FolderID:     folderID,
		ObjectKey:    objectKey,
		OriginalName: source.OriginalName,
		ContentType:  source.ContentType,
		SizeBytes:    source.SizeBytes,
		ETag:         info.ETag,
	}

	if err := s.repository.Create(ctx, copyFile); err != nil {
		_ = s.objectStore.Delete(ctx, objectKey)
		return nil, err
	}

	return copyFile, nil
}

func (s *Service) Delete(
	ctx context.Context,
	id string,
	userID string,
) error {
	access, err := s.authorization.Access(
		ctx,
		userID,
		authorization.ResourceFile,
		id,
	)
	if err != nil {
		return mapAuthorizationError(err)
	}

	if !access.Owner {
		return authorization.ErrForbidden
	}

	file, err := s.repository.GetByIDAndOwner(
		ctx,
		id,
		userID,
	)
	if err != nil {
		return err
	}

	if err := s.objectStore.Delete(ctx, file.ObjectKey); err != nil {
		return fmt.Errorf("delete stored object: %w", err)
	}

	if err := s.repository.DeleteMetadata(
		ctx,
		id,
		userID,
	); err != nil {
		return fmt.Errorf("delete metadata: %w", err)
	}

	return nil
}

func (s *Service) authorize(
	ctx context.Context,
	userID string,
	fileID string,
	action authorization.Action,
) error {
	_, err := s.authorization.Authorize(
		ctx,
		userID,
		authorization.ResourceFile,
		fileID,
		action,
	)

	return mapAuthorizationError(err)
}

func (s *Service) authorizeFolder(
	ctx context.Context,
	userID string,
	folderID string,
	action authorization.Action,
) error {
	_, err := s.authorization.Authorize(
		ctx,
		userID,
		authorization.ResourceFolder,
		folderID,
		action,
	)

	return mapAuthorizationError(err)
}

func mapAuthorizationError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, authorization.ErrResourceNotFound) {
		return ErrNotFound
	}

	return err
}

func (s *Service) validateOwnedFolder(
	ctx context.Context,
	ownerID string,
	folderID *string,
) error {
	if folderID == nil {
		return nil
	}

	_, err := s.folderRepository.GetByIDAndOwner(
		ctx,
		*folderID,
		ownerID,
	)

	if errors.Is(err, folders.ErrNotFound) {
		return ErrNotFound
	}

	return err
}

func normalizeContentType(filename string, contentType string) string {
	contentType = strings.TrimSpace(contentType)

	if contentType != "" &&
		contentType != "application/octet-stream" {
		return contentType
	}

	extension := filepath.Ext(filename)

	if extension != "" {
		if detected := mime.TypeByExtension(extension); detected != "" {
			return detected
		}
	}

	return "application/octet-stream"
}
