package auth

import (
	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/micro/v3/service/auth/client"
)

// DefaultAuth implementation
var DefaultAuth auth.Auth = client.NewAuth()

type (
	// AccessToken is an alias for auth.Token
	AccessToken = auth.Token
	// Account is an alias for auth.Account
	Account = auth.Account
	// Resource is an alias for auth.Resource
	Resource = auth.Resource
	// Rule is an alias for auth.Rule
	Rule = auth.Rule
)

// Generate a new account
func Generate(id string, opts ...auth.GenerateOption) (*Account, error) {
	return DefaultAuth.Generate(id, opts...)
}

// Verify an account has access to a resource using the rules
func Verify(acc *Account, res *Resource, opts ...auth.VerifyOption) error {
	return DefaultAuth.Verify(acc, res, opts...)
}

// Inspect a token
func Inspect(token string) (*Account, error) {
	return DefaultAuth.Inspect(token)
}

// Token generated using refresh token or credentials
func Token(opts ...auth.TokenOption) (*AccessToken, error) {
	return DefaultAuth.Token(opts...)
}

// Grant access to a resource
func Grant(rule *Rule) error {
	return DefaultAuth.Grant(rule)
}

// Revoke access to a resource
func Revoke(rule *Rule) error {
	return DefaultAuth.Revoke(rule)
}

// Rules returns all the rules used to verify requests
func Rules(...auth.RulesOption) ([]*Rule, error) {
	return DefaultAuth.Rules()
}
