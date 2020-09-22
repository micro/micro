package manager

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	gorun "github.com/micro/go-micro/v3/runtime"
	"github.com/micro/go-micro/v3/runtime/builder"
	"github.com/micro/go-micro/v3/store"
	cachest "github.com/micro/go-micro/v3/store/cache"
	filest "github.com/micro/go-micro/v3/store/file"
	"github.com/micro/go-micro/v3/store/memory"
	"github.com/micro/micro/v3/internal/namespace"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/runtime"
	"github.com/micro/micro/v3/service/runtime/util"
)

// Init initializes the runtime
func (m *manager) Init(...gorun.Option) error {
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
	return m.publishEvent(gorun.Create, srv, &options)
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
	srvs, err := m.readServices(options.Namespace, &runtime.Service{
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
	ret := []*runtime.Service{}
	for _, srv := range srvs {
		ret = append(ret, srv.Service)
		md, ok := statuses[srv.Service.Name+":"+srv.Service.Version]
		if !ok {
			continue
		}
		srv.Service.Status = md.Status
		srv.Service.Metadata["error"] = md.Error
		if !md.Updated.IsZero() {
			srv.Service.Metadata["started"] = md.Updated.Format(time.RFC3339)
		}
	}

	return ret, nil
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
	return m.publishEvent(gorun.Update, srv, &runtime.CreateOptions{Namespace: options.Namespace})
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
	return m.publishEvent(gorun.Delete, srv, &runtime.CreateOptions{Namespace: options.Namespace})
}

func (m *manager) CreateNamespace(ns string) error {
	// Do we need to store this locally?
	return runtime.DefaultRuntime.CreateNamespace(ns)
}

func (m *manager) DeleteNamespace(ns string) error {
	return runtime.DefaultRuntime.DeleteNamespace(ns)
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

	// Watch services that were running previously
	go m.watchServices()

	return nil
}

// Logs for a service
func (m *manager) Logs(srv *runtime.Service, opts ...runtime.LogsOption) (runtime.Logs, error) {
	return runtime.Log(srv, opts...)
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
			running[v.Name+":"+v.Version+":"+v.Source] = v
		}

		for _, srv := range srvs {
			if _, ok := running[srv.Service.Name+":"+srv.Service.Version+":"+srv.Service.Source]; ok {
				// already running, don't need to start again
				continue
			}

			// if the source is a blob, we must pull it and save it to a tmp dir so it can be accessed
			// by the local runtime which has no concept of the blob store
			if strings.HasPrefix(srv.Service.Source, "source://") {
				// create a tmp dir to store the source in
				dir, err := ioutil.TempDir(os.TempDir(), fmt.Sprintf("source-%v-*", srv.Service.Name))
				if err != nil {
					if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
						logger.Errorf("Error restarting service %v: %v", srv.Service.Name, err)
					}
					continue
				}

				// pull the source from the blob store
				src, err := util.ReadSource(srv.Service, srv.Options.Secrets, ns)
				if err != nil {
					if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
						logger.Errorf("Error restarting service %v: %v", srv.Service.Name, err)
					}
					continue
				}

				// unarchive the tar into the directory
				if err := util.Unarchive(src, dir); err != nil {
					if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
						logger.Errorf("Error restarting service %v: %v", srv.Service.Name, err)
					}
					continue
				}

				// set the service's source to the tmp dir, but don't write this back to the store. the service
				// might not be at the top level of the source (e.g. in the case of a mono-repo) so we'll join
				// the entrypoint
				srv.Service.Source = filepath.Join(dir, srv.Service.Metadata["entrypoint"])
			}

			// generate an auth account for the service to use
			acc, err := m.generateAccount(srv.Service, ns)
			if err != nil {
				continue
			}

			// construct the options
			options := []runtime.CreateOption{
				runtime.CreateImage(srv.Options.Image),
				runtime.CreateType(srv.Options.Type),
				runtime.CreateNamespace(ns),
				runtime.WithArgs(srv.Options.Args...),
				runtime.WithCommand(srv.Options.Command...),
				runtime.WithEnv(m.runtimeEnv(srv.Service, srv.Options)),
			}

			// inject the credentials into the service if present
			if len(acc.ID) > 0 && len(acc.Secret) > 0 {
				options = append(options, runtime.WithSecret("MICRO_AUTH_ID", acc.ID))
				options = append(options, runtime.WithSecret("MICRO_AUTH_SECRET", acc.Secret))
			}

			// add the secrets provided by the client
			for key, value := range srv.Options.Secrets {
				options = append(options, runtime.WithSecret(key, value))
			}

			// create the service
			if err := runtime.Create(srv.Service, options...); err != nil {
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
	// cache is a memory store which is used to store any information we don't want to write to the
	// global store, e.g. service status / errors (these will change depending on the
	// managed runtime and hence won't be the same globally).
	cache store.Store
	// fileCache is a cache store used to store any information we don't want to write to the
	// global store but want to persist across restarts, e.g. events consumed
	fileCache store.Store
	// builder is used to prebuild source. it can be nil.
	builder builder.Builder
}

// New returns a manager for the runtime
func New(opts ...Option) gorun.Runtime {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	return &manager{
		builder:   options.Builder,
		cache:     memory.NewStore(),
		fileCache: cachest.NewStore(filest.NewStore()),
	}
}

type Options struct {
	Builder builder.Builder
}

type Option func(o *Options)

func Builder(b builder.Builder) Option {
	return func(o *Options) {
		o.Builder = b
	}
}
