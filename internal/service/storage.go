package service

import (
	"github.com/moooll/cat-service-mongo/internal/repository"
	rediscache "github.com/moooll/cat-service-mongo/internal/repository/rediscache"
)

// Storage type is for working with the database from endpoints
type Storage struct {
	catalog *repository.MongoCatalog
	cache   *rediscache.Redis
}

// NewStorage creates new *Storage
func NewStorage(cat *repository.MongoCatalog, cache *rediscache.Redis) *Storage {
	return &Storage{
		cat,
		cache,
	}
}
