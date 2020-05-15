package manager

import (
	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/go-micro/v2/runtime"
	"github.com/micro/go-micro/v2/store"
	"github.com/micro/go-micro/v2/store/memory"
	"github.com/micro/micro/v2/internal/namespace"
)

// Init initializes the runtime
func (m *manager) Init(...runtime.Option) error {
	return nil
}

// Create registers a service
func (m *manager) Create(srv *runtime.Service, opts ...runtime.CreateOption) error {
	// parse the options
	var options runtime.CreateOptions
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
	return m.publishEvent(eventTypeCreated, srv, &options)
}

// Read returns the service which matches the criteria provided
func (m *manager) Read(opts ...runtime.ReadOption) ([]*runtime.Service, error) {
	// parse the options
	var options runtime.ReadOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.Namespace) == 0 {
		options.Namespace = namespace.DefaultNamespace
	}

	// query the store. TODO: query by type? (it isn't an attr of srv)
	return m.readServices(options.Namespace, &runtime.Service{
		Name:    options.Service,
		Version: options.Version,
	})
}

// Update the service in place
func (m *manager) Update(srv *runtime.Service, opts ...runtime.UpdateOption) error {
	// parse the options
	var options runtime.UpdateOptions
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
	return m.publishEvent(eventTypeUpdated, srv, &runtime.CreateOptions{Namespace: options.Namespace})
}

// Remove a service
func (m *manager) Delete(srv *runtime.Service, opts ...runtime.DeleteOption) error {
	// parse the options
	var options runtime.DeleteOptions
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
	return m.publishEvent(eventTypeDeleted, srv, &runtime.CreateOptions{Namespace: options.Namespace})
}

// Starts the manager
func (m *manager) Start() error {
	if m.running {
		return nil
	}
	m.running = true

	// start the runtime we're going to manage
	if err := m.Runtime.Start(); err != nil {
		return err
	}

	// watch events written to the store
	go m.watchEvents()

	// todo: periodically load the status of services from the runtime

	// todo: compare the store to the runtime incase we missed any events

	return nil
}

// Stop the manager
func (m *manager) Stop() error {
	if !m.running {
		return nil
	}
	m.running = false

	return m.Runtime.Stop()
}

// String describes runtime
func (m *manager) String() string {
	return "manager"
}

type manager struct {
	// runtime being managed
	runtime.Runtime
	// options passed by the caller
	options Options
	// running is true after Start is called
	running bool
	// eventsConsumed is a memory store which contains the
	// ids of all events recently consumed by the manager
	eventsConsumed store.Store
}

// New returns a manager for the runtime
func New(r runtime.Runtime, opts ...Option) runtime.Runtime {
	// parse the options
	var options Options
	for _, o := range opts {
		o(&options)
	}

	// set the defaults
	if options.Store == nil {
		options.Store = *cmd.DefaultCmd.Options().Store
	}

	return &manager{
		Runtime:        r,
		options:        options,
		eventsConsumed: memory.NewStore(),
	}
}
