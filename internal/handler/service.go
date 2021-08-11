package handler

import (
	"github.com/moooll/cat-service-mongo/internal/service"
	"github.com/moooll/cat-service-mongo/internal/streams"
)

// Service contains *service.Storage to interact with storage from handlers
type Service struct {
	storage *service.Storage
	stream *streams.StreamService
}

func NewService (s *service.Storage, stream *streams.StreamService) *Service {
	return &Service{
		storage: s,
		stream: stream,
	} 
}

