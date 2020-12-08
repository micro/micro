package client

import (
	"io"
	"sync"

	pb "github.com/micro/micro/v3/proto/runtime"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context"
	"github.com/micro/micro/v3/service/runtime"
)

type svc struct {
	sync.RWMutex
	options runtime.Options
	runtime pb.RuntimeService
}

// Init initializes runtime with given options
func (s *svc) Init(opts ...runtime.Option) error {
	s.Lock()
	defer s.Unlock()

	for _, o := range opts {
		o(&s.options)
	}

	s.runtime = pb.NewRuntimeService("runtime", client.DefaultClient)

	return nil
}

// Create a resource
func (s *svc) Create(resource runtime.Resource, opts ...runtime.CreateOption) error {
	var options runtime.CreateOptions
	for _, o := range opts {
		o(&options)
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

		// runtime namespace create request
		req := &pb.CreateRequest{
			Resource: &pb.Resource{
				Namespace: &pb.Namespace{
					Name: namespace.Name,
				},
			},
		}

		if _, err := s.runtime.Create(context.DefaultContext, req, client.WithAuthToken()); err != nil {
			return err
		}

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

		// runtime namespace create request
		req := &pb.CreateRequest{
			Resource: &pb.Resource{
				Networkpolicy: &pb.NetworkPolicy{
					Name:      networkPolicy.Name,
					Namespace: networkPolicy.Namespace,
				},
			},
		}

		if _, err := s.runtime.Create(context.DefaultContext, req, client.WithAuthToken()); err != nil {
			return err
		}

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

		// runtime namespace create request
		req := &pb.CreateRequest{
			Resource: &pb.Resource{
				Resourcequota: &pb.ResourceQuota{
					Requests: &pb.Resources{
						CPU:              int32(resourceQuota.Requests.CPU),
						EphemeralStorage: int32(resourceQuota.Requests.Disk),
						Memory:           int32(resourceQuota.Requests.Mem),
					},
					Limits: &pb.Resources{
						CPU:              int32(resourceQuota.Limits.CPU),
						EphemeralStorage: int32(resourceQuota.Limits.Disk),
						Memory:           int32(resourceQuota.Limits.Mem),
					},
					Name:      resourceQuota.Name,
					Namespace: resourceQuota.Namespace,
				},
			},
		}

		if _, err := s.runtime.Create(context.DefaultContext, req, client.WithAuthToken()); err != nil {
			return err
		}

	case runtime.TypeService:

		// Assert the resource back into a *runtime.Service
		svc, ok := resource.(*runtime.Service)
		if !ok {
			return runtime.ErrInvalidResource
		}

		// set the default source from MICRO_RUNTIME_SOURCE
		if len(svc.Source) == 0 {
			svc.Source = s.options.Source
		}

		// runtime service create request
		req := &pb.CreateRequest{
			Resource: &pb.Resource{
				Service: &pb.Service{
					Name:     svc.Name,
					Version:  svc.Version,
					Source:   svc.Source,
					Metadata: svc.Metadata,
				},
			},
			Options: &pb.CreateOptions{
				Command:    options.Command,
				Args:       options.Args,
				Env:        options.Env,
				Type:       options.Type,
				Image:      options.Image,
				Namespace:  options.Namespace,
				Secrets:    options.Secrets,
				Entrypoint: options.Entrypoint,
				Volumes:    options.Volumes,
			},
		}

		if _, err := s.runtime.Create(context.DefaultContext, req, client.WithAuthToken()); err != nil {
			return err
		}
	default:
		return runtime.ErrInvalidResource
	}

	return nil
}

func (s *svc) Logs(resource runtime.Resource, options ...runtime.LogsOption) (runtime.LogStream, error) {
	var opts runtime.LogsOptions
	for _, o := range options {
		o(&opts)
	}

	// Handle the various different types of resources:
	switch resource.Type() {
	case runtime.TypeNamespace:
		// noop (Namespace is not supported by *kubernetes.Logs())
		return nil, nil
	case runtime.TypeNetworkPolicy:
		// noop (NetworkPolicy is not supported by *kubernetes.Logs())
		return nil, nil
	case runtime.TypeResourceQuota:
		// noop (ResourceQuota is not supported by *kubernetes.Logs())
		return nil, nil
	case runtime.TypeService:
		// Assert the resource back into a *runtime.Service
		service, ok := resource.(*runtime.Service)
		if !ok {
			return nil, runtime.ErrInvalidResource
		}

		ls, err := s.runtime.Logs(context.DefaultContext, &pb.LogsRequest{
			Service: service.Name,
			Stream:  opts.Stream,
			Count:   opts.Count,
			Options: &pb.LogsOptions{
				Namespace: opts.Namespace,
			},
		}, client.WithAuthToken())
		if err != nil {
			return nil, err
		}

		logStream := &serviceLogs{
			service: service.Name,
			stream:  make(chan runtime.Log),
			stop:    make(chan bool),
		}

		go func() {
			for {
				select {
				// @todo this never seems to return, investigate
				case <-ls.Context().Done():
					logStream.Stop()
				}
			}
		}()

		go func() {
			for {
				select {
				// @todo this never seems to return, investigate
				case <-ls.Context().Done():
					close(logStream.stream)
					return
				case <-logStream.stop:
					close(logStream.stream)
					return
				default:
					record := pb.LogRecord{}

					if err := ls.RecvMsg(&record); err != nil {
						if err != io.EOF {
							logStream.err = err
						}
						close(logStream.stream)
						logStream.Stop()
						return
					}

					logStream.stream <- runtime.Log{
						Message:  record.GetMessage(),
						Metadata: record.GetMetadata(),
					}
				}
			}
		}()

		return logStream, nil
	default:
		return nil, runtime.ErrInvalidResource
	}
}

type serviceLogs struct {
	service string
	stream  chan runtime.Log
	sync.Mutex
	stop chan bool
	err  error
}

func (l *serviceLogs) Error() error {
	return l.err
}

func (l *serviceLogs) Chan() chan runtime.Log {
	return l.stream
}

func (l *serviceLogs) Stop() error {
	l.Lock()
	defer l.Unlock()
	select {
	case <-l.stop:
		return nil
	default:
		close(l.stop)
	}
	return nil
}

// Read returns the service with the given name from the runtime
func (s *svc) Read(opts ...runtime.ReadOption) ([]*runtime.Service, error) {
	var options runtime.ReadOptions
	for _, o := range opts {
		o(&options)
	}

	// runtime service create request
	req := &pb.ReadRequest{
		Options: &pb.ReadOptions{
			Service:   options.Service,
			Version:   options.Version,
			Type:      options.Type,
			Namespace: options.Namespace,
		},
	}

	resp, err := s.runtime.Read(context.DefaultContext, req, client.WithAuthToken())
	if err != nil {
		return nil, err
	}

	services := make([]*runtime.Service, 0, len(resp.Services))
	for _, service := range resp.Services {
		svc := &runtime.Service{
			Name:     service.Name,
			Version:  service.Version,
			Source:   service.Source,
			Metadata: service.Metadata,
			Status:   runtime.ServiceStatus(service.Status),
		}
		services = append(services, svc)
	}

	return services, nil
}

// Update a resource
func (s *svc) Update(resource runtime.Resource, opts ...runtime.UpdateOption) error {
	var options runtime.UpdateOptions
	for _, o := range opts {
		o(&options)
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

		// runtime namespace update request
		req := &pb.UpdateRequest{
			Resource: &pb.Resource{
				Namespace: &pb.Namespace{
					Name: namespace.Name,
				},
			},
		}

		if _, err := s.runtime.Update(context.DefaultContext, req, client.WithAuthToken()); err != nil {
			return err
		}

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

		// runtime networkpolicy update request
		req := &pb.UpdateRequest{
			Resource: &pb.Resource{
				Networkpolicy: &pb.NetworkPolicy{
					Name:      networkPolicy.Name,
					Namespace: networkPolicy.Namespace,
				},
			},
		}

		if _, err := s.runtime.Update(context.DefaultContext, req, client.WithAuthToken()); err != nil {
			return err
		}

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

		// runtime resourcequota update request
		req := &pb.UpdateRequest{
			Resource: &pb.Resource{
				Resourcequota: &pb.ResourceQuota{
					Requests: &pb.Resources{
						CPU:              int32(resourceQuota.Requests.CPU),
						EphemeralStorage: int32(resourceQuota.Requests.Disk),
						Memory:           int32(resourceQuota.Requests.Mem),
					},
					Limits: &pb.Resources{
						CPU:              int32(resourceQuota.Limits.CPU),
						EphemeralStorage: int32(resourceQuota.Limits.Disk),
						Memory:           int32(resourceQuota.Limits.Mem),
					},
					Name:      resourceQuota.Name,
					Namespace: resourceQuota.Namespace,
				},
			},
		}

		if _, err := s.runtime.Update(context.DefaultContext, req, client.WithAuthToken()); err != nil {
			return err
		}

	case runtime.TypeService:

		// Assert the resource back into a *runtime.Service
		svc, ok := resource.(*runtime.Service)
		if !ok {
			return runtime.ErrInvalidResource
		}

		// runtime service create request
		req := &pb.UpdateRequest{
			Resource: &pb.Resource{
				Service: &pb.Service{
					Name:     svc.Name,
					Version:  svc.Version,
					Source:   svc.Source,
					Metadata: svc.Metadata,
				},
			},
			Options: &pb.UpdateOptions{
				Namespace:  options.Namespace,
				Entrypoint: options.Entrypoint,
			},
		}

		if _, err := s.runtime.Update(context.DefaultContext, req, client.WithAuthToken()); err != nil {
			return err
		}
	default:
		return runtime.ErrInvalidResource
	}

	return nil
}

// Delete a resource
func (s *svc) Delete(resource runtime.Resource, opts ...runtime.DeleteOption) error {
	var options runtime.DeleteOptions
	for _, o := range opts {
		o(&options)
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

		// runtime namespace delete request
		req := &pb.DeleteRequest{
			Resource: &pb.Resource{
				Namespace: &pb.Namespace{
					Name: namespace.Name,
				},
			},
		}

		if _, err := s.runtime.Delete(context.DefaultContext, req, client.WithAuthToken()); err != nil {
			return err
		}

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

		// runtime namespace delete request
		req := &pb.DeleteRequest{
			Resource: &pb.Resource{
				Networkpolicy: &pb.NetworkPolicy{
					Name:      networkPolicy.Name,
					Namespace: networkPolicy.Namespace,
				},
			},
		}

		if _, err := s.runtime.Delete(context.DefaultContext, req, client.WithAuthToken()); err != nil {
			return err
		}

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

		// runtime resourcequota delete request
		req := &pb.DeleteRequest{
			Resource: &pb.Resource{
				Resourcequota: &pb.ResourceQuota{
					Requests: &pb.Resources{
						CPU:              int32(resourceQuota.Requests.CPU),
						EphemeralStorage: int32(resourceQuota.Requests.Disk),
						Memory:           int32(resourceQuota.Requests.Mem),
					},
					Limits: &pb.Resources{
						CPU:              int32(resourceQuota.Limits.CPU),
						EphemeralStorage: int32(resourceQuota.Limits.Disk),
						Memory:           int32(resourceQuota.Limits.Mem),
					},
					Name:      resourceQuota.Name,
					Namespace: resourceQuota.Namespace,
				},
			},
		}

		if _, err := s.runtime.Delete(context.DefaultContext, req, client.WithAuthToken()); err != nil {
			return err
		}

	case runtime.TypeService:

		// Assert the resource back into a *runtime.Service
		svc, ok := resource.(*runtime.Service)
		if !ok {
			return runtime.ErrInvalidResource
		}

		// runtime service dekete request
		req := &pb.DeleteRequest{
			Resource: &pb.Resource{
				Service: &pb.Service{
					Name:     svc.Name,
					Version:  svc.Version,
					Source:   svc.Source,
					Metadata: svc.Metadata,
				},
			},
			Options: &pb.DeleteOptions{
				Namespace: options.Namespace,
			},
		}

		if _, err := s.runtime.Delete(context.DefaultContext, req, client.WithAuthToken()); err != nil {
			return err
		}
	default:
		return runtime.ErrInvalidResource
	}

	return nil
}

// Start starts the runtime
func (s *svc) Start() error {
	// NOTE: nothing to be done here
	return nil
}

// Stop stops the runtime
func (s *svc) Stop() error {
	// NOTE: nothing to be done here
	return nil
}

// Returns the runtime service implementation
func (s *svc) String() string {
	return "service"
}

// NewRuntime creates new service runtime and returns it
func NewRuntime(opts ...runtime.Option) runtime.Runtime {
	var options runtime.Options
	for _, o := range opts {
		o(&options)
	}

	return &svc{
		options: options,
		runtime: pb.NewRuntimeService("runtime", client.DefaultClient),
	}
}
