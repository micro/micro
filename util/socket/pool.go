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
// Original source: github.com/micro/go-micro/v3/util/socket/pool.go

package socket

import (
	"sync"
)

type Pool struct {
	sync.RWMutex
	pool map[string]*Socket
}

func (p *Pool) Get(id string) (*Socket, bool) {
	// attempt to get existing socket
	p.RLock()
	socket, ok := p.pool[id]
	if ok {
		p.RUnlock()
		return socket, ok
	}
	p.RUnlock()

	// save socket
	p.Lock()
	defer p.Unlock()
	// double checked locking
	socket, ok = p.pool[id]
	if ok {
		return socket, ok
	}
	// create new socket
	socket = New(id)
	p.pool[id] = socket

	// return socket
	return socket, false
}

func (p *Pool) Release(s *Socket) {
	p.Lock()
	defer p.Unlock()

	// close the socket
	s.Close()
	delete(p.pool, s.id)
}

// Close the pool and delete all the sockets
func (p *Pool) Close() {
	p.Lock()
	defer p.Unlock()
	for id, sock := range p.pool {
		sock.Close()
		delete(p.pool, id)
	}
}

// NewPool returns a new socket pool
func NewPool() *Pool {
	return &Pool{
		pool: make(map[string]*Socket),
	}
}
