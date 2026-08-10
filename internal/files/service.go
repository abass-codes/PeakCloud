package files

import (
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"

	"github.com/abass-codes/peakcloud/internal/folders"
	"github.com/abass-codes/peakcloud/internal/storage"
)

type Service struct {
	repository       *Repository
	folderRepository *folders.Repository
	objectStore      *storage.ObjectStore
	maxUploadSize    int64
}

func NewService(
	repository *Repository,
	folderRepository *folders.Repository,
	objectStore *storage.ObjectStore,
	maxUploadSize int64,
) *Service {
	return &Service{
		repository:       repository,
		folderRepository: folderRepository,
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

	if err := s.validateFolder(ctx, ownerID, folderID); err != nil {
		return nil, err
	}

	contentType = normalizeContentType(filename, contentType)
	objectKey := NewObjectKey(ownerID)

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
		OwnerID:      ownerID,
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
	ownerID string,
	folderID *string,
) ([]File, error) {
	if err := s.validateFolder(ctx, ownerID, folderID); err != nil {
		return nil, err
	}

	return s.repository.ListByFolder(ctx, ownerID, folderID)
}

func (s *Service) Get(
	ctx context.Context,
	id string,
	ownerID string,
) (*File, error) {
	return s.repository.GetByIDAndOwner(ctx, id, ownerID)
}

func (s *Service) Download(
	ctx context.Context,
	id string,
	ownerID string,
) (*File, *minio.Object, error) {
	file, err := s.repository.GetByIDAndOwner(ctx, id, ownerID)
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
	ownerID string,
	name string,
) (*File, error) {
	name = strings.TrimSpace(name)

	if err := ValidateFilename(name); err != nil {
		return nil, err
	}

	return s.repository.Rename(ctx, id, ownerID, name)
}

func (s *Service) Move(
	ctx context.Context,
	id string,
	ownerID string,
	folderID *string,
) (*File, error) {
	if err := s.validateFolder(ctx, ownerID, folderID); err != nil {
		return nil, err
	}

	return s.repository.Move(ctx, id, ownerID, folderID)
}

func (s *Service) Copy(
	ctx context.Context,
	id string,
	ownerID string,
	folderID *string,
) (*File, error) {
	if err := s.validateFolder(ctx, ownerID, folderID); err != nil {
		return nil, err
	}

	source, err := s.repository.GetByIDAndOwner(ctx, id, ownerID)
	if err != nil {
		return nil, err
	}

	objectKey := NewObjectKey(ownerID)

	info, err := s.objectStore.Copy(
		ctx,
		source.ObjectKey,
		objectKey,
	)
	if err != nil {
		return nil, err
	}

	copyFile := &File{
		OwnerID:      ownerID,
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
	ownerID string,
) error {
	file, err := s.repository.GetByIDAndOwner(ctx, id, ownerID)
	if err != nil {
		return err
	}

	if err := s.objectStore.Delete(ctx, file.ObjectKey); err != nil {
		return fmt.Errorf("delete stored object: %w", err)
	}

	if err := s.repository.DeleteMetadata(ctx, id, ownerID); err != nil {
		return fmt.Errorf("delete metadata: %w", err)
	}

	return nil
}

func (s *Service) validateFolder(
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

	if err == folders.ErrNotFound {
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
