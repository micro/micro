// Package runtime is the micro runtime
package runtime

import (
	"github.com/micro/go-micro/v3/runtime"
	"github.com/micro/micro/v2/service/runtime/client"
)

var (
	// DefaultRuntime implementation
	DefaultRuntime runtime.Runtime = client.NewRuntime()
)
