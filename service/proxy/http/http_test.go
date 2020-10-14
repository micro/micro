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
// Original source: github.com/micro/go-micro/v3/proxy/http/http_test.go

package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/micro/micro/v3/service/client"
	cmucp "github.com/micro/micro/v3/service/client/mucp"
	"github.com/micro/micro/v3/service/registry/memory"
	"github.com/micro/micro/v3/service/router"
	"github.com/micro/micro/v3/service/router/registry"
	"github.com/micro/micro/v3/service/server"
	"github.com/micro/micro/v3/service/server/mucp"
)

type testHandler struct{}

func (t *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"hello": "world"}`))
}

func TestHTTPProxy(t *testing.T) {
	c, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	addr := c.Addr().String()

	url := fmt.Sprintf("http://%s", addr)

	testCases := []struct {
		// http endpoint to call e.g /foo/bar
		httpEp string
		// rpc endpoint called e.g Foo.Bar
		rpcEp string
		// should be an error
		err bool
	}{
		{"/", "Foo.Bar", false},
		{"/", "Foo.Baz", false},
		{"/helloworld", "Hello.World", true},
	}

	// handler
	http.Handle("/", new(testHandler))

	// new proxy
	p := NewSingleHostProxy(url)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reg := memory.NewRegistry()
	rtr := registry.NewRouter(
		router.Registry(reg),
	)

	// new micro service
	service := mucp.NewServer(
		server.Context(ctx),
		server.Name("foobar"),
		server.Registry(reg),
		server.WithRouter(p),
	)

	service.Start()
	defer service.Stop()

	// run service
	// server
	go http.Serve(c, nil)

	cl := cmucp.NewClient(
		client.Router(rtr),
	)

	for _, test := range testCases {
		req := cl.NewRequest("foobar", test.rpcEp, map[string]string{"foo": "bar"}, client.WithContentType("application/json"))
		var rsp map[string]string
		err := cl.Call(ctx, req, &rsp)
		if err != nil && test.err == false {
			t.Fatal(err)
		}
		if v := rsp["hello"]; v != "world" {
			t.Fatalf("Expected hello world got %s from %s", v, test.rpcEp)
		}
	}
}
