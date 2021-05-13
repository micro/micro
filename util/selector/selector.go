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
// Original source: github.com/micro/go-micro/v3/selector/selector.go

// Package selector is for node selection and load balancing
package selector

import (
	"errors"
)

var (
	// ErrNoneAvailable is returned by select when no routes were provided to select from
	ErrNoneAvailable = errors.New("none available")
)

// Selector selects a route from a pool
type Selector interface {
	// Select a route from the pool using the strategy
	Select([]string, ...SelectOption) (Next, error)
	// Record the error returned from a route to inform future selection
	Record(string, error) error
	// Reset the selector
	Reset() error
	// String returns the name of the selector
	String() string
}

// Next returns the next node
type Next func() string
