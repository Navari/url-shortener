package repository

import (
	"time"

	"github.com/shortener/internal/model"
	"gorm.io/gorm"
)

type ShortURLRepository interface {
	Create(shortURL *model.ShortURL) error
	FindByCode(shortCode string) (*model.ShortURL, error)
	FindByID(id uint) (*model.ShortURL, error)
	Delete(shortCode string) error
}

type shortURLRepository struct {
	db *gorm.DB
}

func NewShortURLRepository(db *gorm.DB) ShortURLRepository {
	return &shortURLRepository{db: db}
}

func (r *shortURLRepository) Create(shortURL *model.ShortURL) error {
	return r.db.Create(shortURL).Error
}

func (r *shortURLRepository) FindByCode(shortCode string) (*model.ShortURL, error) {
	var shortURL model.ShortURL
	err := r.db.Where("short_code = ?", shortCode).First(&shortURL).Error
	if err != nil {
		return nil, err
	}
	return &shortURL, nil
}

func (r *shortURLRepository) FindByID(id uint) (*model.ShortURL, error) {
	var shortURL model.ShortURL
	err := r.db.First(&shortURL, id).Error
	if err != nil {
		return nil, err
	}
	return &shortURL, nil
}

func (r *shortURLRepository) Delete(shortCode string) error {
	return r.db.Where("short_code = ?", shortCode).Delete(&model.ShortURL{}).Error
}

func (r *shortURLRepository) IsExpired(shortURL *model.ShortURL) bool {
	if shortURL.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*shortURL.ExpiresAt)
}
