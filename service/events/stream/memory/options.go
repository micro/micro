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
// Original source: github.com/micro/go-micro/v3/events/stream/memory/options.go

package memory

import "github.com/micro/micro/v3/service/store"

// Options which are used to configure the in-memory stream
type Options struct {
	Store store.Store
}

// Option is a function which configures options
type Option func(o *Options)

// Store sets the store to use
func Store(s store.Store) Option {
	return func(o *Options) {
		o.Store = s
	}
}
