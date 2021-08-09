package handler

import (
	"github.com/moooll/cat-service-mongo/internal/repository"
	rediscache "github.com/moooll/cat-service-mongo/internal/repository/rediscache"
)

// Service type is for working with the database from endpoints
type Service struct {
	catalog *repository.MongoCatalog
	cache   *rediscache.Redis
}

// NewService creates new *Service
func NewService(cat *repository.MongoCatalog, cache *rediscache.Redis) *Service {
	return &Service{
		cat,
		cache,
	}
}
