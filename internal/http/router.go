package http

import (
	"github.com/abass-codes/peakcloud/internal/health"
	"github.com/gin-gonic/gin"
)

func NewRouter(healthHandler *health.Handler) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/health", healthHandler.Health)

	return router
}
