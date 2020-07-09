package handler

import (
	"context"

	foobarfn "foobarfn/proto/foobarfn"
)

type Foobarfn struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Foobarfn) Call(ctx context.Context, req *foobarfn.Request, rsp *foobarfn.Response) error {
	rsp.Msg = "Hello " + req.Name
	return nil
}
