package drive

import (
	"errors"
	"net/http"

	"github.com/abass-codes/peakcloud/internal/auth"
	"github.com/abass-codes/peakcloud/internal/authorization"
	"github.com/abass-codes/peakcloud/internal/files"
	"github.com/abass-codes/peakcloud/internal/folders"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	files         *files.Service
	folders       *folders.Service
	authorization *authorization.Service
}

func NewHandler(
	fileService *files.Service,
	folderService *folders.Service,
	authorizationService *authorization.Service,
) *Handler {
	return &Handler{
		files:         fileService,
		folders:       folderService,
		authorization: authorizationService,
	}
}

func (h *Handler) Contents(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var folderID *string

	if value := c.Query("folder_id"); value != "" {
		folderID = &value
	}

	folderList, err := h.folders.List(
		c.Request.Context(),
		user.ID,
		folderID,
	)
	if err != nil {
		writeDriveError(c, err, "folder not found")
		return
	}

	fileList, err := h.files.List(
		c.Request.Context(),
		user.ID,
		folderID,
	)
	if err != nil {
		writeDriveError(c, err, "unable to list drive contents")
		return
	}

	breadcrumbs := []folders.Folder{}

	if folderID != nil {
		breadcrumbs, err = h.folders.Breadcrumbs(
			c.Request.Context(),
			*folderID,
			user.ID,
		)
		if err != nil {
			writeDriveError(c, err, "folder not found")
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"folder_id":   folderID,
		"breadcrumbs": breadcrumbs,
		"folders":     folderList,
		"files":       fileList,
	})
}

type bulkMoveRequest struct {
	Items    []BulkItem `json:"items"`
	FolderID *string    `json:"folder_id"`
}

type bulkDeleteRequest struct {
	Items []BulkItem `json:"items"`
}

type bulkDownloadRequest struct {
	FileIDs []string `json:"file_ids"`
}

func (h *Handler) BulkMoveHandler(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request bulkMoveRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if len(request.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no items selected"})
		return
	}

	if err := h.BulkMove(
		c.Request.Context(),
		user.ID,
		request.Items,
		request.FolderID,
	); err != nil {
		writeDriveError(c, err, "bulk move failed")
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) BulkDeleteHandler(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request bulkDeleteRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if len(request.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no items selected"})
		return
	}

	if err := h.BulkDelete(
		c.Request.Context(),
		user.ID,
		request.Items,
	); err != nil {
		writeDriveError(c, err, "bulk delete failed")
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) BulkDownload(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request bulkDownloadRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if len(request.FileIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no files selected"})
		return
	}

	// Do not write response headers until authorization succeeds.
	for _, id := range request.FileIDs {
		if _, err := h.authorization.Authorize(
			c.Request.Context(),
			user.ID,
			authorization.ResourceFile,
			id,
			authorization.ActionDownload,
		); err != nil {
			writeDriveError(c, err, "bulk download failed")
			return
		}
	}

	c.Header("Content-Type", "application/zip")
	c.Header(
		"Content-Disposition",
		`attachment; filename="peakcloud-files.zip"`,
	)

	if err := h.WriteFilesZip(
		c.Request.Context(),
		user.ID,
		request.FileIDs,
		c.Writer,
	); err != nil {
		return
	}
}

func writeDriveError(
	c *gin.Context,
	err error,
	fallback string,
) {
	switch {
	case errors.Is(err, authorization.ErrResourceNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"error": "resource not found",
		})

	case errors.Is(err, authorization.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{
			"error": "forbidden",
		})

	case errors.Is(err, files.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"error": "file not found",
		})

	case errors.Is(err, folders.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"error": "folder not found",
		})

	case errors.Is(err, ErrUnsupportedItemType):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "unsupported item type",
		})

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fallback,
		})
	}
}
