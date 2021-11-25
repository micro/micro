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
// Original source: github.com/micro/go-micro/v3/debug/profile/profile.go

// Package profile is for profilers
package profile

type Profile interface {
	// Start the profiler
	Start() error
	// Stop the profiler
	Stop() error
	// Name of the profiler
	String() string
}

var (
	DefaultProfile Profile = new(noop)
)

type noop struct{}

func (p *noop) Start() error {
	return nil
}

func (p *noop) Stop() error {
	return nil
}

func (p *noop) String() string {
	return "noop"
}

type Options struct {
	// Name to use for the profile
	Name string
}

type Option func(o *Options)

// Name of the profile
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}
