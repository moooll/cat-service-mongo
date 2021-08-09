package cache

import (
	"context"

	"github.com/go-redis/cache/v8"
	"github.com/google/uuid"
	"github.com/moooll/cat-service-mongo/internal/models"
)

// Redis contains  redis *cache.Cache var
type Redis struct {
	cache *cache.Cache
}

// NewRedisCache returns new cache
func NewRedisCache(c *cache.Cache) *Redis {
	return &Redis{
		cache: c,
	}
}

// GetFromCache gets record from the redis storage
func (c *Redis) GetFromCache(uid uuid.UUID) (cat models.Cat, err error) {
	id := uid.String()
	err = c.cache.Get(context.Background(), id, &cat)
	if err != nil {
		return models.Cat{}, err
	}

	return cat, nil
}

// SetToHash puts the record to redis storage
func (c *Redis) SetToHash(cat models.Cat) (err error) {
	err = c.cache.Set(&cache.Item{
		Key:   cat.ID.String(),
		Value: cat,
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteFromCache deletes record from redis cache
func (c *Redis) DeleteFromCache(uid uuid.UUID) (err error) {
	id := uid.String()
	err = c.cache.Delete(context.Background(), id)
	if err != nil {
		return err
	}

	return nil
}
