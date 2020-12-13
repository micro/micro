package wrapper

import (
	"context"

	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context/metadata"
	"github.com/micro/micro/v3/service/server"
)

type fromServiceWrapper struct {
	client.Client
}

var (
	HeaderPrefix = "Micro-"
)

func (f *fromServiceWrapper) setHeaders(ctx context.Context) context.Context {
	return metadata.MergeContext(ctx, metadata.Metadata{
		HeaderPrefix + "From-Service": server.DefaultServer.Options().Name,
	}, false)
}

func (f *fromServiceWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	ctx = f.setHeaders(ctx)
	return f.Client.Call(ctx, req, rsp, opts...)
}

func (f *fromServiceWrapper) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	ctx = f.setHeaders(ctx)
	return f.Client.Stream(ctx, req, opts...)
}

func (f *fromServiceWrapper) Publish(ctx context.Context, p client.Message, opts ...client.PublishOption) error {
	ctx = f.setHeaders(ctx)
	return f.Client.Publish(ctx, p, opts...)
}

// FromService wraps a client to inject service and auth metadata
func FromService(c client.Client) client.Client {
	return &fromServiceWrapper{c}
}
