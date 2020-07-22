// Package profile is for specific profiles
// @todo this package is the definition of cruft and
// should be rewritten in a more elegant way
package profile

import (
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/registry/memory"
	"github.com/micro/micro/v2/service"
)

// Profile configures an environment
type Profile []service.Option

// Test profile is used for the go test suite
var Test Profile = []service.Option{
	service.Auth(auth.NewAuth()),
	service.Registry(memory.NewRegistry()),
}

// TODO: Convert the below profiles to use the new type

// Local is a profile for local environments
func Local() []string {
	return []string{}
}

// Server is a profile for running things through micro server
// eg runtime config etc will use actual services.
func Server() []string {
	return []string{
		"MICRO_AUTH=service",
		"MICRO_BROKER=service",
		"MICRO_REGISTRY=service",
		"MICRO_RUNTIME=service",
		"MICRO_STORE=service",
		"MICRO_CONFIG=service",
	}
}

func ServerCLI() []string {
	return []string{
		"MICRO_AUTH=service",
		"MICRO_BROKER=service",
		"MICRO_REGISTRY=service",
		"MICRO_RUNTIME=service",
		"MICRO_ROUTER=service",
		"MICRO_STORE=service",
		"MICRO_CONFIG=service",
	}
}

// Kubernetes is a profile for kubernetes
func Kubernetes() []string {
	return []string{
		"MICRO_AUTH=service",
		"MICRO_BROKER=service",
		"MICRO_CONFIG=service",
		"MICRO_NETWORK=service",
		"MICRO_REGISTRY=service",
		"MICRO_RUNTIME=service",
		"MICRO_STORE=service",
		"MICRO_AUTH_ADDRESS=micro-auth:8010",
		"MICRO_BROKER_ADDRESS=micro-store:8001",
		"MICRO_NETWORK_ADDRESS=micro-network:8080",
		"MICRO_REGISTRY_ADDRESS=micro-registry:8000",
		"MICRO_ROUTER_ADDRESS=micro-runtime:8084",
		"MICRO_RUNTIME_ADDRESS=micro-runtime:8088",
		"MICRO_STORE_ADDRESS=micro-store:8002",
	}
}

// Platform is a platform profile
func Platform() []string {
	return []string{
		// TODO: debug, monitor, etc
		"MICRO_AUTH=service",
		"MICRO_BROKER=service",
		"MICRO_CONFIG=service",
		"MICRO_NETWORK=service",
		"MICRO_REGISTRY=service",
		"MICRO_RUNTIME=service",
		"MICRO_ROUTER=service",
		"MICRO_STORE=service",
		// now set the addresses
		"MICRO_AUTH_ADDRESS=micro-auth.default.svc:8010",
		"MICRO_BROKER_ADDRESS=micro-store.default.svc:8001",
		"MICRO_NETWORK_ADDRESS=micro-network.default.svc:8080",
		"MICRO_REGISTRY_ADDRESS=micro-registry.default.svc:8000",
		"MICRO_ROUTER_ADDRESS=micro-runtime.default.svc:8084",
		"MICRO_RUNTIME_ADDRESS=micro-runtime.default.svc:8088",
		"MICRO_STORE_ADDRESS=micro-store.default.svc:8002",
	}
}

// Platform is a platform profile
func PlatformCLI() []string {
	return []string{
		// TODO: debug, monitor, etc
		"MICRO_AUTH=service",
		"MICRO_BROKER=service",
		"MICRO_CONFIG=service",
		"MICRO_REGISTRY=service",
		"MICRO_ROUTER=service",
		"MICRO_RUNTIME=service",
		"MICRO_STORE=service",
	}
}
