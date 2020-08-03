package cmd

import (
	"time"

	"github.com/micro/cli/v2"
	goauth "github.com/micro/go-micro/v3/auth"
	"github.com/micro/micro/v3/client/cli/namespace"
	clitoken "github.com/micro/micro/v3/client/cli/token"
	"github.com/micro/micro/v3/client/cli/util"
	"github.com/micro/micro/v3/service/auth"
)

// setupAuthForCLI handles exchanging refresh tokens to access tokens
// The structure of the local micro userconfig file is the following:
// micro.auth.[envName].token: temporary access token
// micro.auth.[envName].refresh-token: long lived refresh token
// micro.auth.[envName].expiry: expiration time of the access token, seconds since Unix epoch.
func setupAuthForCLI(ctx *cli.Context) error {
	env := util.GetEnv(ctx)
	ns, err := namespace.Get(env.Name)
	if err != nil {
		return err
	}

	tok, err := clitoken.Get(env.Name)
	if err != nil {
		return err
	}

	// If there is no refresh token, do not try to refresh it
	if len(tok.RefreshToken) == 0 {
		return nil
	}

	// Check if token is valid
	if time.Now().Before(tok.Expiry.Add(-15 * time.Second)) {
		auth.DefaultAuth.Init(goauth.ClientToken(tok))
		return nil
	}

	// Get new access token from refresh token if it's close to expiry
	tok, err = auth.Token(
		goauth.WithToken(tok.RefreshToken),
		goauth.WithTokenIssuer(ns),
	)
	if err != nil {
		clitoken.Remove(env.Name)
		return nil
	}

	// Save the token to user config file
	auth.DefaultAuth.Init(goauth.ClientToken(tok))
	return clitoken.Save(env.Name, tok)
}
