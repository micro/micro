package client

import (
	"context"

	ccli "github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/metadata"
	cliutil "github.com/micro/micro/v2/client/cli/util"
	"github.com/micro/micro/v2/internal/config"
)

// New returns a wrapped grpc client which will inject the
// token found in config into each request
func New(ctx *ccli.Context) client.Client {
	env := cliutil.GetEnv(ctx)
	token, _ := config.Get("micro", "auth", env.Name, "token")
	return &wrapper{grpc.NewClient(), token, env.Name, ctx}
}

type wrapper struct {
	client.Client
	token string
	env   string
	ctx   *ccli.Context
}

func (a *wrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	if len(a.token) > 0 {
		ctx = metadata.Set(ctx, "Authorization", auth.BearerScheme+a.token)
	}

	// network will be used by the router to filter available routes.
	// var network string
	// if util.IsLocal(a.ctx) || util.IsServer(a.ctx) {
	// 	network = router.DefaultNetwork
	// } else {
	// 	network = strings.ReplaceAll(a.env, "/", "-")
	// }

	// ctx = metadata.Set(ctx, "Micro-Network", network)
	return a.Client.Call(ctx, req, rsp, opts...)
}
