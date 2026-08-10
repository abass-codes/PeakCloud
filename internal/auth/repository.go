package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotFound = errors.New("user not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(
	ctx context.Context,
	email string,
	displayName string,
	passwordHash string,
) (*User, error) {
	const query = `
		INSERT INTO users (email, display_name, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email, display_name, password_hash, created_at, updated_at
	`

	var user User

	err := r.db.QueryRow(
		ctx,
		query,
		strings.ToLower(strings.TrimSpace(email)),
		strings.TrimSpace(displayName),
		passwordHash,
	).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &user, nil
}

func (r *Repository) FindUserByEmail(
	ctx context.Context,
	email string,
) (*User, error) {
	const query = `
		SELECT id, email, display_name, password_hash, created_at, updated_at
		FROM users
		WHERE lower(email) = lower($1)
	`

	var user User

	err := r.db.QueryRow(
		ctx,
		query,
		strings.TrimSpace(email),
	).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}

	return &user, nil
}

func (r *Repository) CreateSession(
	ctx context.Context,
	userID string,
	tokenHash string,
	expiresAt time.Time,
) error {
	const query = `
		INSERT INTO sessions (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`

	if _, err := r.db.Exec(
		ctx,
		query,
		userID,
		tokenHash,
		expiresAt,
	); err != nil {
		return fmt.Errorf("create session: %w", err)
	}

	return nil
}

func (r *Repository) FindUserBySession(
	ctx context.Context,
	tokenHash string,
) (*User, error) {
	const query = `
		SELECT
			u.id,
			u.email,
			u.display_name,
			u.password_hash,
			u.created_at,
			u.updated_at
		FROM sessions s
		JOIN users u ON u.id = s.user_id
		WHERE s.token_hash = $1
		  AND s.expires_at > NOW()
	`

	var user User

	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("find user by session: %w", err)
	}

	return &user, nil
}

func (r *Repository) DeleteSession(
	ctx context.Context,
	tokenHash string,
) error {
	const query = `DELETE FROM sessions WHERE token_hash = $1`

	if _, err := r.db.Exec(ctx, query, tokenHash); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}
