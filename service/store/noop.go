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
// Original source: github.com/micro/go-micro/v3/store/noop.go

package store

type noopStore struct{}

func (n *noopStore) Init(opts ...Option) error {
	return nil
}

func (n *noopStore) Options() Options {
	return Options{}
}

func (n *noopStore) String() string {
	return "noop"
}

func (n *noopStore) Read(key string, opts ...ReadOption) ([]*Record, error) {
	return []*Record{}, nil
}

func (n *noopStore) Write(r *Record, opts ...WriteOption) error {
	return nil
}

func (n *noopStore) Delete(key string, opts ...DeleteOption) error {
	return nil
}

func (n *noopStore) List(opts ...ListOption) ([]string, error) {
	return []string{}, nil
}

func (n *noopStore) Close() error {
	return nil
}
