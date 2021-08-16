package handler

import (
	"github.com/moooll/cat-service-mongo/internal/service"
	"github.com/moooll/cat-service-mongo/internal/streams"
)

// Service helps to interact with storage from handlers
type Service struct {
	storage *service.Storage
	stream  *streams.StreamService
}

// NewService creates a new instance of *Service
func NewService(s *service.Storage, stream *streams.StreamService) *Service {
	return &Service{
		storage: s,
		stream:  stream,
	}
}
