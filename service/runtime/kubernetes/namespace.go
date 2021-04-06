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
// Original source: github.com/micro/go-micro/v3/runtime/kubernetes/namespace.go

package kubernetes

import (
	"strings"

	"github.com/micro/micro/v3/internal/kubernetes/client"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/runtime"
)

func (k *kubernetes) ensureNamepaceExists(ns string) error {
	namespace := client.Format(ns)
	if namespace == client.DefaultNamespace {
		return nil
	}

	exist, err := k.namespaceExists(namespace)
	if err == nil && exist {
		return nil
	}
	if err != nil {
		if logger.V(logger.WarnLevel, logger.DefaultLogger) {
			logger.Warnf("Error checking namespace %v exists: %v", namespace, err)
		}
		return err
	}

	if err := k.autoCreateNamespace(namespace); err != nil {
		if logger.V(logger.WarnLevel, logger.DefaultLogger) {
			logger.Warnf("Error creating namespace %v: %v", namespace, err)
		}
		return err
	}

	return nil
}

// namespaceExists returns a boolean indicating if a namespace exists in the cache
func (k *kubernetes) namespaceExists(name string) (bool, error) {
	// populate the cache
	if k.namespaces == nil {
		if logger.V(logger.DebugLevel, logger.DefaultLogger) {
			logger.Debugf("Populating namespace cache")
		}

		namespaceList := new(client.NamespaceList)
		resource := &client.Resource{Kind: "namespace", Value: namespaceList}
		if err := k.client.List(resource); err != nil {
			return false, err
		}

		if logger.V(logger.DebugLevel, logger.DefaultLogger) {
			logger.Debugf("Populated namespace cache successfully with %v items", len(namespaceList.Items))
		}
		k.namespaces = namespaceList.Items
	}

	// check if the namespace exists in the cache
	for _, n := range k.namespaces {
		if n.Metadata.Name == name {
			return true, nil
		}
	}

	return false, nil
}

// autoCreateNamespace creates a new k8s namespace
func (k *kubernetes) autoCreateNamespace(namespace string) error {
	ns := client.Namespace{Metadata: &client.Metadata{Name: namespace}}
	err := k.client.Create(&client.Resource{Kind: "namespace", Value: ns})

	// ignore err already exists
	if err != nil && strings.Contains(err.Error(), "already exists") {
		logger.Debugf("Ignoring ErrAlreadyExists for namespace %v: %v", namespace, err)
		err = nil
	}

	// add to cache and create networkpolicy
	if err == nil && k.namespaces != nil {
		k.namespaces = append(k.namespaces, ns)
	}

	return err
}

// createNamespace creates a namespace resource
func (k *kubernetes) createNamespace(namespace *runtime.Namespace) error {
	err := k.client.Create(&client.Resource{
		Kind: "namespace",
		Name: namespace.Name,
		Value: &client.Namespace{
			Metadata: &client.Metadata{
				Name: namespace.Name,
			},
		},
	}, client.CreateNamespace(namespace.Name))
	if err != nil {
		if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
			logger.Errorf("Error creating namespace %s: %v", namespace.String(), err)
		}
	}
	return err
}

// deleteNamespace deletes a namespace resource
func (k *kubernetes) deleteNamespace(namespace *runtime.Namespace) error {
	err := k.client.Delete(&client.Resource{
		Kind: "namespace",
		Name: namespace.Name,
		Value: &client.Namespace{
			Metadata: &client.Metadata{
				Name: namespace.Name,
			},
		},
	}, client.DeleteNamespace(namespace.Name))
	if err != nil {
		if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
			logger.Errorf("Error deleting namespace %s: %v", namespace.String(), err)
		}
	}
	return err
}
