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
// Original source: github.com/micro/go-micro/v3/network/resolver/resolver.go

// Package resolver resolves network names to addresses
package resolver

// Resolver is network resolver. It's used to find network nodes
// via the name to connect to. This is done based on Network.Name().
// Before we can be part of any network, we have to connect to it.
type Resolver interface {
	// Resolve returns a list of addresses for a name
	Resolve(name string) ([]*Record, error)
}

// A resolved record
type Record struct {
	Address  string `json:"address"`
	Priority int64  `json:"priority"`
}
