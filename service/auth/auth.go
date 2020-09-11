package auth

import (
	"context"

	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/micro/v3/service/auth/client"
)

const (
	// BearerScheme used for Authorization header
	BearerScheme = "Bearer "
)

var (
	// DefaultAuth implementation
	DefaultAuth auth.Auth = client.NewAuth()
	// ErrInvalidToken is when the token provided is not valid
	ErrInvalidToken = auth.ErrInvalidToken
	// ErrForbidden is when a user does not have the necessary scope to access a resource
	ErrForbidden = auth.ErrForbidden
)

type (
	// AccountToken is an alias for auth.Token
	AccountToken = auth.Token
	// Account is an alias for auth.Account
	Account = auth.Account
	// Resource is an alias for auth.Resource
	Resource = auth.Resource
	// Rule is an alias for auth.Rule
	Rule = auth.Rule
)

// Generate a new account
func Generate(id string, opts ...GenerateOption) (*Account, error) {
	return DefaultAuth.Generate(id, opts...)
}

// Verify an account has access to a resource using the rules
func Verify(acc *Account, res *Resource, opts ...VerifyOption) error {
	return DefaultAuth.Verify(acc, res, opts...)
}

// Inspect a token
func Inspect(token string) (*Account, error) {
	return DefaultAuth.Inspect(token)
}

// Token generated using refresh token or credentials
func Token(opts ...TokenOption) (*AccountToken, error) {
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
func Rules(...RulesOption) ([]*Rule, error) {
	return DefaultAuth.Rules()
}

type accountKey struct{}

// AccountFromContext gets the account from the context, which
// is set by the auth wrapper at the start of a call. If the account
// is not set, a nil account will be returned. The error is only returned
// when there was a problem retrieving an account
func AccountFromContext(ctx context.Context) (*Account, bool) {
	acc, ok := ctx.Value(accountKey{}).(*Account)
	return acc, ok
}

// ContextWithAccount sets the account in the context
func ContextWithAccount(ctx context.Context, account *Account) context.Context {
	return context.WithValue(ctx, accountKey{}, account)
}
