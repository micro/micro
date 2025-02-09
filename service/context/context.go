// Package context provides a context for accessing services
package context

import (
	"context"
)

var (
	// DefaultContext is a context which can be used to access micro services
	DefaultContext = WithNamespace("micro")

	NamespaceKey = "Micro-Namespace"
)

// WithNamespace creates a new context with the given namespace
func WithNamespace(ns string) context.Context {
	return SetNamespace(context.TODO(), ns)
}

// SetNamespace sets the namespace for a context
func SetNamespace(ctx context.Context, ns string) context.Context {
	return Set(ctx, NamespaceKey, ns)
}

// SetMetadata sets the metadata within the context
func SetMetadata(ctx context.Context, k, v string) context.Context {
	return Set(ctx, k, v)
}

// GetMetadata returns metadata from the context
func GetMetadata(ctx context.Context, k string) (string, bool) {
	return Get(ctx, k)
}
