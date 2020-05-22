package client

import (
	"context"
	"strings"

	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/micro/v2/client/cli/util"
	"github.com/micro/micro/v2/internal/config"
)

// New returns a wrapped grpc client which will inject the
// token found in config into each request
func New() client.Client {
	env, _ := config.Get("env")
	token, _ := config.Get(env, "auth", "token")
	return &wrapper{grpc.NewClient(), token, env}
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
