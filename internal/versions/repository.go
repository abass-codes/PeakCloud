package versions

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

func (r *Repository) Create(
	ctx context.Context,
	version *Version,
) error {
	const query = `
		INSERT INTO file_versions (
			file_id,
			version_number,
			object_key,
			size_bytes,
			content_type,
			etag,
			created_by
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		version.FileID,
		version.VersionNumber,
		version.ObjectKey,
		version.SizeBytes,
		version.ContentType,
		nullString(version.ETag),
		version.CreatedBy,
	).Scan(
		&version.ID,
		&version.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create file version: %w", err)
	}

	return nil
}

func (r *Repository) GetByNumber(
	ctx context.Context,
	fileID string,
	versionNumber int,
) (*Version, error) {
	const query = `
		SELECT
			id,
			file_id,
			version_number,
			object_key,
			size_bytes,
			content_type,
			COALESCE(etag, ''),
			created_by,
			created_at
		FROM file_versions
		WHERE file_id = $1
		  AND version_number = $2
	`

	version := &Version{}

	err := r.db.QueryRow(
		ctx,
		query,
		fileID,
		versionNumber,
	).Scan(
		&version.ID,
		&version.FileID,
		&version.VersionNumber,
		&version.ObjectKey,
		&version.SizeBytes,
		&version.ContentType,
		&version.ETag,
		&version.CreatedBy,
		&version.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get file version: %w", err)
	}

	return version, nil
}

func (r *Repository) ListByFile(
	ctx context.Context,
	fileID string,
) ([]Version, error) {
	const query = `
		SELECT
			id,
			file_id,
			version_number,
			object_key,
			size_bytes,
			content_type,
			COALESCE(etag, ''),
			created_by,
			created_at
		FROM file_versions
		WHERE file_id = $1
		ORDER BY version_number DESC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		fileID,
	)
	if err != nil {
		return nil, fmt.Errorf("list file versions: %w", err)
	}
	defer rows.Close()

	result := make([]Version, 0)

	for rows.Next() {
		var version Version

		if err := rows.Scan(
			&version.ID,
			&version.FileID,
			&version.VersionNumber,
			&version.ObjectKey,
			&version.SizeBytes,
			&version.ContentType,
			&version.ETag,
			&version.CreatedBy,
			&version.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf(
				"scan file version: %w",
				err,
			)
		}

		result = append(result, version)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate file versions: %w",
			err,
		)
	}

	return result, nil
}

func (r *Repository) NextVersionNumber(
	ctx context.Context,
	fileID string,
) (int, error) {
	const query = `
		SELECT COALESCE(MAX(version_number), 0) + 1
		FROM file_versions
		WHERE file_id = $1
	`

	var next int

	if err := r.db.QueryRow(
		ctx,
		query,
		fileID,
	).Scan(&next); err != nil {
		return 0, fmt.Errorf(
			"get next version number: %w",
			err,
		)
	}

	return next, nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	id string,
) (*Version, error) {
	const query = `
		SELECT
			id,
			file_id,
			version_number,
			object_key,
			size_bytes,
			content_type,
			COALESCE(etag, ''),
			created_by,
			created_at
		FROM file_versions
		WHERE id = $1
	`

	version := &Version{}

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&version.ID,
		&version.FileID,
		&version.VersionNumber,
		&version.ObjectKey,
		&version.SizeBytes,
		&version.ContentType,
		&version.ETag,
		&version.CreatedBy,
		&version.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf(
			"get file version by id: %w",
			err,
		)
	}

	return version, nil
}

func (r *Repository) GetLatest(
	ctx context.Context,
	fileID string,
) (*Version, error) {
	const query = `
		SELECT
			id,
			file_id,
			version_number,
			object_key,
			size_bytes,
			content_type,
			COALESCE(etag, ''),
			created_by,
			created_at
		FROM file_versions
		WHERE file_id = $1
		ORDER BY version_number DESC
		LIMIT 1
	`

	version := &Version{}

	err := r.db.QueryRow(
		ctx,
		query,
		fileID,
	).Scan(
		&version.ID,
		&version.FileID,
		&version.VersionNumber,
		&version.ObjectKey,
		&version.SizeBytes,
		&version.ContentType,
		&version.ETag,
		&version.CreatedBy,
		&version.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf(
			"get latest file version: %w",
			err,
		)
	}

	return version, nil
}

func nullString(value string) any {
	if value == "" {
		return nil
	}

	return value
}
