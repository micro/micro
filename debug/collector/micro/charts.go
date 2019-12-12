package micro

import "github.com/netdata/go-orchestrator/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

const (
	chartServiceStarted = "micro_service_started"
	chartServiceUptime  = "micro_service_uptime"
	chartServiceMemory  = "micro_service_memory"
	chartServiceThreads = "micro_service_threads"
	chartServiceGC      = "micro_service_gc"
	chartServiceGCRate  = "micro_service_gcrate"
)

var charts = Charts{
	{
		ID:    chartServiceStarted,
		Title: "Start Time",
		Units: "timestamp",
		Fam:   "uptime",
		Ctx:   "micro.service.started",
	},
	{
		ID:    chartServiceUptime,
		Title: "Uptime",
		Units: "seconds",
		Fam:   "uptime",
		Ctx:   "micro.service.uptime",
	},
	{
		ID:    chartServiceMemory,
		Title: "Heap Allocated",
		Units: "B",
		Fam:   "memory",
		Ctx:   "micro.service.memory",
	},
	{
		ID:    chartServiceThreads,
		Title: "goroutines",
		Units: "goroutines",
		Fam:   "threads",
		Ctx:   "micro.service.threads",
	},
	{
		ID:    chartServiceGC,
		Title: "Cumulative GC Pause Total",
		Units: "nanoseconds",
		Fam:   "gc",
		Ctx:   "micro.service.gc",
	},
	{
		ID:    chartServiceGCRate,
		Title: "GC Pause rate",
		Units: "ns/s",
		Fam:   "gc",
		Ctx:   "micro.service.gcrate",
	},
	// TODO: debug_metrics when design is finalised.
}

