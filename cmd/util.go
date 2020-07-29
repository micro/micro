package cmd

import (
	"time"

	"github.com/google/uuid"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/micro/v3/client/cli/namespace"
	clitoken "github.com/micro/micro/v3/client/cli/token"
	"github.com/micro/micro/v3/client/cli/util"
	muauth "github.com/micro/micro/v3/service/auth"
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
		muauth.DefaultAuth.Init(auth.ClientToken(tok))
		return nil
	}

	// Get new access token from refresh token if it's close to expiry
	tok, err = muauth.DefaultAuth.Token(
		auth.WithToken(tok.RefreshToken),
		auth.WithTokenIssuer(ns),
	)
	if err != nil {
		clitoken.Remove(env.Name)
		return nil
	}

	// Save the token to user config file
	muauth.DefaultAuth.Init(auth.ClientToken(tok))
	return clitoken.Save(env.Name, tok)
}

// setupAuthForService generates auth credentials for the service
func setupAuthForService() error {
	a := muauth.DefaultAuth

	// extract the account creds from options, these can be set by flags
	accID := muauth.DefaultAuth.Options().ID
	accSecret := a.Options().Secret

	// if no credentials were provided, self generate an account
	if len(accID) == 0 && len(accSecret) == 0 {
		opts := []auth.GenerateOption{
			auth.WithType("service"),
			auth.WithScopes("service"),
		}

		acc, err := a.Generate(uuid.New().String(), opts...)
		if err != nil {
			return err
		}
		if logger.V(logger.DebugLevel, logger.DefaultLogger) {
			logger.Debugf("Auth [%v] Generated an auth account", a.String())
		}

		accID = acc.ID
		accSecret = acc.Secret
	}

	// generate the first token
	token, err := a.Token(
		auth.WithCredentials(accID, accSecret),
		auth.WithExpiry(time.Minute*10),
	)
	if err != nil {
		return err
	}

	// set the credentials and token in auth options
	a.Init(
		auth.ClientToken(token),
		auth.Credentials(accID, accSecret),
	)
	return nil
}

// refreshAuthToken if it is close to expiring
func refreshAuthToken(stop chan bool) {
	t := time.NewTicker(time.Second * 15)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			// don't refresh the token if it's not close to expiring
			tok := muauth.DefaultAuth.Options().Token
			if tok.Expiry.Unix() > time.Now().Add(time.Minute).Unix() {
				continue
			}

			// generate the first token
			tok, err := muauth.DefaultAuth.Token(
				auth.WithToken(tok.RefreshToken),
				auth.WithExpiry(time.Minute*10),
			)
			if err != nil {
				logger.Warnf("[Auth] Error refreshing token: %v", err)
				continue
			}

			// set the token
			muauth.DefaultAuth.Init(auth.ClientToken(tok))
		case <-stop:
			return
		}
	}
}
