package wrapper

import (
	"context"

	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/server"
)

type logWrapper struct {
	client.Client
}

func (l *logWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	logger.Debugf("Calling service %s endpoint %s", req.Service(), req.Endpoint())
	return l.Client.Call(ctx, req, rsp, opts...)
}

func LogClient(c client.Client) client.Client {
	return &logWrapper{c}
}

func LogHandler() server.HandlerWrapper {
	// return a handler wrapper
	return func(h server.HandlerFunc) server.HandlerFunc {
		// return a function that returns a function
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			logger.Debugf("Serving request for service %s endpoint %s", req.Service(), req.Endpoint())
			return h(ctx, req, rsp)
		}
	}
}
