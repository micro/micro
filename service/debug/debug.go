package debug

import (
	"github.com/micro/go-micro/v3/debug/profile"
	"github.com/micro/go-micro/v3/debug/stats"
	"github.com/micro/go-micro/v3/debug/trace"
	"github.com/micro/go-micro/v3/debug/trace/memory"
)

var (
	DefaultTracer   trace.Tracer    = memory.NewTracer()
	DefaultStats    stats.Stats     = stats.NewStats()
	DefaultProfiler profile.Profile = nil
)
