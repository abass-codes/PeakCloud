package sharing

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindUserByEmail(
	ctx context.Context,
	email string,
) (string, string, error) {
	var id string
	var foundEmail string

	err := r.db.QueryRow(
		ctx,
		`SELECT id, email
		 FROM users
		 WHERE lower(email) = lower($1)`,
		email,
	).Scan(&id, &foundEmail)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", ErrRecipientNotFound
	}

	return id, foundEmail, err
}

func (r *Repository) FileOwnedBy(
	ctx context.Context,
	fileID string,
	ownerID string,
) (string, bool, error) {
	var name string

	err := r.db.QueryRow(
		ctx,
		`SELECT original_name
		 FROM files
		 WHERE id = $1
		   AND owner_id = $2`,
		fileID,
		ownerID,
	).Scan(&name)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", false, nil
	}

	if err != nil {
		return "", false, err
	}

	return name, true, nil
}

func (r *Repository) FolderOwnedBy(
	ctx context.Context,
	folderID string,
	ownerID string,
) (string, bool, error) {
	var name string

	err := r.db.QueryRow(
		ctx,
		`SELECT name
		 FROM folders
		 WHERE id = $1
		   AND owner_id = $2`,
		folderID,
		ownerID,
	).Scan(&name)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", false, nil
	}

	if err != nil {
		return "", false, err
	}

	return name, true, nil
}

func (r *Repository) GetPublicFile(
	ctx context.Context,
	fileID string,
	ownerID string,
) (*PublicFile, error) {
	var file PublicFile

	err := r.db.QueryRow(
		ctx,
		`SELECT
			id,
			object_key,
			original_name,
			content_type,
			size_bytes
		 FROM files
		 WHERE id = $1
		   AND owner_id = $2`,
		fileID,
		ownerID,
	).Scan(
		&file.ID,
		&file.ObjectKey,
		&file.OriginalName,
		&file.ContentType,
		&file.SizeBytes,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (r *Repository) CreateShare(
	ctx context.Context,
	share *Share,
) error {
	return r.db.QueryRow(
		ctx,
		`INSERT INTO resource_shares (
			id,
			owner_id,
			recipient_id,
			resource_type,
			resource_id,
			permission,
			allow_download
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING created_at, updated_at`,
		share.ID,
		share.OwnerID,
		share.RecipientID,
		share.ResourceType,
		share.ResourceID,
		share.Permission,
		share.AllowDownload,
	).Scan(&share.CreatedAt, &share.UpdatedAt)
}

func (r *Repository) ListOwnedShares(
	ctx context.Context,
	ownerID string,
) ([]Share, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT
			s.id,
			s.owner_id,
			s.recipient_id,
			u.email,
			s.resource_type,
			s.resource_id,
			CASE
				WHEN s.resource_type = 'file'
				THEN (
					SELECT f.original_name
					FROM files f
					WHERE f.id = s.resource_id
				)
				ELSE (
					SELECT fo.name
					FROM folders fo
					WHERE fo.id = s.resource_id
				)
			END,
			s.permission,
			s.allow_download,
			s.created_at,
			s.updated_at
		 FROM resource_shares s
		 JOIN users u ON u.id = s.recipient_id
		 WHERE s.owner_id = $1
		 ORDER BY s.created_at DESC`,
		ownerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	shares := make([]Share, 0)

	for rows.Next() {
		var share Share

		if err := rows.Scan(
			&share.ID,
			&share.OwnerID,
			&share.RecipientID,
			&share.RecipientEmail,
			&share.ResourceType,
			&share.ResourceID,
			&share.ResourceName,
			&share.Permission,
			&share.AllowDownload,
			&share.CreatedAt,
			&share.UpdatedAt,
		); err != nil {
			return nil, err
		}

		shares = append(shares, share)
	}

	return shares, rows.Err()
}

func (r *Repository) ListSharedWithMe(
	ctx context.Context,
	userID string,
) ([]Share, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT
			s.id,
			s.owner_id,
			s.recipient_id,
			s.resource_type,
			s.resource_id,
			CASE
				WHEN s.resource_type = 'file'
				THEN (
					SELECT f.original_name
					FROM files f
					WHERE f.id = s.resource_id
				)
				ELSE (
					SELECT fo.name
					FROM folders fo
					WHERE fo.id = s.resource_id
				)
			END,
			s.permission,
			s.allow_download,
			s.created_at,
			s.updated_at
		 FROM resource_shares s
		 WHERE s.recipient_id = $1
		 ORDER BY s.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	shares := make([]Share, 0)

	for rows.Next() {
		var share Share

		if err := rows.Scan(
			&share.ID,
			&share.OwnerID,
			&share.RecipientID,
			&share.ResourceType,
			&share.ResourceID,
			&share.ResourceName,
			&share.Permission,
			&share.AllowDownload,
			&share.CreatedAt,
			&share.UpdatedAt,
		); err != nil {
			return nil, err
		}

		shares = append(shares, share)
	}

	return shares, rows.Err()
}

func (r *Repository) UpdateShare(
	ctx context.Context,
	id string,
	ownerID string,
	permission Permission,
	allowDownload bool,
) error {
	tag, err := r.db.Exec(
		ctx,
		`UPDATE resource_shares
		 SET permission = $3,
		     allow_download = $4,
		     updated_at = NOW()
		 WHERE id = $1
		   AND owner_id = $2`,
		id,
		ownerID,
		permission,
		allowDownload,
	)

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *Repository) DeleteShare(
	ctx context.Context,
	id string,
	ownerID string,
) error {
	tag, err := r.db.Exec(
		ctx,
		`DELETE FROM resource_shares
		 WHERE id = $1
		   AND owner_id = $2`,
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

type publicLinkRecord struct {
	PublicLink
	TokenHash    string
	PasswordHash string
}

func (r *Repository) CreatePublicLink(
	ctx context.Context,
	link *publicLinkRecord,
) error {
	return r.db.QueryRow(
		ctx,
		`INSERT INTO public_share_links (
			id,
			owner_id,
			resource_type,
			resource_id,
			token_hash,
			password_hash,
			permission,
			allow_download,
			expires_at
		)
		VALUES ($1,$2,$3,$4,$5,NULLIF($6,''),$7,$8,$9)
		RETURNING created_at, updated_at`,
		link.ID,
		link.OwnerID,
		link.ResourceType,
		link.ResourceID,
		link.TokenHash,
		link.PasswordHash,
		link.Permission,
		link.AllowDownload,
		link.ExpiresAt,
	).Scan(&link.CreatedAt, &link.UpdatedAt)
}

func (r *Repository) FindPublicLink(
	ctx context.Context,
	tokenHash string,
) (*publicLinkRecord, error) {
	var link publicLinkRecord

	err := r.db.QueryRow(
		ctx,
		`SELECT
			id,
			owner_id,
			resource_type,
			resource_id,
			token_hash,
			COALESCE(password_hash, ''),
			permission,
			allow_download,
			expires_at,
			revoked_at,
			created_at,
			updated_at
		 FROM public_share_links
		 WHERE token_hash = $1`,
		tokenHash,
	).Scan(
		&link.ID,
		&link.OwnerID,
		&link.ResourceType,
		&link.ResourceID,
		&link.TokenHash,
		&link.PasswordHash,
		&link.Permission,
		&link.AllowDownload,
		&link.ExpiresAt,
		&link.RevokedAt,
		&link.CreatedAt,
		&link.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &link, nil
}

func (r *Repository) ListPublicLinks(
	ctx context.Context,
	ownerID string,
) ([]PublicLink, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT
			p.id,
			p.owner_id,
			p.resource_type,
			p.resource_id,
			CASE
				WHEN p.resource_type = 'file'
				THEN (
					SELECT f.original_name
					FROM files f
					WHERE f.id = p.resource_id
				)
				ELSE (
					SELECT fo.name
					FROM folders fo
					WHERE fo.id = p.resource_id
				)
			END,
			p.permission,
			p.allow_download,
			p.password_hash IS NOT NULL,
			p.expires_at,
			p.revoked_at,
			p.created_at,
			p.updated_at
		 FROM public_share_links p
		 WHERE p.owner_id = $1
		 ORDER BY p.created_at DESC`,
		ownerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := make([]PublicLink, 0)

	for rows.Next() {
		var link PublicLink

		if err := rows.Scan(
			&link.ID,
			&link.OwnerID,
			&link.ResourceType,
			&link.ResourceID,
			&link.ResourceName,
			&link.Permission,
			&link.AllowDownload,
			&link.PasswordSet,
			&link.ExpiresAt,
			&link.RevokedAt,
			&link.CreatedAt,
			&link.UpdatedAt,
		); err != nil {
			return nil, err
		}

		links = append(links, link)
	}

	return links, rows.Err()
}

func (r *Repository) RevokePublicLink(
	ctx context.Context,
	id string,
	ownerID string,
) error {
	now := time.Now().UTC()

	tag, err := r.db.Exec(
		ctx,
		`UPDATE public_share_links
		 SET revoked_at = $3,
		     updated_at = $3
		 WHERE id = $1
		   AND owner_id = $2
		   AND revoked_at IS NULL`,
		id,
		ownerID,
		now,
	)

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
