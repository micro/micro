package handler

import (
	"context"
	"fmt"

	example "example-service/proto"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) Call(ctx context.Context, req *example.Request, rsp *example.Response) error {
	fmt.Println("Received Example.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}
