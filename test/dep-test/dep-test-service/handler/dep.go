package handler

import (
	"context"

	log "github.com/micro/go-micro/v3/logger"

	dep "dep-test-service/proto/dep"
)

type Dep struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Dep) Call(ctx context.Context, req *dep.Request, rsp *dep.Response) error {
	log.Info("Received Dep.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Dep) Stream(ctx context.Context, req *dep.StreamingRequest, stream dep.Dep_StreamStream) error {
	log.Infof("Received Dep.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&dep.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Dep) PingPong(ctx context.Context, stream dep.Dep_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&dep.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
