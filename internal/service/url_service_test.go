package service

import (
	"testing"
	"time"

	"github.com/shortener/internal/config"
	"github.com/shortener/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Mock Repository
type MockShortURLRepository struct {
	mock.Mock
}

func (m *MockShortURLRepository) Create(shortURL *model.ShortURL) error {
	args := m.Called(shortURL)
	return args.Error(0)
}

func (m *MockShortURLRepository) FindByCode(shortCode string) (*model.ShortURL, error) {
	args := m.Called(shortCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ShortURL), args.Error(1)
}

func (m *MockShortURLRepository) FindByID(id uint) (*model.ShortURL, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ShortURL), args.Error(1)
}

func (m *MockShortURLRepository) Delete(shortCode string) error {
	args := m.Called(shortCode)
	return args.Error(0)
}

// MockRedisClient implements cache.CacheInterface
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Set(key, value string, expiration time.Duration) error {
	args := m.Called(key, value, expiration)
	return args.Error(0)
}

func (m *MockRedisClient) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisClient) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockRedisClient) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRedisClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestCreateShortURL(t *testing.T) {
	// Setup
	mockRepo := new(MockShortURLRepository)
	mockCache := new(MockRedisClient)
	cfg := &config.Config{
		App: config.AppConfig{
			BaseURL:         "http://localhost:8080",
			CacheTTL:        3600,
			ShortCodeLength: 6,
		},
	}

	service := NewURLService(mockRepo, mockCache, cfg)

	// Test data
	req := &model.CreateShortURLRequest{
		URL: "https://example.com",
	}

	// Mock expectations
	mockRepo.On("Create", mock.AnythingOfType("*model.ShortURL")).Return(nil)
	mockRepo.On("FindByCode", mock.AnythingOfType("string")).Return(nil, gorm.ErrRecordNotFound)
	mockCache.On("Set", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil)

	// Execute
	response, err := service.CreateShortURL(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.ShortCode)
	assert.Equal(t, req.URL, response.OriginalURL)
	assert.Contains(t, response.ShortURL, response.ShortCode)

	// Verify mock calls
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestCreateShortURLWithoutExpiry(t *testing.T) {
	// Setup
	mockRepo := new(MockShortURLRepository)
	mockCache := new(MockRedisClient)
	cfg := &config.Config{
		App: config.AppConfig{
			BaseURL:         "http://localhost:8080",
			CacheTTL:        3600,
			ShortCodeLength: 6,
		},
	}

	service := NewURLService(mockRepo, mockCache, cfg)

	// Test data - expires_at parametresi YOK
	req := &model.CreateShortURLRequest{
		URL: "https://example.com/without-expiry",
	}

	// Mock expectations
	mockRepo.On("Create", mock.AnythingOfType("*model.ShortURL")).Return(nil)
	mockRepo.On("FindByCode", mock.AnythingOfType("string")).Return(nil, gorm.ErrRecordNotFound)
	mockCache.On("Set", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil)

	// Execute
	response, err := service.CreateShortURL(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.ShortCode)
	assert.Equal(t, req.URL, response.OriginalURL)
	assert.Contains(t, response.ShortURL, response.ShortCode)
	assert.Nil(t, response.ExpiresAt) // expires_at nil olmalı

	// Verify mock calls
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestGetOriginalURL_FromCache(t *testing.T) {
	// Setup
	mockRepo := new(MockShortURLRepository)
	mockCache := new(MockRedisClient)
	cfg := &config.Config{
		App: config.AppConfig{
			BaseURL:         "http://localhost:8080",
			CacheTTL:        3600,
			ShortCodeLength: 6,
		},
	}

	service := NewURLService(mockRepo, mockCache, cfg)

	// Test data
	shortCode := "abc123"
	originalURL := "https://example.com"
	cacheKey := "short_url:" + shortCode

	// Mock expectations
	mockCache.On("Get", cacheKey).Return(originalURL, nil)

	// Execute
	result, err := service.GetOriginalURL(shortCode)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, originalURL, result)

	// Verify mock calls
	mockCache.AssertExpectations(t)
}

func TestGetOriginalURL_FromDatabase(t *testing.T) {
	// Setup
	mockRepo := new(MockShortURLRepository)
	mockCache := new(MockRedisClient)
	cfg := &config.Config{
		App: config.AppConfig{
			BaseURL:         "http://localhost:8080",
			CacheTTL:        3600,
			ShortCodeLength: 6,
		},
	}

	service := NewURLService(mockRepo, mockCache, cfg)

	// Test data
	shortCode := "abc123"
	originalURL := "https://example.com"
	cacheKey := "short_url:" + shortCode

	shortURL := &model.ShortURL{
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		CreatedAt:   time.Now(),
	}

	// Mock expectations
	mockCache.On("Get", cacheKey).Return("", assert.AnError)
	mockRepo.On("FindByCode", shortCode).Return(shortURL, nil)
	mockCache.On("Set", cacheKey, originalURL, mock.AnythingOfType("time.Duration")).Return(nil)

	// Execute
	result, err := service.GetOriginalURL(shortCode)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, originalURL, result)

	// Verify mock calls
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestGetOriginalURL_Expired(t *testing.T) {
	// Setup
	mockRepo := new(MockShortURLRepository)
	mockCache := new(MockRedisClient)
	cfg := &config.Config{
		App: config.AppConfig{
			BaseURL:         "http://localhost:8080",
			CacheTTL:        3600,
			ShortCodeLength: 6,
		},
	}

	service := NewURLService(mockRepo, mockCache, cfg)

	// Test data
	shortCode := "abc123"
	originalURL := "https://example.com"
	cacheKey := "short_url:" + shortCode
	expiredTime := time.Now().Add(-time.Hour) // 1 hour ago

	shortURL := &model.ShortURL{
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		CreatedAt:   time.Now(),
		ExpiresAt:   &expiredTime,
	}

	// Mock expectations
	mockCache.On("Get", cacheKey).Return("", assert.AnError)
	mockRepo.On("FindByCode", shortCode).Return(shortURL, nil)

	// Execute
	result, err := service.GetOriginalURL(shortCode)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expired")
	assert.Empty(t, result)

	// Verify mock calls
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestGetURLStats(t *testing.T) {
	// Setup
	mockRepo := new(MockShortURLRepository)
	mockCache := new(MockRedisClient)
	cfg := &config.Config{
		App: config.AppConfig{
			BaseURL:         "http://localhost:8080",
			CacheTTL:        3600,
			ShortCodeLength: 6,
		},
	}

	service := NewURLService(mockRepo, mockCache, cfg)

	// Test data
	shortCode := "abc123"
	originalURL := "https://example.com"
	createdAt := time.Now()

	shortURL := &model.ShortURL{
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		CreatedAt:   createdAt,
	}

	// Mock expectations
	mockRepo.On("FindByCode", shortCode).Return(shortURL, nil)

	// Execute
	stats, err := service.GetURLStats(shortCode)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, shortCode, stats.ShortCode)
	assert.Equal(t, originalURL, stats.OriginalURL)
	assert.Equal(t, createdAt, stats.CreatedAt)

	// Verify mock calls
	mockRepo.AssertExpectations(t)
}
