package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/shortener/internal/cache"
	"github.com/shortener/internal/config"
	"github.com/shortener/internal/model"
	"github.com/shortener/internal/repository"
	"gorm.io/gorm"
)

type URLService interface {
	CreateShortURL(req *model.CreateShortURLRequest) (*model.CreateShortURLResponse, error)
	GetOriginalURL(shortCode string) (string, error)
	GetURLStats(shortCode string) (*model.URLStatsResponse, error)
}

type urlService struct {
	repo   repository.ShortURLRepository
	cache  cache.CacheInterface
	config *config.Config
}

func NewURLService(repo repository.ShortURLRepository, cache cache.CacheInterface, cfg *config.Config) URLService {
	return &urlService{
		repo:   repo,
		cache:  cache,
		config: cfg,
	}
}

func (s *urlService) CreateShortURL(req *model.CreateShortURLRequest) (*model.CreateShortURLResponse, error) {
	shortCode, err := s.generateShortCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate short code: %w", err)
	}

	shortURL := &model.ShortURL{
		ShortCode:   shortCode,
		OriginalURL: req.URL,
		ExpiresAt:   req.ExpiresAt,
	}

	if err := s.repo.Create(shortURL); err != nil {
		return nil, fmt.Errorf("failed to create short URL: %w", err)
	}

	// Cache the URL
	cacheKey := fmt.Sprintf("short_url:%s", shortCode)
	expiration := time.Duration(s.config.App.CacheTTL) * time.Second

	if req.ExpiresAt != nil {
		timeUntilExpiry := time.Until(*req.ExpiresAt)
		if timeUntilExpiry < expiration {
			expiration = timeUntilExpiry
		}
	}

	s.cache.Set(cacheKey, req.URL, expiration)

	response := &model.CreateShortURLResponse{
		ShortCode:   shortCode,
		ShortURL:    fmt.Sprintf("%s/%s", s.config.App.BaseURL, shortCode),
		OriginalURL: req.URL,
		ExpiresAt:   req.ExpiresAt,
	}

	return response, nil
}

func (s *urlService) GetOriginalURL(shortCode string) (string, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("short_url:%s", shortCode)
	if cachedURL, err := s.cache.Get(cacheKey); err == nil {
		return cachedURL, nil
	}

	// Fallback to database
	shortURL, err := s.repo.FindByCode(shortCode)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("short URL not found")
		}
		return "", fmt.Errorf("failed to find short URL: %w", err)
	}

	// Check if expired
	if shortURL.ExpiresAt != nil && time.Now().After(*shortURL.ExpiresAt) {
		return "", fmt.Errorf("short URL has expired")
	}

	// Cache the result
	expiration := time.Duration(s.config.App.CacheTTL) * time.Second
	if shortURL.ExpiresAt != nil {
		timeUntilExpiry := time.Until(*shortURL.ExpiresAt)
		if timeUntilExpiry < expiration {
			expiration = timeUntilExpiry
		}
	}
	s.cache.Set(cacheKey, shortURL.OriginalURL, expiration)

	return shortURL.OriginalURL, nil
}

func (s *urlService) GetURLStats(shortCode string) (*model.URLStatsResponse, error) {
	shortURL, err := s.repo.FindByCode(shortCode)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("short URL not found")
		}
		return nil, fmt.Errorf("failed to find short URL: %w", err)
	}

	response := &model.URLStatsResponse{
		ShortCode:   shortURL.ShortCode,
		OriginalURL: shortURL.OriginalURL,
		CreatedAt:   shortURL.CreatedAt,
		ExpiresAt:   shortURL.ExpiresAt,
	}

	return response, nil
}

func (s *urlService) generateShortCode() (string, error) {
	const maxRetries = 5

	for i := 0; i < maxRetries; i++ {
		// Generate random bytes
		randomBytes := make([]byte, s.config.App.ShortCodeLength)
		if _, err := rand.Read(randomBytes); err != nil {
			return "", err
		}

		// Encode to base64 and clean up
		encoded := base64.URLEncoding.EncodeToString(randomBytes)
		shortCode := strings.ReplaceAll(encoded, "-", "")
		shortCode = strings.ReplaceAll(shortCode, "_", "")

		if len(shortCode) > s.config.App.ShortCodeLength {
			shortCode = shortCode[:s.config.App.ShortCodeLength]
		}

		// Check if already exists
		_, err := s.repo.FindByCode(shortCode)
		if err == gorm.ErrRecordNotFound {
			return shortCode, nil
		} else if err != nil {
			return "", err
		}
		// If exists, retry
	}

	return "", fmt.Errorf("failed to generate unique short code after %d retries", maxRetries)
}
