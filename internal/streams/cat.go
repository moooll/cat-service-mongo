// Package streams helps to work with Redis Streams, push and read messages from streams
package streams

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// StreamService represents the entity for working with redis streams
type StreamService struct {
	client *redis.Client
}

// NewStreamService creates new instance of StreamService
func NewStreamService(c *redis.Client) *StreamService {
	return &StreamService{
		client: c,
	}
}

// Push adds data to the stream
func (s *StreamService) Push(ctx context.Context, data interface{}) error {
	err := s.client.XAdd(ctx, &redis.XAddArgs{
		Stream: "delete-cats",
		MaxLen: 0,
		Values: data,
	}).Err()
	if err != nil {
		return err
	}

	return nil
}

// Read reads data from the stream
func (s *StreamService) Read(ctx context.Context, id string) (data interface{}, err error) {
	data, err = s.client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{"delete-cats", id},
		Count:   1,
		Block:   0,
	}).Result()
	if err != nil {
		return nil, err
	}

	return data, nil
}
