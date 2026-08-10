package http

import (
	"time"

	"github.com/abass-codes/peakcloud/internal/auth"
	"github.com/abass-codes/peakcloud/internal/drive"
	"github.com/abass-codes/peakcloud/internal/files"
	"github.com/abass-codes/peakcloud/internal/folders"
	"github.com/abass-codes/peakcloud/internal/health"
	"github.com/abass-codes/peakcloud/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	healthHandler *health.Handler,
	authHandler *auth.Handler,
	fileHandler *files.Handler,
	folderHandler *folders.Handler,
	driveHandler *drive.Handler,
	authService *auth.Service,
	cookieName string,
	webURL string,
) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{webURL},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/health", healthHandler.Health)

	authRateLimiter := middleware.NewRateLimiter(
		10,
		time.Minute,
	)

	api := router.Group("/api/v1")

	authRoutes := api.Group("/auth")
	authRoutes.Use(authRateLimiter.Middleware())
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}

	authRoutes.POST("/logout", authHandler.Logout)

	protected := api.Group("")
	protected.Use(auth.Middleware(authService, cookieName))
	{
		protected.GET("/me", authHandler.Me)

		protected.POST("/files", fileHandler.Upload)
		protected.GET("/files", fileHandler.List)
		protected.GET("/files/:id", fileHandler.Get)
		protected.GET("/files/:id/download", fileHandler.Download)
		protected.GET("/files/:id/preview", fileHandler.Preview)
		protected.GET("/files/:id/content", fileHandler.Content)
		protected.PATCH("/files/:id/name", fileHandler.Rename)
		protected.PATCH("/files/:id/location", fileHandler.Move)
		protected.POST("/files/:id/copy", fileHandler.Copy)
		protected.DELETE("/files/:id", fileHandler.Delete)

		protected.POST("/folders", folderHandler.Create)
		protected.GET("/folders", folderHandler.List)
		protected.GET("/folders/:id", folderHandler.Get)
		protected.GET("/folders/:id/breadcrumbs", folderHandler.Breadcrumbs)
		protected.PATCH("/folders/:id/name", folderHandler.Rename)
		protected.PATCH("/folders/:id/location", folderHandler.Move)
		protected.DELETE("/folders/:id", folderHandler.Delete)

		protected.GET("/drive", driveHandler.Contents)
		protected.POST("/drive/bulk/move", driveHandler.BulkMoveHandler)
		protected.POST("/drive/bulk/delete", driveHandler.BulkDeleteHandler)
		protected.POST("/drive/bulk/download", driveHandler.BulkDownload)
	}

	return router
}
