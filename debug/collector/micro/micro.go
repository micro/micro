package micro

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/netdata/go-orchestrator/module"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	stats "github.com/micro/micro/debug/stats/proto"
)

// Config is the go-micro collector configuration
type Config struct{}

// func init() {
// 	creator := module.Creator{
// 		Defaults: module.Defaults{
// 			Disabled: false,
// 		},
// 		Create: func() module.Module { return New() },
// 	}

// 	module.Register("go_micro_services", creator)
// }

// New creates the micro module with default values
func New() *Micro {
	return &Micro{
		Config:   Config{},
		services: make(map[string]bool),
		indexes:  make(map[string]map[string]bool),
		metrics:  make(map[string]int64),
		client:   client.DefaultClient,
	}
}

// WithClient sets the micro client
func (m *Micro) WithClient(c client.Client) *Micro {
	m.client = c
	return m
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

// Micro is a netdata module that collects metrics from go-micro services
type Micro struct {
	// netdata
	module.Base
	Config `yaml:",inline"`
	// charts  *Charts
	metrics map[string]int64

	// go-micro
	client client.Client

	// internal
	sync.RWMutex
	cached   []*registry.Service
	services map[string]bool
	indexes  map[string]map[string]bool
}

// Cleanup makes cleanup
func (Micro) Cleanup() {}

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

// Check makes check
func (m *Micro) Check() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := m.collect(ctx)
	return err == nil
}

// Charts creates Charts
func (m *Micro) Charts() *Charts {
	m.RLock()
	c := charts.Copy()
	m.RUnlock()
	return c
}

// Collect collects metrics
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
	sort.Sort(sortableSnapshot(snapshots))
	getIndex := func(s *stats.Snapshot) string {
		if _, found := m.indexes[s.Service.Name]; !found {
			m.indexes[s.Service.Name] = make(map[string]bool)
		}
		m.indexes[s.Service.Name][key(s)] = true
		return strconv.Itoa(len(m.indexes[s.Service.Name]))
	}
	m.Lock()
	defer m.Unlock()
	for _, snap := range snapshots {
		svc := key(snap)
		if _, found := m.services[svc]; !found {
			m.services[svc] = true
			for _, ch := range charts {
				if ch.ID == chartServiceGCRate {
					ch.AddDim(&module.Dim{
						ID:   svc + "_" + ch.ID,
						Name: strings.TrimPrefix(strings.ReplaceAll(snap.Service.Name, ".", "_"), "go_micro_") + getIndex(snap),
						Algo: module.Incremental,
					})
				} else {
					ch.AddDim(&module.Dim{
						ID:   svc + "_" + ch.ID,
						Name: strings.TrimPrefix(strings.ReplaceAll(snap.Service.Name, ".", "_"), "go_micro_") + getIndex(snap),
						Algo: module.Absolute,
					})
				}
				m.Logger.Debug("Added dimension" + svc + "_" + ch.ID)
				ch.MarkNotCreated()
			}
		}
	}
	return nil
}

// Collect contacts the Debug service to retrieve snapshots of stats
func (m *Micro) collect(ctx context.Context) error {
	// Grab snapshots from the Debug service
	req := &stats.ReadRequest{}
	rsp := &stats.ReadResponse{}
	err := m.client.Call(ctx, client.NewRequest("go.micro.debug", "Stats.Read", req), rsp)
	if err != nil {
		return err
	}

	// If we don't already have a dimension for the service, create it
	if err := m.updateCharts(rsp.Stats); err != nil {
		m.Logger.Error(err)
		return err
	}

	// Populate metrics map
	m.Lock()
	for _, s := range rsp.Stats {
		k := key(s)
		m.metrics[k+"_"+chartServiceStarted] = int64(s.Started)
		m.metrics[k+"_"+chartServiceUptime] = int64(s.Uptime)
		m.metrics[k+"_"+chartServiceMemory] = int64(s.Memory)
		m.metrics[k+"_"+chartServiceThreads] = int64(s.Threads)
		m.metrics[k+"_"+chartServiceGC] = int64(s.Gc)
		m.metrics[k+"_"+chartServiceGCRate] = int64(s.Gc)
	}
	m.Unlock()
	return nil
}

func key(s *stats.Snapshot) string {
	return strings.ReplaceAll(s.Service.Node.Id+s.Service.Version, ".", "_")
}

type sortableSnapshot []*stats.Snapshot

func (s sortableSnapshot) Len() int      { return len(s) }
func (s sortableSnapshot) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s sortableSnapshot) Less(i, j int) bool {
	return s[i].Service.Node.Id+s[i].Service.Version < s[j].Service.Node.Id+s[j].Service.Version
}

