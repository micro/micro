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
				EphemeralStorage: fmt.Sprintf("%dMi", resourceQuota.Requests.Disk),
				Memory:           fmt.Sprintf("%dMi", resourceQuota.Requests.Mem),
			},
			Limits: &client.ResourceLimits{
				CPU:              fmt.Sprintf("%dm", resourceQuota.Limits.CPU),
				EphemeralStorage: fmt.Sprintf("%dMi", resourceQuota.Limits.Disk),
				Memory:           fmt.Sprintf("%dMi", resourceQuota.Limits.Mem),
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
				EphemeralStorage: fmt.Sprintf("%dMi", resourceQuota.Requests.Disk),
				Memory:           fmt.Sprintf("%dMi", resourceQuota.Requests.Mem),
			},
			Limits: &client.ResourceLimits{
				CPU:              fmt.Sprintf("%dm", resourceQuota.Limits.CPU),
				EphemeralStorage: fmt.Sprintf("%dMi", resourceQuota.Limits.Disk),
				Memory:           fmt.Sprintf("%dMi", resourceQuota.Limits.Mem),
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
				EphemeralStorage: fmt.Sprintf("%dMi", resourceQuota.Requests.Disk),
				Memory:           fmt.Sprintf("%dMi", resourceQuota.Requests.Mem),
			},
			Limits: &client.ResourceLimits{
				CPU:              fmt.Sprintf("%dm", resourceQuota.Limits.CPU),
				EphemeralStorage: fmt.Sprintf("%dMi", resourceQuota.Limits.Disk),
				Memory:           fmt.Sprintf("%dMi", resourceQuota.Limits.Mem),
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
