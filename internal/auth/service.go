package auth

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Service struct {
	repository *Repository
	sessionTTL time.Duration
}

func NewService(repository *Repository, sessionTTL time.Duration) *Service {
	return &Service{
		repository: repository,
		sessionTTL: sessionTTL,
	}
}

func (s *Service) Register(
	ctx context.Context,
	email string,
	displayName string,
	password string,
) (*User, string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	displayName = strings.TrimSpace(displayName)

	if !validEmail(email) ||
		len(displayName) < 2 ||
		len(displayName) > 100 ||
		len(password) < 12 ||
		len(password) > 128 {
		return nil, "", ErrInvalidInput
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, "", fmt.Errorf("hash password: %w", err)
	}

	user, err := s.repository.CreateUser(
		ctx,
		email,
		displayName,
		passwordHash,
	)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, "", ErrEmailAlreadyExists
		}

		return nil, "", err
	}

	token, err := s.createSession(ctx, user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *Service) Login(
	ctx context.Context,
	email string,
	password string,
) (*User, string, error) {
	user, err := s.repository.FindUserByEmail(ctx, email)
	if errors.Is(err, ErrUserNotFound) {
		return nil, "", ErrInvalidCredentials
	}

	if err != nil {
		return nil, "", err
	}

	valid, err := VerifyPassword(password, user.PasswordHash)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	if !valid {
		return nil, "", ErrInvalidCredentials
	}

	token, err := s.createSession(ctx, user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *Service) Authenticate(
	ctx context.Context,
	token string,
) (*User, error) {
	if token == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := s.repository.FindUserBySession(
		ctx,
		HashSessionToken(token),
	)
	if errors.Is(err, ErrUserNotFound) {
		return nil, ErrInvalidCredentials
	}

	return user, err
}

func (s *Service) Logout(
	ctx context.Context,
	token string,
) error {
	if token == "" {
		return nil
	}

	return s.repository.DeleteSession(
		ctx,
		HashSessionToken(token),
	)
}

func (s *Service) createSession(
	ctx context.Context,
	userID string,
) (string, error) {
	token, err := GenerateSessionToken()
	if err != nil {
		return "", err
	}

	err = s.repository.CreateSession(
		ctx,
		userID,
		HashSessionToken(token),
		time.Now().Add(s.sessionTTL),
	)
	if err != nil {
		return "", err
	}

	return token, nil
}

func validEmail(value string) bool {
	address, err := mail.ParseAddress(value)
	if err != nil {
		return false
	}

	return strings.EqualFold(address.Address, value)
}
