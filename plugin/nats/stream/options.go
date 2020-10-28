// Copyright 2020 Asim Aslam
//
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
// Original source: github.com/micro/go-plugins/v3/events/stream/nats/options.go

package nats

import "crypto/tls"

// Options which are used to configure the nats stream
type Options struct {
	ClusterID string
	ClientID  string
	Address   string
	TLSConfig *tls.Config
}

// Option is a function which configures options
type Option func(o *Options)

// ClusterID sets the cluster id for the nats connection
func ClusterID(id string) Option {
	return func(o *Options) {
		o.ClusterID = id
	}
}

// ClientID sets the client id for the nats connection
func ClientID(id string) Option {
	return func(o *Options) {
		o.ClientID = id
	}
}

// Address of the nats cluster
func Address(addr string) Option {
	return func(o *Options) {
		o.Address = addr
	}
}

// TLSConfig to use when connecting to the cluster
func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}
