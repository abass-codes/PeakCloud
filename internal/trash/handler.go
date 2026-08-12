package trash

import (
	"errors"
	"net/http"

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

func (h *Handler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "unauthorized"},
		)
		return
	}

	contents, err := h.service.List(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(
		http.StatusOK,
		contents,
	)
}

func (h *Handler) Trash(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "unauthorized"},
		)
		return
	}

	resourceType := ResourceType(
		c.Param("type"),
	)
	resourceID := c.Param("id")

	if err := h.service.Trash(
		c.Request.Context(),
		userID,
		resourceType,
		resourceID,
	); err != nil {
		writeError(c, err)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "resource moved to trash",
		},
	)
}

func (h *Handler) Restore(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "unauthorized"},
		)
		return
	}

	resourceType := ResourceType(
		c.Param("type"),
	)
	resourceID := c.Param("id")

	if err := h.service.Restore(
		c.Request.Context(),
		userID,
		resourceType,
		resourceID,
	); err != nil {
		writeError(c, err)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "resource restored",
		},
	)
}

func (h *Handler) DeletePermanently(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "unauthorized"},
		)
		return
	}

	resourceType := ResourceType(
		c.Param("type"),
	)
	resourceID := c.Param("id")

	if err := h.service.DeletePermanently(
		c.Request.Context(),
		userID,
		resourceType,
		resourceID,
	); err != nil {
		writeError(c, err)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "resource permanently deleted",
		},
	)
}

func writeError(
	c *gin.Context,
	err error,
) {
	switch {
	case errors.Is(err, ErrInvalidType):
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)

	case errors.Is(err, ErrNotFound):
		c.JSON(
			http.StatusNotFound,
			gin.H{"error": err.Error()},
		)

	case errors.Is(err, ErrForbidden):
		c.JSON(
			http.StatusForbidden,
			gin.H{"error": err.Error()},
		)

	case errors.Is(err, ErrAlreadyDeleted):
		c.JSON(
			http.StatusConflict,
			gin.H{"error": err.Error()},
		)

	default:
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "internal server error"},
		)
	}
}
