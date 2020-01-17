package runtime

import (
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro/runtime"
	"github.com/micro/go-micro/store"
	"github.com/micro/go-micro/util/log"
)

type manager struct {
	Runtime runtime.Runtime
	Store   store.Store

	sync.RWMutex
	// internal cache of services
	services map[string]*runtimeService

	running bool
	exit    chan bool
	// env to inject into the service
	// TODO: use profiles not env vars
	env []string
}

// stored in store
type runtimeService struct {
	Service *runtime.Service       `json:"service"`
	Options *runtime.CreateOptions `json:"options"`
}

var (
	eventTick = time.Second * 10
)

func key(s *runtime.Service) string {
	return s.Name + ":" + s.Version
}

func (m *manager) Init(opts ...runtime.Option) error {
	return nil
}

func (m *manager) Create(s *runtime.Service, opts ...runtime.CreateOption) error {
	m.Lock()
	defer m.Unlock()

	// we need to parse options to get the env vars
	var options runtime.CreateOptions
	for _, o := range opts {
		o(&options)
	}

	// setup the runtime env
	opts = append(opts, runtime.WithEnv(append(options.Env, m.env...)))

	if s.Metadata == nil {
		s.Metadata = make(map[string]string)
	}

	// make copy
	cp := new(runtime.Service)
	cp.Name = s.Name
	cp.Version = s.Version
	cp.Source = s.Source
	cp.Metadata = make(map[string]string)

	for k, v := range s.Metadata {
		cp.Metadata[k] = v
	}

	// create the service
	if err := m.Runtime.Create(cp, opts...); err != nil {
		return err
	}

	// set the initial status
	s.Metadata["status"] = "started"

	// create service key
	k := key(s)

	rs := &runtimeService{
		Service: s,
		Options: &options,
	}

	// save
	m.services[k] = rs

	// marshall the content
	b, err := json.Marshal(rs)
	if err != nil {
		return err
	}

	// save the record
	if err := m.Store.Write(&store.Record{
		Key:   k,
		Value: b,
	}); err != nil {
		return err
	}

	return nil
}

func (m *manager) Read(opts ...runtime.ReadOption) ([]*runtime.Service, error) {
	var options runtime.ReadOptions
	for _, o := range opts {
		o(&options)
	}

	var services []*runtime.Service

	m.RLock()
	defer m.RUnlock()

	for _, rs := range m.services {
		srv := options.Service
		ver := options.Version

		if len(srv) > 0 && rs.Service.Name != srv {
			continue
		}

		if len(ver) > 0 && rs.Service.Version != ver {
			continue
		}

		// TODO: service type options.Type
		cp := new(runtime.Service)
		*cp = *rs.Service
		services = append(services, cp)
	}

	return services, nil
}

func (m *manager) Update(s *runtime.Service) error {
	// create the service
	if err := m.Runtime.Update(s); err != nil {
		return err
	}

	k := key(s)
	// read the existing record
	r, err := m.Store.Read(k)
	if err != nil {
		return err
	}

	if len(r) == 0 {
		return errors.New("service not found")
	}

	var rs runtimeService
	if err := json.Unmarshal(r[0].Value, &rs); err != nil {
		return err
	}
	// set the service
	rs.Service = s
	// TODO: allow setting opts

	// marshall the content
	b, err := json.Marshal(rs)
	if err != nil {
		return err
	}

	// save the record
	return m.Store.Write(&store.Record{
		Key:   k,
		Value: b,
	})
}

func (m *manager) Delete(s *runtime.Service) error {
	m.Lock()
	defer m.Unlock()

	k := key(s)

	// save local status
	v, ok := m.services[k]
	if ok {
		v.Service.Metadata["status"] = "stopping"
	}

	// delete from runtime
	if err := m.Runtime.Delete(s); err != nil {
		v.Service.Metadata["status"] = "error"
		v.Service.Metadata["error"] = err.Error()
		return err
	}

	// delete from store
	return m.Store.Delete(k)
}

func (m *manager) List() ([]*runtime.Service, error) {
	m.RLock()
	defer m.RUnlock()

	services := make([]*runtime.Service, 0, len(m.services))

	for _, service := range m.services {
		cp := new(runtime.Service)
		*cp = *service.Service
		services = append(services, cp)
	}

	return services, nil
}

// TODO: watch events rather than poll
func (m *manager) run() {
	//
	t := time.NewTicker(eventTick)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			// list the keys from store
			records, err := m.Store.List()
			if err != nil {
				continue
			}

			// list whats already runnning
			services, err := m.Runtime.List()
			if err != nil {
				continue
			}

			// generate service map of running things
			running := make(map[string]*runtime.Service)

			for _, service := range services {
				k := key(service)
				running[k] = service
			}

			// create a map of services that should actually run
			shouldRun := make(map[string]*runtimeService)

			// iterate through and see what we need to run
			for _, record := range records {
				// decode the record
				var rs *runtimeService
				if err := json.Unmarshal(record.Value, &rs); err != nil {
					continue
				}

				// things to run
				shouldRun[record.Key] = rs

				// check if its already running
				if v, ok := running[record.Key]; ok {
					log.Logf("%v %v", v.Metadata, rs.Service.Metadata)

					// set the status
					storeStatus := rs.Service.Metadata["status"]
					runningStatus := v.Metadata["status"]

					if storeStatus != "running" {
						rs.Service.Metadata["status"] = "running"
						// write the updated status
						b, _ := json.Marshal(rs)
						record.Value = b
						// store the record
						m.Store.Write(record)
					} else if storeStatus != runningStatus {
						log.Logf("Setting %v %v", v.Metadata, rs.Service.Metadata)
						rs.Service.Metadata["status"] = v.Metadata["status"]
						if len(v.Metadata["error"]) > 0 {
							rs.Service.Metadata["error"] = v.Metadata["error"]
						}
						// write the updated status
						b, _ := json.Marshal(rs)
						record.Value = b
						// store the record
						m.Store.Write(record)
					}

					continue
				}

				opts := []runtime.CreateOption{
					runtime.WithCommand(rs.Options.Command...),
					runtime.WithEnv(rs.Options.Env),
					runtime.CreateType(rs.Options.Type),
				}

				log.Logf("Creating service %s version %s source %s", rs.Service.Name, rs.Service.Version, rs.Service.Source)

				// set the status to starting
				rs.Service.Metadata["status"] = "starting"

				// service does not exist so start it
				if err := m.Runtime.Create(rs.Service, opts...); err != nil {
					log.Logf("Erroring running %s: %v", rs.Service.Name, err)

					// an error is already recorded
					if rs.Service.Metadata["status"] == "error" {
						continue
					}

					// save the error
					rs.Service.Metadata["status"] = "error"
					rs.Service.Metadata["error"] = err.Error()
				}

				// write the updated status
				b, _ := json.Marshal(rs)
				record.Value = b
				// store the record
				m.Store.Write(record)
			}

			// check what we need to stop from the running list
			for _, service := range services {
				k := key(service)

				// check if it should be running
				if _, ok := shouldRun[k]; ok {
					continue
				}

				log.Logf("Stopping %s", k)

				// should not be running
				m.Runtime.Delete(service)
			}

			// save the current list of running things
			m.services = shouldRun
		case <-m.exit:
			return
		}
	}
}

func (m *manager) Start() error {
	m.Lock()
	defer m.Unlock()

	if m.running {
		return nil
	}

	// reset the exit channel
	m.exit = make(chan bool)

	// start the runtime
	if err := m.Runtime.Start(); err != nil {
		return err
	}

	// start the internal manager
	go m.run()

	// set to running
	m.running = true

	return nil
}

func (m *manager) Stop() error {
	m.Lock()
	defer m.Unlock()

	if !m.running {
		return nil
	}

	select {
	case <-m.exit:
		return nil
	default:
		close(m.exit)
		m.Runtime.Stop()
		m.running = false
	}

	return nil
}

func newManager(ctx *cli.Context, r runtime.Runtime, s store.Store) *manager {
	var env []string
	// peel out the env
	for _, ev := range ctx.StringSlice("env") {
		for _, val := range strings.Split(ev, ",") {
			env = append(env, strings.TrimSpace(val))
		}
	}

	return &manager{
		Runtime:  r,
		Store:    s,
		env:      env,
		services: make(map[string]*runtimeService),
		exit:     make(chan bool),
	}
}
