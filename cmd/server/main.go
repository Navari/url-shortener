package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shortener/internal/cache"
	"github.com/shortener/internal/config"
	"github.com/shortener/internal/handler"
	"github.com/shortener/internal/logger"
	"github.com/shortener/internal/repository"
	"github.com/shortener/internal/service"
	"go.uber.org/zap"
)

// @title URL Shortener API
// @version 1.0
// @description Bu bir URL kısaltma servisi API'sidir
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token authentication
func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	if err := logger.Init(cfg.Server.Env); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting URL Shortener Service", zap.String("env", cfg.Server.Env))

	// Initialize database
	db, err := repository.NewDatabase(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Initialize Redis
	redisClient := cache.NewRedisClient(cfg)
	if err := redisClient.Ping(); err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	logger.Info("Database and Redis connections established successfully")

	// Initialize repositories
	shortURLRepo := repository.NewShortURLRepository(db)

	// Initialize services
	urlService := service.NewURLService(shortURLRepo, redisClient, cfg)

	// Setup routes
	router := handler.SetupRoutes(urlService, cfg)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server starting", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	// Close database connection
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
	}

	// Close Redis connection
	redisClient.Close()

	logger.Info("Server exited")
}
