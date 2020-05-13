package namespace

import (
	"context"

	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/metadata"
)

const (
	// TODO: Move default namespace out of go-micro
	DefaultNamespace = auth.DefaultNamespace
	// NamespaceKey is used to set/get the namespace from the
	// context
	NamespaceKey = "Micro-Namespace"
)

// FromContext gets the namespace from the context
func FromContext(ctx context.Context) string {
	// if there is an account, we use its namespace
	if acc, ok := auth.AccountFromContext(ctx); ok {
		return acc.Namespace
	}

	// next check for the namespace key set by micro web or api
	if ns, ok := metadata.Get(ctx, NamespaceKey); ok {
		return ns
	}

	// fallback to the default namespace
	return DefaultNamespace
}
