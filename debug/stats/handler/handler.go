// Package handler is the handler for the `micro debug stats` service
package handler

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/config/cmd"
	debug "github.com/micro/go-micro/debug/service/proto"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/cache"
	"github.com/micro/go-micro/util/log"
	stats "github.com/micro/micro/debug/stats/proto"
)

// New initialises and returns a new Stats service handler
func New() (*Stats, error) {
	s := &Stats{
		registry: cache.New(*cmd.DefaultOptions().Registry),
		client:   *cmd.DefaultOptions().Client,
	}

	if err := s.scan(); err != nil {
		return nil, err
	}
	done := make(chan bool)
	s.Start(done)
	defer close(done)
	return s, nil
}

// Stats is the Debug.Stats handler
type Stats struct {
	registry  registry.Registry
	client    client.Client
	snapshots []*stats.Snapshot
	next      []*stats.Snapshot
	cached    []*registry.Service
	sync.RWMutex
}

// Read returns gets a snapshot of all current stats
func (s *Stats) Read(ctx context.Context, req *stats.ReadRequest, rsp *stats.ReadResponse) error {
	s.RLock()
	rsp.Stats = s.snapshots
	s.RUnlock()
	return nil
}

// Stream starts streaming stats
func (s *Stats) Stream(ctx context.Context, req *stats.StreamRequest, rsp stats.Stats_StreamStream) error {
	return errors.New("Not Implemented")
}

// Start Starts scraping other services
// close the returned channel to stop scraping
func (s *Stats) Start(done <-chan bool) {
	go func(s *Stats) {
		for {
			select {
			case <-done:
				return
			default:
				s.scrape()
				time.Sleep(time.Second)
			}
		}
	}(s)
	go func(s *Stats) {
		rescan := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-done:
				rescan.Stop()
				return
			case <-rescan.C:
				s.scan()
			}
		}
	}(s)
}

func (s *Stats) scan() error {
	services, err := s.registry.ListServices()
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	serviceMap := make(map[string]*registry.Service)

	// check each service has nodes
	for _, service := range services {
		if len(service.Nodes) > 0 {
			serviceMap[service.Name+service.Version] = service
			continue
		}

		// get nodes that does not exist
		newServices, err := s.registry.GetService(service.Name)
		if err != nil {
			continue
		}

		// store service by version
		for _, service := range newServices {
			serviceMap[service.Name+service.Version] = service
		}
	}

	// flatten the map
	var serviceList []*registry.Service

	for _, service := range serviceMap {
		serviceList = append(serviceList, service)
	}

	// save the list
	s.Lock()
	s.cached = serviceList
	s.Unlock()
	return nil
}

func (s *Stats) scrape() {
	s.RLock()
	// Create a local copy of cached services
	services := make([]*registry.Service, len(s.cached))
	copy(services, s.cached)
	s.RUnlock()

	// Start building the next list of snapshots
	s.Lock()
	s.next = make([]*stats.Snapshot, 0)
	s.Unlock()

	// Call each node of each service in goroutines
	var wg sync.WaitGroup
	for _, svc := range services {

		// Ignore nodeless and non mucp services
		if len(svc.Nodes) == 0 {
			continue
		}
		if svc.Nodes[0].Metadata["protocol"] != "mucp" {
			continue
		}

		// Call every node
		for _, node := range svc.Nodes {
			go func(st *Stats, service *registry.Service, node *registry.Node) {
				wg.Add(1)
				defer wg.Done()

				// create new context to cancel within a few seconds
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()

				req := s.client.NewRequest(svc.Name, "Debug.Stats", &debug.StatsRequest{})
				rsp := new(debug.StatsResponse)
				if err := s.client.Call(ctx, req, rsp, client.WithAddress(node.Address)); err != nil {
					log.Errorf("Error calling %s@%s (%s)", svc.Name, node.Address, err.Error())
					return
				}

				// Append the new snapshot
				snap := stats.Snapshot{
					Service: &stats.Service{
						Name:    service.Name,
						Version: service.Version,
						Node: &stats.Node{
							Id: node.Address,
						},
					},
					Started: int64(rsp.Started),
					Uptime:  rsp.Uptime,
					Memory:  rsp.Memory,
					Threads: rsp.Threads,
					Gc:      rsp.Gc,
				}
				s.Lock()
				s.next = append(s.next, &snap)
				s.Unlock()
			}(s, svc, node)
		}
	}
	wg.Wait()

	// Swap in the snapshots
	s.Lock()
	s.snapshots = s.next
	s.Unlock()
}
