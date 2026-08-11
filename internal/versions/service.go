package versions

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/abass-codes/peakcloud/internal/authorization"
	"github.com/abass-codes/peakcloud/internal/files"
	"github.com/abass-codes/peakcloud/internal/storage"
)

type Service struct {
	repository     *Repository
	fileRepository *files.Repository
	authorization  *authorization.Service
	objectStore    *storage.ObjectStore
	maxUploadSize  int64
}

func NewService(
	repository *Repository,
	fileRepository *files.Repository,
	authorizationService *authorization.Service,
	objectStore *storage.ObjectStore,
	maxUploadSize int64,
) *Service {
	return &Service{
		repository:     repository,
		fileRepository: fileRepository,
		authorization:  authorizationService,
		objectStore:    objectStore,
		maxUploadSize:  maxUploadSize,
	}
}

func (s *Service) Create(
	ctx context.Context,
	userID string,
	fileID string,
	contentType string,
	reader io.Reader,
	size int64,
) (*Version, error) {
	if size < 0 || size > s.maxUploadSize {
		return nil, files.ErrFileTooLarge
	}

	if _, err := s.authorization.Authorize(
		ctx,
		userID,
		authorization.ResourceFile,
		fileID,
		authorization.ActionEdit,
	); err != nil {
		return nil, mapAuthorizationError(err)
	}

	file, err := s.fileRepository.GetByID(ctx, fileID)
	if err != nil {
		if errors.Is(err, files.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	if contentType == "" {
		contentType = file.ContentType
	}

	versionNumber, err := s.repository.NextVersionNumber(
		ctx,
		fileID,
	)
	if err != nil {
		return nil, err
	}

	objectKey := NewObjectKey(
		file.OwnerID,
		file.ID,
		versionNumber,
	)

	objectInfo, err := s.objectStore.Put(
		ctx,
		objectKey,
		reader,
		size,
		contentType,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"store version object: %w",
			err,
		)
	}

	version := &Version{
		FileID:        file.ID,
		VersionNumber: versionNumber,
		ObjectKey:     objectKey,
		SizeBytes:     objectInfo.Size,
		ContentType:   contentType,
		ETag:          objectInfo.ETag,
		CreatedBy:     userID,
	}

	if err := s.repository.Create(ctx, version); err != nil {
		_ = s.objectStore.Delete(ctx, objectKey)
		return nil, fmt.Errorf(
			"persist version metadata: %w",
			err,
		)
	}

	if _, err := s.fileRepository.UpdateContent(
		ctx,
		file.ID,
		objectKey,
		contentType,
		objectInfo.Size,
		objectInfo.ETag,
	); err != nil {
		return nil, fmt.Errorf(
			"update current file content: %w",
			err,
		)
	}

	return version, nil
}

func (s *Service) List(
	ctx context.Context,
	userID string,
	fileID string,
) ([]Version, error) {
	if _, err := s.authorization.Authorize(
		ctx,
		userID,
		authorization.ResourceFile,
		fileID,
		authorization.ActionRead,
	); err != nil {
		return nil, mapAuthorizationError(err)
	}

	return s.repository.ListByFile(ctx, fileID)
}

func (s *Service) Get(
	ctx context.Context,
	userID string,
	fileID string,
	versionNumber int,
) (*Version, error) {
	if err := ValidateVersionNumber(versionNumber); err != nil {
		return nil, err
	}

	if _, err := s.authorization.Authorize(
		ctx,
		userID,
		authorization.ResourceFile,
		fileID,
		authorization.ActionRead,
	); err != nil {
		return nil, mapAuthorizationError(err)
	}

	return s.repository.GetByNumber(
		ctx,
		fileID,
		versionNumber,
	)
}

func (s *Service) Content(
	ctx context.Context,
	userID string,
	fileID string,
	versionNumber int,
) (*Version, io.ReadCloser, error) {
	version, err := s.Get(
		ctx,
		userID,
		fileID,
		versionNumber,
	)
	if err != nil {
		return nil, nil, err
	}

	object, err := s.objectStore.Get(
		ctx,
		version.ObjectKey,
	)
	if err != nil {
		return nil, nil, err
	}

	return version, object, nil
}

func (s *Service) Download(
	ctx context.Context,
	userID string,
	fileID string,
	versionNumber int,
) (*Version, io.ReadCloser, error) {
	if err := ValidateVersionNumber(versionNumber); err != nil {
		return nil, nil, err
	}

	if _, err := s.authorization.Authorize(
		ctx,
		userID,
		authorization.ResourceFile,
		fileID,
		authorization.ActionDownload,
	); err != nil {
		return nil, nil, mapAuthorizationError(err)
	}

	version, err := s.repository.GetByNumber(
		ctx,
		fileID,
		versionNumber,
	)
	if err != nil {
		return nil, nil, err
	}

	object, err := s.objectStore.Get(
		ctx,
		version.ObjectKey,
	)
	if err != nil {
		return nil, nil, err
	}

	return version, object, nil
}

func mapAuthorizationError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(
		err,
		authorization.ErrResourceNotFound,
	) {
		return ErrNotFound
	}

	if errors.Is(err, authorization.ErrForbidden) {
		return ErrForbidden
	}

	return err
}

func (s *Service) Restore(
	ctx context.Context,
	userID string,
	fileID string,
	versionNumber int,
) (*Version, error) {
	if err := ValidateVersionNumber(versionNumber); err != nil {
		return nil, err
	}

	if _, err := s.authorization.Authorize(
		ctx,
		userID,
		authorization.ResourceFile,
		fileID,
		authorization.ActionEdit,
	); err != nil {
		return nil, mapAuthorizationError(err)
	}

	file, err := s.fileRepository.GetByID(ctx, fileID)
	if err != nil {
		if errors.Is(err, files.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	source, err := s.repository.GetByNumber(
		ctx,
		fileID,
		versionNumber,
	)
	if err != nil {
		return nil, err
	}

	object, err := s.objectStore.Get(
		ctx,
		source.ObjectKey,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get source version object: %w",
			err,
		)
	}
	defer object.Close()

	nextVersion, err := s.repository.NextVersionNumber(
		ctx,
		fileID,
	)
	if err != nil {
		return nil, err
	}

	objectKey := NewObjectKey(
		file.OwnerID,
		file.ID,
		nextVersion,
	)

	objectInfo, err := s.objectStore.Put(
		ctx,
		objectKey,
		object,
		source.SizeBytes,
		source.ContentType,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"store restored version object: %w",
			err,
		)
	}

	version := &Version{
		FileID:        file.ID,
		VersionNumber: nextVersion,
		ObjectKey:     objectKey,
		SizeBytes:     objectInfo.Size,
		ContentType:   source.ContentType,
		ETag:          objectInfo.ETag,
		CreatedBy:     userID,
	}

	if err := s.repository.Create(ctx, version); err != nil {
		_ = s.objectStore.Delete(ctx, objectKey)

		return nil, fmt.Errorf(
			"persist restored version metadata: %w",
			err,
		)
	}

	if _, err := s.fileRepository.UpdateContent(
		ctx,
		file.ID,
		objectKey,
		source.ContentType,
		objectInfo.Size,
		objectInfo.ETag,
	); err != nil {
		return nil, fmt.Errorf(
			"update restored file content: %w",
			err,
		)
	}

	return version, nil
}
