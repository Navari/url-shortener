package cache

import "time"

// CacheInterface defines methods for cache operations
type CacheInterface interface {
	Set(key, value string, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	Ping() error
	Close() error
}
