package kubernetes

import (
	"github.com/micro/micro/v3/internal/kubernetes/client"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/runtime"
)

// createResourceQuota creates a resourcequota resource
func (k *kubernetes) createResourceQuota(resourceQuota *runtime.ResourceQuota) error {
	err := k.client.Create(&client.Resource{
		Kind:  "resourcequota",
		Value: client.NewResourceQuota(resourceQuota),
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
		Kind:  "resourcequota",
		Name:  resourceQuota.Name,
		Value: client.NewResourceQuota(resourceQuota),
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
		Kind:  "resourcequota",
		Name:  resourceQuota.Name,
		Value: client.NewResourceQuota(resourceQuota),
	}, client.DeleteNamespace(resourceQuota.Namespace))
	if err != nil {
		if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
			logger.Errorf("Error deleting resource %s: %v", resourceQuota.String(), err)
		}
	}
	return err
}
