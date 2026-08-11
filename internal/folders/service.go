package folders

import (
	"context"
	"errors"
	"strings"

	"github.com/abass-codes/peakcloud/internal/authorization"
)

type Service struct {
	repository    *Repository
	authorization *authorization.Service
}

func NewService(
	repository *Repository,
	authorizationService *authorization.Service,
) *Service {
	return &Service{
		repository:    repository,
		authorization: authorizationService,
	}
}

func (s *Service) Create(
	ctx context.Context,
	userID string,
	parentID *string,
	name string,
) (*Folder, error) {
	name = strings.TrimSpace(name)

	if err := ValidateName(name); err != nil {
		return nil, err
	}

	ownerID := userID

	if parentID != nil {
		if err := s.authorize(
			ctx,
			userID,
			*parentID,
			authorization.ActionEdit,
		); err != nil {
			return nil, err
		}

		parent, err := s.repository.GetByID(
			ctx,
			*parentID,
		)
		if err != nil {
			return nil, err
		}

		ownerID = parent.OwnerID
	}

	folder := &Folder{
		OwnerID:  ownerID,
		ParentID: parentID,
		Name:     name,
	}

	if err := s.repository.Create(ctx, folder); err != nil {
		return nil, err
	}

	return folder, nil
}

func (s *Service) Get(
	ctx context.Context,
	id string,
	userID string,
) (*Folder, error) {
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

func (s *Service) List(
	ctx context.Context,
	userID string,
	parentID *string,
) ([]Folder, error) {
	if parentID == nil {
		return s.repository.ListChildren(
			ctx,
			userID,
			nil,
		)
	}

	if err := s.authorize(
		ctx,
		userID,
		*parentID,
		authorization.ActionRead,
	); err != nil {
		return nil, err
	}

	return s.repository.ListChildrenByFolder(
		ctx,
		*parentID,
	)
}

func (s *Service) Rename(
	ctx context.Context,
	id string,
	userID string,
	name string,
) (*Folder, error) {
	name = strings.TrimSpace(name)

	if err := ValidateName(name); err != nil {
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

	return s.repository.RenameByID(
		ctx,
		id,
		name,
	)
}

func (s *Service) Move(
	ctx context.Context,
	id string,
	userID string,
	parentID *string,
) (*Folder, error) {
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

	if parentID == nil {
		if source.OwnerID != userID {
			return nil, authorization.ErrForbidden
		}

		return s.repository.MoveByID(
			ctx,
			id,
			nil,
		)
	}

	if id == *parentID {
		return nil, ErrInvalidMove
	}

	if err := s.authorize(
		ctx,
		userID,
		*parentID,
		authorization.ActionEdit,
	); err != nil {
		return nil, err
	}

	destination, err := s.repository.GetByID(
		ctx,
		*parentID,
	)
	if err != nil {
		return nil, err
	}

	if destination.OwnerID != source.OwnerID {
		return nil, authorization.ErrForbidden
	}

	descendant, err := s.repository.IsDescendant(
		ctx,
		source.OwnerID,
		id,
		*parentID,
	)
	if err != nil {
		return nil, err
	}

	if descendant {
		return nil, ErrInvalidMove
	}

	return s.repository.MoveByID(
		ctx,
		id,
		parentID,
	)
}

func (s *Service) Breadcrumbs(
	ctx context.Context,
	id string,
	userID string,
) ([]Folder, error) {
	if err := s.authorize(
		ctx,
		userID,
		id,
		authorization.ActionRead,
	); err != nil {
		return nil, err
	}

	folder, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.repository.Breadcrumbs(
		ctx,
		id,
		folder.OwnerID,
	)
}

func (s *Service) Delete(
	ctx context.Context,
	id string,
	userID string,
) error {
	access, err := s.authorization.Access(
		ctx,
		userID,
		authorization.ResourceFolder,
		id,
	)
	if err != nil {
		return mapAuthorizationError(err)
	}

	if !access.Owner {
		return authorization.ErrForbidden
	}

	return s.repository.Delete(
		ctx,
		id,
		userID,
	)
}

func (s *Service) authorize(
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
