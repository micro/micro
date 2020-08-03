package template

var (
	HandlerSRV = `package handler

import (
	"context"

	log "github.com/micro/micro/v3/service/logger"

	{{dehyphen .Alias}} "{{.Dir}}/proto"
)

type {{title .Alias}} struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *{{title .Alias}}) Call(ctx context.Context, req *{{dehyphen .Alias}}.Request, rsp *{{dehyphen .Alias}}.Response) error {
	log.Info("Received {{title .Alias}}.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *{{title .Alias}}) Stream(ctx context.Context, req *{{dehyphen .Alias}}.StreamingRequest, stream {{dehyphen .Alias}}.{{title .Alias}}_StreamStream) error {
	log.Infof("Received {{title .Alias}}.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&{{dehyphen .Alias}}.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *{{title .Alias}}) PingPong(ctx context.Context, stream {{dehyphen .Alias}}.{{title .Alias}}_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&{{dehyphen .Alias}}.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
`

	SubscriberSRV = `package subscriber

import (
	"context"
	log "github.com/micro/micro/v3/service/logger"

	{{dehyphen .Alias}} "{{.Dir}}/proto"
)

type {{title .Alias}} struct{}

func (e *{{title .Alias}}) Handle(ctx context.Context, msg *{{dehyphen .Alias}}.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *{{dehyphen .Alias}}.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
`
)
