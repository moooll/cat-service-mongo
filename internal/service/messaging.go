package service

import (
	"context"
	"log"
	"time"

	"github.com/moooll/cat-service-mongo/internal/streams"
	"github.com/segmentio/kafka-go"
)

// MessageToRedisOnDelete sends message mes to the given stream
func MessageToRedisOnDelete(ctx context.Context, ss *streams.StreamService, stream string, mes string) error {
	data := map[string]interface{}{"act": "delete", "message": mes}
	err := ss.Push(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// // ListenRedisOnDelete listens
// func ListenRedisOnDelete(ctx context.Context, ss *streams.StreamService, id string) error {
// 	msg, err := ss.Read(ctx, id)
// 	if err != nil {
// 		return err
// 	}

// 	log.Println("message on delete: ", fmt.Sprint(msg))
// 	return nil
// }

// FromRedisToKafka reads messages from Redis stream and writes to kafka with kfk.Read() func.
func FromRedisToKafka(ctx context.Context, ss *streams.StreamService, conn *kafka.Conn) error {
	data, err := ss.Read(ctx, "$")
	if err != nil {
		return err
	}

	_, err = conn.WriteMessages(kafka.Message{
		Key:   []byte("delete-cats:"),
		Value: data.([]byte),
		Time:  time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}

// ReadFromKafka is a wrapper around internal/kafka.Read. Logs messages from Kafka.
func ReadFromKafka(ctx context.Context, conn *kafka.Conn, batch *kafka.Batch, b []byte) error {
	m, err := batch.Read(b)
	if err != nil {
		return err
	}

	log.Println(string(m))
	return nil
}

