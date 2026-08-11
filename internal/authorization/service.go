package authorization

import "context"

type AccessRepository interface {
	FileAccess(
		ctx context.Context,
		userID string,
		fileID string,
	) (Access, error)

	FolderAccess(
		ctx context.Context,
		userID string,
		folderID string,
	) (Access, error)
}

type Service struct {
	repository AccessRepository
}

func NewService(repository AccessRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Access(
	ctx context.Context,
	userID string,
	resourceType ResourceType,
	resourceID string,
) (Access, error) {
	switch resourceType {
	case ResourceFile:
		return s.repository.FileAccess(
			ctx,
			userID,
			resourceID,
		)

	case ResourceFolder:
		return s.repository.FolderAccess(
			ctx,
			userID,
			resourceID,
		)

	default:
		return Access{}, ErrInvalidResource
	}
}

func (s *Service) Authorize(
	ctx context.Context,
	userID string,
	resourceType ResourceType,
	resourceID string,
	action Action,
) (Access, error) {
	access, err := s.Access(
		ctx,
		userID,
		resourceType,
		resourceID,
	)
	if err != nil {
		return Access{}, err
	}

	if !access.Allows(action) {
		return Access{}, ErrForbidden
	}

	return access, nil
}

func (s *Service) Can(
	ctx context.Context,
	userID string,
	resourceType ResourceType,
	resourceID string,
	action Action,
) (bool, error) {
	access, err := s.Access(
		ctx,
		userID,
		resourceType,
		resourceID,
	)
	if err != nil {
		return false, err
	}

	return access.Allows(action), nil
}
