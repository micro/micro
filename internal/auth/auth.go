package auth

import "github.com/micro/go-micro/v2/auth"

// SystemRules are the default rules which are applied to the runtime services
var SystemRules = map[string][]*auth.Resource{
	"service": {
		&auth.Resource{Namespace: auth.DefaultNamespace, Type: "*", Name: "*", Endpoint: "*"},
	},
	"admin": {
		&auth.Resource{Namespace: auth.DefaultNamespace, Type: "*", Name: "*", Endpoint: "*"},
	},
	"developer": {
		&auth.Resource{Namespace: auth.DefaultNamespace, Type: "*", Name: "*", Endpoint: "*"},
	},
	"*": {
		&auth.Resource{Namespace: auth.DefaultNamespace, Type: "service", Name: "go.micro.auth", Endpoint: "Auth.Generate"},
		&auth.Resource{Namespace: auth.DefaultNamespace, Type: "service", Name: "go.micro.auth", Endpoint: "Auth.Token"},
		&auth.Resource{Namespace: auth.DefaultNamespace, Type: "service", Name: "go.micro.registry", Endpoint: "Registry.GetService"},
	},
}
