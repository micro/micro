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
// Original source: github.com/micro/go-micro/v3/network/transport/http/options.go

package http

import (
	"context"
	"net/http"

	"github.com/micro/micro/v3/internal/network/transport"
)

// Handle registers the handler for the given pattern.
func Handle(pattern string, handler http.Handler) transport.Option {
	return func(o *transport.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		handlers, ok := o.Context.Value("http_handlers").(map[string]http.Handler)
		if !ok {
			handlers = make(map[string]http.Handler)
		}
		handlers[pattern] = handler
		o.Context = context.WithValue(o.Context, "http_handlers", handlers)
	}
}
