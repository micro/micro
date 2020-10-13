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
// Original source: github.com/micro/go-micro/v3/debug/profile/http/http.go

// Package http enables the http profiler
package http

import (
	"context"
	"net/http"
	"net/http/pprof"
	"sync"

	"github.com/micro/micro/v3/internal/debug/profile"
)

type httpProfile struct {
	sync.Mutex
	running bool
	server  *http.Server
}

var (
	DefaultAddress = ":6060"
)

// Start the profiler
func (h *httpProfile) Start() error {
	h.Lock()
	defer h.Unlock()

	if h.running {
		return nil
	}

	go func() {
		if err := h.server.ListenAndServe(); err != nil {
			h.Lock()
			h.running = false
			h.Unlock()
		}
	}()

	h.running = true

	return nil
}

// Stop the profiler
func (h *httpProfile) Stop() error {
	h.Lock()
	defer h.Unlock()

	if !h.running {
		return nil
	}

	h.running = false

	return h.server.Shutdown(context.TODO())
}

func (h *httpProfile) String() string {
	return "http"
}

func NewProfile(opts ...profile.Option) profile.Profile {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	return &httpProfile{
		server: &http.Server{
			Addr:    DefaultAddress,
			Handler: mux,
		},
	}
}
