package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shortener/internal/config"
	"github.com/shortener/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(urlService service.URLService, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize handlers
	urlHandler := NewURLHandler(urlService)

	// Health check endpoint
	router.GET("/healthz", urlHandler.HealthCheck)

	// API routes
	api := router.Group("/api/v1")
	{
		// Protected route - requires auth token
		api.POST("/shorten", AuthMiddleware(cfg), urlHandler.CreateShortURL)
		// Public route - no auth required
		api.GET("/stats/:code", urlHandler.GetURLStats)
	}

	// Redirect route (short URL resolution)
	router.GET("/:code", urlHandler.RedirectToOriginalURL)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
