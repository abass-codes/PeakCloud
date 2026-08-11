package authorization

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

func (r *Repository) FileAccess(
	ctx context.Context,
	userID string,
	fileID string,
) (Access, error) {
	var ownerID string
	var folderID *string

	err := r.db.QueryRow(
		ctx,
		`SELECT owner_id, folder_id
		 FROM files
		 WHERE id = $1`,
		fileID,
	).Scan(&ownerID, &folderID)

	if errors.Is(err, pgx.ErrNoRows) {
		return Access{}, ErrResourceNotFound
	}

	if err != nil {
		return Access{}, err
	}

	if ownerID == userID {
		return ownerAccess(), nil
	}

	best := Access{}

	direct, err := r.directAccess(
		ctx,
		userID,
		ResourceFile,
		fileID,
	)
	if err != nil {
		return Access{}, err
	}

	best = strongerAccess(best, direct)

	if folderID != nil {
		inherited, err := r.folderTreeAccess(
			ctx,
			userID,
			*folderID,
		)
		if err != nil {
			return Access{}, err
		}

		best = strongerAccess(best, inherited)
	}

	return best, nil
}

func (r *Repository) FolderAccess(
	ctx context.Context,
	userID string,
	folderID string,
) (Access, error) {
	var ownerID string

	err := r.db.QueryRow(
		ctx,
		`SELECT owner_id
		 FROM folders
		 WHERE id = $1`,
		folderID,
	).Scan(&ownerID)

	if errors.Is(err, pgx.ErrNoRows) {
		return Access{}, ErrResourceNotFound
	}

	if err != nil {
		return Access{}, err
	}

	if ownerID == userID {
		return ownerAccess(), nil
	}

	return r.folderTreeAccess(
		ctx,
		userID,
		folderID,
	)
}

func (r *Repository) directAccess(
	ctx context.Context,
	userID string,
	resourceType ResourceType,
	resourceID string,
) (Access, error) {
	var permission Permission
	var allowDownload bool

	err := r.db.QueryRow(
		ctx,
		`SELECT permission, allow_download
		 FROM resource_shares
		 WHERE recipient_id = $1
		   AND resource_type = $2
		   AND resource_id = $3`,
		userID,
		resourceType,
		resourceID,
	).Scan(
		&permission,
		&allowDownload,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return Access{}, nil
	}

	if err != nil {
		return Access{}, err
	}

	return Access{
		Permission:    permission,
		AllowDownload: allowDownload,
	}, nil
}

func (r *Repository) folderTreeAccess(
	ctx context.Context,
	userID string,
	folderID string,
) (Access, error) {
	rows, err := r.db.Query(
		ctx,
		`
		WITH RECURSIVE ancestry AS (
			SELECT
				id,
				parent_id,
				0 AS depth
			FROM folders
			WHERE id = $2

			UNION ALL

			SELECT
				f.id,
				f.parent_id,
				a.depth + 1
			FROM folders f
			JOIN ancestry a
			  ON a.parent_id = f.id
		)
		SELECT
			s.permission,
			s.allow_download,
			a.depth
		FROM ancestry a
		JOIN resource_shares s
		  ON s.resource_type = 'folder'
		 AND s.resource_id = a.id
		 AND s.recipient_id = $1
		ORDER BY a.depth ASC
		`,
		userID,
		folderID,
	)
	if err != nil {
		return Access{}, err
	}
	defer rows.Close()

	best := Access{}

	for rows.Next() {
		var candidate Access
		var depth int

		if err := rows.Scan(
			&candidate.Permission,
			&candidate.AllowDownload,
			&depth,
		); err != nil {
			return Access{}, err
		}

		best = strongerAccess(best, candidate)
	}

	if err := rows.Err(); err != nil {
		return Access{}, err
	}

	return best, nil
}

func ownerAccess() Access {
	return Access{
		Owner:         true,
		Permission:    PermissionOwner,
		AllowDownload: true,
	}
}

func strongerAccess(current Access, candidate Access) Access {
	currentRank := permissionRank(current)
	candidateRank := permissionRank(candidate)

	if candidateRank > currentRank {
		return candidate
	}

	if candidateRank < currentRank {
		return current
	}

	if candidateRank == 0 {
		return current
	}

	current.AllowDownload =
		current.AllowDownload || candidate.AllowDownload

	return current
}

func permissionRank(access Access) int {
	if access.Owner || access.Permission == PermissionOwner {
		return 3
	}

	switch access.Permission {
	case PermissionEditor:
		return 2
	case PermissionViewer:
		return 1
	default:
		return 0
	}
}
