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
// Original source: github.com/micro/go-micro/v3/debug/trace/options.go

package trace

type Options struct {
	// Size is the size of ring buffer
	Size int
}

type Option func(o *Options)

type ReadOptions struct {
	// Trace id
	Trace string
}

type ReadOption func(o *ReadOptions)

// Read the given trace
func ReadTrace(t string) ReadOption {
	return func(o *ReadOptions) {
		o.Trace = t
	}
}

const (
	// DefaultSize of the buffer
	DefaultSize = 64
)

// DefaultOptions returns default options
func DefaultOptions() Options {
	return Options{
		Size: DefaultSize,
	}
}
