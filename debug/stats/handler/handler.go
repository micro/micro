// Package handler is the handler for the `micro debug stats` service
package handler

import (
	"context"
	"sync"
	"time"

	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/config/cmd"
	debug "github.com/micro/go-micro/debug/service/proto"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/cache"
	"github.com/micro/go-micro/util/log"
	stats "github.com/micro/micro/debug/stats/proto"
)

// New initialises and returns a new Stats service handler
func New(done <-chan bool) (*Stats, error) {
	s := &Stats{
		registry: cache.New(*cmd.DefaultOptions().Registry),
		client:   *cmd.DefaultOptions().Client,
	}

	if err := s.scan(); err != nil {
		return nil, err
	}

	s.Start(done)
	return s, nil
}

// Stats is the Debug.Stats handler
type Stats struct {
	registry registry.Registry
	client   client.Client

	sync.RWMutex
	snapshots []*stats.Snapshot
	cached    []*registry.Service
}

// Read returns gets a snapshot of all current stats
func (s *Stats) Read(ctx context.Context, req *stats.ReadRequest, rsp *stats.ReadResponse) error {
	if req.Service == nil {
		s.RLock()
		rsp.Stats = s.snapshots
		s.RUnlock()
		return nil
	}

	filter := func(a, b string) bool {
		if len(b) == 0 {
			return true
		}
		return a == b
	}

	s.RLock()
	for _, s := range s.snapshots {
		if !filter(s.Service.Name, req.Service.Name) {
			continue
		}
		if !filter(s.Service.Version, req.Service.Version) {
			continue
		}
		// append snapshot
		rsp.Stats = append(rsp.Stats, s)
	}
	s.RUnlock()

	return nil
}

func (s *Stats) Write(ctx context.Context, req *stats.WriteRequest, rsp *stats.WriteResponse) error {
	return errors.BadRequest("go.micro.debug.stats", "not implemented")
}

// Stream starts streaming stats
func (s *Stats) Stream(ctx context.Context, req *stats.StreamRequest, rsp stats.Stats_StreamStream) error {
	return errors.BadRequest("go.micro.debug.stats", "not implemented")
}

// Start Starts scraping other services until the provided channel is closed
func (s *Stats) Start(done <-chan bool) {
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				s.scrape()
				time.Sleep(time.Second)
			}
		}
	}()

	go func() {
		t := time.NewTicker(10 * time.Second)
		defer t.Stop()

		for {
			select {
			case <-done:
				return
			case <-t.C:
				if err := s.scan(); err != nil {
					log.Debug(err)
				}
			}
		}
	}()
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
	var mtx sync.Mutex
	next := make([]*stats.Snapshot, 0)

	// Call each node of each service in goroutines
	var wg sync.WaitGroup

	protocol := s.client.String()
	transport := s.client.Options().Transport.String()

	for _, svc := range services {
		// Ignore nodeless and non mucp services
		if len(svc.Nodes) == 0 {
			continue
		}
		// Call every node
		for _, node := range svc.Nodes {
			if node.Metadata["protocol"] != protocol {
				continue
			}
			if node.Metadata["transport"] != transport {
				continue
			}

			wg.Add(1)

			go func(st *Stats, service *registry.Service, node *registry.Node) {
				defer wg.Done()

				// create new context to cancel within a few seconds
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()

				req := s.client.NewRequest(service.Name, "Debug.Stats", &debug.StatsRequest{})
				rsp := new(debug.StatsResponse)
				if err := s.client.Call(ctx, req, rsp, client.WithAddress(node.Address)); err != nil {
					log.Errorf("Error calling %s@%s (%s)", service.Name, node.Address, err.Error())
					return
				}

				// Append the new snapshot
				snap := &stats.Snapshot{
					Service: &stats.Service{
						Name:    service.Name,
						Version: service.Version,
						Node: &stats.Node{
							Id:      node.Id,
							Address: node.Address,
						},
					},
					Started: int64(rsp.Started),
					Uptime:  rsp.Uptime,
					Memory:  rsp.Memory,
					Threads: rsp.Threads,
					Gc:      rsp.Gc,
				}

				mtx.Lock()
				next = append(next, snap)
				mtx.Unlock()
			}(s, svc, node)
		}
	}
	wg.Wait()

	// Swap in the snapshots
	s.Lock()
	s.snapshots = next
	s.Unlock()
}
