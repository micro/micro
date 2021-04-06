package auth

import (
	"github.com/micro/micro/v3/service/auth"
)

const (
	// BearerScheme used for Authorization header
	BearerScheme = "Bearer "
	// TokenCookieName is the name of the cookie which stores the auth token
	TokenCookieName = "micro-token"
)

// SystemRules are the default rules which are applied to the runtime services
var SystemRules = []*auth.Rule{
	&auth.Rule{
		ID:       "default",
		Scope:    auth.ScopePublic,
		Access:   auth.AccessGranted,
		Resource: &auth.Resource{Type: "*", Name: "*", Endpoint: "*"},
	},
}
