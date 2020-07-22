// Package runtime is the micro runtime
package runtime

import (
	"github.com/micro/go-micro/v2/runtime"
	"github.com/micro/go-micro/v2/runtime/service"
)

var (
	// DefaultRuntime implementation
	DefaultRuntime runtime.Runtime = service.NewRuntime()
)
