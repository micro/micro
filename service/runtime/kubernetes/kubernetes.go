// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Original source: github.com/micro/go-micro/v3/runtime/kubernetes/kubernetes.go

// Package kubernetes implements kubernetes micro runtime
package kubernetes

import (
	"fmt"
	"sync"
	"time"

	"github.com/micro/micro/v3/internal/kubernetes/api"
	"github.com/micro/micro/v3/internal/kubernetes/client"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/runtime"
)

var (
	DefaultServiceResources = &runtime.Resources{
		CPU:  200,
		Mem:  200,
		Disk: 2000,
	}

	DefaultImage = "micro/cells:v3"
)

// action to take on runtime service
type action int

type kubernetes struct {
	sync.Mutex
	// options configure runtime
	options runtime.Options
	// client is kubernetes client
	client client.Client
	// namespaces which exist
	namespaces []client.Namespace
}

// Init initializes runtime options
func (k *kubernetes) Init(opts ...runtime.Option) error {
	for _, o := range opts {
		o(&k.options)
	}
	return nil
}

func (k *kubernetes) Logs(resource runtime.Resource, options ...runtime.LogsOption) (runtime.LogStream, error) {
	// Handle the various different types of resources:
	switch resource.Type() {
	case runtime.TypeNamespace:
		// noop (Namespace is not supported by *kubernetes.Logs())
		return nil, nil
	case runtime.TypeNetworkPolicy:
		// noop (NetworkPolicy is not supported by *kubernetes.Logs()))
		return nil, nil
	case runtime.TypeResourceQuota:
		// noop (ResourceQuota is not supported by *kubernetes.Logs()))
		return nil, nil
	case runtime.TypeService:
		// Assert the resource back into a *runtime.Service
		s, ok := resource.(*runtime.Service)
		if !ok {
			return nil, runtime.ErrInvalidResource
		}

		klo := newLog(k.client, s.Name, options...)

		// if its not a stream then read the records, return and close
		if !klo.options.Stream {
			records, err := klo.Read()
			if err != nil {
				logger.Errorf("Failed to get logs for service '%v' from k8s: %v", s.Name, err)
				return nil, err
			}

			// create a stream even though we're not streaming
			kstream := &kubeStream{
				// make a stream buffer of size records
				stream: make(chan runtime.Log, len(records)),
				stop:   make(chan bool),
			}

			// load the records
			for _, record := range records {
				kstream.stream <- record
			}

			// close the stream so it doesn't block
			close(kstream.stream)

			return kstream, nil
		}

		// otherwise stream the logs
		stream, err := klo.Stream()
		if err != nil {
			return nil, err
		}

		return stream, nil
	default:
		return nil, runtime.ErrInvalidResource
	}
}

type kubeStream struct {
	// the k8s log stream
	stream chan runtime.Log
	// the stop chan
	sync.Mutex
	stop chan bool
	err  error
}

func (k *kubeStream) Error() error {
	return k.err
}

func (k *kubeStream) Chan() chan runtime.Log {
	return k.stream
}

func (k *kubeStream) Stop() error {
	k.Lock()
	defer k.Unlock()

	select {
	case <-k.stop:
		return nil
	default:
		close(k.stop)
	}

	return nil
}

// Create a resource
func (k *kubernetes) Create(resource runtime.Resource, opts ...runtime.CreateOption) error {
	k.Lock()
	defer k.Unlock()
	return k.create(resource, opts...)
}

func (k *kubernetes) create(resource runtime.Resource, opts ...runtime.CreateOption) error {
	// parse the options
	options := &runtime.CreateOptions{
		Type:      k.options.Type,
		Image:     k.options.Image,
		Namespace: client.DefaultNamespace,
	}
	for _, o := range opts {
		o(options)
	}

	// Handle the various different types of resources:
	switch resource.Type() {
	case runtime.TypeNamespace:
		// Assert the resource back into a *runtime.Namespace
		namespace, ok := resource.(*runtime.Namespace)
		if !ok {
			return runtime.ErrInvalidResource
		}
		return k.createNamespace(namespace)
	case runtime.TypeNetworkPolicy:
		// Assert the resource back into a *runtime.NetworkPolicy
		networkPolicy, ok := resource.(*runtime.NetworkPolicy)
		if !ok {
			return runtime.ErrInvalidResource
		}
		return k.createNetworkPolicy(networkPolicy)
	case runtime.TypeResourceQuota:
		// Assert the resource back into a *runtime.ResourceQuota
		resourceQuota, ok := resource.(*runtime.ResourceQuota)
		if !ok {
			return runtime.ErrInvalidResource
		}
		return k.createResourceQuota(resourceQuota)
	case runtime.TypeService:
		// Assert the resource back into a *runtime.Service
		s, ok := resource.(*runtime.Service)
		if !ok {
			return runtime.ErrInvalidResource
		}

		// default the service's source and version
		if len(s.Source) == 0 {
			s.Source = k.options.Source
		}
		if len(s.Version) == 0 {
			s.Version = "latest"
		}

		// ensure the namespace exists
		if err := k.ensureNamepaceExists(options.Namespace); err != nil {
			return nil
		}

		// create a secret for the deployment
		if err := k.createCredentials(s, options); err != nil {
			return err
		}

		// create some default resource requests
		if options.Resources == nil && options.Namespace != "micro" {
			options.Resources = DefaultServiceResources
		}

		if len(options.Image) == 0 {
			options.Image = DefaultImage
		}

		// create the deployment and set the runtime class name if provided
		dep := client.NewDeployment(s, options)
		if rcn := getRuntimeClassName(k.options.Context); len(rcn) > 0 {
			dep.Value.(*client.Deployment).Spec.Template.PodSpec.RuntimeClassName = rcn
			logger.Infof("Setting runtime class name to %v", rcn)
		}

		// create the deployment
		if err := k.client.Create(dep, client.CreateNamespace(options.Namespace)); err != nil {
			if parseError(err).Reason == "AlreadyExists" {
				return runtime.ErrAlreadyExists
			}
			if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
				logger.Errorf("Runtime failed to create deployment: %v", err)
			}
			return err
		}

		// create the service, one could already exist for another version so ignore ErrAlreadyExists
		if err := k.client.Create(client.NewService(s, options), client.CreateNamespace(options.Namespace)); err != nil {
			if parseError(err).Reason == "AlreadyExists" {
				return nil
			}
			if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
				logger.Errorf("Runtime failed to create service: %v", err)
			}
			return err
		}

		return nil
	default:
		return runtime.ErrInvalidResource
	}
}

// Read returns all instances of given service
func (k *kubernetes) Read(opts ...runtime.ReadOption) ([]*runtime.Service, error) {
	k.Lock()
	defer k.Unlock()

	// parse the options
	options := runtime.ReadOptions{
		Namespace: client.DefaultNamespace,
	}
	for _, o := range opts {
		o(&options)
	}

	// construct the query
	labels := map[string]string{}
	if len(options.Service) > 0 {
		labels["name"] = client.Format(options.Service)
	}
	if len(options.Version) > 0 {
		labels["version"] = client.Format(options.Version)
	}
	if len(options.Type) > 0 {
		labels["micro"] = client.Format(options.Type)
	}

	// lookup all the serivces which match this query, if one service has two different versions,
	// they'll be returned as two seperate resullts
	return k.getServices(client.GetNamespace(options.Namespace), client.GetLabels(labels))
}

// Update a resource in place
func (k *kubernetes) Update(resource runtime.Resource, opts ...runtime.UpdateOption) error {
	k.Lock()
	defer k.Unlock()

	// parse the options
	options := runtime.UpdateOptions{
		Namespace: client.DefaultNamespace,
	}
	for _, o := range opts {
		o(&options)
	}

	// Handle the various different types of resources:
	switch resource.Type() {
	case runtime.TypeNamespace:
		// noop (Namespace is not supported by *kubernetes.Update())
		return nil
	case runtime.TypeNetworkPolicy:
		// Assert the resource back into a *runtime.NetworkPolicy
		networkPolicy, ok := resource.(*runtime.NetworkPolicy)
		if !ok {
			return runtime.ErrInvalidResource
		}
		return k.updateNetworkPolicy(networkPolicy)
	case runtime.TypeResourceQuota:
		// Assert the resource back into a *runtime.ResourceQuota
		resourceQuota, ok := resource.(*runtime.ResourceQuota)
		if !ok {
			return runtime.ErrInvalidResource
		}
		return k.updateResourceQuota(resourceQuota)
	case runtime.TypeService:

		// Assert the resource back into a *runtime.Service
		s, ok := resource.(*runtime.Service)
		if !ok {
			return runtime.ErrInvalidResource
		}

		// construct the query
		labels := map[string]string{}
		if len(s.Name) > 0 {
			labels["name"] = client.Format(s.Name)
		}
		if len(s.Version) > 0 {
			labels["version"] = client.Format(s.Version)
		}

		// get the existing deployments
		depList := new(client.DeploymentList)
		d := &client.Resource{
			Kind:  "deployment",
			Value: depList,
		}
		depOpts := []client.GetOption{
			client.GetNamespace(options.Namespace),
			client.GetLabels(labels),
		}
		if err := k.client.Get(d, depOpts...); err != nil {
			return err
		} else if len(depList.Items) == 0 {
			return runtime.ErrNotFound
		}

		// update the deployments which match the query
		for _, dep := range depList.Items {
			// the service wan't created by the k8s runtime
			if dep.Metadata == nil || dep.Metadata.Annotations == nil {
				continue
			}

			// update metadata
			for k, v := range s.Metadata {
				dep.Metadata.Annotations[k] = v
			}

			if rcn := getRuntimeClassName(k.options.Context); len(rcn) > 0 {
				dep.Spec.Template.PodSpec.RuntimeClassName = rcn
				logger.Infof("Setting runtime class name to %v", rcn)
			}

			// update build time annotation
			dep.Spec.Template.Metadata.Annotations["updated"] = fmt.Sprintf("%d", time.Now().Unix())

			// update the deployment
			res := &client.Resource{
				Kind:  "deployment",
				Name:  resourceName(s),
				Value: &dep,
			}
			if err := k.client.Update(res, client.UpdateNamespace(options.Namespace)); err != nil {
				if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
					logger.Errorf("Runtime failed to update deployment: %v", err)
				}
				return err
			}
		}

		return nil
	default:
		return runtime.ErrInvalidResource
	}
}

// Delete removes a resource
func (k *kubernetes) Delete(resource runtime.Resource, opts ...runtime.DeleteOption) error {
	k.Lock()
	defer k.Unlock()

	options := runtime.DeleteOptions{
		Namespace: client.DefaultNamespace,
	}
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
		return k.deleteNamespace(namespace)
	case runtime.TypeNetworkPolicy:
		// Assert the resource back into a *runtime.NetworkPolicy
		networkPolicy, ok := resource.(*runtime.NetworkPolicy)
		if !ok {
			return runtime.ErrInvalidResource
		}
		return k.deleteNetworkPolicy(networkPolicy)
	case runtime.TypeResourceQuota:
		// Assert the resource back into a *runtime.ResourceQuota
		resourceQuota, ok := resource.(*runtime.ResourceQuota)
		if !ok {
			return runtime.ErrInvalidResource
		}
		return k.deleteResourceQuota(resourceQuota)
	case runtime.TypeService:

		// Assert the resource back into a *runtime.Service
		s, ok := resource.(*runtime.Service)
		if !ok {
			return runtime.ErrInvalidResource
		}

		// delete the deployment
		dep := client.NewDeployment(s, &runtime.CreateOptions{
			Type:      k.options.Type,
			Namespace: options.Namespace,
		})
		if err := k.client.Delete(dep, client.DeleteNamespace(options.Namespace)); err != nil {
			if err == api.ErrNotFound {
				return runtime.ErrNotFound
			}
			if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
				logger.Errorf("Runtime failed to delete deployment: %v", err)
			}
			return err
		}

		// delete the credentials
		if err := k.deleteCredentials(s, &runtime.CreateOptions{Namespace: options.Namespace}); err != nil {
			return err
		}

		// if there are more deployments for this service, then don't delete it
		labels := map[string]string{}
		if len(s.Name) > 0 {
			labels["name"] = client.Format(s.Name)
		}

		// get the existing services. todo: refactor to just get the deployments
		services, err := k.getServices(client.GetNamespace(options.Namespace), client.GetLabels(labels))
		if err != nil || len(services) > 0 {
			return err
		}

		// delete the service
		srv := client.NewService(s, &runtime.CreateOptions{
			Type:      k.options.Type,
			Namespace: options.Namespace,
		})
		if err := k.client.Delete(srv, client.DeleteNamespace(options.Namespace)); err != nil {
			if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
				logger.Errorf("Runtime failed to delete service: %v", err)
			}
			return err
		}

		return nil
	default:
		return runtime.ErrInvalidResource
	}
}

// Start starts the runtime
func (k *kubernetes) Start() error {
	return nil
}

// Stop shuts down the runtime
func (k *kubernetes) Stop() error {
	return nil
}

// String implements stringer interface
func (k *kubernetes) String() string {
	return "kubernetes"
}

// NewRuntime creates new kubernetes runtime
func NewRuntime(opts ...runtime.Option) runtime.Runtime {
	// get default options
	options := runtime.Options{
		// Create labels with type "micro": "service"
		Type: "service",
	}

	// apply requested options
	for _, o := range opts {
		o(&options)
	}

	// kubernetes client
	client := client.NewClusterClient()

	return &kubernetes{
		options: options,
		client:  client,
	}
}
