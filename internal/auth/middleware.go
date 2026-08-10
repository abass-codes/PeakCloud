package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Middleware(service *Service, cookieName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(cookieName)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "unauthorized"},
			)
			return
		}

		user, err := service.Authenticate(
			c.Request.Context(),
			token,
		)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "unauthorized"},
			)
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func UserFromContext(c *gin.Context) (*User, bool) {
	value, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	user, ok := value.(*User)
	if !ok || user == nil {
		return nil, false
	}

	return user, true
}
