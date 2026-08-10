package files

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/abass-codes/peakcloud/internal/auth"
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

	header, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	if header.Size > h.maxUploadSize {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"error": "file exceeds maximum upload size",
		})
		return
	}

	source, err := header.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to read uploaded file",
		})
		return
	}
	defer source.Close()

	contentType := header.Header.Get("Content-Type")

	file, err := h.service.Upload(
		c.Request.Context(),
		user.ID,
		header.Filename,
		contentType,
		source,
		header.Size,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidFilename):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filename"})
		case errors.Is(err, ErrFileTooLarge):
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "file exceeds maximum upload size",
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

	files, err := h.service.List(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to list files",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
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

	if _, err := io.Copy(c.Writer, object); err != nil {
		return
	}
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
