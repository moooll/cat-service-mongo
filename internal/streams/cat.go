package streams

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type PushArgs struct {
	Stream string
	Data   interface{}
}

type ReadArgs struct {
	Name string
	Id   string
	Mes chan interface{}
}

type StreamService struct {
	client *redis.Client
}

func NewStreamService(c *redis.Client) *StreamService {
	return &StreamService{
		client: c,
	}
}

// Push adds data to the stream
func (s *StreamService) Push(ctx context.Context, args PushArgs) error {
	err := s.client.XAdd(ctx, &redis.XAddArgs{
		Stream: args.Stream,
		Values: args.Data,
	}).Err()
	if err != nil {
		return err
	}

	return nil
}

// Read reads data from the stream
func (s *StreamService) Read(ctx context.Context, args ReadArgs) error {
	data, err := s.client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{args.Name, args.Id},
		Count: 1,
	}).Result()
	if err != nil {
		return err
	}

	for _, v := range data {
		for _, f := range v.Messages {
			args.Mes <- f
		}
	}

	return nil
}
