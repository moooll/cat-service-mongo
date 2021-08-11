package service

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/moooll/cat-service-mongo/internal/streams"
)

type ListenOnDeleteArgs struct {
	Ctx context.Context
	Wg *sync.WaitGroup
	Ss streams.StreamService
	Args streams.ReadArgs
	ErChan chan error
}
func MessageOnDelete(ctx context.Context, ss *streams.StreamService, stream string, mes string) error {
	data := map[string]interface{}{"act": "delete", "message": mes}
	args := &streams.PushArgs{
		Stream: stream,
		Data:   data,
	}
	err := ss.Push(ctx, *args)
	if err != nil {
		return err
	}

	return nil
}

func ListenOnDelete(args ListenOnDeleteArgs) {
	defer args.Wg.Done()
	var resp interface{}
	err := args.Ss.Read(args.Ctx, args.Args)
	if err != nil {
		args.ErChan <- err
	}

	log.Println("message on delete: ", fmt.Sprint(resp))
}
