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
// Original source: github.com/micro/go-micro/v3/store/file/options.go

package file

import (
	"context"

	"github.com/micro/micro/v3/service/store"
)

type dirKey struct{}

// WithDir sets the directory to store the files in
func WithDir(dir string) store.StoreOption {
	return func(o *store.StoreOptions) {
		if o.Context == nil {
			o.Context = context.WithValue(context.Background(), dirKey{}, dir)
		} else {
			o.Context = context.WithValue(o.Context, dirKey{}, dir)
		}
	}
}
