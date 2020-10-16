// Copyright 2020 Asim Aslam
//
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
// Original source: github.com/micro/go-micro/v3/api/server/cors/cors.go

package cors

import (
	"net/http"
)

// CombinedCORSHandler wraps a server and provides CORS headers
func CombinedCORSHandler(h http.Handler) http.Handler {
	return corsHandler{h}
}

type corsHandler struct {
	handler http.Handler
}

func (c corsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w, r)

	if r.Method == "OPTIONS" {
		return
	}

	c.handler.ServeHTTP(w, r)
}

// SetHeaders sets the CORS headers
func SetHeaders(w http.ResponseWriter, r *http.Request) {
	set := func(w http.ResponseWriter, k, v string) {
		if v := w.Header().Get(k); len(v) > 0 {
			return
		}
		w.Header().Set(k, v)
	}

	if origin := r.Header.Get("Origin"); len(origin) > 0 {
		set(w, "Access-Control-Allow-Origin", origin)
	} else {
		set(w, "Access-Control-Allow-Origin", "*")
	}

	set(w, "Access-Control-Allow-Credentials", "true")
	set(w, "Access-Control-Allow-Methods", "POST, PATCH, GET, OPTIONS, PUT, DELETE")
	set(w, "Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Micro-Namespace")
}
