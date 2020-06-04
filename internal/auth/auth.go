package auth

import "github.com/micro/go-micro/v2/auth"

// TokenCookieName is the name of the cookie which stores the auth token
const TokenCookieName = "micro-token"

// SystemRules are the default rules which are applied to the runtime services
var SystemRules = []*auth.Rule{
	&auth.Rule{
		ID:       "default",
		Scope:    "*",
		Resource: &auth.Resource{Type: "*", Name: "*", Endpoint: "*"},
	},
	&auth.Rule{
		ID:       "auth-public",
		Scope:    "",
		Resource: &auth.Resource{Type: "service", Name: "go.micro.auth", Endpoint: "*"},
	},
	&auth.Rule{
		ID:       "registry-get",
		Scope:    "",
		Resource: &auth.Resource{Type: "service", Name: "go.micro.registry", Endpoint: "Registry.GetService"},
	},
	&auth.Rule{
		ID:       "registry-list",
		Scope:    "",
		Resource: &auth.Resource{Type: "service", Name: "go.micro.registry", Endpoint: "Registry.ListServices"},
	},
}
