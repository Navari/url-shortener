package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shortener/internal/logger"
	"github.com/shortener/internal/model"
	"github.com/shortener/internal/service"
	"go.uber.org/zap"
)

type URLHandler struct {
	urlService service.URLService
}

func NewURLHandler(urlService service.URLService) *URLHandler {
	return &URLHandler{
		urlService: urlService,
	}
}

// CreateShortURL creates a new short URL
// @Summary Create short URL
// @Description Create a new short URL from original URL (requires Bearer token)
// @Tags urls
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param request body model.CreateShortURLRequest true "URL creation request (expires_at optional)"
// @Success 201 {object} model.CreateShortURLResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/shorten [post]
func (h *URLHandler) CreateShortURL(c *gin.Context) {
	var req model.CreateShortURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz istek formatı"})
		return
	}

	response, err := h.urlService.CreateShortURL(&req)
	if err != nil {
		logger.Error("Failed to create short URL", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kısa URL oluşturulamadı"})
		return
	}

	logger.Info("Short URL created successfully", zap.String("short_code", response.ShortCode))
	c.JSON(http.StatusCreated, response)
}

// RedirectToOriginalURL redirects to the original URL
// @Summary Redirect to original URL
// @Description Redirect to the original URL using short code
// @Tags urls
// @Param code path string true "Short code"
// @Success 302 "Redirect to original URL"
// @Failure 404 {object} model.ErrorResponse
// @Failure 410 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /{code} [get]
func (h *URLHandler) RedirectToOriginalURL(c *gin.Context) {
	shortCode := c.Param("code")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kısa kod gereklidir"})
		return
	}

	originalURL, err := h.urlService.GetOriginalURL(shortCode)
	if err != nil {
		if err.Error() == "short URL not found" {
			logger.Warn("Short URL not found", zap.String("short_code", shortCode))
			c.JSON(http.StatusNotFound, gin.H{"error": "Kısa URL bulunamadı"})
			return
		}
		if err.Error() == "short URL has expired" {
			logger.Warn("Short URL expired", zap.String("short_code", shortCode))
			c.JSON(http.StatusGone, gin.H{"error": "Kısa URL'in süresi dolmuş"})
			return
		}
		logger.Error("Failed to get original URL", zap.Error(err), zap.String("short_code", shortCode))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sunucu hatası"})
		return
	}

	logger.Info("Redirecting to original URL", zap.String("short_code", shortCode), zap.String("original_url", originalURL))
	c.Redirect(http.StatusFound, originalURL)
}

// GetURLStats returns statistics for a short URL
// @Summary Get URL statistics
// @Description Get statistics for a short URL
// @Tags urls
// @Param code path string true "Short code"
// @Success 200 {object} model.URLStatsResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/stats/{code} [get]
func (h *URLHandler) GetURLStats(c *gin.Context) {
	shortCode := c.Param("code")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kısa kod gereklidir"})
		return
	}

	stats, err := h.urlService.GetURLStats(shortCode)
	if err != nil {
		if err.Error() == "short URL not found" {
			logger.Warn("Short URL not found for stats", zap.String("short_code", shortCode))
			c.JSON(http.StatusNotFound, gin.H{"error": "Kısa URL bulunamadı"})
			return
		}
		logger.Error("Failed to get URL stats", zap.Error(err), zap.String("short_code", shortCode))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "İstatistikler alınamadı"})
		return
	}

	logger.Info("URL stats retrieved", zap.String("short_code", shortCode))
	c.JSON(http.StatusOK, stats)
}

// HealthCheck returns service health status
// @Summary Health check
// @Description Returns the health status of the service
// @Tags health
// @Produce json
// @Success 200 {object} model.HealthResponse
// @Router /healthz [get]
func (h *URLHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "url-shortener",
		"version": "1.0.0",
	})
}
