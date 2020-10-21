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

package kubernetes

import (
	"fmt"

	"github.com/micro/micro/v3/internal/kubernetes/client"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/runtime"
)

// createResourceQuota creates a resourcequota resource
func (k *kubernetes) createResourceQuota(resourceQuota *runtime.ResourceQuota) error {
	err := k.client.Create(&client.Resource{
		Kind: "resourcequota",
		Value: client.ResourceQuota{
			Requests: &client.ResourceLimits{
				CPU:              fmt.Sprintf("%dm", resourceQuota.Requests.CPU),
				EphemeralStorage: fmt.Sprintf("%dm", resourceQuota.Requests.Disk),
				Memory:           fmt.Sprintf("%dm", resourceQuota.Requests.Mem),
			},
			Limits: &client.ResourceLimits{
				CPU:              fmt.Sprintf("%dm", resourceQuota.Limits.CPU),
				EphemeralStorage: fmt.Sprintf("%dm", resourceQuota.Limits.Disk),
				Memory:           fmt.Sprintf("%dm", resourceQuota.Limits.Mem),
			},
			Metadata: &client.Metadata{
				Name:      resourceQuota.Name,
				Namespace: resourceQuota.Namespace,
			},
		},
	}, client.CreateNamespace(resourceQuota.Namespace))
	if err != nil {
		if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
			logger.Errorf("Error creating resource %s: %v", resourceQuota.String(), err)
		}
	}
	return err
}

// updateResourceQuota updates a resourcequota resource in-place
func (k *kubernetes) updateResourceQuota(resourceQuota *runtime.ResourceQuota) error {
	err := k.client.Update(&client.Resource{
		Kind: "resourcequota",
		Value: client.ResourceQuota{
			Requests: &client.ResourceLimits{
				CPU:              fmt.Sprintf("%dm", resourceQuota.Requests.CPU),
				EphemeralStorage: fmt.Sprintf("%dm", resourceQuota.Requests.Disk),
				Memory:           fmt.Sprintf("%dm", resourceQuota.Requests.Mem),
			},
			Limits: &client.ResourceLimits{
				CPU:              fmt.Sprintf("%dm", resourceQuota.Limits.CPU),
				EphemeralStorage: fmt.Sprintf("%dm", resourceQuota.Limits.Disk),
				Memory:           fmt.Sprintf("%dm", resourceQuota.Limits.Mem),
			},
			Metadata: &client.Metadata{
				Name:      resourceQuota.Name,
				Namespace: resourceQuota.Namespace,
			},
		},
	}, client.UpdateNamespace(resourceQuota.Namespace))
	if err != nil {
		if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
			logger.Errorf("Error updating resource %s: %v", resourceQuota.String(), err)
		}
	}
	return err
}

// deleteResourcequota deletes a resourcequota resource
func (k *kubernetes) deleteResourceQuota(resourceQuota *runtime.ResourceQuota) error {
	err := k.client.Delete(&client.Resource{
		Kind: "resourcequota",
		Value: client.ResourceQuota{
			Requests: &client.ResourceLimits{
				CPU:              fmt.Sprintf("%dm", resourceQuota.Requests.CPU),
				EphemeralStorage: fmt.Sprintf("%dm", resourceQuota.Requests.Disk),
				Memory:           fmt.Sprintf("%dm", resourceQuota.Requests.Mem),
			},
			Limits: &client.ResourceLimits{
				CPU:              fmt.Sprintf("%dm", resourceQuota.Limits.CPU),
				EphemeralStorage: fmt.Sprintf("%dm", resourceQuota.Limits.Disk),
				Memory:           fmt.Sprintf("%dm", resourceQuota.Limits.Mem),
			},
			Metadata: &client.Metadata{
				Name:      resourceQuota.Name,
				Namespace: resourceQuota.Namespace,
			},
		},
	}, client.DeleteNamespace(resourceQuota.Namespace))
	if err != nil {
		if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
			logger.Errorf("Error deleting resource %s: %v", resourceQuota.String(), err)
		}
	}
	return err
}
