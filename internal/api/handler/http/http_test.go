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
// Original source: github.com/micro/go-micro/v3/api/handler/http/http_test.go

package http

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/micro/micro/v3/internal/api/handler"
	"github.com/micro/micro/v3/internal/api/resolver"
	"github.com/micro/micro/v3/internal/api/resolver/vpath"
	"github.com/micro/micro/v3/internal/api/router"
	regRouter "github.com/micro/micro/v3/internal/api/router/registry"
	"github.com/micro/micro/v3/service/registry"
	"github.com/micro/micro/v3/service/registry/memory"
)

func testHttp(t *testing.T, path, service, ns string) {
	r := memory.NewRegistry()

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	s := &registry.Service{
		Name: service,
		Nodes: []*registry.Node{
			{
				Id:      service + "-1",
				Address: l.Addr().String(),
			},
		},
	}

	r.Register(s)
	defer r.Deregister(s)

	// setup the test handler
	m := http.NewServeMux()
	m.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`you got served`))
	})

	// start http test serve
	go http.Serve(l, m)

	// create new request and writer
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	// initialise the handler
	rt := regRouter.NewRouter(
		router.WithHandler("http"),
		router.WithRegistry(r),
		router.WithResolver(vpath.NewResolver(
			resolver.WithServicePrefix(ns),
		)),
	)

	p := NewHandler(handler.WithRouter(rt))

	// execute the handler
	p.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("Expected 200 response got %d %s", w.Code, w.Body.String())
	}

	if w.Body.String() != "you got served" {
		t.Fatalf("Expected body: you got served. Got: %s", w.Body.String())
	}
}

func TestHttpHandler(t *testing.T) {
	testData := []struct {
		path      string
		service   string
		namespace string
	}{
		{
			"/test/foo",
			"go.micro.api.test",
			"go.micro.api",
		},
		{
			"/test/foo/baz",
			"go.micro.api.test",
			"go.micro.api",
		},
		{
			"/v1/foo",
			"go.micro.api.v1.foo",
			"go.micro.api",
		},
		{
			"/v1/foo/bar",
			"go.micro.api.v1.foo",
			"go.micro.api",
		},
		{
			"/v2/baz",
			"go.micro.api.v2.baz",
			"go.micro.api",
		},
		{
			"/v2/baz/bar",
			"go.micro.api.v2.baz",
			"go.micro.api",
		},
		{
			"/v2/baz/bar",
			"v2.baz",
			"",
		},
	}

	for _, d := range testData {
		t.Run(d.service, func(t *testing.T) {
			testHttp(t, d.path, d.service, d.namespace)
		})
	}
}
