package folders

import (
	"context"
	"strings"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	ownerID string,
	parentID *string,
	name string,
) (*Folder, error) {
	name = strings.TrimSpace(name)

	if err := ValidateName(name); err != nil {
		return nil, err
	}

	if parentID != nil {
		if _, err := s.repository.GetByIDAndOwner(
			ctx,
			*parentID,
			ownerID,
		); err != nil {
			if err == ErrNotFound {
				return nil, ErrInvalidParent
			}

			return nil, err
		}
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
	ownerID string,
) (*Folder, error) {
	return s.repository.GetByIDAndOwner(ctx, id, ownerID)
}

func (s *Service) List(
	ctx context.Context,
	ownerID string,
	parentID *string,
) ([]Folder, error) {
	if parentID != nil {
		if _, err := s.repository.GetByIDAndOwner(
			ctx,
			*parentID,
			ownerID,
		); err != nil {
			return nil, err
		}
	}

	return s.repository.ListChildren(ctx, ownerID, parentID)
}

func (s *Service) Rename(
	ctx context.Context,
	id string,
	ownerID string,
	name string,
) (*Folder, error) {
	name = strings.TrimSpace(name)

	if err := ValidateName(name); err != nil {
		return nil, err
	}

	return s.repository.Rename(ctx, id, ownerID, name)
}

func (s *Service) Move(
	ctx context.Context,
	id string,
	ownerID string,
	parentID *string,
) (*Folder, error) {
	if parentID == nil {
		return s.repository.Move(ctx, id, ownerID, nil)
	}

	if id == *parentID {
		return nil, ErrInvalidMove
	}

	if _, err := s.repository.GetByIDAndOwner(
		ctx,
		*parentID,
		ownerID,
	); err != nil {
		if err == ErrNotFound {
			return nil, ErrInvalidParent
		}

		return nil, err
	}

	descendant, err := s.repository.IsDescendant(
		ctx,
		ownerID,
		id,
		*parentID,
	)
	if err != nil {
		return nil, err
	}

	if descendant {
		return nil, ErrInvalidMove
	}

	return s.repository.Move(ctx, id, ownerID, parentID)
}

func (s *Service) Breadcrumbs(
	ctx context.Context,
	id string,
	ownerID string,
) ([]Folder, error) {
	return s.repository.Breadcrumbs(ctx, id, ownerID)
}

func (s *Service) Delete(
	ctx context.Context,
	id string,
	ownerID string,
) error {
	if _, err := s.repository.GetByIDAndOwner(
		ctx,
		id,
		ownerID,
	); err != nil {
		return err
	}

	return s.repository.Delete(ctx, id, ownerID)
}
