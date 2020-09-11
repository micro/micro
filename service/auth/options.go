package auth

import (
	"context"
	"time"

	"github.com/micro/go-micro/v3/auth"
)

type (
	// GenerateOption is an alias for auth.GenerateOption
	GenerateOption = auth.GenerateOption
	// TokenOption is an alias for auth.TokenOption
	TokenOption = auth.TokenOption
	// VerifyOption is an alias for auth.VerifyOption
	VerifyOption = auth.VerifyOption
	// RulesOption is an alias for auth.RulesOption
	RulesOption = auth.RulesOption
)

// WithSecret for the generated account
func WithSecret(s string) GenerateOption {
	return func(o *auth.GenerateOptions) {
		o.Secret = s
	}
}

// WithType for the generated account
func WithType(t string) GenerateOption {
	return func(o *auth.GenerateOptions) {
		o.Type = t
	}
}

// WithMetadata for the generated account
func WithMetadata(md map[string]string) GenerateOption {
	return func(o *auth.GenerateOptions) {
		o.Metadata = md
	}
}

// WithProvider for the generated account
func WithProvider(p string) GenerateOption {
	return func(o *auth.GenerateOptions) {
		o.Provider = p
	}
}

// WithScopes for the generated account
func WithScopes(s ...string) GenerateOption {
	return func(o *auth.GenerateOptions) {
		o.Scopes = s
	}
}

// WithIssuer for the generated account
func WithIssuer(i string) GenerateOption {
	return func(o *auth.GenerateOptions) {
		o.Issuer = i
	}
}

// WithName for the generated account
func WithName(n string) GenerateOption {
	return func(o *auth.GenerateOptions) {
		o.Name = n
	}
}

// WithExpiry for the token
func WithExpiry(ex time.Duration) TokenOption {
	return func(o *auth.TokenOptions) {
		o.Expiry = ex
	}
}

func WithCredentials(id, secret string) TokenOption {
	return func(o *auth.TokenOptions) {
		o.ID = id
		o.Secret = secret
	}
}

func WithToken(rt string) TokenOption {
	return func(o *auth.TokenOptions) {
		o.RefreshToken = rt
	}
}

func WithTokenIssuer(iss string) TokenOption {
	return func(o *auth.TokenOptions) {
		o.Issuer = iss
	}
}

func VerifyContext(ctx context.Context) VerifyOption {
	return func(o *auth.VerifyOptions) {
		o.Context = ctx
	}
}
func VerifyNamespace(ns string) VerifyOption {
	return func(o *auth.VerifyOptions) {
		o.Namespace = ns
	}
}

func RulesContext(ctx context.Context) RulesOption {
	return func(o *auth.RulesOptions) {
		o.Context = ctx
	}
}

func RulesNamespace(ns string) RulesOption {
	return func(o *auth.RulesOptions) {
		o.Namespace = ns
	}
}
