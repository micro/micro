// Package handler is the handler for the `micro debug trace` service
package handler

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/config/cmd"
	debug "github.com/micro/go-micro/v2/debug/service/proto"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/cache"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/go-micro/v2/util/ring"
	trace "github.com/micro/micro/v2/debug/trace/proto"
)

// New initialises and returns a new trace service handler
func New(done <-chan bool, windowSize int) (*Trace, error) {
	s := &Trace{
		registry:            cache.New(*cmd.DefaultOptions().Registry),
		client:              *cmd.DefaultOptions().Client,
		historicalSnapshots: ring.New(windowSize),
	}

	if err := s.scan(); err != nil {
		return nil, err
	}

	s.Start(done)
	return s, nil
}

// trace is the Debug.trace handler
type Trace struct {
	registry registry.Registry
	client   client.Client

	sync.RWMutex
	// current snapshots for each service
	snapshots []*trace.Snapshot
	// historical snapshots from the start
	historicalSnapshots *ring.Buffer
	cached              []*registry.Service
}

// Filters out all spans that are part of a trace that hits a given service.
func filterServiceSpans(service string, snapshots []*trace.Snapshot) []*trace.Span {
	// trace id -> span id -> span
	groupByTrace := map[string]map[string]*trace.Span{}
	for _, snapshot := range snapshots {
		for _, span := range snapshot.GetSpans() {
			_, ok := groupByTrace[span.GetTrace()]
			if !ok {
				groupByTrace[span.GetTrace()] = map[string]*trace.Span{}
			}
			groupByTrace[span.GetTrace()][span.GetId()] = span
		}
	}
	ret := []*trace.Span{}
	for _, spanMap := range groupByTrace {
		spans := []*trace.Span{}
		shouldAppend := false
		for _, span := range spanMap {
			spans = append(spans, span)
			if strings.Contains(span.GetName(), service) {
				shouldAppend = true
			}
			if shouldAppend {
				ret = append(ret, spans...)
			}
		}
	}
	return ret
}

func dedupeSpans(spans []*trace.Span) []*trace.Span {
	m := map[string]*trace.Span{}
	for _, span := range spans {
		m[span.GetId()] = span
	}
	ret := []*trace.Span{}
	for _, span := range m {
		ret = append(ret, span)
	}
	return ret
}

func snapshotsToSpans(snapshots []*trace.Snapshot) []*trace.Span {
	ret := []*trace.Span{}
	for _, snapshot := range snapshots {
		ret = append(ret, snapshot.GetSpans()...)
	}
	return ret
}

// Read returns gets a snapshot of all current trace3
func (s *Trace) Read(ctx context.Context, req *trace.ReadRequest, rsp *trace.ReadResponse) error {
	allSnapshots := []*trace.Snapshot{}
	func() {
		s.RLock()
		defer s.RUnlock()
		if req.Past {
			entries := s.historicalSnapshots.Get(3600)
			for _, entry := range entries {
				allSnapshots = append(allSnapshots, entry.Value.([]*trace.Snapshot)...)
			}
		} else {
			// Using an else since the latest snapshot is already in the ring buffer
			allSnapshots = append(allSnapshots, s.snapshots...)
		}
	}()
	if req.Service == nil {
		rsp.Spans = dedupeSpans(snapshotsToSpans(allSnapshots))
		return nil
	}
	spans := filterServiceSpans(req.GetService().GetName(), allSnapshots)
	if req.GetLimit() == 0 {
		rsp.Spans = spans
	} else {
		lim := req.GetLimit()
		if lim >= int64(len(spans)) {
			lim = int64(len(spans))
		}
		rsp.Spans = spans[0:lim]
	}
	return nil
}

func (s *Trace) Write(ctx context.Context, req *trace.WriteRequest, rsp *trace.WriteResponse) error {
	return errors.BadRequest("go.micro.debug.trace", "not implemented")
}

// Stream starts streaming trace
func (s *Trace) Stream(ctx context.Context, req *trace.StreamRequest, rsp trace.Trace_StreamStream) error {
	return errors.BadRequest("go.micro.debug.trace", "not implemented")
}

// Start Starts scraping other services until the provided channel is closed
func (s *Trace) Start(done <-chan bool) {
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				s.scrape()
				time.Sleep(5 * time.Second)
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
					log.Error(err)
				}
			}
		}
	}()
}

func (s *Trace) scan() error {
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

func (s *Trace) scrape() {
	s.RLock()
	// Create a local copy of cached services
	services := make([]*registry.Service, len(s.cached))
	copy(services, s.cached)
	s.RUnlock()

	// Start building the next list of snapshots
	var mtx sync.Mutex
	next := make([]*trace.Snapshot, 0)

	// Call each node of each service in goroutines
	var wg sync.WaitGroup

	protocol := s.client.String()

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

			wg.Add(1)

			go func(st *Trace, service *registry.Service, node *registry.Node) {
				defer wg.Done()

				// create new context to cancel within a few seconds
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				req := s.client.NewRequest(service.Name, "Debug.Trace", &debug.TraceResponse{})
				rsp := new(debug.TraceResponse)
				if err := s.client.Call(ctx, req, rsp, client.WithAddress(node.Address)); err != nil {
					log.Errorf("Error calling %s@%s (%s)", service.Name, node.Address, err.Error())
					return
				}
				spans := []*trace.Span{}
				for _, v := range rsp.GetSpans() {
					spans = append(spans, &trace.Span{
						Trace:    v.GetTrace(),
						Id:       v.GetId(),
						Parent:   v.GetParent(),
						Name:     v.GetName(),
						Started:  v.GetStarted(),
						Duration: v.GetDuration(),
						Metadata: v.GetMetadata(),
					})
				}
				// Append the new snapshot
				snap := &trace.Snapshot{
					Service: &trace.Service{
						Name:    service.Name,
						Version: service.Version,
						Node: &trace.Node{
							Id:      node.Id,
							Address: node.Address,
						},
					},
					Spans: spans,
				}
				//timestamp := time.Now().Unix()
				// snap.Timestamp = uint64(timestamp)
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
	s.historicalSnapshots.Put(next)
	s.Unlock()
}
