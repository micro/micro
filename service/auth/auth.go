package auth

import (
	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/micro/v3/service/auth/client"
)

// DefaultAuth implementation
var DefaultAuth auth.Auth = client.NewAuth()

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
