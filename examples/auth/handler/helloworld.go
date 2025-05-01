package handler

import (
	"context"

	pb "github.com/micro/micro/examples/auth/proto"
)

type Helloworld struct{}

func (e *Helloworld) Call(ctx context.Context, req *pb.CallRequest, rsp *pb.CallResponse) error {
	rsp.Msg = "Hello " + req.Name

	return nil
}
