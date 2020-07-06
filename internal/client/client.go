package client

import (
	"context"
	"fmt"
	"strconv"
	"time"

	ccli "github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/auth"
	srvAuth "github.com/micro/go-micro/v2/auth/service"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/config/cmd"
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
	client := &wrapper{grpc.NewClient(), "", ns, env.ProxyAddress, env.Name, ctx}
	return client
}

type wrapper struct {
	client.Client
	token        string
	ns           string
	proxyAddress string
	envName      string
	context      *ccli.Context
}

func (a *wrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	err := a.getAccessToken(a.envName, a.context)
	if err != nil {
		return err
	}
	if len(a.token) > 0 {
		ctx = metadata.Set(ctx, "Authorization", auth.BearerScheme+a.token)
	}
	if len(a.proxyAddress) > 0 {
		opts = append(opts, client.WithAddress(a.proxyAddress))
	}
	ctx = metadata.Set(ctx, "Micro-Namespace", a.ns)
	return a.Client.Call(ctx, req, rsp, opts...)
}

func (a *wrapper) authFromContext(ctx *ccli.Context) auth.Auth {
	if cliutil.IsLocal(ctx) {
		return *cmd.DefaultCmd.Options().Auth
	}
	return srvAuth.NewAuth(
		auth.WithClient(a),
	)
}

// getAccessToken handles exchanging refresh tokens to access tokens
// The structure of the local micro userconfig file is the following:
// micro.auth.[envName].token: temporary access token
// micro.auth.[envName].refresh-token: long lived refresh token
// micro.auth.[envName].expiry: expiration time of the access token, seconds since Unix epoch.
func (a *wrapper) getAccessToken(envName string, ctx *ccli.Context) error {
	path := []string{"micro", "auth", envName}
	accessToken, _ := config.Get(append(path, "token")...)

	// Save the access token so it's usable for calls
	a.token = accessToken

	refreshToken, err := config.Get(append(path, "refresh-token")...)
	if err != nil {
		// Gracefully degrading here in case the user only has a temporary access token at hand.
		// The call will fail on the receiving end.
		return nil
	}

	// See if the access token has expired
	expiry, _ := config.Get("micro", "auth", envName, "refresh-token")
	if len(expiry) == 0 {
		return nil
	}
	expiryInt, err := strconv.ParseInt(expiry, 10, 64)
	if err != nil {
		return nil
	}
	if time.Now().Before(time.Unix(expiryInt, 0).Add(-15 * time.Second)) {
		return nil
	}
	// Get new access token from refresh token
	tok, err := a.authFromContext(a.context).Token(auth.WithToken(refreshToken))
	if err != nil {
		return err
	}

	// Save the token to user config file
	return SaveToken(envName, tok)
}

// SaveToken saves the auth token to the user's local config file
func SaveToken(envName string, token *auth.Token) error {
	if err := config.Set(token.AccessToken, "micro", "auth", envName, "token"); err != nil {
		return err
	}
	// Store the refresh token in micro config
	if err := config.Set(token.RefreshToken, "micro", "auth", envName, "refresh-token"); err != nil {
		return err
	}
	// Store the refresh token in micro config
	return config.Set(fmt.Sprintf("%v", token.Expiry.Unix()), "micro", "auth", envName, "expiry")
}
