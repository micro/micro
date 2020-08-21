package handler

import (
	"context"

	log "github.com/micro/micro/v3/service/logger"

	idiomatic "idiomatic/proto"
)

type Idiomatic struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Idiomatic) Call(ctx context.Context, req *idiomatic.Request, rsp *idiomatic.Response) error {
	log.Info("Received Idiomatic.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Idiomatic) Stream(ctx context.Context, req *idiomatic.StreamingRequest, stream idiomatic.Idiomatic_StreamStream) error {
	log.Infof("Received Idiomatic.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&idiomatic.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Idiomatic) PingPong(ctx context.Context, stream idiomatic.Idiomatic_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&idiomatic.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
