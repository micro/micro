package client

import (
	"context"
	"strings"

	ccli "github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/micro/v2/client/cli/util"
	cliutil "github.com/micro/micro/v2/client/cli/util"
	"github.com/micro/micro/v2/internal/config"
)

// New returns a wrapped grpc client which will inject the
// token found in config into each request
func New(ctx *ccli.Context) client.Client {
	var env cliutil.Env
	if len(ctx.String("env")) > 0 {
		env = cliutil.GetEnvByName(ctx.String("env"))
	} else {
		env = cliutil.GetEnv()
	}
	token, _ := config.Get("micro", "auth", env.Name, "token")
	return &wrapper{grpc.NewClient(), token, env.Name}
}

type wrapper struct {
	client.Client

	token string
	env   string
}

func (a *wrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	if len(a.token) > 0 {
		ctx = metadata.Set(ctx, "Authorization", auth.BearerScheme+a.token)
	}
	if len(a.env) > 0 && !util.IsLocal() && !util.IsServer() {
		env := strings.ReplaceAll(a.env, "/", "-")
		ctx = metadata.Set(ctx, "Micro-Namespace", env)
	}
	return a.Client.Call(ctx, req, rsp, opts...)
}
