package files

import (
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"

	"github.com/abass-codes/peakcloud/internal/storage"
)

type Service struct {
	repository    *Repository
	objectStore   *storage.ObjectStore
	maxUploadSize int64
}

func NewService(
	repository *Repository,
	objectStore *storage.ObjectStore,
	maxUploadSize int64,
) *Service {
	return &Service{
		repository:    repository,
		objectStore:   objectStore,
		maxUploadSize: maxUploadSize,
	}
}

func (s *Service) Upload(
	ctx context.Context,
	ownerID string,
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

func (s *Service) List(ctx context.Context, ownerID string) ([]File, error) {
	return s.repository.ListByOwner(ctx, ownerID)
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
		return nil, nil, fmt.Errorf("load object: %w", err)
	}

	return file, object, nil
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
		return fmt.Errorf("delete object: %w", err)
	}

	if err := s.repository.DeleteByIDAndOwner(ctx, id, ownerID); err != nil {
		return fmt.Errorf("delete metadata: %w", err)
	}

	return nil
}

func normalizeContentType(filename, supplied string) string {
	supplied = strings.TrimSpace(supplied)

	if supplied != "" && supplied != "application/octet-stream" {
		return supplied
	}

	extension := filepath.Ext(filename)
	if extension != "" {
		if detected := mime.TypeByExtension(extension); detected != "" {
			return detected
		}
	}

	return "application/octet-stream"
}
