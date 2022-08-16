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
	"net"
	"net/http"

	"github.com/micro/micro/v3/service/network/transport"
)

type netListener struct{}

// getNetListener Get net.Listener from ListenOptions
func getNetListener(o *transport.ListenOptions) net.Listener {
	if o.Context == nil {
		return nil
	}

	if l, ok := o.Context.Value(netListener{}).(net.Listener); ok && l != nil {
		return l
	}

	return nil
}

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

// NetListener Set net.Listener for httpTransport
func NetListener(customListener net.Listener) transport.ListenOption {
	return func(o *transport.ListenOptions) {
		if customListener == nil {
			return
		}
		if o.Context == nil {
			o.Context = context.TODO()
		}
		o.Context = context.WithValue(o.Context, netListener{}, customListener)
	}
}
