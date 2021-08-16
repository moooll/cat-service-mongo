// Package cache contains functions to store, retrieve and delete data from redis cache
package cache

import (
	"context"
	"fmt"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"github.com/moooll/cat-service-mongo/internal/models"
)

// Redis contains  redis *cache.Cache var
type Redis struct {
	cache  *cache.Cache
	client *redis.Client
}

// NewRedisCache returns new cache
func NewRedisCache(c *cache.Cache, cl *redis.Client) *Redis {
	return &Redis{
		cache:  c,
		client: cl,
	}
}

// GetFromCache gets record from the redis storage
func (c *Redis) GetFromCache(ctx context.Context, uid uuid.UUID) (cat models.Cat, err error) {
	id := uid.String()
	err = c.cache.Get(ctx, id, &cat)
	if err != nil {
		return models.Cat{}, err
	}

	return cat, nil
}

// GetAllFromCache gets all records from the redis storage
func (c *Redis) GetAllFromCache(ctx context.Context) (cats []models.Cat, err error) {
	intK, err := c.client.Do(ctx, "KEYS", "*").Result()
	if err != nil {
		return []models.Cat{}, err
	}
	k := intK.([]interface{})
	for _, v := range k {
		key := fmt.Sprint(v)
		var cat models.Cat
		err = c.cache.Get(context.Background(), key, &cat)
		if err != nil {
			return []models.Cat{}, err
		}

		cats = append(cats, cat)
	}
	return cats, nil
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

// SetAllToHash puts all records to redis storage
func (c *Redis) SetAllToHash(cats []models.Cat) (err error) {
	for _, v := range cats {
		err = c.cache.Set(&cache.Item{
			Key:   v.ID.String(),
			Value: v,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteFromCache deletes record from redis cache
func (c *Redis) DeleteFromCache(ctx context.Context, uid uuid.UUID) (err error) {
	id := uid.String()
	err = c.cache.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
