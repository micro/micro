package namespace

import (
	"context"
	"errors"

	"github.com/micro/micro/v3/service/auth"
)

var (
	// ErrUnauthorized is returned by Authorize when a context without a blank account tries to access
	// a restricted namespace
	ErrUnauthorized = errors.New("An account is required")
	// ErrForbidden is returned by Authorize when a context is trying to access a namespace it doesn't
	// have access to
	ErrForbidden = errors.New("Access denied to namespace")
)

const (
	// DefaultNamespace used by the server
	DefaultNamespace = "micro"
)

// Authorize will return an error if the context cannot access the given namespace
func Authorize(ctx context.Context, namespace string, opts ...AuthorizeOption) error {
	// parse the options
	var options AuthorizeOptions
	for _, o := range opts {
		o(&options)
	}

	// check to see if the namespace was made public
	if namespace == options.PublicNamespace {
		return nil
	}

	// accounts are always required so we can identify the caller. If auth is not configured, the noop
	// auth implementation will return a blank account with the default namespace set, allowing the caller
	// access to all resources
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return ErrUnauthorized
	}

	// the server can access all namespaces
	if acc.Issuer == DefaultNamespace {
		return nil
	}

	// ensure the account is requesing access to it's own namespace
	if acc.Issuer != namespace {
		return ErrForbidden
	}

	return nil
}

// AuthorizeOptions are used to configure the Authorize method
type AuthorizeOptions struct {
	PublicNamespace string
}

// AuthorizeOption sets an attribute on AuthorizeOptions
type AuthorizeOption func(o *AuthorizeOptions)

// Public indicates a namespace is public and can be accessed by anyone
func Public(ns string) AuthorizeOption {
	return func(o *AuthorizeOptions) {
		o.PublicNamespace = ns
	}
}
