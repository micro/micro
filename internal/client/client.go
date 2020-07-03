package client

import (
	"context"

	ccli "github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/micro/v2/client/cli/namespace"
	cliutil "github.com/micro/micro/v2/client/cli/util"
	"github.com/micro/micro/v2/internal/config"
)

// New returns a wrapped grpc client which will inject the
// token found in config into each request
func New(ctx *ccli.Context) client.Client {
	env := cliutil.GetEnv(ctx)
	ns, _ := namespace.Get(env.Name)
	token, _ := config.Get("micro", "auth", env.Name, "token")
	return &wrapper{grpc.NewClient(), token, ns, env.ProxyAddress}
}

type wrapper struct {
	client.Client
	token        string
	ns           string
	proxyAddress string
}

func (a *wrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	if len(a.token) > 0 {
		ctx = metadata.Set(ctx, "Authorization", auth.BearerScheme+a.token)
	}
	if len(a.proxyAddress) > 0 {
		opts = append(opts, client.WithAddress(a.proxyAddress))
	}
	ctx = metadata.Set(ctx, "Micro-Namespace", a.ns)
	return a.Client.Call(ctx, req, rsp, opts...)
}
