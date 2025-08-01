package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shortener/internal/config"
	"github.com/shortener/internal/logger"
	"go.uber.org/zap"
)

// AuthMiddleware validates the Bearer token for protected endpoints
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header gereklidir"})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz Authorization header formatı"})
			c.Abort()
			return
		}

		token := parts[1]
		if token != cfg.Auth.Token {
			logger.Warn("Invalid token provided", zap.String("provided_token", token))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz token"})
			c.Abort()
			return
		}

		logger.Debug("Token validated successfully")
		c.Next()
	}
}
