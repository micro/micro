package infra

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Platform defines a complete platform
type Platform struct {
	Name    string
	Domain  string
	Gslb    string
	Kv      string
	Regions []struct {
		Provider string
		Region   string
		Control  []string
		Resource []string
		Network  []string
	}
}

// Steps generates an action plan from a Platform description
func (p *Platform) Steps() ([]Step, error) {
	// Not secure random, it doesn't matter as it's only to generate non colliding directory names
	rand.Seed(time.Now().UnixNano())
	runID := rand.Int31()
	var steps []Step
	// 1: Ensure Remote state is available
	steps = append(steps, Step{&RemoteState{ID: p.Name + "-check-remote-state", Name: p.Name + "-check-remote-state"}})

	// 2: Set up KV namespace
	steps = append(steps, Step{
		&TerraformModule{
			ID:     p.Name + "-global-kv",
			Name:   p.Name + "-global-kv",
			Source: "./infra/kv/" + p.Kv,
			Path:   fmt.Sprintf("/tmp/%s-%d", p.Name+"-kv", runID),
		},
	})

	for _, r := range p.Regions {
		// 2.1 Create Kubernetes cluster
		k := &Kubernetes{
			Name:     p.Name,
			Region:   r.Region,
			Provider: r.Provider,
		}
		cluster, err := k.Steps(runID)
		if err != nil {
			return steps, errors.Wrap(err, "Kubernetes cluster steps failed")
		}
		steps = append(steps, cluster...)

		// 2.2 Create namespaces
		vars := make(map[string]string)
		env := make(map[string]string)
		vars["control_namespace"] = strings.ToLower(fmt.Sprintf("%s-control", p.Name))
		vars["resource_namespace"] = strings.ToLower(fmt.Sprintf("%s-resource", p.Name))
		vars["network_namespace"] = strings.ToLower(fmt.Sprintf("%s-network", p.Name))
		env["KUBECONFIG"] = fmt.Sprintf("/tmp/%s-%s-%s-kubeconfig-%d/kubeconfig", p.Name, r.Region, r.Provider, runID)
		steps = append(steps, Step{
			&TerraformModule{
				ID:        p.Name + "-" + r.Region + "-" + r.Provider + "-namespaces",
				Name:      p.Name + "-" + r.Region + "-" + r.Provider + "-namespaces",
				Source:    "./infra/kubernetes/namespaces",
				Path:      fmt.Sprintf("/tmp/%s-%s-%s-namespaces-%d", p.Name, r.Region, r.Provider, runID),
				Variables: vars,
				Env:       env,
			},
		})

		// 2.3 Create shared resources
		vars = make(map[string]string)
		env = make(map[string]string)
		remoteStates := make(map[string]string)
		if r.Provider == "aws" {
			vars["in_aws"] = "true"
		} else {
			vars["in_aws"] = "false"
		}
		env["KUBECONFIG"] = fmt.Sprintf("/tmp/%s-%s-%s-kubeconfig-%d/kubeconfig", p.Name, r.Region, r.Provider, runID)
		remoteStates["namespaces"] = p.Name + "-" + r.Region + "-" + r.Provider + "-namespaces"
		steps = append(steps, Step{
			&TerraformModule{
				ID:           p.Name + "-" + r.Region + "-" + r.Provider + "-resource",
				Name:         p.Name + "-" + r.Region + "-" + r.Provider + "-resource",
				Source:       "./infra/resource",
				Path:         fmt.Sprintf("/tmp/%s-%s-%s-resource-%d", p.Name, r.Region, r.Provider, runID),
				Variables:    vars,
				Env:          env,
				RemoteStates: remoteStates,
			},
		})

		// 2.4 Create control plane
		vars = make(map[string]string)
		env = make(map[string]string)
		remoteStates = make(map[string]string)
		vars["domain_name"] = p.Domain
		env["KUBECONFIG"] = fmt.Sprintf("/tmp/%s-%s-%s-kubeconfig-%d/kubeconfig", p.Name, r.Region, r.Provider, runID)
		remoteStates["namespaces"] = p.Name + "-" + r.Region + "-" + r.Provider + "-namespaces"
		steps = append(steps, Step{
			&TerraformModule{
				ID:           p.Name + "-" + r.Region + "-" + r.Provider + "-control",
				Name:         p.Name + "-" + r.Region + "-" + r.Provider + "-control",
				Source:       "./infra/control",
				Path:         fmt.Sprintf("/tmp/%s-%s-%s-control-%d", p.Name, r.Region, r.Provider, runID),
				Variables:    vars,
				Env:          env,
				RemoteStates: remoteStates,
			},
		})

		// 2.5 Create network
		vars = make(map[string]string)
		env = make(map[string]string)
		remoteStates = make(map[string]string)
		vars["domain_name"] = p.Domain
		vars["cloudflare_account_id"] = "TODO"
		vars["cloudflare_dns_zone_id"] = "TODO"
		vars["cloudflare_api_token"] = "TODO"
		vars["region_slug"] = r.Region + "-" + r.Provider
		env["KUBECONFIG"] = fmt.Sprintf("/tmp/%s-%s-%s-kubeconfig-%d/kubeconfig", p.Name, r.Region, r.Provider, runID)
		remoteStates["namespaces"] = p.Name + "-" + r.Region + "-" + r.Provider + "-namespaces"
		remoteStates["kv"] = p.Name + "-global-kv"
		steps = append(steps, Step{
			&TerraformModule{
				ID:           p.Name + "-" + r.Region + "-" + r.Provider + "-network",
				Name:         p.Name + "-" + r.Region + "-" + r.Provider + "-network",
				Source:       "./infra/network",
				Path:         fmt.Sprintf("/tmp/%s-%s-%s-network-%d", p.Name, r.Region, r.Provider, runID),
				Variables:    vars,
				Env:          env,
				RemoteStates: remoteStates,
			},
		})
	}

	return steps, nil
}
