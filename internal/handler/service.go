package handler

import (
	"github.com/moooll/cat-service-mongo/internal/service"
)

// Service contains *service.Storage to interact with storage from handlers
type Service struct {
	storage *service.Storage
}

func NewService (s *service.Storage) *Service {
	return &Service{
		storage: s,
	} 
}

