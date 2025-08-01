package model

import (
	"time"

	"gorm.io/gorm"
)

type ShortURL struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ShortCode   string         `gorm:"uniqueIndex;size:10;not null" json:"short_code"`
	OriginalURL string         `gorm:"not null" json:"original_url"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	ExpiresAt   *time.Time     `json:"expires_at,omitempty"`
}

type CreateShortURLRequest struct {
	URL       string     `json:"url" binding:"required,url"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type CreateShortURLResponse struct {
	ShortCode   string     `json:"short_code"`
	ShortURL    string     `json:"short_url"`
	OriginalURL string     `json:"original_url"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type URLStatsResponse struct {
	ShortCode   string     `json:"short_code"`
	OriginalURL string     `json:"original_url"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}
