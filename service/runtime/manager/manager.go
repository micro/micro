package manager

import (
	gorun "github.com/micro/go-micro/v3/runtime"
	"github.com/micro/go-micro/v3/store"
	cachest "github.com/micro/go-micro/v3/store/cache"
	filest "github.com/micro/go-micro/v3/store/file"
	"github.com/micro/go-micro/v3/store/memory"
	"github.com/micro/micro/v3/internal/namespace"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/runtime"
)

// Init initializes the runtime
func (m *manager) Init(...gorun.Option) error {
	return nil
}

// Create registers a service
func (m *manager) Create(srv *gorun.Service, opts ...gorun.CreateOption) error {
	// parse the options
	var options gorun.CreateOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.Namespace) == 0 {
		options.Namespace = namespace.DefaultNamespace
	}

	// set defaults
	if srv.Metadata == nil {
		srv.Metadata = make(map[string]string)
	}
	if len(srv.Version) == 0 {
		srv.Version = "latest"
	}

	// write the object to the store
	if err := m.createService(srv, &options); err != nil {
		return err
	}

	// publish the event, this will apply it aysnc to the runtime
	return m.publishEvent(gorun.Create, srv, &options)
}

// Read returns the service which matches the criteria provided
func (m *manager) Read(opts ...gorun.ReadOption) ([]*gorun.Service, error) {
	// parse the options
	var options gorun.ReadOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.Namespace) == 0 {
		options.Namespace = namespace.DefaultNamespace
	}

	// query the store. TODO: query by type? (it isn't an attr of srv)
	srvs, err := m.readServices(options.Namespace, &gorun.Service{
		Name:    options.Service,
		Version: options.Version,
	})
	if err != nil {
		return nil, err
	}
	// add the metadata to the service from the local runtime (e.g. status, err)
	statuses, err := m.listStatuses(options.Namespace)
	if err != nil {
		return nil, err
	}
	ret := []*gorun.Service{}
	for _, srv := range srvs {
		ret = append(ret, srv.Service)
		md, ok := statuses[srv.Service.Name+":"+srv.Service.Version]
		if !ok {
			continue
		}
		srv.Service.Metadata["status"] = md.Status
		srv.Service.Metadata["error"] = md.Error
	}

	return ret, nil
}

// Update the service in place
func (m *manager) Update(srv *gorun.Service, opts ...gorun.UpdateOption) error {
	// parse the options
	var options gorun.UpdateOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.Namespace) == 0 {
		options.Namespace = namespace.DefaultNamespace
	}

	// set defaults
	if len(srv.Version) == 0 {
		srv.Version = "latest"
	}

	// publish the update event which will trigger an update in the runtime
	return m.publishEvent(gorun.Update, srv, &gorun.CreateOptions{Namespace: options.Namespace})
}

// Remove a service
func (m *manager) Delete(srv *gorun.Service, opts ...gorun.DeleteOption) error {
	// parse the options
	var options gorun.DeleteOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.Namespace) == 0 {
		options.Namespace = namespace.DefaultNamespace
	}

	// set defaults
	if len(srv.Version) == 0 {
		srv.Version = "latest"
	}

	// delete from the store
	if err := m.deleteService(options.Namespace, srv); err != nil {
		return err
	}

	// publish the event which will trigger a delete in the runtime
	return m.publishEvent(gorun.Delete, srv, &gorun.CreateOptions{Namespace: options.Namespace})
}

// Starts the manager
func (m *manager) Start() error {
	if m.running {
		return nil
	}
	m.running = true

	// start the runtime we're going to manage
	if err := runtime.DefaultRuntime.Start(); err != nil {
		return err
	}

	// watch events written to the store
	go m.watchEvents()

	// periodically load the status of services from the runtime
	go m.watchStatus()

	// todo: compare the store to the runtime incase we missed any events

	// Resurrect services that were running previously
	go m.resurrectServices()

	return nil
}

// Logs for a service
func (m *manager) Logs(srv *gorun.Service, opts ...gorun.LogsOption) (gorun.LogStream, error) {
	return runtime.Logs(srv, opts...)
}

func (m *manager) resurrectServices() {
	nss, err := m.listNamespaces()
	if err != nil {
		logger.Warnf("Error listing namespaces: %v", err)
		return
	}

	for _, ns := range nss {
		srvs, err := m.readServices(ns, &gorun.Service{})
		if err != nil {
			logger.Warnf("Error reading services from the %v namespace: %v", ns, err)
			return
		}

		running := map[string]*gorun.Service{}
		curr, _ := runtime.Read(gorun.ReadNamespace(ns))
		for _, v := range curr {
			running[v.Name+":"+v.Version+":"+v.Source] = v
		}

		for _, srv := range srvs {
			if _, ok := running[srv.Service.Name+":"+srv.Service.Version+":"+srv.Service.Source]; ok {
				// already running, don't need to start again
				continue
			}

			// generate an auth account for the service to use
			acc, err := m.generateAccount(srv.Service, ns)
			if err != nil {
				continue
			}

			runtime.Create(srv.Service,
				gorun.CreateImage(srv.Options.Image),
				gorun.CreateType(srv.Options.Type),
				gorun.CreateNamespace(ns),
				gorun.WithArgs(srv.Options.Args...),
				gorun.WithCommand(srv.Options.Command...),
				gorun.WithEnv(m.runtimeEnv(srv.Service, srv.Options)),
				gorun.CreateCredentials(acc.ID, acc.Secret),
			)
		}
	}
}

// Stop the manager
func (m *manager) Stop() error {
	if !m.running {
		return nil
	}
	m.running = false

	return runtime.DefaultRuntime.Stop()
}

// String describes runtime
func (m *manager) String() string {
	return "manager"
}

type manager struct {
	// running is true after Start is called
	running bool
	// cache is a memory store which is used to store any information we don't want to write to the
	// global store, e.g. service status / errors (these will change depending on the
	// managed runtime and hence won't be the same globally).
	cache store.Store
	// fileCache is a cache store used to store any information we don't want to write to the
	// global store but want to persist across restarts, e.g. events consumed
	fileCache store.Store
}

// New returns a manager for the runtime
func New() gorun.Runtime {
	return &manager{
		cache:     memory.NewStore(),
		fileCache: cachest.NewStore(filest.NewStore()),
	}
}
