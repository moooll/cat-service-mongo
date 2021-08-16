// Package service helps handlers to interact with Storage
package service

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/moooll/cat-service-mongo/internal/models"
)

// SaveToStorage saves the record both in redis cache and mongoDB
func (s *Storage) SaveToStorage(ctx context.Context, cat models.Cat) error {
	err := s.cache.SetToHash(cat)
	if err != nil {
		return err
	}

	err = s.catalog.Save(ctx, cat)
	if err != nil {
		return err
	}

	return nil
}

// DeleteFromStorage deletes the record both from redis cache and mongoDB
func (s *Storage) DeleteFromStorage(ctx context.Context, id uuid.UUID) error {
	err := s.cache.DeleteFromCache(ctx, id)
	if err != nil {
		log.Println("delete from cache error ", err)
		return err
	}

	_, err = s.catalog.Delete(ctx, id)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// GetAllFromStorage retrieves all records from redis cache and if not found from mongoDB
func (s *Storage) GetAllFromStorage(ctx context.Context) (cats []models.Cat, err error) {
	cats, err = s.cache.GetAllFromCache(ctx)
	if err != nil || len(cats) == 0 {
		cats, err = s.catalog.GetAll(ctx)
		if err != nil {
			return []models.Cat{}, err
		}

		err = s.cache.SetAllToHash(cats)
		if err != nil {
			return []models.Cat{}, err
		}
	}

	return cats, nil
}

// GetFromStorage retrieves a record from redis cache and if not found from mongoDB
func (s *Storage) GetFromStorage(ctx context.Context, id uuid.UUID) (models.Cat, error) {
	cat, err := s.cache.GetFromCache(ctx, id)
	if err != nil || (models.Cat{}) == cat {
		cat, err = s.catalog.Get(ctx, id)
		if err != nil {
			return models.Cat{}, err
		}

		err = s.cache.SetToHash(cat)
		if err != nil {
			return models.Cat{}, err
		}
	}

	return cat, nil
}

// UpdateStorage updates record both in redis cache and mongoDB
func (s *Storage) UpdateStorage(ctx context.Context, cat models.Cat) error {
	err := s.cache.SetToHash(cat)
	if err != nil {
		return err
	}

	err = s.catalog.Update(ctx, cat)
	if err != nil {
		return err
	}

	return nil
}
