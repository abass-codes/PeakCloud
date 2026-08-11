package sharing

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/abass-codes/peakcloud/internal/auth"
	"github.com/abass-codes/peakcloud/internal/storage"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service     *Service
	objectStore *storage.ObjectStore
}

func NewHandler(
	service *Service,
	objectStore *storage.ObjectStore,
) *Handler {
	return &Handler{
		service:     service,
		objectStore: objectStore,
	}
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNotFound),
		errors.Is(err, ErrRecipientNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

	case errors.Is(err, ErrForbidden):
		c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})

	case errors.Is(err, ErrAlreadyShared):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

	case errors.Is(err, ErrPasswordRequired):
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":             err.Error(),
			"password_required": true,
		})

	case errors.Is(err, ErrInvalidPassword):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

	case errors.Is(err, ErrExpired),
		errors.Is(err, ErrRevoked):
		c.JSON(http.StatusGone, gin.H{"error": err.Error()})

	case errors.Is(err, ErrDownloadDenied):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})

	case errors.Is(err, ErrInvalidPermission),
		errors.Is(err, ErrInvalidResource),
		errors.Is(err, ErrCannotShareSelf):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
	}
}

func currentUser(c *gin.Context) (string, bool) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return "", false
	}

	return user.ID, true
}

func (h *Handler) Create(c *gin.Context) {
	userID, ok := currentUser(c)
	if !ok {
		return
	}

	var input CreateShareRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	share, err := h.service.CreateShare(
		c.Request.Context(),
		userID,
		input,
	)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"share": share})
}

func (h *Handler) List(c *gin.Context) {
	userID, ok := currentUser(c)
	if !ok {
		return
	}

	shares, err := h.service.ListOwnedShares(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"shares": shares})
}

func (h *Handler) SharedWithMe(c *gin.Context) {
	userID, ok := currentUser(c)
	if !ok {
		return
	}

	shares, err := h.service.SharedWithMe(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"shares": shares})
}

func (h *Handler) Update(c *gin.Context) {
	userID, ok := currentUser(c)
	if !ok {
		return
	}

	var input UpdateShareRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.service.UpdateShare(
		c.Request.Context(),
		userID,
		c.Param("id"),
		input,
	); err != nil {
		writeError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) Delete(c *gin.Context) {
	userID, ok := currentUser(c)
	if !ok {
		return
	}

	if err := h.service.DeleteShare(
		c.Request.Context(),
		userID,
		c.Param("id"),
	); err != nil {
		writeError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) CreatePublicLink(c *gin.Context) {
	userID, ok := currentUser(c)
	if !ok {
		return
	}

	var input CreatePublicLinkRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	link, err := h.service.CreatePublicLink(
		c.Request.Context(),
		userID,
		input,
	)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"link": link})
}

func (h *Handler) ListPublicLinks(c *gin.Context) {
	userID, ok := currentUser(c)
	if !ok {
		return
	}

	links, err := h.service.ListPublicLinks(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"links": links})
}

func (h *Handler) RevokePublicLink(c *gin.Context) {
	userID, ok := currentUser(c)
	if !ok {
		return
	}

	if err := h.service.RevokePublicLink(
		c.Request.Context(),
		userID,
		c.Param("id"),
	); err != nil {
		writeError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) ResolvePublic(c *gin.Context) {
	var input PublicAccessRequest
	_ = c.ShouldBindJSON(&input)

	link, err := h.service.ResolvePublicLink(
		c.Request.Context(),
		c.Param("token"),
		input.Password,
	)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"link": link})
}

func (h *Handler) PublicContent(c *gin.Context) {
	password := c.Query("password")

	file, err := h.service.ResolvePublicFile(
		c.Request.Context(),
		c.Param("token"),
		password,
		false,
	)
	if err != nil {
		writeError(c, err)
		return
	}

	object, err := h.objectStore.Get(
		c.Request.Context(),
		file.ObjectKey,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to load shared file",
		})
		return
	}
	defer object.Close()

	c.Header("Content-Type", file.ContentType)
	c.Header("Content-Length", stringInt64(file.SizeBytes))
	c.Header(
		"Content-Disposition",
		`inline; filename*=UTF-8''`+url.PathEscape(file.OriginalName),
	)

	c.Status(http.StatusOK)

	_, _ = io.Copy(c.Writer, object)
}

func (h *Handler) PublicDownload(c *gin.Context) {
	password := c.Query("password")

	file, err := h.service.ResolvePublicFile(
		c.Request.Context(),
		c.Param("token"),
		password,
		true,
	)
	if err != nil {
		writeError(c, err)
		return
	}

	object, err := h.objectStore.Get(
		c.Request.Context(),
		file.ObjectKey,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to download shared file",
		})
		return
	}
	defer object.Close()

	c.Header("Content-Type", file.ContentType)
	c.Header("Content-Length", stringInt64(file.SizeBytes))
	c.Header(
		"Content-Disposition",
		`attachment; filename*=UTF-8''`+url.PathEscape(file.OriginalName),
	)

	c.Status(http.StatusOK)

	_, _ = io.Copy(c.Writer, object)
}

func stringInt64(value int64) string {
	return fmt.Sprintf("%d", value)
}
