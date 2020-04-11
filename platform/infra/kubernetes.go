package infra

import (
	"fmt"

	"github.com/spf13/viper"
)

// Kubernetes represents a Kube Cluster
type Kubernetes struct {
	Name     string
	Region   string
	Provider string
}

// Steps generates steps that provision a Kubernetes cluster
func (k *Kubernetes) Steps(runID int32) ([]Step, error) {
	var s []Step
	k8sName := k.internalName("k8s")
	configName := k.internalName("kubeconfig")
	vars := make(map[string]string)
	vars["name"] = k.Name
	vars["kubernetes"] = k.Provider
	vars["region"] = k.Region
	vars["args"] = fmt.Sprintf(`["%s","%s"]`, k8sName, viper.GetString("aws-region"))
	remoteStates := make(map[string]string)
	remoteStates["k8s"] = k8sName
	s = append(s,
		// Provision the cluster
		Step{
			&TerraformModule{
				ID:        k8sName,
				Name:      k8sName,
				Source:    "./infra/kubernetes/" + k.Provider,
				Path:      fmt.Sprintf("/tmp/%s-%d", k8sName, runID),
				Variables: vars,
			},
		},
		// Grab the Kubernetes Config
		Step{
			&TerraformModule{
				ID:           configName,
				Name:         configName,
				Source:       "./infra/kubernetes/kubeconfig",
				Path:         fmt.Sprintf("/tmp/%s-%d", configName, runID),
				Variables:    vars,
				RemoteStates: remoteStates,
			},
		},
	)
	return s, nil
}

// Config returns steps to save a Kubernetes config
func (k *Kubernetes) Config(runID int32, path string) ([]Step, error) {
	k8sName := k.internalName("k8s")
	configName := k.internalName("kubeconfig")
	vars := make(map[string]string)
	vars["name"] = k.Name
	vars["kubernetes"] = k.Provider
	vars["region"] = k.Region
	vars["args"] = fmt.Sprintf(`["%s","%s"]`, k8sName, viper.GetString("aws-region"))
	vars["output_path"] = path
	remoteStates := make(map[string]string)
	remoteStates["k8s"] = k8sName
	return []Step{
		Step{
			&TerraformModule{
				ID:           configName,
				Name:         configName,
				Source:       "./infra/kubernetes/kubeconfig",
				Path:         fmt.Sprintf("/tmp/%s-%d", configName, runID),
				Variables:    vars,
				RemoteStates: remoteStates,
			},
		},
	}, nil
}

func (k *Kubernetes) internalName(module string) string {
	return fmt.Sprintf("%s-%s-%s-%s", k.Name, k.Region, k.Provider, module)
}
