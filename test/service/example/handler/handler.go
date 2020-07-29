package handler

import (
	"context"
	"fmt"

	pb "example/proto"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) Call(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	fmt.Println("Received Example.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}
