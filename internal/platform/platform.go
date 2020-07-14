// Package platform manages the runtime services as a platform
package platform

var (
	// list of services managed
	Services = []string{
		// runtime services
		"config",   // ????
		"network",  // :8085
		"runtime",  // :8088
		"registry", // :8000
		"broker",   // :8001
		"store",    // :8002
		"router",   // :8084
		"debug",    // :????
		"proxy",    // :8081
		"api",      // :8080
		"auth",     // :8010
		"web",      // :8082
	}
)
