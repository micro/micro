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
// Original source: github.com/micro/go-micro/v3/network/network.go

// Package network implements micro network node
package network

import (
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/server"
)

var (
	// DefaultNetwork implementation
	DefaultNetwork Network
)

// Error is network node errors
type Error interface {
	// Count is current count of errors
	Count() int
	// Msg is last error message
	Msg() string
}

// Status is node status
type Status interface {
	// Error reports error status
	Error() Error
}

// Node is network node
type Node interface {
	// Id is node id
	Id() string
	// Address is node bind address
	Address() string
	// Peers returns node peers
	Peers() []Node
	// Network is the network node is in
	Network() Network
	// Status returns node status
	Status() Status
}

// Network is micro network
type Network interface {
	// Node is network node
	Node
	// Initialise options
	Init(...Option) error
	// Options returns the network options
	Options() Options
	// Name of the network
	Name() string
	// Connect starts the resolver and tunnel server
	Connect() error
	// Close stops the tunnel and resolving
	Close() error
	// Client is micro client
	Client() client.Client
	// Server is micro server
	Server() server.Server
}
