package sharing

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) resourceOwnedBy(
	ctx context.Context,
	resourceType ResourceType,
	resourceID string,
	ownerID string,
) (string, error) {
	switch resourceType {
	case ResourceFile:
		name, ok, err := s.repository.FileOwnedBy(
			ctx,
			resourceID,
			ownerID,
		)
		if err != nil {
			return "", err
		}
		if !ok {
			return "", ErrForbidden
		}

		return name, nil

	case ResourceFolder:
		name, ok, err := s.repository.FolderOwnedBy(
			ctx,
			resourceID,
			ownerID,
		)
		if err != nil {
			return "", err
		}
		if !ok {
			return "", ErrForbidden
		}

		return name, nil

	default:
		return "", ErrInvalidResource
	}
}

func (s *Service) CreateShare(
	ctx context.Context,
	ownerID string,
	input CreateShareRequest,
) (*Share, error) {
	input.RecipientEmail = NormalizeEmail(input.RecipientEmail)

	if input.RecipientEmail == "" {
		return nil, ErrRecipientNotFound
	}

	if !ValidResourceType(input.ResourceType) {
		return nil, ErrInvalidResource
	}

	if !ValidPermission(input.Permission) {
		return nil, ErrInvalidPermission
	}

	resourceName, err := s.resourceOwnedBy(
		ctx,
		input.ResourceType,
		input.ResourceID,
		ownerID,
	)
	if err != nil {
		return nil, err
	}

	recipientID, recipientEmail, err :=
		s.repository.FindUserByEmail(ctx, input.RecipientEmail)
	if err != nil {
		return nil, err
	}

	if recipientID == ownerID {
		return nil, ErrCannotShareSelf
	}

	share := &Share{
		ID:             uuid.NewString(),
		OwnerID:        ownerID,
		RecipientID:    recipientID,
		RecipientEmail: recipientEmail,
		ResourceType:   input.ResourceType,
		ResourceID:     input.ResourceID,
		ResourceName:   resourceName,
		Permission:     input.Permission,
		AllowDownload:  input.AllowDownload,
	}

	if err := s.repository.CreateShare(ctx, share); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrAlreadyShared
		}

		return nil, err
	}

	return share, nil
}

func (s *Service) ListOwnedShares(
	ctx context.Context,
	ownerID string,
) ([]Share, error) {
	return s.repository.ListOwnedShares(ctx, ownerID)
}

func (s *Service) SharedWithMe(
	ctx context.Context,
	userID string,
) ([]Share, error) {
	return s.repository.ListSharedWithMe(ctx, userID)
}

func (s *Service) UpdateShare(
	ctx context.Context,
	ownerID string,
	id string,
	input UpdateShareRequest,
) error {
	if !ValidPermission(input.Permission) {
		return ErrInvalidPermission
	}

	return s.repository.UpdateShare(
		ctx,
		id,
		ownerID,
		input.Permission,
		input.AllowDownload,
	)
}

func (s *Service) DeleteShare(
	ctx context.Context,
	ownerID string,
	id string,
) error {
	return s.repository.DeleteShare(ctx, id, ownerID)
}

func (s *Service) CreatePublicLink(
	ctx context.Context,
	ownerID string,
	input CreatePublicLinkRequest,
) (*PublicLinkCreated, error) {
	if !ValidResourceType(input.ResourceType) {
		return nil, ErrInvalidResource
	}

	if !ValidPermission(input.Permission) {
		return nil, ErrInvalidPermission
	}

	resourceName, err := s.resourceOwnedBy(
		ctx,
		input.ResourceType,
		input.ResourceID,
		ownerID,
	)
	if err != nil {
		return nil, err
	}

	if input.ExpiresAt != nil &&
		!input.ExpiresAt.After(time.Now()) {
		return nil, ErrExpired
	}

	token, err := GenerateToken()
	if err != nil {
		return nil, err
	}

	passwordHash, err := HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	record := &publicLinkRecord{
		PublicLink: PublicLink{
			ID:            uuid.NewString(),
			OwnerID:       ownerID,
			ResourceType:  input.ResourceType,
			ResourceID:    input.ResourceID,
			ResourceName:  resourceName,
			Permission:    input.Permission,
			AllowDownload: input.AllowDownload,
			PasswordSet:   passwordHash != "",
			ExpiresAt:     input.ExpiresAt,
		},
		TokenHash:    HashToken(token),
		PasswordHash: passwordHash,
	}

	if err := s.repository.CreatePublicLink(ctx, record); err != nil {
		return nil, err
	}

	return &PublicLinkCreated{
		PublicLink: record.PublicLink,
		Token:      token,
	}, nil
}

func (s *Service) resolveRecord(
	ctx context.Context,
	token string,
	password string,
) (*publicLinkRecord, error) {
	record, err := s.repository.FindPublicLink(
		ctx,
		HashToken(token),
	)
	if err != nil {
		return nil, err
	}

	if record.RevokedAt != nil {
		return nil, ErrRevoked
	}

	if record.ExpiresAt != nil &&
		time.Now().After(*record.ExpiresAt) {
		return nil, ErrExpired
	}

	if record.PasswordHash != "" {
		if password == "" {
			return nil, ErrPasswordRequired
		}

		if !VerifyPassword(record.PasswordHash, password) {
			return nil, ErrInvalidPassword
		}
	}

	return record, nil
}

func (s *Service) ResolvePublicLink(
	ctx context.Context,
	token string,
	password string,
) (*PublicLink, error) {
	record, err := s.resolveRecord(ctx, token, password)
	if err != nil {
		return nil, err
	}

	name, err := s.resourceOwnedBy(
		ctx,
		record.ResourceType,
		record.ResourceID,
		record.OwnerID,
	)
	if err != nil {
		return nil, ErrNotFound
	}

	record.ResourceName = name
	record.PasswordSet = record.PasswordHash != ""

	result := record.PublicLink
	result.OwnerID = ""

	return &result, nil
}

func (s *Service) ResolvePublicFile(
	ctx context.Context,
	token string,
	password string,
	requireDownload bool,
) (*PublicFile, error) {
	record, err := s.resolveRecord(ctx, token, password)
	if err != nil {
		return nil, err
	}

	if record.ResourceType != ResourceFile {
		return nil, ErrInvalidResource
	}

	if requireDownload && !record.AllowDownload {
		return nil, ErrDownloadDenied
	}

	return s.repository.GetPublicFile(
		ctx,
		record.ResourceID,
		record.OwnerID,
	)
}

func (s *Service) ListPublicLinks(
	ctx context.Context,
	ownerID string,
) ([]PublicLink, error) {
	return s.repository.ListPublicLinks(ctx, ownerID)
}

func (s *Service) RevokePublicLink(
	ctx context.Context,
	ownerID string,
	id string,
) error {
	return s.repository.RevokePublicLink(ctx, id, ownerID)
}
