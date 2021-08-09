package handler

import (
	"cat-service/internal/repository"
)

type Service struct {
	catalog *repository.Catalog
}

func NewService(cat *repository.Catalog) *Service {
	return &Service{
		cat,
	}
}