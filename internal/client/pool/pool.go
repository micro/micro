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
// Original source: github.com/micro/go-micro/v3/util/pool/pool.go

// Package pool is a connection pool
package pool

import (
	"time"

	"github.com/micro/micro/v3/internal/network/transport"
)

// Pool is an interface for connection pooling
type Pool interface {
	// Close the pool
	Close() error
	// Get a connection
	Get(addr string, opts ...transport.DialOption) (Conn, error)
	// Release the connection
	Release(c Conn, status error) error
}

type Conn interface {
	// unique id of connection
	Id() string
	// time it was created
	Created() time.Time
	// embedded connection
	transport.Client
}

func NewPool(opts ...Option) Pool {
	var options Options
	for _, o := range opts {
		o(&options)
	}
	return newPool(options)
}
