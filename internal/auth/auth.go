package auth

import "github.com/micro/go-micro/v2/auth"

// TokenCookieName is the name of the cookie which stores the auth token
const TokenCookieName = "micro-token"

// SystemRules are the default rules which are applied to the runtime services
var SystemRules = map[string][]*auth.Resource{
	"*": {
		&auth.Resource{Type: "*", Name: "*", Endpoint: "*"},
	},
	"": {
		&auth.Resource{Type: "service", Name: "go.micro.auth", Endpoint: "Auth.Generate"},
		&auth.Resource{Type: "service", Name: "go.micro.auth", Endpoint: "Auth.Token"},
		&auth.Resource{Type: "service", Name: "go.micro.auth", Endpoint: "Auth.Inspect"},
		&auth.Resource{Type: "service", Name: "go.micro.registry", Endpoint: "Registry.GetService"},
		&auth.Resource{Type: "service", Name: "go.micro.registry", Endpoint: "Registry.ListServices"},
	},
}
