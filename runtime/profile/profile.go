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
		"MICRO_PROXY=go.micro.proxy",
		"MICRO_PROXY_ADDRESS=micro-proxy",
	}
}
