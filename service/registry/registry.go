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
// Original source: github.com/micro/go-micro/v3/registry/registry.go

// Package registry is the micro registry
package registry

import (
	"errors"
)

var (
	// DefaultRegistry implementation
	DefaultRegistry Registry
	// ErrNotFound error when GetService is called
	ErrNotFound = errors.New("service not found")
	// ErrWatcherStopped error when watcher is stopped
	ErrWatcherStopped = errors.New("watcher stopped")
)

const (
	// WildcardDomain indicates any domain
	WildcardDomain = "*"
	// DefaultDomain to use if none was provided in options
	DefaultDomain = "micro"
)

// Registry provides an interface for service discovery
// and an abstraction over varying implementations
// {consul, etcd, zookeeper, ...}
type Registry interface {
	Init(...Option) error
	Options() Options
	Register(*Service, ...RegisterOption) error
	Deregister(*Service, ...DeregisterOption) error
	GetService(string, ...GetOption) ([]*Service, error)
	ListServices(...ListOption) ([]*Service, error)
	Watch(...WatchOption) (Watcher, error)
	String() string
}

type Service struct {
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Metadata  map[string]string `json:"metadata"`
	Endpoints []*Endpoint       `json:"endpoints"`
	Nodes     []*Node           `json:"nodes"`
}

type Node struct {
	Id       string            `json:"id"`
	Address  string            `json:"address"`
	Metadata map[string]string `json:"metadata"`
}

type Endpoint struct {
	Name     string            `json:"name"`
	Request  *Value            `json:"request"`
	Response *Value            `json:"response"`
	Metadata map[string]string `json:"metadata"`
}

type Value struct {
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Values []*Value `json:"values"`
}

// GetService from the registry
func GetService(service string) ([]*Service, error) {
	return DefaultRegistry.GetService(service)
}

// ListServices in the registry
func ListServices() ([]*Service, error) {
	return DefaultRegistry.ListServices()
}

// Watch the registry for updates
func Watch() (Watcher, error) {
	return DefaultRegistry.Watch()
}
