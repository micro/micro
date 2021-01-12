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
// Original source: github.com/micro/go-micro/v3/runtime/resource.go

// Package runtime is a service runtime manager
package runtime

import "fmt"

const (
	TypeNamespace     = "namespace"
	TypeNetworkPolicy = "networkpolicy"
	TypeResourceQuota = "resourcequota"
	TypeService       = "service"
)

var (
	defaultRequests = &Resources{
		Mem:  2000, // 2GB
		Disk: 5000, // 5GB
	}
	defaultLimits = &Resources{
		Mem:  4000,  // 4GB
		Disk: 10000, // 10GB
	}
)

// Resource represents any resource handled by runtime
type Resource interface {
	String() string
	Type() string
}

// Namespace represents a logical namespace for organising resources
type Namespace struct {
	// Name of the namespace
	Name string
}

// NewNamespace mints a new namespace
func NewNamespace(name string) (*Namespace, error) {
	if name == "" {
		return nil, ErrInvalidResource
	}
	return &Namespace{
		Name: name,
	}, nil
}

// String implements Resource
func (r *Namespace) String() string {
	return r.Name
}

// Type implements Resource
func (*Namespace) Type() string {
	return TypeNamespace
}

// NetworkPolicy represents an ACL of label pairs allowing ignress to a namespace
type NetworkPolicy struct {
	// The labels allowed ingress by this policy
	AllowedLabels map[string]string
	// Name of the network policy
	Name string
	// Namespace the network policy belongs to
	Namespace string
}

// NewNetworkPolicy mints a new networkpolicy
func NewNetworkPolicy(name, namespace string, allowedLabels map[string]string) (*NetworkPolicy, error) {
	if name == "" || namespace == "" {
		return nil, ErrInvalidResource
	}
	if allowedLabels == nil {
		allowedLabels = map[string]string{
			"origin": "micro",
		}
	}
	return &NetworkPolicy{
		AllowedLabels: allowedLabels,
		Name:          name,
		Namespace:     namespace,
	}, nil
}

// String implements Resource
func (r *NetworkPolicy) String() string {
	return fmt.Sprintf("%s.%s", r.Namespace, r.Name)
}

// Type implements Resource
func (*NetworkPolicy) Type() string {
	return TypeNetworkPolicy
}

// ResourceQuota represents an ACL of label pairs allowing ignress to a namespace
type ResourceQuota struct {
	// Name of the resource quota
	Name string
	// Namespace the resource quota belongs to
	Namespace string
	// Quota for resource REQUESTS
	Requests *Resources
	// Quota for resource LIMITS
	Limits *Resources
}

// NewResourceQuota mints a new resourcequota
func NewResourceQuota(name, namespace string, requests, limits *Resources) (*ResourceQuota, error) {
	if name == "" || namespace == "" {
		return nil, ErrInvalidResource
	}

	rq := &ResourceQuota{
		Name:      name,
		Namespace: namespace,
		Requests:  defaultRequests,
		Limits:    defaultLimits,
	}

	if requests != nil {
		rq.Requests = requests
	}

	if limits != nil {
		rq.Limits = limits
	}

	return rq, nil
}

// String implements Resource
func (r *ResourceQuota) String() string {
	return fmt.Sprintf("%s.%s", r.Namespace, r.Name)
}

// Type implements Resource
func (*ResourceQuota) Type() string {
	return TypeResourceQuota
}

// Service represents a Micro service running within a namespace
type Service struct {
	// Name of the service
	Name string
	// Version of the service
	Version string
	// url location of source
	Source string
	// Metadata stores metadata
	Metadata map[string]string
	// Status of the service
	Status ServiceStatus
}

// NewService mints a new service
func NewService(name, version string) (*Service, error) {
	if name == "" {
		return nil, ErrInvalidResource
	}
	if version == "" {
		version = "latest"
	}
	return &Service{
		Name:    name,
		Version: version,
	}, nil
}

// String implements Resource
func (r *Service) String() string {
	return fmt.Sprintf("service://%s@%s:%s", r.Metadata["namespace"], r.Name, r.Version)
}

// Type implements Resource
func (*Service) Type() string {
	return TypeService
}
