package drive

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/abass-codes/peakcloud/internal/authorization"
)

var ErrUnsupportedItemType = errors.New("unsupported item type")

type BulkItem struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

func (h *Handler) BulkMove(
	ctx context.Context,
	userID string,
	items []BulkItem,
	folderID *string,
) error {
	// Preflight the destination before mutating anything.
	if folderID != nil {
		if _, err := h.authorization.Authorize(
			ctx,
			userID,
			authorization.ResourceFolder,
			*folderID,
			authorization.ActionEdit,
		); err != nil {
			return err
		}
	}

	// Preflight every source item before moving anything.
	for _, item := range items {
		resourceType, err := resourceTypeForItem(item.Type)
		if err != nil {
			return err
		}

		if _, err := h.authorization.Authorize(
			ctx,
			userID,
			resourceType,
			item.ID,
			authorization.ActionEdit,
		); err != nil {
			return err
		}
	}

	// Authorization passed for the entire batch.
	for _, item := range items {
		switch item.Type {
		case "file":
			if _, err := h.files.Move(
				ctx,
				item.ID,
				userID,
				folderID,
			); err != nil {
				return err
			}

		case "folder":
			if _, err := h.folders.Move(
				ctx,
				item.ID,
				userID,
				folderID,
			); err != nil {
				return err
			}

		default:
			return ErrUnsupportedItemType
		}
	}

	return nil
}

func (h *Handler) BulkDelete(
	ctx context.Context,
	userID string,
	items []BulkItem,
) error {
	// Deletion is owner-only. Validate the entire batch first.
	for _, item := range items {
		resourceType, err := resourceTypeForItem(item.Type)
		if err != nil {
			return err
		}

		access, err := h.authorization.Access(
			ctx,
			userID,
			resourceType,
			item.ID,
		)
		if err != nil {
			return err
		}

		if !access.Owner {
			return authorization.ErrForbidden
		}
	}

	// Every item passed owner validation.
	for _, item := range items {
		switch item.Type {
		case "file":
			if err := h.files.Delete(
				ctx,
				item.ID,
				userID,
			); err != nil {
				return err
			}

		case "folder":
			if err := h.folders.Delete(
				ctx,
				item.ID,
				userID,
			); err != nil {
				return err
			}

		default:
			return ErrUnsupportedItemType
		}
	}

	return nil
}

func (h *Handler) WriteFilesZip(
	ctx context.Context,
	userID string,
	ids []string,
	writer io.Writer,
) error {
	// Preflight every file before writing any ZIP bytes.
	for _, id := range ids {
		if _, err := h.authorization.Authorize(
			ctx,
			userID,
			authorization.ResourceFile,
			id,
			authorization.ActionDownload,
		); err != nil {
			return err
		}
	}

	archive := zip.NewWriter(writer)

	for _, id := range ids {
		file, object, err := h.files.Download(
			ctx,
			id,
			userID,
		)
		if err != nil {
			_ = archive.Close()
			return err
		}

		entry, err := archive.Create(file.OriginalName)
		if err != nil {
			_ = object.Close()
			_ = archive.Close()
			return err
		}

		if _, err := io.Copy(entry, object); err != nil {
			_ = object.Close()
			_ = archive.Close()
			return err
		}

		if err := object.Close(); err != nil {
			_ = archive.Close()
			return err
		}
	}

	return archive.Close()
}

func resourceTypeForItem(
	itemType string,
) (authorization.ResourceType, error) {
	switch itemType {
	case "file":
		return authorization.ResourceFile, nil
	case "folder":
		return authorization.ResourceFolder, nil
	default:
		return "", fmt.Errorf(
			"%w: %s",
			ErrUnsupportedItemType,
			itemType,
		)
	}
}
