package manager

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/build"
	"github.com/micro/micro/v3/service/build/util/tar"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/runtime"
	kclient "github.com/micro/micro/v3/service/runtime/kubernetes/client"
	"github.com/micro/micro/v3/service/runtime/source/git"
	"github.com/micro/micro/v3/service/store"
	"github.com/micro/micro/v3/util/namespace"
)

const (
	// servicePrefix is prefixed to the key for service records
	servicePrefix = "service:"
)

// service is the object persisted in the store
type service struct {
	Service   *runtime.Service       `json:"service"`
	Options   *runtime.CreateOptions `json:"options"`
	Status    runtime.ServiceStatus  `json:"status"`
	UpdatedAt time.Time              `json:"last_updated"`
	Error     string                 `json:"error"`
}

// key to write the service to the store under, e.g:
// "service/foo/bar:latest"
func (s *service) Key() string {
	return servicePrefix + s.Options.Namespace + ":" + s.Service.Name + ":" + s.Service.Version
}

// unique is a helper method to filter a slice of strings
// down to unique entries
func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func (m *manager) checkServices() {
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

// writeService to the store
func (m *manager) writeService(srv *service) error {
	bytes, err := json.Marshal(srv)
	if err != nil {
		return err
	}

	return store.Write(&store.Record{Key: srv.Key(), Value: bytes})
}

// deleteService from the store
func (m *manager) deleteService(srv *service) error {
	return store.Delete(srv.Key())
}

// readServices returns all the services in a given namespace. If a service name and
// version are provided it will filter using these as well
func (m *manager) readServices(namespace string, srv *runtime.Service) ([]*service, error) {
	prefix := servicePrefix + namespace + ":"
	if len(srv.Name) > 0 {
		prefix += srv.Name + ":"
	}
	if len(srv.Name) > 0 && len(srv.Version) > 0 {
		prefix += srv.Version
	}

	recs, err := store.Read(prefix, store.ReadPrefix())
	if err != nil {
		return nil, err
	} else if len(recs) == 0 {
		return make([]*service, 0), nil
	}

	srvs := make([]*service, 0, len(recs))
	for _, r := range recs {
		var s *service
		if err := json.Unmarshal(r.Value, &s); err != nil {
			return nil, err
		}
		srvs = append(srvs, s)
	}

	return srvs, nil
}

// listNamespaces of the services in the store. todo: remove this and have the watchServices func
// query the store directly
func (m *manager) listNamespaces() ([]string, error) {
	recs, err := store.Read(servicePrefix, store.ReadPrefix())
	if err != nil {
		return nil, err
	}
	if len(recs) == 0 {
		return []string{namespace.DefaultNamespace}, nil
	}

	namespaces := make([]string, 0, len(recs))
	for _, rec := range recs {
		// key is formatted 'prefix:namespace:name:version'
		if comps := strings.Split(rec.Key, ":"); len(comps) == 4 {
			namespaces = append(namespaces, comps[1])
		} else {
			return nil, fmt.Errorf("Invalid key: %v", rec.Key)
		}
	}

	return unique(namespaces), nil
}

func (m *manager) buildAndRun(srv *service) {
	if err := m.build(srv); err != nil {
		return
	}

	srv.Status = runtime.Starting
	m.writeService(srv)

	if err := m.createServiceInRuntime(srv); err != nil {
		srv.Status = runtime.Error
		srv.Error = fmt.Sprintf("Error creating service: %v", err)
		m.writeService(srv)
	}
}

func (m *manager) buildAndUpdate(srv *service) {
	if err := m.build(srv); err != nil {
		return
	}

	srv.Status = runtime.Starting
	m.writeService(srv)

	if err := m.updateServiceInRuntime(srv); err != nil {
		srv.Status = runtime.Error
		srv.Error = fmt.Sprintf("Error updating service: %v", err)
		m.writeService(srv)
	}
}

func (m *manager) build(srv *service) error {
	logger.Infof("Preparing to build %v:%v", srv.Service.Name, srv.Service.Version)

	// set the service status to building
	srv.Status = runtime.Building
	m.writeService(srv)

	// handleError will set the error status on the service
	handleError := func(err error, msg string) {
		logger.Warnf("Build failed %v:%v: %v %v", srv.Service.Name, srv.Service.Version, msg, err)
		srv.Status = runtime.Error
		srv.Error = fmt.Sprintf("%v: %v", msg, err)
		m.writeService(srv)
	}

	// load the source
	var source io.Reader
	var err error
	if strings.HasPrefix(srv.Service.Source, "source://") {
		// if the source was uploaded to the blob store, it'll have source:// as the prefix
		nsOpt := store.BlobNamespace(srv.Options.Namespace)
		source, err = store.DefaultBlobStore.Read(srv.Service.Source, nsOpt)
	} else {
		// the source will otherwise be a git remote, we'll clone it and then tar archive the result
		gitSrc, err := git.ParseSource(srv.Service.Source)
		if err != nil {
			handleError(err, "Error parsing git source")
			return err
		}
		if len(srv.Options.Entrypoint) == 0 {
			srv.Options.Entrypoint = gitSrc.Folder
		}

		// checkout the source
		gitSrc.Ref = srv.Service.Version
		dir, err := git.CheckoutSource(gitSrc, srv.Options.Secrets)
		if err != nil {
			handleError(err, "Error fetching git source")
			return err
		}

		// archive the source so it can be passed to the build
		source, err = tar.Archive(dir)
	}
	if err != nil {
		handleError(err, "Error loading source")
		return err
	}

	// build the source
	logger.Infof("Build starting %v:%v", srv.Service.Name, srv.Service.Version)
	build, err := build.DefaultBuilder.Build(source,
		build.Archive("tar"),
		build.Entrypoint(srv.Options.Entrypoint),
	)
	logger.Infof("Build finished %v:%v %v", srv.Service.Name, srv.Service.Version, err)
	if err != nil {
		handleError(err, "Error building service")
		return err
	}

	// for the kubernetes runtime, the container needs to pull the source (it's not got access to the
	// local filesystem like the local runtime does). hence we'll upload the source to the blob store
	// which the cell (container) can then pull via the Runtime.Build.Read RPC.
	if m.Runtime.String() != "local" {
		logger.Infof("Uploading build %v:%v", srv.Service.Name, srv.Service.Version)
		nsOpt := store.BlobNamespace(srv.Options.Namespace)
		key := fmt.Sprintf("build://%v:%v", srv.Service.Name, srv.Service.Version)
		if err := store.DefaultBlobStore.Write(key, build, nsOpt); err != nil {
			handleError(err, "Error uploading build")
			return err
		}
	}

	return nil
}

func (m *manager) updateServiceInRuntime(srv *service) error {
	// construct the options
	options := []runtime.UpdateOption{
		runtime.UpdateEntrypoint(srv.Options.Entrypoint),
		runtime.UpdateNamespace(srv.Options.Namespace),
	}

	// add the secrets
	for key, value := range srv.Options.Secrets {
		options = append(options, runtime.UpdateSecret(key, value))
	}

	// update the service
	return m.Runtime.Update(srv.Service, options...)
}

// createServiceInRuntime will add all the required env vars and secrets and then create the service
func (m *manager) createServiceInRuntime(srv *service) error {
	// generate an auth account for the service to use
	acc, err := m.generateAccount(srv)
	if err != nil {
		return err
	}

	// construct the options
	options := []runtime.CreateOption{
		runtime.CreateEntrypoint(srv.Options.Entrypoint),
		runtime.CreateImage(srv.Options.Image),
		runtime.CreateType(srv.Options.Type),
		runtime.CreateNamespace(srv.Options.Namespace),
		runtime.WithArgs(srv.Options.Args...),
		runtime.WithCommand(srv.Options.Command...),
		runtime.WithEnv(m.runtimeEnv(srv.Service, srv.Options)),
		runtime.CreateInstances(srv.Options.Instances),
	}

	// add the secrets
	for key, value := range srv.Options.Secrets {
		options = append(options, runtime.WithSecret(key, value))
	}

	// inject the credentials into the service if present
	if len(acc.ID) > 0 && len(acc.Secret) > 0 {
		options = append(options, runtime.WithSecret("MICRO_AUTH_ID", acc.ID))
		options = append(options, runtime.WithSecret("MICRO_AUTH_SECRET", acc.Secret))
	}

	// create the service
	return m.Runtime.Create(srv.Service, options...)
}

// checkoutSource will take a service and download the source into a temp directory. This source
// could be a git remote or a reference to some source in the blob store (used for running local
// code on the server).
func (m *manager) checkoutSource(srv *service) (string, error) {
	if strings.HasPrefix(srv.Service.Source, "source://") {
		return m.checkoutBlobSource(srv)
	} else {
		return m.checkoutGitSource(srv)
	}
}

// checkoutBlobSource will checkout source from the blob store using the key in the service's source
// attribute. It will then unarchive the source into a temp directory and return the location of
// said directory.
func (m *manager) checkoutBlobSource(srv *service) (string, error) {
	nsOpt := store.BlobNamespace(srv.Options.Namespace)
	source, err := store.DefaultBlobStore.Read(srv.Service.Source, nsOpt)
	if err != nil {
		return "", err
	}

	dir, err := ioutil.TempDir(os.TempDir(), "blob-*")
	if err != nil {
		return "", err
	}

	if err := tar.Unarchive(source, dir); err != nil {
		return "", err
	}

	return dir, nil
}

// checkoutGitSource will download source from a git remote into a temp dir and then return the
// location of that temp directory
func (m *manager) checkoutGitSource(srv *service) (string, error) {
	gitSrc, err := git.ParseSource(srv.Service.Source)
	if err != nil {
		return "", err
	}
	gitSrc.Ref = srv.Service.Version

	dir, err := git.CheckoutSource(gitSrc, srv.Options.Secrets)
	if err != nil {
		return "", err
	}

	// the dir will contain the entire repo, however the use could've specified a subfolder within
	// that repo. this is the case for mono-repos
	if len(srv.Options.Entrypoint) == 0 {
		srv.Options.Entrypoint = gitSrc.Folder
	}

	return dir, nil
}

// runtimeEnv returns the environment variables which should  be used when creating a service.
func (m *manager) runtimeEnv(srv *runtime.Service, options *runtime.CreateOptions) []string {
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
	env := map[string]string{
		// ensure a profile for the services isn't set, they
		// should use the default RPC clients
		"MICRO_PROFILE": "service",
		// pass the service's name and version
		"MICRO_SERVICE_NAME":    srv.Name,
		"MICRO_SERVICE_VERSION": srv.Version,
		// set the proxy for the service to use (e.g. micro network)
		// using the proxy which has been configured for the runtime
		"MICRO_PROXY": client.DefaultClient.Options().Proxy,
	}

	// bind to port 8080, this is what the k8s tcp readiness check will use
	if runtime.DefaultRuntime.String() == "kubernetes" {
		env["MICRO_SERVICE_ADDRESS"] = ":8080"
	}

	// set the env vars provided
	setEnv(options.Env, env)

	// set the service namespace
	if len(options.Namespace) > 0 {
		env["MICRO_NAMESPACE"] = options.Namespace
	}

	// create a new env
	var vars []string
	for k, v := range env {
		vars = append(vars, k+"="+v)
	}

	// setup the runtime env
	return vars
}

func (m *manager) generateAccount(srv *service) (*auth.Account, error) {
	accName := srv.Service.Name + "-" + srv.Service.Version

	opts := []auth.GenerateOption{
		auth.WithIssuer(srv.Options.Namespace),
		auth.WithScopes("service"),
		auth.WithType("service"),
	}

	acc, err := auth.Generate(accName, opts...)
	if err != nil {
		if logger.V(logger.WarnLevel, logger.DefaultLogger) {
			logger.Warnf("Error generating account %v: %v", accName, err)
		}
		return nil, err
	}
	if logger.V(logger.DebugLevel, logger.DefaultLogger) {
		logger.Debugf("Generated auth account: %v, secret: [len: %v]", acc.ID, len(acc.Secret))
	}

	return acc, nil
}

// cleanupBlobStore deletes the source code and build from the blob store once the service finishes
// running.
func (m *manager) cleanupBlobStore(srv *service) {
	// delete the raw source code
	opt := store.BlobNamespace(srv.Options.Namespace)
	srcKey := fmt.Sprintf("source://%v:%v", srv.Service.Name, srv.Service.Version)
	if err := store.DefaultBlobStore.Delete(srcKey, opt); err != nil && err != store.ErrNotFound {
		logger.Warnf("Error deleting source %v: %v", srcKey, err)
	}

	// if there is no build enabled, there won't be any build to delete
	if build.DefaultBuilder == nil {
		return
	}

	// delete the binary
	buildKey := fmt.Sprintf("build://%v:%v", srv.Service.Name, srv.Service.Version)
	if err := store.DefaultBlobStore.Delete(buildKey, opt); err != nil && err != store.ErrNotFound {
		logger.Warnf("Error deleting build %v: %v", srcKey, err)
	}
}

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
		rs, ok := rSrvMap[kclient.Format(s.Service.Name)+":"+kclient.Format(s.Service.Version)]
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
		if options.Instances > 0 {
			service.Options.Instances = options.Instances
		}
		if len(options.Entrypoint) > 0 {
			service.Options.Entrypoint = options.Entrypoint
		}
		if len(options.Secrets) > 0 {
			service.Options.Secrets = options.Secrets
		}

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

// watchServices periodically checks services and whether they need to be recreated
func (m *manager) watchServices() {
	t := time.NewTicker(time.Second * 10)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			m.checkServices()
		case <-m.exit:
			return
		}
	}
}

// Stop the manager
func (m *manager) Stop() error {
	if !m.running {
		return nil
	}
	m.running = false

	// ping to exit
	select {
	case m.exit <- true:
	default:
	}

	return runtime.DefaultRuntime.Stop()
}

// String describes runtime
func (m *manager) String() string {
	return "manager"
}

type manager struct {
	// running is true after Start is called
	running bool
	exit    chan bool

	runtime.Runtime
}

// New returns a manager for the runtime
func New() runtime.Runtime {
	return &manager{
		exit:    make(chan bool, 1),
		Runtime: NewCache(runtime.DefaultRuntime),
	}
}
