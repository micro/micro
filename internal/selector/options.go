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
// Original source: github.com/micro/go-micro/v3/selector/options.go

package selector

// Options used to configure a selector
type Options struct{}

// Option updates the options
type Option func(*Options)

// SelectOptions used to configure selection
type SelectOptions struct{}

// SelectOption updates the select options
type SelectOption func(*SelectOptions)

// NewSelectOptions parses select options
func NewSelectOptions(opts ...SelectOption) SelectOptions {
	var options SelectOptions
	for _, o := range opts {
		o(&options)
	}

	return options
}
