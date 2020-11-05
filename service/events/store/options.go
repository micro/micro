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
// Original source: github.com/micro/go-micro/v3/events/store/options.go

package store

import (
	"time"

	"github.com/micro/micro/v3/service/store"
)

type Options struct {
	Store store.Store
	TTL   time.Duration
}

type Option func(o *Options)

// WithStore sets the underlying store to use
func WithStore(s store.Store) Option {
	return func(o *Options) {
		o.Store = s
	}
}

// WithTTL sets the default TTL
func WithTTL(ttl time.Duration) Option {
	return func(o *Options) {
		o.TTL = ttl
	}
}
