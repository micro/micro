// Package profile is for specific profiles
package profile

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
		"MICRO_ROUTER=service",
		"MICRO_RUNTIME=service",
		"MICRO_STORE=service",
		"MICRO_STORE_BACKEND=file",
		"MICRO_PROXY=service",
	}
}

// Kubernetes is a profile for kubernetes
func Kubernetes() []string {
	return []string{}
}

// Platform is a platform profile
func Platform() []string {
	return []string{
		// TODO: debug, monitor, etc
		"MICRO_AUTH=service",
		"MICRO_BROKER=service",
		"MICRO_REGISTRY=service",
		"MICRO_ROUTER=service",
		"MICRO_RUNTIME=service",
		"MICRO_STORE=service",
		"MICRO_PROXY=service",
		"MICRO_CONFIG=service",
		// now set the addresses
		"MICRO_BROKER_ADDRESS=micro-store:8001",
		"MICRO_REGISTRY_ADDRESS=micro-registry:8000",
		"MICRO_PROXY_ADDRESS=micro-proxy:8081",
		"MICRO_ROUTER_ADDRESS=micro-runtime:8084",
		"MICRO_RUNTIME_ADDRESS=micro-runtime:8088",
		"MICRO_STORE_ADDRESS=micro-store:8002",
		// set the athens proxy to speedup builds
		"GOPROXY=http://athens-proxy",
	}
}
