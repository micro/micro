package auth

import (
	"github.com/micro/go-micro/v3/auth"
)

// TokenCookieName is the name of the cookie which stores the auth token
const TokenCookieName = "micro-token"

// SystemRules are the default rules which are applied to the runtime services
var SystemRules = []*auth.Rule{
	&auth.Rule{
		ID:       "default",
		Scope:    auth.ScopePublic,
		Access:   auth.AccessGranted,
		Resource: &auth.Resource{Type: "*", Name: "*", Endpoint: "*"},
	},
}
