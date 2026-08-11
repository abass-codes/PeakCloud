package versions

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/abass-codes/peakcloud/internal/auth"
	"github.com/abass-codes/peakcloud/internal/files"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Create(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "unauthorized"},
		)
		return
	}

	fileID := c.Param("id")

	upload, err := c.FormFile("file")
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "file is required"},
		)
		return
	}

	reader, err := upload.Open()
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "open upload"},
		)
		return
	}
	defer reader.Close()

	version, err := h.service.Create(
		c.Request.Context(),
		user.ID,
		fileID,
		upload.Header.Get("Content-Type"),
		reader,
		upload.Size,
	)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{"version": version},
	)
}

func (h *Handler) List(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "unauthorized"},
		)
		return
	}

	result, err := h.service.List(
		c.Request.Context(),
		user.ID,
		c.Param("id"),
	)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"versions": result,
	})
}

func (h *Handler) Get(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "unauthorized"},
		)
		return
	}

	versionNumber, err := parseVersion(c)
	if err != nil {
		writeError(c, err)
		return
	}

	version, err := h.service.Get(
		c.Request.Context(),
		user.ID,
		c.Param("id"),
		versionNumber,
	)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"version": version,
	})
}

func (h *Handler) Content(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "unauthorized"},
		)
		return
	}

	versionNumber, err := parseVersion(c)
	if err != nil {
		writeError(c, err)
		return
	}

	version, object, err := h.service.Content(
		c.Request.Context(),
		user.ID,
		c.Param("id"),
		versionNumber,
	)
	if err != nil {
		writeError(c, err)
		return
	}
	defer object.Close()

	c.Header("Content-Type", version.ContentType)
	c.Header(
		"Content-Length",
		strconv.FormatInt(version.SizeBytes, 10),
	)

	c.DataFromReader(
		http.StatusOK,
		version.SizeBytes,
		version.ContentType,
		object,
		nil,
	)
}

func (h *Handler) Download(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "unauthorized"},
		)
		return
	}

	versionNumber, err := parseVersion(c)
	if err != nil {
		writeError(c, err)
		return
	}

	version, object, err := h.service.Download(
		c.Request.Context(),
		user.ID,
		c.Param("id"),
		versionNumber,
	)
	if err != nil {
		writeError(c, err)
		return
	}
	defer object.Close()

	filename := fmt.Sprintf(
		"version-%d",
		version.VersionNumber,
	)

	c.Header(
		"Content-Disposition",
		fmt.Sprintf(
			`attachment; filename*=UTF-8''%s`,
			url.QueryEscape(filename),
		),
	)

	c.DataFromReader(
		http.StatusOK,
		version.SizeBytes,
		version.ContentType,
		object,
		nil,
	)
}

func parseVersion(c *gin.Context) (int, error) {
	value, err := strconv.Atoi(c.Param("version"))
	if err != nil {
		return 0, ErrInvalidVersion
	}

	if err := ValidateVersionNumber(value); err != nil {
		return 0, err
	}

	return value, nil
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		c.JSON(
			http.StatusNotFound,
			gin.H{"error": "version not found"},
		)

	case errors.Is(err, ErrForbidden):
		c.JSON(
			http.StatusForbidden,
			gin.H{"error": "forbidden"},
		)

	case errors.Is(err, ErrInvalidVersion):
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "invalid version"},
		)

	case errors.Is(err, files.ErrFileTooLarge):
		c.JSON(
			http.StatusRequestEntityTooLarge,
			gin.H{"error": err.Error()},
		)

	default:
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "internal server error"},
		)
	}
}

func (h *Handler) Restore(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "unauthorized"},
		)
		return
	}

	versionNumber, err := parseVersion(c)
	if err != nil {
		writeError(c, err)
		return
	}

	version, err := h.service.Restore(
		c.Request.Context(),
		user.ID,
		c.Param("id"),
		versionNumber,
	)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{"version": version},
	)
}
