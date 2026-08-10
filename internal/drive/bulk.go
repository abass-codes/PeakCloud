package drive

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
)

type BulkItem struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

func (h *Handler) BulkMove(
	ctx context.Context,
	ownerID string,
	items []BulkItem,
	folderID *string,
) error {
	for _, item := range items {
		switch item.Type {
		case "file":
			if _, err := h.files.Move(
				ctx,
				item.ID,
				ownerID,
				folderID,
			); err != nil {
				return err
			}

		case "folder":
			if _, err := h.folders.Move(
				ctx,
				item.ID,
				ownerID,
				folderID,
			); err != nil {
				return err
			}

		default:
			return fmt.Errorf("unsupported item type")
		}
	}

	return nil
}

func (h *Handler) BulkDelete(
	ctx context.Context,
	ownerID string,
	items []BulkItem,
) error {
	for _, item := range items {
		switch item.Type {
		case "file":
			if err := h.files.Delete(
				ctx,
				item.ID,
				ownerID,
			); err != nil {
				return err
			}

		case "folder":
			if err := h.folders.Delete(
				ctx,
				item.ID,
				ownerID,
			); err != nil {
				return err
			}

		default:
			return fmt.Errorf("unsupported item type")
		}
	}

	return nil
}

func (h *Handler) WriteFilesZip(
	ctx context.Context,
	ownerID string,
	ids []string,
	writer io.Writer,
) error {
	archive := zip.NewWriter(writer)
	defer archive.Close()

	for _, id := range ids {
		file, object, err := h.files.Download(
			ctx,
			id,
			ownerID,
		)
		if err != nil {
			return err
		}

		entry, err := archive.Create(file.OriginalName)
		if err != nil {
			_ = object.Close()
			return err
		}

		if _, err := io.Copy(entry, object); err != nil {
			_ = object.Close()
			return err
		}

		if err := object.Close(); err != nil {
			return err
		}
	}

	return nil
}
