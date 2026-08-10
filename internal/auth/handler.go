package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service    *Service
	cookieName string
	sessionTTL time.Duration
	secure     bool
}

type registerRequest struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewHandler(
	service *Service,
	cookieName string,
	sessionTTL time.Duration,
	secure bool,
) *Handler {
	return &Handler{
		service:    service,
		cookieName: cookieName,
		sessionTTL: sessionTTL,
		secure:     secure,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var request registerRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, token, err := h.service.Register(
		c.Request.Context(),
		request.Email,
		request.DisplayName,
		request.Password,
	)

	switch {
	case errors.Is(err, ErrInvalidInput):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "provide a valid email, display name, and password of at least 12 characters",
		})
		return

	case errors.Is(err, ErrEmailAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{
			"error": "an account with that email already exists",
		})
		return

	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to create account",
		})
		return
	}

	h.setSessionCookie(c, token)

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func (h *Handler) Login(c *gin.Context) {
	var request loginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, token, err := h.service.Login(
		c.Request.Context(),
		request.Email,
		request.Password,
	)

	if errors.Is(err, ErrInvalidCredentials) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to login",
		})
		return
	}

	h.setSessionCookie(c, token)

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *Handler) Logout(c *gin.Context) {
	token, _ := c.Cookie(h.cookieName)

	if err := h.service.Logout(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to logout",
		})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     h.cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.secure,
		SameSite: http.SameSiteLaxMode,
	})

	c.Status(http.StatusNoContent)
}

func (h *Handler) Me(c *gin.Context) {
	value, exists := c.Get("user")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, ok := value.(*User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user context"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *Handler) setSessionCookie(c *gin.Context, token string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     h.cookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(h.sessionTTL.Seconds()),
		HttpOnly: true,
		Secure:   h.secure,
		SameSite: http.SameSiteLaxMode,
	})
}
