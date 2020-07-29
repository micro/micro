package subscriber

import (
	"context"
	log "github.com/micro/micro/v3/service/logger"

	dep "dep-test-service/proto/dep"
)

type Dep struct{}

func (e *Dep) Handle(ctx context.Context, msg *dep.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *dep.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
