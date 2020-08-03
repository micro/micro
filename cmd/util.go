package cmd

import (
	"time"

	"github.com/google/uuid"
	"github.com/micro/cli/v2"
	goauth "github.com/micro/go-micro/v3/auth"
	"github.com/micro/go-micro/v3/logger"
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

// setupAuthForService generates auth credentials for the service
func setupAuthForService() error {
	opts := auth.DefaultAuth.Options()

	// extract the account creds from options, these can be set by flags
	accID := opts.ID
	accSecret := opts.Secret

	// if no credentials were provided, self generate an account
	if len(accID) == 0 && len(accSecret) == 0 {
		opts := []goauth.GenerateOption{
			goauth.WithType("service"),
			goauth.WithScopes("service"),
		}

		acc, err := auth.Generate(uuid.New().String(), opts...)
		if err != nil {
			return err
		}
		if logger.V(logger.DebugLevel, logger.DefaultLogger) {
			logger.Debugf("Auth [%v] Generated an auth account", auth.DefaultAuth.String())
		}

		accID = acc.ID
		accSecret = acc.Secret
	}

	// generate the first token
	token, err := auth.Token(
		goauth.WithCredentials(accID, accSecret),
		goauth.WithExpiry(time.Minute*10),
	)
	if err != nil {
		return err
	}

	// set the credentials and token in auth options
	auth.DefaultAuth.Init(
		goauth.ClientToken(token),
		goauth.Credentials(accID, accSecret),
	)
	return nil
}

// refreshAuthToken if it is close to expiring
func refreshAuthToken(stop chan bool) {
	// can't refresh a token we dno't have
	if auth.DefaultAuth.Options().Token == nil {
		return
	}

	t := time.NewTicker(time.Second * 15)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			// don't refresh the token if it's not close to expiring
			tok := auth.DefaultAuth.Options().Token
			if tok.Expiry.Unix() > time.Now().Add(time.Minute).Unix() {
				continue
			}

			// generate the first token
			tok, err := auth.Token(
				goauth.WithToken(tok.RefreshToken),
				goauth.WithExpiry(time.Minute*10),
			)
			if err != nil {
				logger.Warnf("[Auth] Error refreshing token: %v", err)
				continue
			}

			// set the token
			auth.DefaultAuth.Init(goauth.ClientToken(tok))
		case <-stop:
			return
		}
	}
}
