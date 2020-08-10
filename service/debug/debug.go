package debug

import (
	"github.com/micro/go-micro/v3/debug/log"
	memLog "github.com/micro/go-micro/v3/debug/log/memory"
	"github.com/micro/go-micro/v3/debug/profile"
	"github.com/micro/go-micro/v3/debug/stats"
	memStats "github.com/micro/go-micro/v3/debug/stats/memory"
	"github.com/micro/go-micro/v3/debug/trace"
	memTrace "github.com/micro/go-micro/v3/debug/trace/memory"
)

var (
	DefaultLog      log.Log         = memLog.NewLog()
	DefaultTracer   trace.Tracer    = memTrace.NewTracer()
	DefaultStats    stats.Stats     = memStats.NewStats()
	DefaultProfiler profile.Profile = nil
)
