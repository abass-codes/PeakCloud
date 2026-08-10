package files

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/abass-codes/peakcloud/internal/auth"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service       *Service
	maxUploadSize int64
}

func NewHandler(service *Service, maxUploadSize int64) *Handler {
	return &Handler{
		service:       service,
		maxUploadSize: maxUploadSize,
	}
}

type renameRequest struct {
	Name string `json:"name"`
}

type locationRequest struct {
	FolderID *string `json:"folder_id"`
}

func (h *Handler) Upload(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.Request.Body = http.MaxBytesReader(
		c.Writer,
		c.Request.Body,
		h.maxUploadSize+(1<<20),
	)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	if fileHeader.Size > h.maxUploadSize {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"error": "file exceeds upload limit",
		})
		return
	}

	source, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to read file"})
		return
	}
	defer source.Close()

	var folderID *string

	if value := c.PostForm("folder_id"); value != "" {
		folderID = &value
	}

	file, err := h.service.Upload(
		c.Request.Context(),
		user.ID,
		folderID,
		fileHeader.Filename,
		fileHeader.Header.Get("Content-Type"),
		source,
		fileHeader.Size,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrFileTooLarge):
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "file exceeds upload limit",
			})

		case errors.Is(err, ErrInvalidFilename):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid filename",
			})

		case errors.Is(err, ErrNotFound):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid destination folder",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "unable to upload file",
			})
		}

		return
	}

	c.JSON(http.StatusCreated, gin.H{"file": file})
}

func (h *Handler) List(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var folderID *string

	if value := c.Query("folder_id"); value != "" {
		folderID = &value
	}

	result, err := h.service.List(
		c.Request.Context(),
		user.ID,
		folderID,
	)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to list files",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": result})
}

func (h *Handler) Get(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	file, err := h.service.Get(
		c.Request.Context(),
		c.Param("id"),
		user.ID,
	)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to retrieve file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"file": file})
}

func (h *Handler) Download(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	file, object, err := h.service.Download(
		c.Request.Context(),
		c.Param("id"),
		user.ID,
	)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to download file",
		})
		return
	}
	defer object.Close()

	c.Header("Content-Type", file.ContentType)
	c.Header(
		"Content-Disposition",
		fmt.Sprintf(
			"attachment; filename*=UTF-8''%s",
			url.PathEscape(file.OriginalName),
		),
	)
	c.Header("Content-Length", fmt.Sprintf("%d", file.SizeBytes))
	c.Status(http.StatusOK)

	_, _ = io.Copy(c.Writer, object)
}

func (h *Handler) Rename(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request renameRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	file, err := h.service.Rename(
		c.Request.Context(),
		c.Param("id"),
		user.ID,
		request.Name,
	)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}

		if errors.Is(err, ErrInvalidFilename) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filename"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to rename file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"file": file})
}

func (h *Handler) Move(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request locationRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	file, err := h.service.Move(
		c.Request.Context(),
		c.Param("id"),
		user.ID,
		request.FolderID,
	)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "file or destination folder not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to move file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"file": file})
}

func (h *Handler) Copy(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request locationRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	file, err := h.service.Copy(
		c.Request.Context(),
		c.Param("id"),
		user.ID,
		request.FolderID,
	)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "file or destination folder not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to copy file",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"file": file})
}

func (h *Handler) Delete(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err := h.service.Delete(
		c.Request.Context(),
		c.Param("id"),
		user.ID,
	)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to delete file",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) Content(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	file, object, err := h.service.Download(
		c.Request.Context(),
		c.Param("id"),
		user.ID,
	)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to preview file",
		})
		return
	}
	defer object.Close()

	preview := ClassifyPreview(
		file.ContentType,
		file.OriginalName,
	)

	if !preview.Previewable {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{
			"error": "file type does not support preview",
		})
		return
	}

	if (preview.Kind == PreviewText || preview.Kind == PreviewCode) &&
		file.SizeBytes > MaxTextPreviewBytes {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"error": "text file is too large to preview",
		})
		return
	}

	c.Header("Content-Type", file.ContentType)
	c.Header(
		"Content-Disposition",
		`inline; filename="`+sanitizeContentDispositionFilename(file.OriginalName)+`"`,
	)
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Cache-Control", "private, no-store")

	if _, err := io.Copy(c.Writer, object); err != nil {
		return
	}
}

func sanitizeContentDispositionFilename(name string) string {
	name = strings.ReplaceAll(name, `\`, "_")
	name = strings.ReplaceAll(name, `"`, "_")
	name = strings.ReplaceAll(name, "\r", "_")
	name = strings.ReplaceAll(name, "\n", "_")

	return name
}

func (h *Handler) Preview(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	file, err := h.service.Get(
		c.Request.Context(),
		c.Param("id"),
		user.ID,
	)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to load file preview",
		})
		return
	}

	preview := ClassifyPreview(
		file.ContentType,
		file.OriginalName,
	)

	c.JSON(http.StatusOK, gin.H{
		"file":    file,
		"preview": preview,
	})
}
