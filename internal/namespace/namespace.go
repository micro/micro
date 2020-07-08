package namespace

import (
	"context"
	"errors"

	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/metadata"
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
	// NamespaceKey is used to set/get the namespace from the context
	NamespaceKey = "Micro-Namespace"
)

// FromContext gets the namespace from the context
func FromContext(ctx context.Context) string {
	// get the namespace which is set at ingress by micro web / api / proxy etc. The go-micro auth
	// wrapper will ensure the account making the request has the necessary issuer.
	ns, _ := metadata.Get(ctx, NamespaceKey)
	return ns
}

// ContextWithNamespace sets the namespace in the context
func ContextWithNamespace(ctx context.Context, ns string) context.Context {
	return metadata.Set(ctx, NamespaceKey, ns)
}

// Authorize will return an error if the context cannot access the given namespace
func Authorize(ctx context.Context, namespace string, options ...Option) error {
	opts := &Options{}
	for _, o := range options {
		o(opts)
	}
	if len(opts.Public) > 0 {
		for _, public := range opts.Public {
			if namespace == public {
				return nil
			}
		}
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

// Options for namespace operations
type Options struct {
	// Public namespaces
	Public []string
}

// Option sets an option
type Option func(*Options)

func WithPublic(namespaces ...string) Option {
	return func(o *Options) {
		o.Public = namespaces
	}
}
