package auth

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/micro/v3/internal/cmd"
)

func init() {
	cmd.Init(func(ctx *cli.Context) error {
		// set the services namespace as the issuer
		if ns := ctx.String("namespace"); len(ns) > 0 {
			DefaultAuth.Init(auth.Issuer(ns))
		}

		// setup JWT private / public keys
		if len(ctx.String("auth_public_key")) > 0 {
			DefaultAuth.Init(auth.PublicKey(ctx.String("auth_public_key")))
		}
		if len(ctx.String("auth_private_key")) > 0 {
			DefaultAuth.Init(auth.PrivateKey(ctx.String("auth_private_key")))
		}

		// start a goroutine to refresh the access token periodically
		go refreshAccessToken()

		// don't setup auth credentials if we already have them
		if DefaultAuth.Options().Token != nil {
			return nil
		}

		// use the credentials passed to the service to authenticate
		if len(ctx.String("auth_id")) > 0 || len(ctx.String("auth_secret")) > 0 {
			return authUsingCreds(ctx.String("namespace"), ctx.String("auth_id"), ctx.String("auth_secret"))
		}

		// no credentials were provided, attempt to self-generate creds
		return generateCreds(ctx.String("namespace"))
	})
}

// DefaultAuth implementation
var DefaultAuth auth.Auth

// Generate a new account
func Generate(id string, opts ...auth.GenerateOption) (*auth.Account, error) {
	return DefaultAuth.Generate(id, opts...)
}

// Verify an account has access to a resource using the rules
func Verify(acc *auth.Account, res *auth.Resource, opts ...auth.VerifyOption) error {
	return DefaultAuth.Verify(acc, res, opts...)
}

// Inspect a token
func Inspect(token string) (*auth.Account, error) {
	return DefaultAuth.Inspect(token)
}

// Token generated using refresh token or credentials
func Token(opts ...auth.TokenOption) (*auth.Token, error) {
	return DefaultAuth.Token(opts...)
}

// Grant access to a resource
func Grant(rule *auth.Rule) error {
	return DefaultAuth.Grant(rule)
}

// Revoke access to a resource
func Revoke(rule *auth.Rule) error {
	return DefaultAuth.Revoke(rule)
}

// Rules returns all the rules used to verify requests
func Rules(...auth.RulesOption) ([]*auth.Rule, error) {
	return DefaultAuth.Rules()
}
