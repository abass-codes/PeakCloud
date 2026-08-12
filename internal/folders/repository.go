package folders

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, folder *Folder) error {
	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO folders (
			owner_id,
			parent_id,
			name
		)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
		`,
		folder.OwnerID,
		folder.ParentID,
		folder.Name,
	).Scan(
		&folder.ID,
		&folder.CreatedAt,
		&folder.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicateName
		}

		return err
	}

	return nil
}

func (r *Repository) GetByIDAndOwner(
	ctx context.Context,
	id string,
	ownerID string,
) (*Folder, error) {
	var folder Folder

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			owner_id,
			parent_id,
			name,
			created_at,
			updated_at
		FROM folders
		WHERE id = $1
		  AND owner_id = $2
              AND deleted_at IS NULL
		`,
		id,
		ownerID,
	).Scan(
		&folder.ID,
		&folder.OwnerID,
		&folder.ParentID,
		&folder.Name,
		&folder.CreatedAt,
		&folder.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &folder, nil
}

func (r *Repository) ListChildren(
	ctx context.Context,
	ownerID string,
	parentID *string,
) ([]Folder, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id,
			owner_id,
			parent_id,
			name,
			created_at,
			updated_at
		FROM folders
		WHERE owner_id = $1
		  AND parent_id IS NOT DISTINCT FROM $2::uuid
                  AND deleted_at IS NULL
		ORDER BY lower(name), created_at
		`,
		ownerID,
		parentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]Folder, 0)

	for rows.Next() {
		var folder Folder

		if err := rows.Scan(
			&folder.ID,
			&folder.OwnerID,
			&folder.ParentID,
			&folder.Name,
			&folder.CreatedAt,
			&folder.UpdatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, folder)
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
) (*Folder, error) {
	var folder Folder

	err := r.db.QueryRow(
		ctx,
		`
		UPDATE folders
		SET
			name = $3,
			updated_at = NOW()
		WHERE id = $1
		  AND owner_id = $2
              AND deleted_at IS NULL
		RETURNING
			id,
			owner_id,
			parent_id,
			name,
			created_at,
			updated_at
		`,
		id,
		ownerID,
		name,
	).Scan(
		&folder.ID,
		&folder.OwnerID,
		&folder.ParentID,
		&folder.Name,
		&folder.CreatedAt,
		&folder.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrDuplicateName
		}

		return nil, err
	}

	return &folder, nil
}

func (r *Repository) Move(
	ctx context.Context,
	id string,
	ownerID string,
	parentID *string,
) (*Folder, error) {
	var folder Folder

	err := r.db.QueryRow(
		ctx,
		`
		UPDATE folders
		SET
			parent_id = $3,
			updated_at = NOW()
		WHERE id = $1
		  AND owner_id = $2
              AND deleted_at IS NULL
		RETURNING
			id,
			owner_id,
			parent_id,
			name,
			created_at,
			updated_at
		`,
		id,
		ownerID,
		parentID,
	).Scan(
		&folder.ID,
		&folder.OwnerID,
		&folder.ParentID,
		&folder.Name,
		&folder.CreatedAt,
		&folder.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrDuplicateName
		}

		return nil, err
	}

	return &folder, nil
}

func (r *Repository) IsDescendant(
	ctx context.Context,
	ownerID string,
	ancestorID string,
	candidateID string,
) (bool, error) {
	var exists bool

	err := r.db.QueryRow(
		ctx,
		`
		WITH RECURSIVE descendants AS (
			SELECT id
			FROM folders
			WHERE parent_id = $1
			  AND owner_id = $2
                      AND deleted_at IS NULL

			UNION ALL

			SELECT f.id
			FROM folders f
			JOIN descendants d
			  ON f.parent_id = d.id
			WHERE f.owner_id = $2
                      AND f.deleted_at IS NULL
		)
		SELECT EXISTS (
			SELECT 1
			FROM descendants
			WHERE id = $3
		)
		`,
		ancestorID,
		ownerID,
		candidateID,
	).Scan(&exists)

	return exists, err
}

func (r *Repository) Breadcrumbs(
	ctx context.Context,
	id string,
	ownerID string,
) ([]Folder, error) {
	rows, err := r.db.Query(
		ctx,
		`
		WITH RECURSIVE path AS (
			SELECT
				id,
				owner_id,
				parent_id,
				name,
				created_at,
				updated_at,
				0 AS depth
			FROM folders
			WHERE id = $1
			  AND owner_id = $2
                      AND deleted_at IS NULL

			UNION ALL

			SELECT
				f.id,
				f.owner_id,
				f.parent_id,
				f.name,
				f.created_at,
				f.updated_at,
				p.depth + 1
			FROM folders f
			JOIN path p
			  ON p.parent_id = f.id
			WHERE f.owner_id = $2
                      AND f.deleted_at IS NULL
		)
		SELECT
			id,
			owner_id,
			parent_id,
			name,
			created_at,
			updated_at
		FROM path
		ORDER BY depth DESC
		`,
		id,
		ownerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]Folder, 0)

	for rows.Next() {
		var folder Folder

		if err := rows.Scan(
			&folder.ID,
			&folder.OwnerID,
			&folder.ParentID,
			&folder.Name,
			&folder.CreatedAt,
			&folder.UpdatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, folder)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, ErrNotFound
	}

	return result, nil
}

func (r *Repository) Delete(
	ctx context.Context,
	id string,
	ownerID string,
) error {
	tag, err := r.db.Exec(
		ctx,
		`
		DELETE FROM folders
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

func (r *Repository) GetByID(
	ctx context.Context,
	id string,
) (*Folder, error) {
	var folder Folder

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			owner_id,
			parent_id,
			name,
			created_at,
			updated_at
		FROM folders
		WHERE id = $1
                  AND deleted_at IS NULL
		`,
		id,
	).Scan(
		&folder.ID,
		&folder.OwnerID,
		&folder.ParentID,
		&folder.Name,
		&folder.CreatedAt,
		&folder.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &folder, nil
}

func (r *Repository) ListChildrenByFolder(
	ctx context.Context,
	parentID string,
) ([]Folder, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id,
			owner_id,
			parent_id,
			name,
			created_at,
			updated_at
		FROM folders
		WHERE parent_id = $1
                  AND deleted_at IS NULL
		ORDER BY lower(name), created_at
		`,
		parentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	folders := make([]Folder, 0)

	for rows.Next() {
		var folder Folder

		if err := rows.Scan(
			&folder.ID,
			&folder.OwnerID,
			&folder.ParentID,
			&folder.Name,
			&folder.CreatedAt,
			&folder.UpdatedAt,
		); err != nil {
			return nil, err
		}

		folders = append(folders, folder)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return folders, nil
}

func (r *Repository) RenameByID(
	ctx context.Context,
	id string,
	name string,
) (*Folder, error) {
	var folder Folder

	err := r.db.QueryRow(
		ctx,
		`
		UPDATE folders
		SET
			name = $2,
			updated_at = NOW()
		WHERE id = $1
                  AND deleted_at IS NULL
		RETURNING
			id,
			owner_id,
			parent_id,
			name,
			created_at,
			updated_at
		`,
		id,
		name,
	).Scan(
		&folder.ID,
		&folder.OwnerID,
		&folder.ParentID,
		&folder.Name,
		&folder.CreatedAt,
		&folder.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrDuplicateName
		}

		return nil, err
	}

	return &folder, nil
}

func (r *Repository) MoveByID(
	ctx context.Context,
	id string,
	parentID *string,
) (*Folder, error) {
	var folder Folder

	err := r.db.QueryRow(
		ctx,
		`
		UPDATE folders
		SET
			parent_id = $2,
			updated_at = NOW()
		WHERE id = $1
                  AND deleted_at IS NULL
		RETURNING
			id,
			owner_id,
			parent_id,
			name,
			created_at,
			updated_at
		`,
		id,
		parentID,
	).Scan(
		&folder.ID,
		&folder.OwnerID,
		&folder.ParentID,
		&folder.Name,
		&folder.CreatedAt,
		&folder.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrDuplicateName
		}

		return nil, err
	}

	return &folder, nil
}

func (r *Repository) Trash(
	ctx context.Context,
	id string,
	ownerID string,
) (*Folder, error) {
	var folder Folder

	err := r.db.QueryRow(
		ctx,
		`
		UPDATE folders
		SET
			deleted_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
		  AND owner_id = $2
		  AND deleted_at IS NULL
		RETURNING
			id,
			owner_id,
			parent_id,
			name,
			created_at,
			updated_at,
			deleted_at
		`,
		id,
		ownerID,
	).Scan(
		&folder.ID,
		&folder.OwnerID,
		&folder.ParentID,
		&folder.Name,
		&folder.CreatedAt,
		&folder.UpdatedAt,
		&folder.DeletedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &folder, nil
}

func (r *Repository) RestoreFromTrash(
	ctx context.Context,
	id string,
	ownerID string,
) (*Folder, error) {
	var folder Folder

	err := r.db.QueryRow(
		ctx,
		`
		UPDATE folders
		SET
			deleted_at = NULL,
			updated_at = NOW()
		WHERE id = $1
		  AND owner_id = $2
		  AND deleted_at IS NOT NULL
		RETURNING
			id,
			owner_id,
			parent_id,
			name,
			created_at,
			updated_at,
			deleted_at
		`,
		id,
		ownerID,
	).Scan(
		&folder.ID,
		&folder.OwnerID,
		&folder.ParentID,
		&folder.Name,
		&folder.CreatedAt,
		&folder.UpdatedAt,
		&folder.DeletedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &folder, nil
}

func (r *Repository) DeletePermanently(
	ctx context.Context,
	id string,
	ownerID string,
) (*Folder, error) {
	var folder Folder

	err := r.db.QueryRow(
		ctx,
		`
		DELETE FROM folders
		WHERE id = $1
		  AND owner_id = $2
		  AND deleted_at IS NOT NULL
		RETURNING
			id,
			owner_id,
			parent_id,
			name,
			created_at,
			updated_at,
			deleted_at
		`,
		id,
		ownerID,
	).Scan(
		&folder.ID,
		&folder.OwnerID,
		&folder.ParentID,
		&folder.Name,
		&folder.CreatedAt,
		&folder.UpdatedAt,
		&folder.DeletedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &folder, nil
}

func (r *Repository) ListTrashed(
	ctx context.Context,
	ownerID string,
) ([]Folder, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id,
			owner_id,
			parent_id,
			name,
			created_at,
			updated_at,
			deleted_at
		FROM folders
		WHERE owner_id = $1
		  AND deleted_at IS NOT NULL
		ORDER BY deleted_at DESC
		`,
		ownerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]Folder, 0)

	for rows.Next() {
		var folder Folder

		if err := rows.Scan(
			&folder.ID,
			&folder.OwnerID,
			&folder.ParentID,
			&folder.Name,
			&folder.CreatedAt,
			&folder.UpdatedAt,
			&folder.DeletedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, folder)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
