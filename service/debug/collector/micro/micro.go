package micro

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/micro/go-micro/v2/client"
	stats "github.com/micro/micro/v2/service/debug/stats/proto"
	"github.com/netdata/go-orchestrator/module"
)

// Config is the go-micro collector configuration
type Config struct{}

// New creates the micro module with default values
func New(c client.Client) *Micro {
	return &Micro{
		client:    c,
		Config:    Config{},
		charts:    charts(),
		metrics:   make(map[string]int64),
		services:  make(map[string]time.Time),
		snapshots: make(map[string][]*stats.Snapshot),
	}
}

// Register Registers with the orchestrator
func (m *Micro) Register() {
	module.Register("go_micro_services", module.Creator{
		Defaults: module.Defaults{
			Disabled: false,
		},
		Create: func() module.Module { return m },
	})
}

// Micro is a netdata module that collects metrics from the go.micro.debug service
type Micro struct {
	// netdata
	module.Base
	Config  `yaml:",inline"`
	charts  Charts
	metrics map[string]int64

	// go-micro
	client client.Client

	// internal
	sync.RWMutex
	// a cache of services we've seen
	services map[string]time.Time
	// the current list of service snapshots
	snapshots map[string][]*stats.Snapshot
}

// Cleanup is a no-op, called by netdata's orchestrator before shutdown
func (m *Micro) Cleanup() {}

// Init ensures the client is working, then starts saving snapshots.
func (m *Micro) Init() bool {
	// do initial scrape
	if err := m.collect(context.Background()); err != nil {
		m.Logger.Error(err)
	}

	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			if err := m.collect(ctx); err != nil {
				m.Logger.Error(err)
			}
			cancel()
			time.Sleep(time.Second)
		}
	}()

	return true
}

// Check can be called by netdata's orchestrator to ensure the collector is working
func (m *Micro) Check() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := m.collect(ctx)
	return err == nil
}

// Charts passes the module's charts to the orchestrator
func (m *Micro) Charts() *Charts {
	m.RLock()
	c := m.charts.Copy()
	m.RUnlock()
	return c
}

// Collect returns the current view of the metrics to the orchestrator
func (m *Micro) Collect() map[string]int64 {
	metrics := make(map[string]int64)
	m.RLock()
	for k, v := range m.metrics {
		metrics[k] = v
	}
	m.RUnlock()
	return metrics
}

func (m *Micro) updateCharts(snapshots []*stats.Snapshot) error {
	updateCharts := func(svc *stats.Snapshot, key string, idx int) {
		// create as many dimensions as we have
		for _, ch := range m.charts {
			name := strings.TrimPrefix(key, "go_micro_")
			id := fmt.Sprintf("%s_%d_%s", key, idx, ch.ID)

			// update our cached state
			m.services[id] = time.Now()

			switch ch.ID {
			case chartServiceGCRate, chartServiceRequests, chartServiceErrors:
				ch.AddDim(&module.Dim{
					Algo: module.Incremental,
					ID:   id,
					Name: fmt.Sprintf("%s.%d", name, idx),
				})
			default:
				ch.AddDim(&module.Dim{
					Algo: module.Absolute,
					ID:   id,
					Name: fmt.Sprintf("%s.%d", name, idx),
				})
			}

			ch.MarkNotCreated()
		}
	}

	// Construct a new Micro.services map based on the list of incoming service snapshots.
	// The map is keyed on Name_Version and sorted by Node ID for consistent graphs
	newServices := make(map[string][]*stats.Snapshot)

	// generate a new service map based on the snapshot
	for _, snap := range snapshots {
		key := key(snap.Service)
		if srv, found := newServices[key]; found {
			newServices[key] = append(srv, snap)
		} else {
			newServices[key] = []*stats.Snapshot{snap}
		}
	}

	for _, services := range newServices {
		sort.Slice(services, func(i, j int) bool {
			return services[i].Service.Node.Id < services[j].Service.Node.Id
		})
	}

	// Check for any services that we used to have that disappeared
	m.Lock()
	defer m.Unlock()

	// Create / remove chart dimensions based on the previous state.
	for key, services := range newServices {
		// update the charts
		for i, service := range services {
			updateCharts(service, key, i)
		}
	}

	// check if we have any stale charts that need to be removed
	for id, lastSeen := range m.services {
		if time.Since(lastSeen) > time.Hour {
			// Service was in old map, isn't in new map, so remove the dimensions for it.
			for _, ch := range m.charts {
				if ch.HasDim(id) {
					ch.MarkDimRemove(id, true)
				}
			}
		}
	}

	// swap in the new services, then return (m.Unlock was deferred)
	m.snapshots = newServices

	return nil
}

// Collect contacts the Debug service to retrieve snapshots of stats
func (m *Micro) collect(ctx context.Context) error {
	// Grab snapshots from the Debug service
	req := &stats.ReadRequest{}
	rsp := &stats.ReadResponse{}

	err := m.client.Call(ctx, m.client.NewRequest("go.micro.debug", "Stats.Read", req), rsp)
	if err != nil {
		return err
	}

	// If we don't already have a dimension for the service, create it
	if err := m.updateCharts(rsp.Stats); err != nil {
		m.Logger.Error(err)
		return err
	}

	// Populate metrics map
	m.RLock()

	for name, services := range m.snapshots {
		for i, s := range services {
			m.metrics[fmt.Sprintf("%s_%d_%s", name, i, chartServiceStarted)] = int64(s.Started)
			m.metrics[fmt.Sprintf("%s_%d_%s", name, i, chartServiceUptime)] = int64(s.Uptime)
			m.metrics[fmt.Sprintf("%s_%d_%s", name, i, chartServiceMemory)] = int64(s.Memory)
			m.metrics[fmt.Sprintf("%s_%d_%s", name, i, chartServiceThreads)] = int64(s.Threads)
			m.metrics[fmt.Sprintf("%s_%d_%s", name, i, chartServiceGC)] = int64(s.Gc)
			m.metrics[fmt.Sprintf("%s_%d_%s", name, i, chartServiceGCRate)] = int64(s.Gc)
			m.metrics[fmt.Sprintf("%s_%d_%s", name, i, chartServiceRequests)] = int64(s.Requests)
			m.metrics[fmt.Sprintf("%s_%d_%s", name, i, chartServiceErrors)] = int64(s.Errors)
		}
	}

	m.RUnlock()

	return nil
}

func format(v string) string {
	return strings.ReplaceAll(v, ".", "_")
}

func key(s *stats.Service) string {
	return format(s.Name + "-" + s.Version)
}
