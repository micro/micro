package subscriber

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	foobarfn "foobarfn/proto/foobarfn"
)

type Foobarfn struct{}

func (e *Foobarfn) Handle(ctx context.Context, msg *foobarfn.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}
