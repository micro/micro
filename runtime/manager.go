package runtime

import (
	"encoding/json"
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
	services map[string]*runtime.Service

	running bool
	exit    chan bool
	// env to inject into the service
	// TODO: use profiles not env vars
	env []string
}

var (
	eventTick = time.Second * 10
)

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

	// set the initial status
	s.Metadata["status"] = "started"

	// make copy
	cp := new(runtime.Service)
	*cp = *s
	for k, v := range s.Metadata {
		cp.Metadata[k] = v
	}

	// create the service
	err := m.Runtime.Create(cp, opts...)
	if err != nil {
		s.Metadata["status"] = "error"
		s.Metadata["error"] = err.Error()
	}

	// save
	m.services[s.Name+s.Version] = s

	// marshall the content
	b, _ := json.Marshal(s)

	// save the record
	err = m.Store.Write(&store.Record{
		Key:   s.Name + s.Version,
		Value: b,
	})
	if err != nil {
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

	for _, service := range m.services {
		srv := options.Service
		ver := options.Version

		if len(srv) > 0 && service.Name != srv {
			continue
		}

		if len(ver) > 0 && service.Version != ver {
			continue
		}

		// TODO: service type options.Type
		cp := new(runtime.Service)
		*cp = *service
		services = append(services, cp)
	}

	return services, nil
}

func (m *manager) Update(s *runtime.Service) error {
	// create the service
	if err := m.Runtime.Update(s); err != nil {
		return err
	}

	// marshall the content
	b, _ := json.Marshal(s)

	// save the record
	return m.Store.Write(&store.Record{
		Key:   s.Name + s.Version,
		Value: b,
	})
}

func (m *manager) Delete(s *runtime.Service) error {
	m.Lock()
	defer m.Unlock()

	// save local status
	v, ok := m.services[s.Name+s.Version]
	if ok {
		v.Metadata["status"] = "stopping"
	}

	// delete from runtime
	if err := m.Runtime.Delete(s); err != nil {
		v.Metadata["status"] = "error"
		v.Metadata["error"] = err.Error()
		return err
	}

	// delete from store
	return m.Store.Delete(s.Name + s.Version)
}

func (m *manager) List() ([]*runtime.Service, error) {
	m.RLock()
	defer m.RUnlock()

	services := make([]*runtime.Service, 0, len(m.services))

	for _, service := range m.services {
		cp := new(runtime.Service)
		*cp = *service
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
				running[service.Name+service.Version] = service
			}

			// create a map of services that should actually run
			shouldRun := make(map[string]*runtime.Service)

			// iterate through and see what we need to run
			for _, record := range records {
				// decode the record
				var service *runtime.Service
				if err := json.Unmarshal(record.Value, &service); err != nil {
					continue
				}

				// things to run
				shouldRun[record.Key] = service

				// check if its already running
				if v, ok := running[service.Name+service.Version]; ok {
					// if the runtime status is different use that
					if v.Metadata["status"] != service.Metadata["status"] {
						service.Metadata["status"] = v.Metadata["status"]
						if len(v.Metadata["error"]) > 0 {
							service.Metadata["error"] = v.Metadata["error"]
						}
					}
					continue
				}

				// service does not exist so start it
				if err := m.Runtime.Create(service); err != nil {
					log.Logf("Erroring running %s: %v", service.Name, err)

					// an error is already recorded
					if len(service.Metadata["error"]) > 0 {
						continue
					}

					// save the error
					service.Metadata["status"] = "error"
					service.Metadata["error"] = err.Error()
					b, _ := json.Marshal(service)
					record.Value = b
					// store the record
					m.Store.Write(record)
				}
			}

			// check what we need to stop from the running list
			for _, service := range services {
				// check if it should be running
				if _, ok := shouldRun[service.Name+service.Version]; ok {
					continue
				}

				log.Logf("Stopping service %s version %s", service.Name, service.Version)

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
		services: make(map[string]*runtime.Service),
		exit:     make(chan bool),
	}
}
