package auth

import "github.com/micro/go-micro/v2/auth"

// SystemRules are the default rules which are applied to the runtime services
var SystemRules = map[string][]*auth.Resource{
	"*": {
		&auth.Resource{Namespace: "*", Type: "*", Name: "*", Endpoint: "*"},
	},
	"": {
		&auth.Resource{Namespace: "*", Type: "service", Name: "go.micro.auth", Endpoint: "Auth.Generate"},
		&auth.Resource{Namespace: "*", Type: "service", Name: "go.micro.auth", Endpoint: "Auth.Token"},
		&auth.Resource{Namespace: "*", Type: "service", Name: "go.micro.auth", Endpoint: "Auth.Inspect"},
		&auth.Resource{Namespace: "*", Type: "service", Name: "go.micro.registry", Endpoint: "Registry.GetService"},
		&auth.Resource{Namespace: "*", Type: "service", Name: "go.micro.registry", Endpoint: "Registry.ListServices"},
	},
}
