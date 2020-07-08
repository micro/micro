// Package client handles calls from the CLI calls
package client

import (
	"context"
	"fmt"
	"os"
	"time"

	ccli "github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/auth"
	srvAuth "github.com/micro/go-micro/v2/auth/service"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/micro/v2/client/cli/namespace"
	clitoken "github.com/micro/micro/v2/client/cli/token"
	cliutil "github.com/micro/micro/v2/client/cli/util"
)

// New returns a wrapped grpc client which will inject the
// token found in config into each request
func New(ctx *ccli.Context) client.Client {
	env := cliutil.GetEnv(ctx)
	ns, _ := namespace.Get(env.Name)
	client := &wrapper{
		Client:       grpc.NewClient(),
		token:        "",
		ns:           ns,
		proxyAddress: env.ProxyAddress,
		envName:      env.Name,
		context:      ctx,
	}
	err := client.getAccessToken(env.Name, ctx)
	if err != nil {
		// @todo this is veeery ugly being here
		fmt.Println(err)
		os.Exit(1)
	}
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
	tok, err := clitoken.Get(envName)
	if err != nil {
		return err
	}

	// Save the access token so it's usable for calls
	a.token = tok.AccessToken

	// If there is no refresh token, do not try to refresh it
	if len(tok.RefreshToken) == 0 {
		return nil
	}

	// Check if token must be refreshed
	if time.Now().Before(tok.Expiry.Add(-15 * time.Second)) {
		return nil
	}

	// Get new access token from refresh token if it's close to expiry
	tok, err = a.authFromContext(a.context).Token(auth.WithToken(tok.RefreshToken))
	if err != nil {
		return err
	}

	// Save the new access token
	a.token = tok.AccessToken

	// Save the token to user config file
	return clitoken.Save(envName, tok)
}
