package files

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, file *File) error {
	const query = `
		INSERT INTO files (
			id,
			owner_id,
			object_key,
			original_name,
			content_type,
			size_bytes,
			etag
		)
		VALUES (
			gen_random_uuid(),
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		file.OwnerID,
		file.ObjectKey,
		file.OriginalName,
		file.ContentType,
		file.SizeBytes,
		file.ETag,
	).Scan(
		&file.ID,
		&file.CreatedAt,
		&file.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create file metadata: %w", err)
	}

	return nil
}

func (r *Repository) ListByOwner(
	ctx context.Context,
	ownerID string,
) ([]File, error) {
	const query = `
		SELECT
			id,
			owner_id,
			object_key,
			original_name,
			content_type,
			size_bytes,
			COALESCE(etag, ''),
			created_at,
			updated_at
		FROM files
		WHERE owner_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("list files: %w", err)
	}
	defer rows.Close()

	result := make([]File, 0)

	for rows.Next() {
		var file File

		if err := rows.Scan(
			&file.ID,
			&file.OwnerID,
			&file.ObjectKey,
			&file.OriginalName,
			&file.ContentType,
			&file.SizeBytes,
			&file.ETag,
			&file.CreatedAt,
			&file.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan file: %w", err)
		}

		result = append(result, file)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate files: %w", err)
	}

	return result, nil
}

func (r *Repository) GetByIDAndOwner(
	ctx context.Context,
	id string,
	ownerID string,
) (*File, error) {
	const query = `
		SELECT
			id,
			owner_id,
			object_key,
			original_name,
			content_type,
			size_bytes,
			COALESCE(etag, ''),
			created_at,
			updated_at
		FROM files
		WHERE id = $1 AND owner_id = $2
	`

	var file File

	err := r.db.QueryRow(ctx, query, id, ownerID).Scan(
		&file.ID,
		&file.OwnerID,
		&file.ObjectKey,
		&file.OriginalName,
		&file.ContentType,
		&file.SizeBytes,
		&file.ETag,
		&file.CreatedAt,
		&file.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get file: %w", err)
	}

	return &file, nil
}

func (r *Repository) DeleteByIDAndOwner(
	ctx context.Context,
	id string,
	ownerID string,
) error {
	const query = `
		DELETE FROM files
		WHERE id = $1 AND owner_id = $2
	`

	tag, err := r.db.Exec(ctx, query, id, ownerID)
	if err != nil {
		return fmt.Errorf("delete file metadata: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
