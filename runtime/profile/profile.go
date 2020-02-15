// Package profile is for specific profiles
package profile

// Local is a profile for local environments
func Local() []string {
	return []string{}
}

// Kubernetes is a profile for kubernetes
func Kubernetes() []string {
	return []string{}
}

// Platform is a platform profile
func Platform() []string {
	return []string{
		"MICRO_STORE=service",
		"MICRO_BROKER=service",
		"MICRO_RUNTIME=service",
		"MICRO_REGISTRY=service",
		// micro proxy routes all requests
		// and expects a k8s service name
		"MICRO_PROXY=go.micro.proxy",
		"MICRO_PROXY_ADDRESS=micro-proxy:8081",
		"GOPROXY=http://athens-proxy",
	}
}
