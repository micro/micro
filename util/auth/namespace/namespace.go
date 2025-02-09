package namespace

import (
	"context"
	"errors"

	"github.com/micro/micro/v5/service/auth"
	merrors "github.com/micro/micro/v5/service/errors"
)

var (
	// ErrUnauthorized is returned by Authorize when a context without a blank account tries to access
	// a restricted namespace
	ErrUnauthorized = errors.New("an account is required")
	// ErrForbidden is returned by Authorize when a context is trying to access a namespace it doesn't
	// have access to
	ErrForbidden = errors.New("access denied to namespace")
)

const (
	// DefaultNamespace used by the server
	DefaultNamespace = "micro"
)

// AuthorizeAdmin returns a service error if the context is not an admin that can access this namespace
// e.g. either an admin for the this namespace or an admin for micro
func AuthorizeAdmin(ctx context.Context, ns, method string) error {
	if err := Authorize(ctx, ns, method); err != nil {
		return err
	}

	adminAcc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return merrors.Unauthorized(method, "Unauthorized")
	}

	// check it's an admin
	if !hasTypeAndScope("user", "admin", adminAcc) && !hasTypeAndScope("service", "service", adminAcc) {
		return merrors.Unauthorized(method, "Unauthorized")
	}
	return nil
}

// Authorize will return a service error if the context cannot access the given namespace
func Authorize(ctx context.Context, namespace, method string, opts ...AuthorizeOption) error {
	if err := authorize(ctx, namespace); err == ErrForbidden {
		return merrors.Forbidden(method, err.Error())
	} else if err == ErrUnauthorized {
		return merrors.Unauthorized(method, err.Error())
	} else if err != nil {
		return merrors.InternalServerError(method, err.Error())
	}
	return nil
}

// authorize will return an error if the context cannot access the given namespace
func authorize(ctx context.Context, namespace string, opts ...AuthorizeOption) error {
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

	// the server and admins can access all namespaces
	if acc.Issuer == DefaultNamespace && (hasTypeAndScope("service", "service", acc) || hasTypeAndScope("user", "admin", acc)) {
		return nil
	}

	// ensure the account is requesting access to it's own namespace
	if acc.Issuer != namespace {
		return ErrForbidden
	}
	// account should be of type user or service
	if acc.Type != "user" && acc.Type != "service" {
		return ErrForbidden
	}

	return nil
}

func hasTypeAndScope(atype, scope string, acc *auth.Account) bool {
	if atype != acc.Type {
		return false
	}
	for _, s := range acc.Scopes {
		if s == scope {
			return true
		}
	}
	return false
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
