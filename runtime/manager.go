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
	mprofile "github.com/micro/micro/runtime/profile"
)

type manager struct {
	Runtime runtime.Runtime
	Store   store.Store

	sync.RWMutex
	// internal cache of services
	services map[string]*runtimeService

	running bool
	exit    chan bool
	// used to propagate events
	events chan *event

	// a runtime profile to set for the service
	profile []string
}

// stored in store
type runtimeService struct {
	Service *runtime.Service       `json:"service"`
	Options *runtime.CreateOptions `json:"options"`
	Status  string                 `json:"status"`
	Error   error                  `json:"error"`
}

type event struct {
	Type    string
	Service *runtime.Service
	Options *runtime.CreateOptions
}

var (
	eventTick = time.Second * 10
)

func copyService(s *runtimeService) *runtime.Service {
	cp := new(runtime.Service)
	cp.Name = s.Service.Name
	cp.Version = s.Service.Version
	cp.Source = s.Service.Source
	cp.Metadata = make(map[string]string)
	for k, v := range s.Service.Metadata {
		cp.Metadata[k] = v
	}
	cp.Metadata["status"] = s.Status
	if s.Error != nil {
		cp.Metadata["error"] = s.Error.Error()
	}
	return cp
}

func key(s *runtime.Service) string {
	return s.Name + ":" + s.Version
}

func (m *manager) sendEvent(ev *event) {
	m.events <- ev
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

	if s.Metadata == nil {
		s.Metadata = make(map[string]string)
	}

	// create service key
	k := key(s)

	rs := &runtimeService{
		Service: s,
		Options: &options,
		Status:  "started",
	}

	// save locally
	m.services[k] = rs

	// send event
	go m.sendEvent(&event{
		Type:    "create",
		Service: s,
		Options: &options,
	})

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

		services = append(services, copyService(rs))
	}

	return services, nil
}

func (m *manager) Update(s *runtime.Service) error {
	m.Lock()
	defer m.Unlock()

	// create the key
	k := key(s)

	// read the existing record
	r, err := m.Store.Read(k)
	if err != nil {
		return err
	}

	// no service
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

	// if not running then run it
	evType := "update"

	// check if it exists
	if _, ok := m.services[k]; !ok {
		// set starting status
		rs.Status = "started"
		evType = "create"
		m.services[k] = &rs
	}

	// fire an update
	go m.sendEvent(&event{
		Type:    evType,
		Service: rs.Service,
		Options: rs.Options,
	})

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
	if !ok {
		return nil
	}

	// set status
	v.Status = "stopped"

	// send event
	go m.sendEvent(&event{
		Type:    "delete",
		Service: v.Service,
	})

	// delete from store
	return m.Store.Delete(k)
}

func (m *manager) List() ([]*runtime.Service, error) {
	m.RLock()
	defer m.RUnlock()

	services := make([]*runtime.Service, 0, len(m.services))

	for _, service := range m.services {
		services = append(services, copyService(service))
	}

	return services, nil
}

func (m *manager) runtimeEnv(options *runtime.CreateOptions) []string {
	setEnv := func(p []string, env map[string]string) {
		for _, v := range p {
			parts := strings.Split(v, "=")
			if len(parts) <= 1 {
				continue
			}
			env[parts[0]] = strings.Join(parts[1:], "=")
		}
	}

	// overwrite any values
	env := map[string]string{}

	// set the env vars provided
	setEnv(options.Env, env)

	// override with vars from the profile
	setEnv(m.profile, env)

	// create a new env
	var vars []string
	for k, v := range env {
		vars = append(vars, k+"="+v)
	}

	// setup the runtime env
	return vars
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
				log.Logf("Failed to list records from store: %v", err)
				continue
			}

			// list whats already runnning
			services, err := m.Runtime.List()
			if err != nil {
				log.Logf("Failed to list runtime services: %v", err)
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
					// TODO: have actual runtime status
					rs.Status = v.Metadata["status"]
					if e := v.Metadata["error"]; len(e) > 0 {
						rs.Error = errors.New(e)
					}
					continue
				}

				// generate the runtime environment
				env := m.runtimeEnv(rs.Options)

				// create a new set of options to use
				opts := []runtime.CreateOption{
					runtime.WithCommand(rs.Options.Command...),
					runtime.WithEnv(env),
					runtime.CreateType(rs.Options.Type),
				}

				log.Logf("Creating service %s version %s source %s", rs.Service.Name, rs.Service.Version, rs.Service.Source)

				// set the status to starting
				rs.Status = "started"

				// service does not exist so start it
				if err := m.Runtime.Create(rs.Service, opts...); err != nil {
					if err != runtime.ErrAlreadyExists {
						log.Logf("Erroring running %s: %v", rs.Service.Name, err)

						// save the error
						rs.Status = "error"
						rs.Error = err
					}
				}
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
		case ev := <-m.events:
			var err error

			switch ev.Type {
			case "delete":
				log.Logf("Deleting %s %s", ev.Service.Name, ev.Service.Version)
				err = m.Runtime.Delete(ev.Service)
			case "update":
				log.Logf("Updating %s %s", ev.Service.Name, ev.Service.Version)
				err = m.Runtime.Update(ev.Service)
			case "create":
				opts := []runtime.CreateOption{
					runtime.WithCommand(ev.Options.Command...),
					runtime.WithEnv(ev.Options.Env),
					runtime.CreateType(ev.Options.Type),
				}

				log.Logf("Creating %s %s", ev.Service.Name, ev.Service.Version)
				err = m.Runtime.Create(ev.Service, opts...)
			}

			if err != nil {
				log.Logf("Erroring executing event %s for %s: %v", ev.Type, ev.Service.Name, err)

				// save the error
				// hacking, its a pointer
				m.Lock()
				v, ok := m.services[key(ev.Service)]
				if ok {
					v.Status = "error"
					v.Error = err
				}
				m.Unlock()
			}
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
	var profile []string
	// peel out the env
	switch ctx.String("profile") {
	case "platform":
		profile = mprofile.Platform()
	}

	return &manager{
		Runtime:  r,
		Store:    s,
		profile:  profile,
		services: make(map[string]*runtimeService),
		exit:     make(chan bool),
		events:   make(chan *event, 8),
	}
}
