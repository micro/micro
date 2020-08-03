package client

import (
	"context"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/client/grpc"
	"github.com/micro/micro/v3/internal/cmd"
)

func init() {
	cmd.Init(func(ctx *cli.Context) error {
		// use the proxy address passed as a flag, e.g. the micro network
		if proxy := ctx.String("proxy_address"); len(proxy) > 0 {
			return DefaultClient.Init(client.Proxy(proxy))
		}

		return nil
	})
}

// DefaultClient for the service
var DefaultClient client.Client = grpc.NewClient()

// NewMessage returns a message which can be published
func NewMessage(topic string, msg interface{}, opts ...client.MessageOption) client.Message {
	return DefaultClient.NewMessage(topic, msg, opts...)
}

// NewRequest returns a request can which be executed using Call or Stream
func NewRequest(service, endpoint string, req interface{}, reqOpts ...client.RequestOption) client.Request {
	return DefaultClient.NewRequest(service, endpoint, req, reqOpts...)
}

// Call performs a request
func Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	return DefaultClient.Call(ctx, req, rsp, opts...)
}

// Stream performs a streaming request
func Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	return DefaultClient.Stream(ctx, req, opts...)
}

// Publish a message
func Publish(ctx context.Context, msg client.Message, opts ...client.PublishOption) error {
	return DefaultClient.Publish(ctx, msg, opts...)
}
