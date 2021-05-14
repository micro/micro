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
// Original source: github.com/micro/go-micro/v3/sync/sync.go

// Package sync is an interface for distributed synchronization
package sync

import (
	"errors"
	"time"
)

var (
	ErrLockTimeout = errors.New("lock timeout")
)

// Sync is an interface for distributed synchronization
type Sync interface {
	// Initialise options
	Init(...Option) error
	// Return the options
	Options() Options
	// Elect a leader
	Leader(id string, opts ...LeaderOption) (Leader, error)
	// Lock acquires a lock
	Lock(id string, opts ...LockOption) error
	// Unlock releases a lock
	Unlock(id string) error
	// Sync implementation
	String() string
}

// Leader provides leadership election
type Leader interface {
	// resign leadership
	Resign() error
	// status returns when leadership is lost
	Status() chan bool
}

type Options struct {
	Nodes  []string
	Prefix string
}

type Option func(o *Options)

type LeaderOptions struct{}

type LeaderOption func(o *LeaderOptions)

type LockOptions struct {
	TTL  time.Duration
	Wait time.Duration
}

type LockOption func(o *LockOptions)
