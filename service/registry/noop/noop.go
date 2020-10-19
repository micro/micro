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
// Original source: github.com/micro/go-micro/v3/registry/noop/noop.go

// Package noop is a registry which does nothing
package noop

import (
	"errors"

	"github.com/micro/micro/v3/service/registry"
)

type noopRegistry struct{}

func (n *noopRegistry) Init(o ...registry.Option) error {
	return nil
}

func (n *noopRegistry) Options() registry.Options {
	return registry.Options{}
}

func (n *noopRegistry) Register(*registry.Service, ...registry.RegisterOption) error {
	return nil
}

func (n *noopRegistry) Deregister(*registry.Service, ...registry.DeregisterOption) error {
	return nil
}

func (n *noopRegistry) GetService(s string, o ...registry.GetOption) ([]*registry.Service, error) {
	return []*registry.Service{}, nil
}

func (n *noopRegistry) ListServices(...registry.ListOption) ([]*registry.Service, error) {
	return []*registry.Service{}, nil
}
func (n *noopRegistry) Watch(...registry.WatchOption) (registry.Watcher, error) {
	return nil, errors.New("not implemented")
}

func (n *noopRegistry) String() string {
	return "noop"
}

// NewRegistry returns a new noop registry
func NewRegistry(opts ...registry.Option) registry.Registry {
	return new(noopRegistry)
}
