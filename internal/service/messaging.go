package service

import (
	"context"
	"fmt"
	"log"

	"github.com/moooll/cat-service-mongo/internal/streams"
)

func MessageOnDelete(ctx context.Context, ss *streams.StreamService, stream string, mes string) error {
	data := map[string]interface{}{"act": "delete", "message": mes}
	err := ss.Push(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func ListenOnDelete(ctx context.Context, ss *streams.StreamService, id string) error {
	msg, err := ss.Read(ctx, id)
	if err != nil {
		return err
	}

	log.Println("message on delete: ", fmt.Sprint(msg))
	return nil
}
