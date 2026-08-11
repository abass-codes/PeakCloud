package folders

import (
	"errors"
	"net/http"

	"github.com/abass-codes/peakcloud/internal/auth"
	"github.com/abass-codes/peakcloud/internal/authorization"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type createRequest struct {
	Name     string  `json:"name"`
	ParentID *string `json:"parent_id"`
}

type renameRequest struct {
	Name string `json:"name"`
}

type moveRequest struct {
	ParentID *string `json:"parent_id"`
}

func (h *Handler) Create(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request createRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	folder, err := h.service.Create(
		c.Request.Context(),
		user.ID,
		request.ParentID,
		request.Name,
	)
	if err != nil {
		if writeAuthorizationError(
			c,
			err,
			"folder not found",
		) {
			return
		}

		writeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"folder": folder})
}

func (h *Handler) Get(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	folder, err := h.service.Get(
		c.Request.Context(),
		c.Param("id"),
		user.ID,
	)
	if err != nil {
		if writeAuthorizationError(
			c,
			err,
			"folder not found",
		) {
			return
		}

		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"folder": folder})
}

func (h *Handler) List(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var parentID *string

	if value := c.Query("parent_id"); value != "" {
		parentID = &value
	}

	result, err := h.service.List(
		c.Request.Context(),
		user.ID,
		parentID,
	)
	if err != nil {
		if writeAuthorizationError(
			c,
			err,
			"folder not found",
		) {
			return
		}

		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"folders": result})
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

	folder, err := h.service.Rename(
		c.Request.Context(),
		c.Param("id"),
		user.ID,
		request.Name,
	)
	if err != nil {
		if writeAuthorizationError(
			c,
			err,
			"folder not found",
		) {
			return
		}

		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"folder": folder})
}

func (h *Handler) Move(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request moveRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	folder, err := h.service.Move(
		c.Request.Context(),
		c.Param("id"),
		user.ID,
		request.ParentID,
	)
	if err != nil {
		if writeAuthorizationError(
			c,
			err,
			"folder not found",
		) {
			return
		}

		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"folder": folder})
}

func (h *Handler) Breadcrumbs(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	result, err := h.service.Breadcrumbs(
		c.Request.Context(),
		c.Param("id"),
		user.ID,
	)
	if err != nil {
		if writeAuthorizationError(
			c,
			err,
			"folder not found",
		) {
			return
		}

		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"breadcrumbs": result})
}

func (h *Handler) Delete(c *gin.Context) {
	user, ok := auth.UserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.Delete(
		c.Request.Context(),
		c.Param("id"),
		user.ID,
	); err != nil {
		writeError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})

	case errors.Is(err, ErrInvalidName):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder name"})

	case errors.Is(err, ErrInvalidParent):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parent folder"})

	case errors.Is(err, ErrDuplicateName):
		c.JSON(http.StatusConflict, gin.H{"error": "folder already exists"})

	case errors.Is(err, ErrInvalidMove):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder move"})

	case errors.Is(err, authorization.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})

	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "folder operation failed",
		})
	}
}

func writeAuthorizationError(
	c *gin.Context,
	err error,
	notFoundMessage string,
) bool {
	switch {
	case errors.Is(err, ErrNotFound),
		errors.Is(err, authorization.ErrResourceNotFound):
		c.JSON(
			http.StatusNotFound,
			gin.H{"error": notFoundMessage},
		)
		return true

	case errors.Is(err, authorization.ErrForbidden):
		c.JSON(
			http.StatusForbidden,
			gin.H{"error": "forbidden"},
		)
		return true

	default:
		return false
	}
}
