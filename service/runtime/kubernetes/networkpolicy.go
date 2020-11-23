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
// Original source: github.com/micro/go-micro/v3/runtime/kubernetes/networkpolicy.go

package kubernetes

import (
	"github.com/micro/micro/v3/internal/kubernetes/client"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/runtime"
)

// createNetworkPolicy creates a networkpolicy resource
func (k *kubernetes) createNetworkPolicy(networkPolicy *runtime.NetworkPolicy) error {
	err := k.client.Create(&client.Resource{
		Kind:  "networkpolicy",
		Value: client.NewNetworkPolicy(networkPolicy.Name, networkPolicy.Namespace, networkPolicy.AllowedLabels),
	}, client.CreateNamespace(networkPolicy.Namespace))
	if err != nil {
		if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
			logger.Errorf("Error creating resource %s: %v", networkPolicy.String(), err)
		}
	}
	return err
}

// updateNetworkPolicy updates a networkpolicy resource in-place
func (k *kubernetes) updateNetworkPolicy(networkPolicy *runtime.NetworkPolicy) error {
	err := k.client.Update(&client.Resource{
		Kind:  "networkpolicy",
		Name:  networkPolicy.Name,
		Value: client.NewNetworkPolicy(networkPolicy.Name, networkPolicy.Namespace, networkPolicy.AllowedLabels),
	}, client.UpdateNamespace(networkPolicy.Namespace))
	if err != nil {
		if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
			logger.Errorf("Error updating resource %s: %v", networkPolicy.String(), err)
		}
	}
	return err
}

// deleteNetworkPolicy deletes a networkpolicy resource
func (k *kubernetes) deleteNetworkPolicy(networkPolicy *runtime.NetworkPolicy) error {
	err := k.client.Delete(&client.Resource{
		Kind:  "networkpolicy",
		Name:  networkPolicy.Name,
		Value: client.NewNetworkPolicy(networkPolicy.Name, networkPolicy.Namespace, networkPolicy.AllowedLabels),
	}, client.DeleteNamespace(networkPolicy.Namespace))
	if err != nil {
		if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
			logger.Errorf("Error deleting resource %s: %v", networkPolicy.String(), err)
		}
	}
	return err
}
