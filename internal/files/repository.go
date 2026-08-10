package files

import (
	"context"
	"errors"

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
	return r.db.QueryRow(
		ctx,
		`
		INSERT INTO files (
			owner_id,
			folder_id,
			object_key,
			original_name,
			content_type,
			size_bytes,
			etag
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
		`,
		file.OwnerID,
		file.FolderID,
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
}

func (r *Repository) GetByIDAndOwner(
	ctx context.Context,
	id string,
	ownerID string,
) (*File, error) {
	var file File

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			owner_id,
			folder_id,
			object_key,
			original_name,
			content_type,
			size_bytes,
			etag,
			created_at,
			updated_at
		FROM files
		WHERE id = $1
		  AND owner_id = $2
		`,
		id,
		ownerID,
	).Scan(
		&file.ID,
		&file.OwnerID,
		&file.FolderID,
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
		return nil, err
	}

	return &file, nil
}

func (r *Repository) ListByFolder(
	ctx context.Context,
	ownerID string,
	folderID *string,
) ([]File, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id,
			owner_id,
			folder_id,
			object_key,
			original_name,
			content_type,
			size_bytes,
			etag,
			created_at,
			updated_at
		FROM files
		WHERE owner_id = $1
		  AND folder_id IS NOT DISTINCT FROM $2::uuid
		ORDER BY lower(original_name), created_at DESC
		`,
		ownerID,
		folderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]File, 0)

	for rows.Next() {
		var file File

		if err := rows.Scan(
			&file.ID,
			&file.OwnerID,
			&file.FolderID,
			&file.ObjectKey,
			&file.OriginalName,
			&file.ContentType,
			&file.SizeBytes,
			&file.ETag,
			&file.CreatedAt,
			&file.UpdatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, file)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repository) Rename(
	ctx context.Context,
	id string,
	ownerID string,
	name string,
) (*File, error) {
	var file File

	err := r.db.QueryRow(
		ctx,
		`
		UPDATE files
		SET
			original_name = $3,
			updated_at = NOW()
		WHERE id = $1
		  AND owner_id = $2
		RETURNING
			id,
			owner_id,
			folder_id,
			object_key,
			original_name,
			content_type,
			size_bytes,
			etag,
			created_at,
			updated_at
		`,
		id,
		ownerID,
		name,
	).Scan(
		&file.ID,
		&file.OwnerID,
		&file.FolderID,
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
		return nil, err
	}

	return &file, nil
}

func (r *Repository) Move(
	ctx context.Context,
	id string,
	ownerID string,
	folderID *string,
) (*File, error) {
	var file File

	err := r.db.QueryRow(
		ctx,
		`
		UPDATE files
		SET
			folder_id = $3,
			updated_at = NOW()
		WHERE id = $1
		  AND owner_id = $2
		RETURNING
			id,
			owner_id,
			folder_id,
			object_key,
			original_name,
			content_type,
			size_bytes,
			etag,
			created_at,
			updated_at
		`,
		id,
		ownerID,
		folderID,
	).Scan(
		&file.ID,
		&file.OwnerID,
		&file.FolderID,
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
		return nil, err
	}

	return &file, nil
}

func (r *Repository) DeleteMetadata(
	ctx context.Context,
	id string,
	ownerID string,
) error {
	tag, err := r.db.Exec(
		ctx,
		`
		DELETE FROM files
		WHERE id = $1
		  AND owner_id = $2
		`,
		id,
		ownerID,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *Repository) ListInFolderTree(
	ctx context.Context,
	ownerID string,
	folderID string,
) ([]File, error) {
	rows, err := r.db.Query(
		ctx,
		`
		WITH RECURSIVE tree AS (
			SELECT id
			FROM folders
			WHERE id = $1
			  AND owner_id = $2

			UNION ALL

			SELECT f.id
			FROM folders f
			JOIN tree t
			  ON f.parent_id = t.id
			WHERE f.owner_id = $2
		)
		SELECT
			fi.id,
			fi.owner_id,
			fi.folder_id,
			fi.object_key,
			fi.original_name,
			fi.content_type,
			fi.size_bytes,
			fi.etag,
			fi.created_at,
			fi.updated_at
		FROM files fi
		JOIN tree t
		  ON fi.folder_id = t.id
		WHERE fi.owner_id = $2
		`,
		folderID,
		ownerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]File, 0)

	for rows.Next() {
		var file File

		if err := rows.Scan(
			&file.ID,
			&file.OwnerID,
			&file.FolderID,
			&file.ObjectKey,
			&file.OriginalName,
			&file.ContentType,
			&file.SizeBytes,
			&file.ETag,
			&file.CreatedAt,
			&file.UpdatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, file)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
