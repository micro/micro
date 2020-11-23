package manager

import (
	"time"

	"github.com/micro/micro/v3/internal/kubernetes/client"
	"github.com/micro/micro/v3/internal/namespace"
	"github.com/micro/micro/v3/service/build"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/runtime"
)

// Init initializes the runtime
func (m *manager) Init(...runtime.Option) error {
	return nil
}

// Create a resource
func (m *manager) Create(resource runtime.Resource, opts ...runtime.CreateOption) error {

	// parse the options
	var options runtime.CreateOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.Namespace) == 0 {
		options.Namespace = namespace.DefaultNamespace
	}

	// Handle the various different types of resources:
	switch resource.Type() {
	case runtime.TypeNamespace:

		// Assert the resource back into a *runtime.Namespace
		namespace, ok := resource.(*runtime.Namespace)
		if !ok {
			return runtime.ErrInvalidResource
		}

		// Allow the options to take precedence
		if options.Namespace != "" {
			namespace.Name = options.Namespace
		}

		// Do we need to store this locally?
		return runtime.DefaultRuntime.Create(namespace)

	case runtime.TypeNetworkPolicy:

		// Assert the resource back into a *runtime.NetworkPolicy
		networkPolicy, ok := resource.(*runtime.NetworkPolicy)
		if !ok {
			return runtime.ErrInvalidResource
		}

		// Allow the options to take precedence
		if options.Namespace != "" {
			networkPolicy.Namespace = options.Namespace
		}

		// Do we need to store this locally?
		return runtime.DefaultRuntime.Create(networkPolicy)

	case runtime.TypeResourceQuota:

		// Assert the resource back into a *runtime.ResourceQuota
		resourceQuota, ok := resource.(*runtime.ResourceQuota)
		if !ok {
			return runtime.ErrInvalidResource
		}

		// Allow the options to take precedence
		if options.Namespace != "" {
			resourceQuota.Namespace = options.Namespace
		}

		// Do we need to store this locally?
		return runtime.DefaultRuntime.Create(resourceQuota)

	case runtime.TypeService:

		// Assert the resource back into a *runtime.Service
		srv, ok := resource.(*runtime.Service)
		if !ok {
			return runtime.ErrInvalidResource
		}

		// set defaults
		if srv.Metadata == nil {
			srv.Metadata = make(map[string]string)
		}
		if len(srv.Version) == 0 {
			srv.Version = "latest"
		}

		// construct the service object
		service := &service{
			Service:   srv,
			Options:   &options,
			UpdatedAt: time.Now(),
		}

		// if there is not a build configured, start the service and then write it to the store
		if build.DefaultBuilder == nil {
			// the source could be a git remote or a reference to the blob store, parse it before we run
			// the service
			var err error
			srv.Source, err = m.checkoutSource(service)
			if err != nil {
				return err
			}

			// create the service in the underlying runtime
			if err := m.createServiceInRuntime(service); err != nil && err != runtime.ErrAlreadyExists {
				return err
			}

			// write the object to the store
			return m.writeService(service)
		}

		// building ths service can take some time so we'll write the service to the store and then
		// perform the build process async
		service.Status = runtime.Pending
		if err := m.writeService(service); err != nil {
			return err
		}

		go m.buildAndRun(service)
		return nil

	default:
		return runtime.ErrInvalidResource
	}
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

	// query the store
	srvs, err := m.readServices(options.Namespace, &runtime.Service{
		Name:    options.Service,
		Version: options.Version,
	})
	if err != nil {
		return nil, err
	}

	// query the runtime and group the resulting services by name:version so they can be queried
	rSrvs, err := m.Runtime.Read(opts...)
	if err != nil {
		return nil, err
	}
	rSrvMap := make(map[string]*runtime.Service, len(rSrvs))
	for _, s := range rSrvs {
		rSrvMap[s.Name+":"+s.Version] = s
	}

	// loop through the services returned from the store and append any info returned by the runtime
	// such as status and error
	result := make([]*runtime.Service, len(srvs))
	for i, s := range srvs {
		result[i] = s.Service

		// check for a status on the service, this could be building, stopping etc
		if s.Status != runtime.Unknown {
			result[i].Status = s.Status
		}
		if len(s.Error) > 0 {
			result[i].Metadata["error"] = s.Error
		}

		// set the last updated, todo: check why this is 'started' and not 'updated'. Consider adding
		// this as an attribute on runtime.Service
		if !s.UpdatedAt.IsZero() {
			result[i].Metadata["started"] = s.UpdatedAt.Format(time.RFC3339)
		}

		// the service might still be building and not have been created in the underlying runtime yet
		rs, ok := rSrvMap[client.Format(s.Service.Name)+":"+client.Format(s.Service.Version)]
		if !ok {
			continue
		}

		// assign the status and error. TODO: make the error an attribute on service
		result[i].Status = rs.Status
		if rs.Metadata != nil && len(rs.Metadata["error"]) > 0 {
			result[i].Metadata["status"] = rs.Metadata["error"]
		}
	}

	return result, nil
}

// Update a resource
func (m *manager) Update(resource runtime.Resource, opts ...runtime.UpdateOption) error {

	// parse the options
	var options runtime.UpdateOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.Namespace) == 0 {
		options.Namespace = namespace.DefaultNamespace
	}

	// Handle the various different types of resources:
	switch resource.Type() {
	case runtime.TypeNamespace:

		// Assert the resource back into a *runtime.Namespace
		namespace, ok := resource.(*runtime.Namespace)
		if !ok {
			return runtime.ErrInvalidResource
		}

		// Allow the options to take precedence
		if options.Namespace != "" {
			namespace.Name = options.Namespace
		}

		// Do we need to store this locally?
		return runtime.DefaultRuntime.Update(namespace)

	case runtime.TypeNetworkPolicy:

		// Assert the resource back into a *runtime.NetworkPolicy
		networkPolicy, ok := resource.(*runtime.NetworkPolicy)
		if !ok {
			return runtime.ErrInvalidResource
		}

		// Allow the options to take precedence
		if options.Namespace != "" {
			networkPolicy.Namespace = options.Namespace
		}

		// Do we need to store this locally?
		return runtime.DefaultRuntime.Update(networkPolicy)

	case runtime.TypeResourceQuota:

		// Assert the resource back into a *runtime.ResourceQuota
		resourceQuota, ok := resource.(*runtime.ResourceQuota)
		if !ok {
			return runtime.ErrInvalidResource
		}

		// Allow the options to take precedence
		if options.Namespace != "" {
			resourceQuota.Namespace = options.Namespace
		}

		// Do we need to store this locally?
		return runtime.DefaultRuntime.Update(resourceQuota)

	case runtime.TypeService:

		// Assert the resource back into a *runtime.Service
		srv, ok := resource.(*runtime.Service)
		if !ok {
			return runtime.ErrInvalidResource
		}

		// set defaults
		if len(srv.Version) == 0 {
			srv.Version = "latest"
		}

		// read the service from the store
		srvs, err := m.readServices(options.Namespace, &runtime.Service{
			Name:    srv.Name,
			Version: srv.Version,
		})
		if err != nil {
			return err
		}
		if len(srvs) == 0 {
			return runtime.ErrNotFound
		}

		// update the service
		service := srvs[0]
		service.Service.Source = srv.Source
		service.UpdatedAt = time.Now()

		// if there is not a build configured, update the service and then write it to the store
		if build.DefaultBuilder == nil {
			// the source could be a git remote or a reference to the blob store, parse it before we run
			// the service
			var err error

			service.Service.Source, err = m.checkoutSource(service)
			if err != nil {
				return err
			}

			// create the service in the underlying runtime
			if err := m.updateServiceInRuntime(service); err != nil {
				return err
			}

			// write the object to the store
			service.Status = runtime.Starting
			service.Error = ""
			return m.writeService(service)
		}

		// building ths service can take some time so we'll write the service to the store and then
		// perform the build process async
		service.Status = runtime.Pending
		if err := m.writeService(service); err != nil {
			return err
		}

		go m.buildAndUpdate(service)
		return nil

	default:
		return runtime.ErrInvalidResource
	}
}

// Delete a resource
func (m *manager) Delete(resource runtime.Resource, opts ...runtime.DeleteOption) error {

	// parse the options
	var options runtime.DeleteOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.Namespace) == 0 {
		options.Namespace = namespace.DefaultNamespace
	}

	// Handle the various different types of resources:
	switch resource.Type() {
	case runtime.TypeNamespace:

		// Assert the resource back into a *runtime.Namespace
		namespace, ok := resource.(*runtime.Namespace)
		if !ok {
			return runtime.ErrInvalidResource
		}

		// Allow the options to take precedence
		if options.Namespace != "" {
			namespace.Name = options.Namespace
		}

		// Do we need to store this locally?
		return runtime.DefaultRuntime.Delete(namespace)

	case runtime.TypeNetworkPolicy:

		// Assert the resource back into a *runtime.NetworkPolicy
		networkPolicy, ok := resource.(*runtime.NetworkPolicy)
		if !ok {
			return runtime.ErrInvalidResource
		}

		// Allow the options to take precedence
		if options.Namespace != "" {
			networkPolicy.Namespace = options.Namespace
		}

		// Do we need to store this locally?
		return runtime.DefaultRuntime.Delete(networkPolicy)

	case runtime.TypeResourceQuota:

		// Assert the resource back into a *runtime.ResourceQuota
		resourceQuota, ok := resource.(*runtime.ResourceQuota)
		if !ok {
			return runtime.ErrInvalidResource
		}

		// Allow the options to take precedence
		if options.Namespace != "" {
			resourceQuota.Namespace = options.Namespace
		}

		// Do we need to store this locally?
		return runtime.DefaultRuntime.Delete(resourceQuota)

	case runtime.TypeService:

		// Assert the resource back into a *runtime.Service
		srv, ok := resource.(*runtime.Service)
		if !ok {
			return runtime.ErrInvalidResource
		}

		// set defaults
		if len(srv.Version) == 0 {
			srv.Version = "latest"
		}

		// read the service from the store
		srvs, err := m.readServices(options.Namespace, &runtime.Service{
			Name:    srv.Name,
			Version: srv.Version,
		})
		if err != nil {
			return err
		}
		if len(srvs) == 0 {
			return runtime.ErrNotFound
		}

		// delete from the underlying runtime
		if err := m.Runtime.Delete(srv, opts...); err != nil && err != runtime.ErrNotFound {
			return err
		}

		// delete from the store
		if err := m.deleteService(srvs[0]); err != nil {
			return err
		}

		// delete the source and binary from the blob store async
		go m.cleanupBlobStore(srvs[0])
		return nil

	default:
		return runtime.ErrInvalidResource
	}
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

	// Watch services that were running previously. TODO: rename and run periodically
	go m.watchServices()

	return nil
}

// Logs for a resource
func (m *manager) Logs(resource runtime.Resource, opts ...runtime.LogsOption) (runtime.LogStream, error) {
	// Handle the various different types of resources:
	switch resource.Type() {
	case runtime.TypeService:

		// Assert the resource back into a *runtime.Service
		srv, ok := resource.(*runtime.Service)
		if !ok {
			return nil, runtime.ErrInvalidResource
		}

		return runtime.Logs(srv, opts...)
	default:
		return nil, runtime.ErrInvalidResource
	}
}

func (m *manager) watchServices() {
	nss, err := m.listNamespaces()
	if err != nil {
		logger.Warnf("Error listing namespaces: %v", err)
		return
	}

	for _, ns := range nss {
		srvs, err := m.readServices(ns, &runtime.Service{})
		if err != nil {
			logger.Warnf("Error reading services from the %v namespace: %v", ns, err)
			return
		}

		running := map[string]*runtime.Service{}
		curr, _ := runtime.Read(runtime.ReadNamespace(ns))
		for _, v := range curr {
			running[v.Name+":"+v.Version] = v
		}

		for _, srv := range srvs {
			// already running, don't need to start again
			if _, ok := running[srv.Service.Name+":"+srv.Service.Version]; ok {
				continue
			}

			// skip services which aren't running for a reason
			if srv.Status == runtime.Error {
				continue
			}
			if srv.Status == runtime.Building {
				continue
			}
			if srv.Status == runtime.Stopped {
				continue
			}

			// create the service
			if err := m.createServiceInRuntime(srv); err != nil {
				if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
					logger.Errorf("Error restarting service: %v", err)
				}
			}
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

	runtime.Runtime
}

// New returns a manager for the runtime
func New() runtime.Runtime {
	return &manager{
		Runtime: NewCache(runtime.DefaultRuntime),
	}
}
